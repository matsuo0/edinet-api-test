package parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// XBRLParser XBRLファイル解析器
type XBRLParser struct{}

// NewXBRLParser 新しいXBRL解析器を作成
func NewXBRLParser() *XBRLParser {
	return &XBRLParser{}
}

// ExtractPublicDocXBRL ZIPからPublicDocのXBRLファイルを抽出
func (x *XBRLParser) ExtractPublicDocXBRL(zipFile string) (string, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", fmt.Errorf("ZIPファイルオープンエラー: %v", err)
	}
	defer r.Close()
	
	for _, f := range r.File {
		if strings.Contains(f.Name, "PublicDoc") && strings.HasSuffix(f.Name, ".xbrl") {
			outPath := filepath.Base(f.Name)
			out, err := os.Create(outPath)
			if err != nil {
				return "", fmt.Errorf("ファイル作成エラー: %v", err)
			}
			
			in, err := f.Open()
			if err != nil {
				out.Close()
				return "", fmt.Errorf("ZIP内ファイルオープンエラー: %v", err)
			}
			
			_, err = io.Copy(out, in)
			in.Close()
			out.Close()
			
			if err != nil {
				return "", fmt.Errorf("ファイルコピーエラー: %v", err)
			}
			
			return outPath, nil
		}
	}
	
	return "", fmt.Errorf("PublicDocのxbrlファイルが見つかりません")
}

// ParseAllXBRL XBRLファイルから全ての値を抽出
func (x *XBRLParser) ParseAllXBRL(xbrlPath string) (map[string]string, error) {
	file, err := os.Open(xbrlPath)
	if err != nil {
		return nil, fmt.Errorf("XBRLファイルオープンエラー: %v", err)
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
			return nil, fmt.Errorf("XMLデコードエラー: %v", err)
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

// ExtractAccountingPeriod XBRLファイルから会計期間の情報を抽出
func (x *XBRLParser) ExtractAccountingPeriod(xbrlPath string) (string, string) {
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
		
		if startDate != "" && endDate != "" {
			break
		}
	}
	
	return startDate, endDate
}

// GetQuarterInfo 会計期間から四半期情報を生成
func (x *XBRLParser) GetQuarterInfo(startDate, endDate string) string {
	if startDate == "" || endDate == "" {
		return "不明"
	}
	
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return "不明"
	}
	
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return "不明"
	}
	
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	
	if days > 300 { // 年次報告書
		year := end.Year()
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

// GetDocTypeName 文書タイプコードを日本語名に変換
func (x *XBRLParser) GetDocTypeName(docTypeCode string) string {
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