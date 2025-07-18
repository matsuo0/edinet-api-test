package config

import (
	"flag"
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
	QuarterOnly  bool
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

	// コマンドライン引数を定義
	var startDate, endDate, targetSecCode, outputFile string
	var quarterOnly bool
	
	flag.StringVar(&startDate, "start", "", "開始日 (YYYY-MM-DD形式)")
	flag.StringVar(&endDate, "end", "", "終了日 (YYYY-MM-DD形式)")
	flag.StringVar(&targetSecCode, "code", "", "対象証券コード（4桁または5桁、空文字列で全企業）")
	flag.StringVar(&outputFile, "output", "", "出力ファイル名")
	flag.BoolVar(&quarterOnly, "quarter", false, "四半期報告書のみを対象にする")
	
	flag.Parse()

	// デフォルト値を設定
	if startDate == "" {
		startDate = "2025-07-10"
	}
	if endDate == "" {
		endDate = "2025-07-16"
	}
	if targetSecCode == "" {
		targetSecCode = "" // 空文字列のままにする（全企業対象）
	}
	if outputFile == "" {
		outputFile = "xbrl_financial_items.csv"
	}

	// 4桁の証券コードの場合は5桁に変換
	if len(targetSecCode) == 4 {
		targetSecCode = targetSecCode + "0"
	}

	return &Config{
		APIKey:        apiKey,
		StartDate:     startDate,
		EndDate:       endDate,
		TargetSecCode: targetSecCode,
		OutputFile:    outputFile,
		QuarterOnly:   quarterOnly,
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