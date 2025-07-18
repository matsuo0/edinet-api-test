package models

// DocInfo EDINET APIの文書情報
type DocInfo struct {
	DocID       string `json:"docID"`
	FilerName   string `json:"filerName"`
	DocTypeCode string `json:"docTypeCode"`
	XbrlFlag    string `json:"xbrlFlag"`
	SecCode     string `json:"secCode"`
}

// DocumentListResponse EDINET APIの文書一覧レスポンス
type DocumentListResponse struct {
	Results []DocInfo `json:"results"`
}

// FinancialData 財務データ
type FinancialData struct {
	Date         string
	SecCode      string
	CompanyName  string
	DocType      string
	Period       string
	// 基本財務データ
	NetSales     string
	GrossProfit  string
	OperatingIncome string
	OrdinaryIncome string
	IncomeBeforeTaxes string
	ProfitLoss   string
	EPS          string
	TotalAssets  string
	CurrentAssets string
	NoncurrentAssets string
	Liabilities  string
	CurrentLiabilities string
	NoncurrentLiabilities string
	NetAssets    string
	CapitalStock string
	RetainedEarnings string
	OperatingCF  string
	InvestmentCF string
	FinancingCF  string
	CashAndEquivalents string
	NetAssetsPerShare string
	EquityRatio  string
	Dividends    string
	// 企業基本情報
	DateOfEstablishment string
	DateOfListing       string
	NumberOfEmployees   string
	ResearchAndDevelopmentExpenses string
	AccountingStandards string
	NameOfIndependentAuditor string
	ConsolidatedOrNonConsolidatedFinancialStatements string
	FiscalYearEnd string
	FiscalYearStart string
	// 収益性指標
	OperatingIncomeRatio string
	OrdinaryIncomeRatio string
	ProfitLossRatio string
	GrossProfitRatio string
	TotalAssetsTurnover string
	NetAssetsTurnover string
	OperatingCashFlowRatio string
	InvestmentCashFlowRatio string
	// 安全性指標
	WorkingCapital string
	DebtRatio string
	FixedRatio string
	FixedLongTermCoverageRatio string
	// 追加財務データ
	CostOfSales string
	SellingGeneralAndAdministrativeExpenses string
	NonOperatingIncome string
	NonOperatingExpenses string
	ExtraordinaryIncome string
	ExtraordinaryLoss string
	IncomeTaxes string
	ProfitLossAttributableToMinorityShareholders string
	ProfitLossAttributableToOwnersOfParent string
	NotesAndAccountsReceivableTrade string
	Inventories string
	PropertyPlantAndEquipment string
	IntangibleAssets string
	InvestmentsAndOtherAssets string
	ShortTermLoansPayable string
	NotesAndAccountsPayableTrade string
	LongTermLoansPayable string
	BondsPayable string
	ProvisionForRetirementBenefits string
	ShareholdersEquity string
	CapitalSurplus string
	ValuationDifferenceOnAvailableForSaleSecurities string
	TreasuryStock string
	// 追加比率指標
	CurrentRatio string
	QuickRatio string
	AccountsReceivableTurnoverDays string
	InventoryTurnoverDays string
	PropertyPlantAndEquipmentTurnover string
	TotalCapitalTurnover string
	OperatingCapitalTurnover string
	InterestCoverageRatio string
	DividendPayoutRatio string
	DividendYield string
	// キャッシュフロー関連
	Depreciation string
	IncreaseDecreaseInProvision string
	IncreaseDecreaseInWorkingCapital string
	ProceedsFromSalesOfInvestmentSecurities string
	PaymentsForPurchaseOfInvestmentSecurities string
	FreeCashFlow string
	CashFlowCoverageRatio string
	// 成長性指標
	NetSalesGrowthRate string
	OperatingIncomeGrowthRate string
	ProfitLossGrowthRate string
	TotalAssetsGrowthRate string
	NetSalesPerEmployee string
	OperatingIncomePerEmployee string
	// メタデータ
	DataCollectionDate string
	DataSource string
	TaxonomyVersion string
	// キャッシュフロー詳細
	IncomeTaxesPaid string
	InterestPaid string
	InterestAndDividendsReceived string
	DividendsReceived string
	DividendsPaid string
	PaymentsForPurchaseOfPropertyPlantAndEquipment string
	ProceedsFromSalesOfPropertyPlantAndEquipment string
	PaymentsForPurchaseOfIntangibleAssets string
	ProceedsFromSalesOfIntangibleAssets string
	ProceedsFromShortTermLoansPayable string
	RepaymentsOfShortTermLoansPayable string
	ProceedsFromLongTermLoansPayable string
	RepaymentsOfLongTermLoansPayable string
	ProceedsFromIssuanceOfBonds string
	RedemptionOfBonds string
} 