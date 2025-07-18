package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"edinet-api-test/internal/config"
	"edinet-api-test/internal/models"
	"edinet-api-test/internal/utils"
)

// CSVWriter CSV出力器
type CSVWriter struct {
	writer    *csv.Writer
	file      *os.File
	headers   []string
	financialTags []string
}

// NewCSVWriter 新しいCSV出力器を作成
func NewCSVWriter(filename string) (*CSVWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("CSV作成エラー: %v", err)
	}
	
	writer := csv.NewWriter(file)
	
	return &CSVWriter{
		writer:        writer,
		file:          file,
		headers:       config.JapaneseHeaders,
		financialTags: config.FinancialTags,
	}, nil
}

// WriteHeader ヘッダーを書き込み
func (c *CSVWriter) WriteHeader() error {
	return c.writer.Write(c.headers)
}

// WriteFinancialData 財務データを書き込み
func (c *CSVWriter) WriteFinancialData(data *models.FinancialData) error {
	row := []string{
		data.Date,
		data.SecCode,
		data.CompanyName,
		data.DocType,
		data.Period,
		// 基本財務データ
		data.NetSales,
		data.GrossProfit,
		data.OperatingIncome,
		data.OrdinaryIncome,
		data.IncomeBeforeTaxes,
		data.ProfitLoss,
		data.EPS,
		data.TotalAssets,
		data.CurrentAssets,
		data.NoncurrentAssets,
		data.Liabilities,
		data.CurrentLiabilities,
		data.NoncurrentLiabilities,
		data.NetAssets,
		data.CapitalStock,
		data.RetainedEarnings,
		data.OperatingCF,
		data.InvestmentCF,
		data.FinancingCF,
		data.CashAndEquivalents,
		data.NetAssetsPerShare,
		data.EquityRatio,
		data.Dividends,
		// 企業基本情報
		data.DateOfEstablishment,
		data.DateOfListing,
		data.NumberOfEmployees,
		data.ResearchAndDevelopmentExpenses,
		"", // 研究開発費比率（計算値）
		data.AccountingStandards,
		data.NameOfIndependentAuditor,
		data.ConsolidatedOrNonConsolidatedFinancialStatements,
		data.FiscalYearEnd,
		data.FiscalYearStart,
		// 収益性指標
		data.OperatingIncomeRatio,
		data.OrdinaryIncomeRatio,
		data.ProfitLossRatio,
		data.GrossProfitRatio,
		data.TotalAssetsTurnover,
		data.NetAssetsTurnover,
		data.OperatingCashFlowRatio,
		data.InvestmentCashFlowRatio,
		// 安全性指標
		data.WorkingCapital,
		data.DebtRatio,
		data.FixedRatio,
		data.FixedLongTermCoverageRatio,
		// 追加財務データ
		data.CostOfSales,
		data.SellingGeneralAndAdministrativeExpenses,
		data.NonOperatingIncome,
		data.NonOperatingExpenses,
		data.ExtraordinaryIncome,
		data.ExtraordinaryLoss,
		data.IncomeTaxes,
		data.ProfitLossAttributableToMinorityShareholders,
		data.ProfitLossAttributableToOwnersOfParent,
		data.NotesAndAccountsReceivableTrade,
		data.Inventories,
		data.PropertyPlantAndEquipment,
		data.IntangibleAssets,
		data.InvestmentsAndOtherAssets,
		data.ShortTermLoansPayable,
		data.NotesAndAccountsPayableTrade,
		data.LongTermLoansPayable,
		data.BondsPayable,
		data.ProvisionForRetirementBenefits,
		data.ShareholdersEquity,
		data.CapitalSurplus,
		data.ValuationDifferenceOnAvailableForSaleSecurities,
		data.TreasuryStock,
		// 追加比率指標
		data.CurrentRatio,
		data.QuickRatio,
		data.AccountsReceivableTurnoverDays,
		data.InventoryTurnoverDays,
		data.PropertyPlantAndEquipmentTurnover,
		data.TotalCapitalTurnover,
		data.OperatingCapitalTurnover,
		data.InterestCoverageRatio,
		data.DividendPayoutRatio,
		data.DividendYield,
		// キャッシュフロー関連
		data.Depreciation,
		data.IncreaseDecreaseInProvision,
		data.IncreaseDecreaseInWorkingCapital,
		data.ProceedsFromSalesOfInvestmentSecurities,
		data.PaymentsForPurchaseOfInvestmentSecurities,
		data.ProceedsFromLongTermLoansPayable,
		data.RepaymentsOfLongTermLoansPayable,
		data.FreeCashFlow,
		data.CashFlowCoverageRatio,
		// 成長性指標
		data.NetSalesGrowthRate,
		data.OperatingIncomeGrowthRate,
		data.ProfitLossGrowthRate,
		data.TotalAssetsGrowthRate,
		data.NetSalesPerEmployee,
		data.OperatingIncomePerEmployee,
		// メタデータ
		data.DataCollectionDate,
		data.DataSource,
		data.TaxonomyVersion,
		// キャッシュフロー詳細
		data.IncomeTaxesPaid,
		data.InterestPaid,
		data.InterestAndDividendsReceived,
		data.DividendsReceived,
		data.DividendsPaid,
		data.PaymentsForPurchaseOfPropertyPlantAndEquipment,
		data.ProceedsFromSalesOfPropertyPlantAndEquipment,
		data.PaymentsForPurchaseOfIntangibleAssets,
		data.ProceedsFromSalesOfIntangibleAssets,
		data.ProceedsFromShortTermLoansPayable,
		data.RepaymentsOfShortTermLoansPayable,
		data.ProceedsFromLongTermLoansPayable,
		data.RepaymentsOfLongTermLoansPayable,
		data.ProceedsFromIssuanceOfBonds,
		data.RedemptionOfBonds,
	}
	
	return c.writer.Write(row)
}

// WriteRow 生の行データを書き込み
func (c *CSVWriter) WriteRow(row []string) error {
	return c.writer.Write(row)
}

// Flush バッファをフラッシュ
func (c *CSVWriter) Flush() {
	c.writer.Flush()
}

// Close ファイルを閉じる
func (c *CSVWriter) Close() error {
	c.Flush()
	return c.file.Close()
}

// ExtractFinancialValues 財務タグから値を抽出（計算値も含む）
func (c *CSVWriter) ExtractFinancialValues(values map[string]string) []string {
	var result []string
	
	// 基本財務値を抽出
	for _, tag := range c.financialTags {
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
		result = append(result, found)
	}
	
	// 計算値を追加
	ratios := utils.CalculateFinancialRatios(values)
	metrics := utils.CalculateAdditionalMetrics(values)
	
	// 計算された比率を追加（順序をヘッダーに合わせる）
	calculatedValues := []string{
		"", // 設立年月日
		"", // 上場年月日
		values["jppfs_cor:NumberOfEmployees"], // 従業員数
		values["jppfs_cor:ResearchAndDevelopmentExpenses"], // 研究開発費
		ratios["jppfs_cor:ResearchAndDevelopmentExpenseRatio"], // 研究開発費比率
		"", // 会計基準
		"", // 監査法人
		"", // 連結・単体
		"", // 決算月
		"", // 年度開始月
		ratios["jppfs_cor:OperatingIncomeRatio"], // 営業利益率
		ratios["jppfs_cor:OrdinaryIncomeRatio"], // 経常利益率
		ratios["jppfs_cor:ProfitLossRatio"], // 当期純利益率
		ratios["jppfs_cor:GrossProfitRatio"], // 売上高総利益率
		ratios["jppfs_cor:TotalAssetsTurnover"], // 総資産回転率
		ratios["jppfs_cor:NetAssetsTurnover"], // 自己資本回転率
		ratios["jppfs_cor:OperatingCashFlowRatio"], // 営業CF比率
		ratios["jppfs_cor:InvestmentCashFlowRatio"], // 投資CF比率
		metrics["jppfs_cor:WorkingCapital"], // 運転資本
		metrics["jppfs_cor:DebtRatio"], // 負債比率
		metrics["jppfs_cor:FixedRatio"], // 固定比率
		metrics["jppfs_cor:FixedLongTermCoverageRatio"], // 固定長期適合率
		utils.GetCurrentTimestamp(), // データ取得日時
		"EDINET", // データソース
		"", // XBRLタクソノミーバージョン
	}
	
	result = append(result, calculatedValues...)
	return result
} 