package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// StockData 株価データ
type StockData struct {
	Price        float64
	MarketCap    float64
	PBR          float64
	PER          float64
	DividendYield float64
	Timestamp    time.Time
}

// CompanyData 企業基本データ
type CompanyData struct {
	EmployeeCount    int
	EstablishmentYear int
	ListingYear      int
	MarketName       string
	IndustryCode     string
}

// IndustryData 業界平均データ
type IndustryData struct {
	AverageROE        float64
	AverageROA        float64
	AverageProfitMargin float64
	AverageDebtRatio  float64
}

// ExternalDataProvider 外部データプロバイダー
type ExternalDataProvider struct {
	client *http.Client
}

// NewExternalDataProvider 新しい外部データプロバイダーを作成
func NewExternalDataProvider() *ExternalDataProvider {
	return &ExternalDataProvider{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetStockData Yahoo Financeから株価データを取得
func (e *ExternalDataProvider) GetStockData(symbol string) (*StockData, error) {
	// Yahoo Finance APIのエンドポイント
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s.JP", symbol)
	
	resp, err := e.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API呼び出しエラー: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}
	
	// JSON解析（簡略化）
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v", err)
	}
	
	// 実際の実装では、Yahoo Financeのレスポンス構造に合わせて解析
	// ここではダミーデータを返す
	return &StockData{
		Price:        5000.0,
		MarketCap:    1000000000000.0, // 1兆円
		PBR:          2.5,
		PER:          15.0,
		DividendYield: 2.0,
		Timestamp:    time.Now(),
	}, nil
}

// GetCompanyData 企業基本データを取得
func (e *ExternalDataProvider) GetCompanyData(companyName string) (*CompanyData, error) {
	// 企業情報ナビAPI（実際のエンドポイントは要確認）
	url := fmt.Sprintf("https://api.company-info.jp/companies/%s", companyName)
	
	resp, err := e.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API呼び出しエラー: %v", err)
	}
	defer resp.Body.Close()
	
	// ダミーデータを返す
	return &CompanyData{
		EmployeeCount:    5000,
		EstablishmentYear: 1889,
		ListingYear:      1949,
		MarketName:       "東証一部",
		IndustryCode:     "G30", // 情報通信業
	}, nil
}

// GetIndustryAverages 業界平均データを取得
func (e *ExternalDataProvider) GetIndustryAverages(industryCode string) (*IndustryData, error) {
	// 統計データAPI（実際のエンドポイントは要確認）
	url := fmt.Sprintf("https://api.statistics.go.jp/industry/%s", industryCode)
	
	resp, err := e.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API呼び出しエラー: %v", err)
	}
	defer resp.Body.Close()
	
	// ダミーデータを返す
	return &IndustryData{
		AverageROE:        8.5,
		AverageROA:        5.2,
		AverageProfitMargin: 12.0,
		AverageDebtRatio:  45.0,
	}, nil
}

// CalculateFinancialRatios 外部データを使用して財務比率を計算
func (e *ExternalDataProvider) CalculateFinancialRatios(
	netSales, netIncome, totalAssets, netAssets float64,
	stockData *StockData,
) map[string]float64 {
	ratios := make(map[string]float64)
	
	// PBR計算
	if netAssets > 0 && stockData.Price > 0 {
		// 1株当たり純資産を計算（簡略化）
		ratios["PBR"] = stockData.Price / (netAssets / 1000000) // 発行済株式数を100万株と仮定
	}
	
	// PER計算
	if netIncome > 0 && stockData.Price > 0 {
		ratios["PER"] = stockData.Price / (netIncome / 1000000) // 発行済株式数を100万株と仮定
	}
	
	// 配当性向計算
	if netIncome > 0 {
		dividendAmount := netIncome * (stockData.DividendYield / 100)
		ratios["DividendPayoutRatio"] = (dividendAmount / netIncome) * 100
	}
	
	// ROE計算
	if netAssets > 0 {
		ratios["ROE"] = (netIncome / netAssets) * 100
	}
	
	// ROA計算
	if totalAssets > 0 {
		ratios["ROA"] = (netIncome / totalAssets) * 100
	}
	
	return ratios
}

// GetCurrentTimestamp 現在のタイムスタンプを取得
func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
} 