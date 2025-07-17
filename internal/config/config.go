package config

import (
	"os"
	"time"
)

// Config アプリケーション設定
type Config struct {
	APIKey       string
	StartDate    string
	EndDate      string
	TargetSecCode string
	OutputFile   string
}

// JapaneseHeaders 日本語ヘッダー
var JapaneseHeaders = []string{
	"日付", "証券コード", "会社名", "文書タイプ", "会計期間",
	"売上高", "売上総利益", "営業利益", "経常利益", "税引前当期純利益", 
	"当期純利益", "1株当たり当期純利益", "総資産", "流動資産", "固定資産", 
	"総負債", "流動負債", "固定負債", "純資産", "資本金", "利益剰余金", 
	"営業CF", "投資CF", "財務CF", "現金及び現金同等物", "1株当たり純資産", 
	"自己資本比率", "配当金",
}

// FinancialTags 財務タグ
var FinancialTags = []string{
	"jppfs_cor:NetSales", "jppfs_cor:GrossProfit", "jppfs_cor:OperatingIncome", 
	"jppfs_cor:OrdinaryIncome", "jppfs_cor:IncomeBeforeIncomeTaxes", "jppfs_cor:ProfitLoss", 
	"jppfs_cor:BasicEarningsLossPerShareSummaryOfBusinessResults", "jppfs_cor:TotalAssets", 
	"jppfs_cor:CurrentAssets", "jppfs_cor:NoncurrentAssets", "jppfs_cor:Liabilities", 
	"jppfs_cor:CurrentLiabilities", "jppfs_cor:NoncurrentLiabilities", "jppfs_cor:NetAssets", 
	"jppfs_cor:CapitalStock", "jppfs_cor:RetainedEarnings", 
	"jppfs_cor:NetCashProvidedByUsedInOperatingActivities", 
	"jppfs_cor:NetCashProvidedByUsedInInvestmentActivities", 
	"jppfs_cor:NetCashProvidedByUsedInFinancingActivities", "jppfs_cor:CashAndCashEquivalents", 
	"jppfs_cor:NetAssetsPerShareSummaryOfBusinessResults", 
	"jppfs_cor:EquityToAssetRatioSummaryOfBusinessResults", "jppfs_cor:DividendsFromSurplus",
}

// LoadConfig 設定を読み込み
func LoadConfig() (*Config, error) {
	apiKey := os.Getenv("EDINET_API_KEY")
	if apiKey == "" {
		return nil, &ConfigError{Message: "EDINET_API_KEYが設定されていません。.envファイルを確認してください。"}
	}

	return &Config{
		APIKey:        apiKey,
		StartDate:     "2025-07-10",
		EndDate:       "2025-07-16",
		TargetSecCode: "40260",
		OutputFile:    "xbrl_financial_items.csv",
	}, nil
}

// GetDateRange 日付範囲を取得
func (c *Config) GetDateRange() (time.Time, time.Time, error) {
	const layout = "2006-01-02"
	start, err := time.Parse(layout, c.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	end, err := time.Parse(layout, c.EndDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	return start, end, nil
}

// ConfigError 設定エラー
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
} 