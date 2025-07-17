package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	
	var docList models.DocumentListResponse
	if err := json.Unmarshal(body, &docList); err != nil {
		return nil, fmt.Errorf("JSONパースエラー: %v", err)
	}
	
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
func FilterDocuments(docs []models.DocInfo, targetSecCode string) []models.DocInfo {
	var filtered []models.DocInfo
	for _, doc := range docs {
		if (doc.DocTypeCode == "120" || doc.DocTypeCode == "130") && 
		   doc.SecCode == targetSecCode && 
		   doc.XbrlFlag == "1" {
			filtered = append(filtered, doc)
		}
	}
	return filtered
} 