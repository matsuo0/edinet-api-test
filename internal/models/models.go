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
} 