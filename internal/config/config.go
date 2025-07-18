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
	// 基本財務データ
	"売上高", "売上総利益", "営業利益", "経常利益", "税引前当期純利益", 
	"当期純利益", "1株当たり当期純利益", "総資産", "流動資産", "固定資産", 
	"総負債", "流動負債", "固定負債", "純資産", "資本金", "利益剰余金", 
	"営業CF", "投資CF", "財務CF", "現金及び現金同等物", "1株当たり純資産", 
	"自己資本比率", "配当金",
	// 企業基本情報
	"設立年月日", "上場年月日", "従業員数", "研究開発費", "研究開発費比率",
	"会計基準", "監査法人", "連結・単体", "決算月", "年度開始月",
	// 収益性指標
	"営業利益率", "経常利益率", "当期純利益率", "売上高総利益率",
	"総資産回転率", "自己資本回転率", "営業CF比率", "投資CF比率",
	// 安全性指標
	"運転資本", "負債比率", "固定比率", "固定長期適合率",
	// 追加財務データ
	"売上原価", "販管費", "営業外収益", "営業外費用", "特別利益", "特別損失",
	"法人税等", "少数株主損益", "親会社株主に帰属する当期純利益",
	"売上債権", "棚卸資産", "有形固定資産", "無形固定資産", "投資その他の資産",
	"短期借入金", "買掛金", "長期借入金", "社債", "退職給付引当金",
	"株主資本", "資本剰余金", "その他有価証券評価差額金", "自己株式",
	// 追加比率指標
	"流動比率", "当座比率", "売上債権回転日数", "棚卸資産回転日数",
	"有形固定資産回転率", "総資本回転率", "営業資本回転率",
	"インタレスト・カバレッジ・レシオ", "配当性向", "配当利回り",
	// キャッシュフロー関連
	"減価償却費", "引当金の増減", "運転資本の増減", "投資活動による収入",
	"投資活動による支出", "財務活動による収入", "財務活動による支出",
	"自由キャッシュフロー", "キャッシュフロー充足率",
	// 成長性指標
	"売上高成長率", "営業利益成長率", "当期純利益成長率", "総資産成長率",
	"従業員一人当たり売上高", "従業員一人当たり営業利益",
	// メタデータ
	"データ取得日時", "データソース", "XBRLタクソノミーバージョン",
	// キャッシュフロー詳細
	"法人税等支払額", "利息支払額", "利息受取額", "配当金受取額", "配当金支払額",
	"有形固定資産取得による支出", "有形固定資産売却による収入", "無形固定資産取得による支出", "無形固定資産売却による収入",
	"短期借入金による収入", "短期借入金返済額", "長期借入金による収入", "長期借入金返済額",
	"社債発行による収入", "社債償還額",
}

// FinancialTags 財務タグ
var FinancialTags = []string{
	// 基本財務データ
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
	// 企業基本情報
	"jppfs_cor:DateOfEstablishment", "jppfs_cor:DateOfListing", 
	"jppfs_cor:NumberOfEmployees", "jppfs_cor:ResearchAndDevelopmentExpenses",
	"jppfs_cor:AccountingStandards", "jppfs_cor:NameOfIndependentAuditor",
	"jppfs_cor:ConsolidatedOrNonConsolidatedFinancialStatements",
	"jppfs_cor:FiscalYearEnd", "jppfs_cor:FiscalYearStart",
	// 追加財務データ
	"jppfs_cor:CostOfSales", "jppfs_cor:SellingGeneralAndAdministrativeExpenses",
	"jppfs_cor:NonOperatingIncome", "jppfs_cor:NonOperatingExpenses",
	"jppfs_cor:ExtraordinaryIncome", "jppfs_cor:ExtraordinaryLoss",
	"jppfs_cor:IncomeTaxes", "jppfs_cor:ProfitLossAttributableToMinorityShareholders",
	"jppfs_cor:ProfitLossAttributableToOwnersOfParent",
	"jppfs_cor:NotesAndAccountsReceivableTrade", "jppfs_cor:Inventories",
	"jppfs_cor:PropertyPlantAndEquipment", "jppfs_cor:IntangibleAssets",
	"jppfs_cor:InvestmentsAndOtherAssets", "jppfs_cor:ShortTermLoansPayable",
	"jppfs_cor:NotesAndAccountsPayableTrade", "jppfs_cor:LongTermLoansPayable",
	"jppfs_cor:BondsPayable", "jppfs_cor:ProvisionForRetirementBenefits",
	"jppfs_cor:ShareholdersEquity", "jppfs_cor:CapitalSurplus",
	"jppfs_cor:ValuationDifferenceOnAvailableForSaleSecurities", "jppfs_cor:TreasuryStock",
	// キャッシュフロー関連
	"jppfs_cor:Depreciation", "jppfs_cor:IncreaseDecreaseInProvision",
	"jppfs_cor:IncreaseDecreaseInWorkingCapital", "jppfs_cor:ProceedsFromSalesOfInvestmentSecurities",
	"jppfs_cor:PaymentsForPurchaseOfInvestmentSecurities", "jppfs_cor:ProceedsFromLongTermLoansPayable",
	"jppfs_cor:RepaymentsOfLongTermLoansPayable",
	// 計算値（EDINETデータから算出）
	"jppfs_cor:OperatingIncomeRatio", "jppfs_cor:OrdinaryIncomeRatio",
	"jppfs_cor:ProfitLossRatio", "jppfs_cor:GrossProfitRatio",
	"jppfs_cor:TotalAssetsTurnover", "jppfs_cor:NetAssetsTurnover",
	"jppfs_cor:OperatingCashFlowRatio", "jppfs_cor:InvestmentCashFlowRatio",
	"jppfs_cor:WorkingCapital", "jppfs_cor:DebtRatio", "jppfs_cor:FixedRatio", 
	"jppfs_cor:FixedLongTermCoverageRatio", "jppfs_cor:CurrentRatio",
	"jppfs_cor:QuickRatio", "jppfs_cor:AccountsReceivableTurnoverDays",
	"jppfs_cor:InventoryTurnoverDays", "jppfs_cor:PropertyPlantAndEquipmentTurnover",
	"jppfs_cor:TotalCapitalTurnover", "jppfs_cor:OperatingCapitalTurnover",
	"jppfs_cor:InterestCoverageRatio", "jppfs_cor:DividendPayoutRatio",
	"jppfs_cor:DividendYield", "jppfs_cor:FreeCashFlow", "jppfs_cor:CashFlowCoverageRatio",
	"jppfs_cor:NetSalesGrowthRate", "jppfs_cor:OperatingIncomeGrowthRate",
	"jppfs_cor:ProfitLossGrowthRate", "jppfs_cor:TotalAssetsGrowthRate",
	"jppfs_cor:NetSalesPerEmployee", "jppfs_cor:OperatingIncomePerEmployee",
	// メタデータ
	"jppfs_cor:DataCollectionDate", "jppfs_cor:DataSource", "jppfs_cor:TaxonomyVersion",
	// キャッシュフロー詳細
	"jppfs_cor:IncomeTaxesPaid", "jppfs_cor:InterestPaid", "jppfs_cor:InterestAndDividendsReceived", "jppfs_cor:DividendsReceived", "jppfs_cor:DividendsPaid",
	"jppfs_cor:PaymentsForPurchaseOfPropertyPlantAndEquipment", "jppfs_cor:ProceedsFromSalesOfPropertyPlantAndEquipment",
	"jppfs_cor:PaymentsForPurchaseOfIntangibleAssets", "jppfs_cor:ProceedsFromSalesOfIntangibleAssets",
	"jppfs_cor:ProceedsFromShortTermLoansPayable", "jppfs_cor:RepaymentsOfShortTermLoansPayable",
	"jppfs_cor:ProceedsFromLongTermLoansPayable", "jppfs_cor:RepaymentsOfLongTermLoansPayable",
	"jppfs_cor:ProceedsFromIssuanceOfBonds", "jppfs_cor:RedemptionOfBonds",
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