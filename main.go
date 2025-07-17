package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/joho/godotenv"
)

// EDINET API レスポンス用構造体
// ... 既存のDocInfo, DocumentListResponse など ...
type DocInfo struct {
	DocID       string `json:"docID"`
	FilerName   string `json:"filerName"`
	DocTypeCode string `json:"docTypeCode"`
	XbrlFlag    string `json:"xbrlFlag"`
	SecCode     string `json:"secCode"`
}
type DocumentListResponse struct {
	Results []DocInfo `json:"results"`
}

func main() {
	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .envファイルを読み込めませんでした: %v", err)
	}
	
	// ここでパラメータ指定
	startDate := "2025-07-10"
	endDate := "2025-07-16"
	targetSecCode := "40260"
	apiKey := os.Getenv("EDINET_API_KEY") // APIキーは環境変数で
	

	// APIキーの確認
	if apiKey == "" {
		log.Fatal("EDINET_API_KEYが設定されていません。.envファイルを確認してください。")
	}

	const layout = "2006-01-02"
	start, _ := time.Parse(layout, startDate)
	end, _ := time.Parse(layout, endDate)

	// 日本語ヘッダー（四半期情報を追加）
	japaneseHeaders := []string{
		"日付","証券コード","会社名","文書タイプ","会計期間",
		"売上高","売上総利益","営業利益","経常利益","税引前当期純利益","当期純利益","1株当たり当期純利益","総資産","流動資産","固定資産","総負債","流動負債","固定負債","純資産","資本金","利益剰余金","営業CF","投資CF","財務CF","現金及び現金同等物","1株当たり純資産","自己資本比率","配当金",
	}
	financialTags := []string{
		"jppfs_cor:NetSales","jppfs_cor:GrossProfit","jppfs_cor:OperatingIncome","jppfs_cor:OrdinaryIncome","jppfs_cor:IncomeBeforeIncomeTaxes","jppfs_cor:ProfitLoss","jppfs_cor:BasicEarningsLossPerShareSummaryOfBusinessResults","jppfs_cor:TotalAssets","jppfs_cor:CurrentAssets","jppfs_cor:NoncurrentAssets","jppfs_cor:Liabilities","jppfs_cor:CurrentLiabilities","jppfs_cor:NoncurrentLiabilities","jppfs_cor:NetAssets","jppfs_cor:CapitalStock","jppfs_cor:RetainedEarnings","jppfs_cor:NetCashProvidedByUsedInOperatingActivities","jppfs_cor:NetCashProvidedByUsedInInvestmentActivities","jppfs_cor:NetCashProvidedByUsedInFinancingActivities","jppfs_cor:CashAndCashEquivalents","jppfs_cor:NetAssetsPerShareSummaryOfBusinessResults","jppfs_cor:EquityToAssetRatioSummaryOfBusinessResults","jppfs_cor:DividendsFromSurplus",
	}

	out, err := os.Create("xbrl_financial_items.csv")
	if err != nil {
		log.Fatalf("CSV作成エラー: %v", err)
	}
	defer out.Close()
	writer := csv.NewWriter(out)
	defer writer.Flush()
	writer.Write(japaneseHeaders)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(layout)
		url := fmt.Sprintf("https://api.edinet-fsa.go.jp/api/v2/documents.json?date=%s&type=2&limit=100", dateStr)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("リクエスト作成エラー: %v", err)
		}
		req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
		req.Header.Set("Accept", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("API取得エラー: %v", err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		var docList DocumentListResponse
		if err := json.Unmarshal(body, &docList); err != nil {
			log.Fatalf("JSONパースエラー: %v", err)
		}
		for _, doc := range docList.Results {
			if (doc.DocTypeCode == "120" || doc.DocTypeCode == "130") && doc.SecCode == targetSecCode && doc.XbrlFlag == "1" {
				// XBRL zipダウンロード
				zipURL := fmt.Sprintf("https://api.edinet-fsa.go.jp/api/v2/documents/%s?type=1", doc.DocID)
				zipFile := doc.DocID + ".zip"
				zreq, _ := http.NewRequest("GET", zipURL, nil)
				zreq.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
				zreq.Header.Set("Accept", "application/zip")
				zresp, err := client.Do(zreq)
				if err != nil {
					log.Printf("zipダウンロード失敗: %v", err)
					continue
				}
				zdata, _ := ioutil.ReadAll(zresp.Body)
				zresp.Body.Close()
				ioutil.WriteFile(zipFile, zdata, 0644)
				// zip解凍→PublicDocのxbrlファイル特定
				xbrlPath, err := extractPublicDocXBRL(zipFile)
				if err != nil {
					log.Printf("XBRL抽出失敗: %v", err)
					continue
				}
				values, err := parseAllXBRL(xbrlPath)
				if err != nil {
					log.Printf("XBRLパース失敗: %v", err)
					continue
				}
				// 文書タイプの日本語名を取得
				docTypeName := getDocTypeName(doc.DocTypeCode)
				
				// 会計期間の情報を抽出
				startDate, endDate := extractAccountingPeriod(xbrlPath)
				quarterInfo := getQuarterInfo(startDate, endDate)
				
				row := []string{
					dateStr,                    // 日付
					doc.SecCode,                // 証券コード
					doc.FilerName,              // 会社名
					docTypeName,                // 文書タイプ
					quarterInfo,                // 会計期間
				}
				
				for _, tag := range financialTags {
					found := ""
					tagOnly := tag
					if strings.Contains(tag, ":") {
						tagOnly = strings.Split(tag, ":")[1]
					}
					for k, v := range values {
						if (strings.Contains(k, ":"+tagOnly+"|") || strings.HasSuffix(k, ":"+tagOnly)) &&
							!strings.Contains(k, "TextBlock") {
							found = v
							break
						}
					}
					row = append(row, found)
				}
				writer.Write(row)
				os.Remove(zipFile)
				os.Remove(xbrlPath)
			}
		}
	}
	fmt.Println("xbrl_financial_items.csv に主要財務項目を出力しました。")
}

// zipからPublicDocのxbrlファイルパスを返す
func extractPublicDocXBRL(zipFile string) (string, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer r.Close()
	for _, f := range r.File {
		if strings.Contains(f.Name, "PublicDoc") && strings.HasSuffix(f.Name, ".xbrl") {
			outPath := filepath.Base(f.Name)
			out, err := os.Create(outPath)
			if err != nil {
				return "", err
			}
			in, err := f.Open()
			if err != nil {
				return "", err
			}
			io.Copy(out, in)
			in.Close()
			out.Close()
			return outPath, nil
		}
	}
	return "", fmt.Errorf("PublicDocのxbrlファイルが見つかりません")
}

// XBRL(XML)から全ての値を抽出（タグ名＋contextRef＋unitRefでユニーク化）
func parseAllXBRL(xbrlPath string) (map[string]string, error) {
	file, err := os.Open(xbrlPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	values := make(map[string]string)
	var currentKey string
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			tagName := se.Name.Local
			if se.Name.Space != "" {
				tagName = se.Name.Space + ":" + se.Name.Local
			}
			contextRef := ""
			unitRef := ""
			for _, attr := range se.Attr {
				if attr.Name.Local == "contextRef" {
					contextRef = attr.Value
				}
				if attr.Name.Local == "unitRef" {
					unitRef = attr.Value
				}
			}
			key := tagName
			if contextRef != "" {
				key += "|contextRef=" + contextRef
			}
			if unitRef != "" {
				key += "|unitRef=" + unitRef
			}
			currentKey = key
		case xml.CharData:
			val := strings.TrimSpace(string(se))
			if val != "" && currentKey != "" {
				values[currentKey] = val
			}
		case xml.EndElement:
			currentKey = ""
		}
	}
	return values, nil
}

// 文書タイプコードを日本語名に変換
func getDocTypeName(docTypeCode string) string {
	switch docTypeCode {
	case "120":
		return "有価証券報告書"
	case "130":
		return "四半期報告書"
	case "140":
		return "有価証券届出書"
	case "150":
		return "発行登録書"
	case "160":
		return "発行登録追補書類"
	case "170":
		return "訂正有価証券報告書"
	case "180":
		return "訂正四半期報告書"
	case "190":
		return "訂正有価証券届出書"
	case "200":
		return "有価証券報告書（外国会社）"
	case "210":
		return "四半期報告書（外国会社）"
	default:
		return "その他"
	}
}

// XBRLファイルから会計期間の情報を抽出
func extractAccountingPeriod(xbrlPath string) (string, string) {
	file, err := os.Open(xbrlPath)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var startDate, endDate string
	
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		
		switch se := tok.(type) {
		case xml.StartElement:
			// startDateまたはendDateタグを探す
			if se.Name.Local == "startDate" || se.Name.Local == "endDate" {
				var date string
				for {
					t, err := decoder.Token()
					if err != nil {
						break
					}
					if cd, ok := t.(xml.CharData); ok {
						date = strings.TrimSpace(string(cd))
						break
					}
					if _, ok := t.(xml.EndElement); ok {
						break
					}
				}
				
				if se.Name.Local == "startDate" {
					startDate = date
				} else if se.Name.Local == "endDate" {
					endDate = date
				}
			}
		}
		
		// 両方の日付が見つかったら終了
		if startDate != "" && endDate != "" {
			break
		}
	}
	
	return startDate, endDate
}

// 会計期間から四半期情報を生成
func getQuarterInfo(startDate, endDate string) string {
	if startDate == "" || endDate == "" {
		return "不明"
	}
	
	// 日付をパース
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return "不明"
	}
	
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return "不明"
	}
	
	// 期間の長さを計算
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	
	// 年次報告書（約1年）か四半期報告書（約3ヶ月）かを判定
	if days > 300 { // 年次報告書
		year := end.Year()
		// 3月決算の場合は前年度
		if end.Month() == 3 {
			year = year - 1
		}
		return fmt.Sprintf("%d年度", year)
	} else { // 四半期報告書
		year := end.Year()
		endMonth := end.Month()
		var quarter string
		
		switch {
		case endMonth <= 3:
			quarter = "Q4"
			year = year - 1
		case endMonth <= 6:
			quarter = "Q1"
		case endMonth <= 9:
			quarter = "Q2"
		case endMonth <= 12:
			quarter = "Q3"
		}
		
		return fmt.Sprintf("%d%s", year, quarter)
	}
}