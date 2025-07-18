package utils

import (
	"strconv"
	"strings"
	"time"
)

// ParseFloat 文字列をfloat64に変換
func ParseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	// カンマを除去
	s = strings.ReplaceAll(s, ",", "")
	// 数値以外の文字を除去
	s = strings.TrimSpace(s)
	
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// FormatFloat float64を文字列に変換
func FormatFloat(f float64) string {
	if f == 0 {
		return ""
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// CalculateFinancialRatios 財務比率を計算
func CalculateFinancialRatios(values map[string]string) map[string]string {
	ratios := make(map[string]string)
	
	// 基本値を取得
	netSales := ParseFloat(values["jppfs_cor:NetSales"])
	grossProfit := ParseFloat(values["jppfs_cor:GrossProfit"])
	operatingIncome := ParseFloat(values["jppfs_cor:OperatingIncome"])
	ordinaryIncome := ParseFloat(values["jppfs_cor:OrdinaryIncome"])
	profitLoss := ParseFloat(values["jppfs_cor:ProfitLoss"])
	totalAssets := ParseFloat(values["jppfs_cor:TotalAssets"])
	netAssets := ParseFloat(values["jppfs_cor:NetAssets"])
	operatingCF := ParseFloat(values["jppfs_cor:NetCashProvidedByUsedInOperatingActivities"])
	investmentCF := ParseFloat(values["jppfs_cor:NetCashProvidedByUsedInInvestmentActivities"])
	researchDevExp := ParseFloat(values["jppfs_cor:ResearchAndDevelopmentExpenses"])
	numberOfEmployees := ParseFloat(values["jppfs_cor:NumberOfEmployees"])
	
	// 研究開発費比率
	if netSales > 0 && researchDevExp > 0 {
		ratios["jppfs_cor:ResearchAndDevelopmentExpenseRatio"] = FormatFloat((researchDevExp / netSales) * 100)
	}
	
	// 営業利益率
	if netSales > 0 {
		ratios["jppfs_cor:OperatingIncomeRatio"] = FormatFloat((operatingIncome / netSales) * 100)
	}
	
	// 経常利益率
	if netSales > 0 {
		ratios["jppfs_cor:OrdinaryIncomeRatio"] = FormatFloat((ordinaryIncome / netSales) * 100)
	}
	
	// 当期純利益率
	if netSales > 0 {
		ratios["jppfs_cor:ProfitLossRatio"] = FormatFloat((profitLoss / netSales) * 100)
	}
	
	// 売上高総利益率
	if netSales > 0 {
		ratios["jppfs_cor:GrossProfitRatio"] = FormatFloat((grossProfit / netSales) * 100)
	}
	
	// 総資産回転率
	if totalAssets > 0 && netSales > 0 {
		ratios["jppfs_cor:TotalAssetsTurnover"] = FormatFloat(netSales / totalAssets)
	}
	
	// 自己資本回転率
	if netAssets > 0 && netSales > 0 {
		ratios["jppfs_cor:NetAssetsTurnover"] = FormatFloat(netSales / netAssets)
	}
	
	// 営業CF比率
	if netSales > 0 && operatingCF != 0 {
		ratios["jppfs_cor:OperatingCashFlowRatio"] = FormatFloat((operatingCF / netSales) * 100)
	}
	
	// 投資CF比率
	if totalAssets > 0 && investmentCF != 0 {
		ratios["jppfs_cor:InvestmentCashFlowRatio"] = FormatFloat((investmentCF / totalAssets) * 100)
	}
	
	// 従業員一人当たり売上高
	if numberOfEmployees > 0 && netSales > 0 {
		ratios["jppfs_cor:NetSalesPerEmployee"] = FormatFloat(netSales / numberOfEmployees)
	}
	
	// 従業員一人当たり営業利益
	if numberOfEmployees > 0 && operatingIncome != 0 {
		ratios["jppfs_cor:OperatingIncomePerEmployee"] = FormatFloat(operatingIncome / numberOfEmployees)
	}
	
	return ratios
}

// CalculateAdditionalMetrics 追加の財務指標を計算
func CalculateAdditionalMetrics(values map[string]string) map[string]string {
	metrics := make(map[string]string)
	
	// 基本値を取得
	currentAssets := ParseFloat(values["jppfs_cor:CurrentAssets"])
	currentLiabilities := ParseFloat(values["jppfs_cor:CurrentLiabilities"])
	noncurrentAssets := ParseFloat(values["jppfs_cor:NoncurrentAssets"])
	liabilities := ParseFloat(values["jppfs_cor:Liabilities"])
	totalAssets := ParseFloat(values["jppfs_cor:TotalAssets"])
	netAssets := ParseFloat(values["jppfs_cor:NetAssets"])
	netSales := ParseFloat(values["jppfs_cor:NetSales"])
	dividends := ParseFloat(values["jppfs_cor:DividendsFromSurplus"])
	profitLoss := ParseFloat(values["jppfs_cor:ProfitLoss"])
	notesReceivable := ParseFloat(values["jppfs_cor:NotesAndAccountsReceivableTrade"])
	inventories := ParseFloat(values["jppfs_cor:Inventories"])
	propertyPlantEquipment := ParseFloat(values["jppfs_cor:PropertyPlantAndEquipment"])
	operatingCF := ParseFloat(values["jppfs_cor:NetCashProvidedByUsedInOperatingActivities"])
	investmentCF := ParseFloat(values["jppfs_cor:NetCashProvidedByUsedInInvestmentActivities"])
	capitalStock := ParseFloat(values["jppfs_cor:CapitalStock"])
	
	// 運転資本
	metrics["jppfs_cor:WorkingCapital"] = FormatFloat(currentAssets - currentLiabilities)
	
	// 負債比率
	if totalAssets > 0 {
		metrics["jppfs_cor:DebtRatio"] = FormatFloat((liabilities / totalAssets) * 100)
	}
	
	// 固定比率
	if netAssets > 0 {
		metrics["jppfs_cor:FixedRatio"] = FormatFloat((noncurrentAssets / netAssets) * 100)
	}
	
	// 固定長期適合率
	if netAssets > 0 {
		metrics["jppfs_cor:FixedLongTermCoverageRatio"] = FormatFloat((noncurrentAssets / netAssets) * 100)
	}
	
	// 流動比率
	if currentLiabilities > 0 {
		metrics["jppfs_cor:CurrentRatio"] = FormatFloat((currentAssets / currentLiabilities) * 100)
	}
	
	// 当座比率（流動資産から棚卸資産を除いたもの）
	if currentLiabilities > 0 {
		quickAssets := currentAssets - inventories
		if quickAssets > 0 {
			metrics["jppfs_cor:QuickRatio"] = FormatFloat((quickAssets / currentLiabilities) * 100)
		}
	}
	
	// 売上債権回転日数
	if netSales > 0 && notesReceivable > 0 {
		metrics["jppfs_cor:AccountsReceivableTurnoverDays"] = FormatFloat((notesReceivable / netSales) * 365)
	}
	
	// 棚卸資産回転日数
	if netSales > 0 && inventories > 0 {
		metrics["jppfs_cor:InventoryTurnoverDays"] = FormatFloat((inventories / netSales) * 365)
	}
	
	// 有形固定資産回転率
	if netSales > 0 && propertyPlantEquipment > 0 {
		metrics["jppfs_cor:PropertyPlantAndEquipmentTurnover"] = FormatFloat(netSales / propertyPlantEquipment)
	}
	
	// 総資本回転率
	if totalAssets > 0 && netSales > 0 {
		metrics["jppfs_cor:TotalCapitalTurnover"] = FormatFloat(netSales / totalAssets)
	}
	
	// 営業資本回転率
	if (currentAssets - currentLiabilities) > 0 && netSales > 0 {
		metrics["jppfs_cor:OperatingCapitalTurnover"] = FormatFloat(netSales / (currentAssets - currentLiabilities))
	}
	
	// 配当性向
	if profitLoss > 0 && dividends > 0 {
		metrics["jppfs_cor:DividendPayoutRatio"] = FormatFloat((dividends / profitLoss) * 100)
	}
	
	// 配当利回り（簡易計算：配当金÷資本金）
	if capitalStock > 0 && dividends > 0 {
		metrics["jppfs_cor:DividendYield"] = FormatFloat((dividends / capitalStock) * 100)
	}
	
	// 自由キャッシュフロー
	if operatingCF != 0 && investmentCF != 0 {
		metrics["jppfs_cor:FreeCashFlow"] = FormatFloat(operatingCF + investmentCF)
	}
	
	// キャッシュフロー充足率
	if investmentCF != 0 && operatingCF != 0 {
		metrics["jppfs_cor:CashFlowCoverageRatio"] = FormatFloat((operatingCF / -investmentCF) * 100)
	}
	
	return metrics
}

// GetCurrentTimestamp 現在のタイムスタンプを取得
func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
} 