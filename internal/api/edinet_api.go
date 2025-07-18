package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"edinet-api-test/internal/models"
)

// EdinetAPI EDINET APIクライアント
type EdinetAPI struct {
	client *http.Client
	apiKey string
}

// NewEdinetAPI 新しいEDINET APIクライアントを作成
func NewEdinetAPI(apiKey string) *EdinetAPI {
	return &EdinetAPI{
		client: &http.Client{},
		apiKey: apiKey,
	}
}

// GetDocuments 指定日の文書一覧を取得
func (e *EdinetAPI) GetDocuments(date time.Time) (*models.DocumentListResponse, error) {
	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("https://api.edinet-fsa.go.jp/api/v2/documents.json?date=%s&type=2&limit=100", dateStr)
	
	fmt.Printf("API URL: %s\n", url)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}
	
	req.Header.Set("Ocp-Apim-Subscription-Key", e.apiKey)
	req.Header.Set("Accept", "application/json")
	
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API取得エラー: %v", err)
	}
	defer resp.Body.Close()
	
	fmt.Printf("API レスポンスステータス: %d\n", resp.StatusCode)
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}
	
	fmt.Printf("API レスポンスサイズ: %d bytes\n", len(body))
	if len(body) < 1000 { // レスポンスが小さい場合は内容を表示
		fmt.Printf("API レスポンス内容: %s\n", string(body))
	}
	
	var docList models.DocumentListResponse
	if err := json.Unmarshal(body, &docList); err != nil {
		return nil, fmt.Errorf("JSONパースエラー: %v", err)
	}
	
	fmt.Printf("取得された文書数: %d\n", len(docList.Results))
	
	return &docList, nil
}

// DownloadXBRLZip XBRLファイルのZIPをダウンロード
func (e *EdinetAPI) DownloadXBRLZip(docID string) ([]byte, error) {
	zipURL := fmt.Sprintf("https://api.edinet-fsa.go.jp/api/v2/documents/%s?type=1", docID)
	
	req, err := http.NewRequest("GET", zipURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ZIPリクエスト作成エラー: %v", err)
	}
	
	req.Header.Set("Ocp-Apim-Subscription-Key", e.apiKey)
	req.Header.Set("Accept", "application/zip")
	
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ZIPダウンロード失敗: %v", err)
	}
	defer resp.Body.Close()
	
	zdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ZIPデータ読み込みエラー: %v", err)
	}
	
	return zdata, nil
}

// FilterDocuments 文書をフィルタリング
func FilterDocuments(docs []models.DocInfo, targetSecCode string, quarterOnly bool) []models.DocInfo {
	var filtered []models.DocInfo
	fmt.Printf("フィルタリング開始: 全%d件の文書を処理\n", len(docs))
	
	for i, doc := range docs {
		fmt.Printf("文書[%d]: DocID=%s, DocTypeCode=%s, SecCode=%s, XbrlFlag=%s, FilerName=%s\n", 
			i, doc.DocID, doc.DocTypeCode, doc.SecCode, doc.XbrlFlag, doc.FilerName)
		
		// 証券コードが指定されている場合は証券コードもチェック
		if targetSecCode != "" {
			if quarterOnly {
				// 四半期報告書のみ
				if doc.DocTypeCode == "130" && 
				   doc.SecCode == targetSecCode && 
				   doc.XbrlFlag == "1" {
					fmt.Printf("  → 四半期報告書として追加\n")
					filtered = append(filtered, doc)
				}
			} else {
				// 有価証券報告書と四半期報告書
				if (doc.DocTypeCode == "120" || doc.DocTypeCode == "130") && 
				   doc.SecCode == targetSecCode && 
				   doc.XbrlFlag == "1" {
					fmt.Printf("  → 対象文書として追加\n")
					filtered = append(filtered, doc)
				}
			}
		} else {
			// 証券コードが指定されていない場合
			if quarterOnly {
				// 四半期報告書のみ
				if doc.DocTypeCode == "130" && 
				   doc.XbrlFlag == "1" {
					fmt.Printf("  → 四半期報告書として追加\n")
					filtered = append(filtered, doc)
				}
			} else {
				// 有価証券報告書と四半期報告書
				if (doc.DocTypeCode == "120" || doc.DocTypeCode == "130") && 
				   doc.XbrlFlag == "1" {
					fmt.Printf("  → 対象文書として追加\n")
					filtered = append(filtered, doc)
				}
			}
		}
	}
	
	fmt.Printf("フィルタリング結果: %d件の文書が対象\n", len(filtered))
	return filtered
} 

// GetSecCodeMapping 証券コードマッピングを取得
func (e *EdinetAPI) GetSecCodeMapping() (map[string]string, error) {
	url := "https://api.edinet-fsa.go.jp/api/v2/companies.json"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}
	
	req.Header.Set("Ocp-Apim-Subscription-Key", e.apiKey)
	req.Header.Set("Accept", "application/json")
	
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API取得エラー: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}
	
	var companyList struct {
		Results []struct {
			SecCode string `json:"secCode"`
			JCN     string `json:"jcn"`
		} `json:"results"`
	}
	
	if err := json.Unmarshal(body, &companyList); err != nil {
		return nil, fmt.Errorf("JSONパースエラー: %v", err)
	}
	
	// 4桁→5桁のマッピングを作成
	mapping := make(map[string]string)
	for _, company := range companyList.Results {
		if len(company.SecCode) == 5 && company.SecCode != "" {
			// 5桁の証券コードから4桁を抽出（末尾の0を除去）
			fourDigit := strings.TrimSuffix(company.SecCode, "0")
			if len(fourDigit) == 4 {
				mapping[fourDigit] = company.SecCode
			}
		}
	}
	
	return mapping, nil
}

// ConvertSecCode 4桁の証券コードを5桁のEDINET証券コードに変換
func (e *EdinetAPI) ConvertSecCode(fourDigitCode string) (string, error) {
	mapping, err := e.GetSecCodeMapping()
	if err != nil {
		return "", err
	}
	
	if fiveDigitCode, exists := mapping[fourDigitCode]; exists {
		return fiveDigitCode, nil
	}
	
	// マッピングにない場合は、末尾に0を付けて試す
	fiveDigitCode := fourDigitCode + "0"
	return fiveDigitCode, nil
} 