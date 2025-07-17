package models

import (
	"encoding/json"
	"testing"
)

func TestDocInfo_JSONMarshal(t *testing.T) {
	doc := DocInfo{
		DocID:       "S100ABCD",
		FilerName:   "テスト株式会社",
		DocTypeCode: "120",
		XbrlFlag:    "1",
		SecCode:     "12345",
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("JSONマーシャルエラー: %v", err)
	}

	var unmarshaled DocInfo
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("JSONアンマーシャルエラー: %v", err)
	}

	if unmarshaled.DocID != doc.DocID {
		t.Errorf("DocID不一致: 期待=%s, 実際=%s", doc.DocID, unmarshaled.DocID)
	}
	if unmarshaled.FilerName != doc.FilerName {
		t.Errorf("FilerName不一致: 期待=%s, 実際=%s", doc.FilerName, unmarshaled.FilerName)
	}
	if unmarshaled.DocTypeCode != doc.DocTypeCode {
		t.Errorf("DocTypeCode不一致: 期待=%s, 実際=%s", doc.DocTypeCode, unmarshaled.DocTypeCode)
	}
	if unmarshaled.XbrlFlag != doc.XbrlFlag {
		t.Errorf("XbrlFlag不一致: 期待=%s, 実際=%s", doc.XbrlFlag, unmarshaled.XbrlFlag)
	}
	if unmarshaled.SecCode != doc.SecCode {
		t.Errorf("SecCode不一致: 期待=%s, 実際=%s", doc.SecCode, unmarshaled.SecCode)
	}
}

func TestDocumentListResponse_JSONMarshal(t *testing.T) {
	response := DocumentListResponse{
		Results: []DocInfo{
			{
				DocID:       "S100ABCD",
				FilerName:   "テスト株式会社A",
				DocTypeCode: "120",
				XbrlFlag:    "1",
				SecCode:     "12345",
			},
			{
				DocID:       "S100EFGH",
				FilerName:   "テスト株式会社B",
				DocTypeCode: "130",
				XbrlFlag:    "1",
				SecCode:     "67890",
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("JSONマーシャルエラー: %v", err)
	}

	var unmarshaled DocumentListResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("JSONアンマーシャルエラー: %v", err)
	}

	if len(unmarshaled.Results) != len(response.Results) {
		t.Errorf("Results数不一致: 期待=%d, 実際=%d", len(response.Results), len(unmarshaled.Results))
	}

	for i, result := range response.Results {
		if unmarshaled.Results[i].DocID != result.DocID {
			t.Errorf("Results[%d].DocID不一致: 期待=%s, 実際=%s", i, result.DocID, unmarshaled.Results[i].DocID)
		}
	}
}

func TestFinancialData_Structure(t *testing.T) {
	data := FinancialData{
		Date:         "2025-01-01",
		SecCode:      "12345",
		CompanyName:  "テスト株式会社",
		DocType:      "有価証券報告書",
		Period:       "2024年度",
		NetSales:     "1000000000",
		GrossProfit:  "500000000",
		OperatingIncome: "200000000",
		OrdinaryIncome: "180000000",
		IncomeBeforeTaxes: "160000000",
		ProfitLoss:   "120000000",
		EPS:          "100.00",
		TotalAssets:  "2000000000",
		CurrentAssets: "800000000",
		NoncurrentAssets: "1200000000",
		Liabilities:  "1000000000",
		CurrentLiabilities: "400000000",
		NoncurrentLiabilities: "600000000",
		NetAssets:    "1000000000",
		CapitalStock: "500000000",
		RetainedEarnings: "300000000",
		OperatingCF:  "150000000",
		InvestmentCF: "-50000000",
		FinancingCF:  "-80000000",
		CashAndEquivalents: "20000000",
		NetAssetsPerShare: "1000.00",
		EquityRatio:  "0.50",
		Dividends:    "50000000",
	}

	// 基本的な構造テスト
	if data.Date == "" {
		t.Error("Dateが空です")
	}
	if data.SecCode == "" {
		t.Error("SecCodeが空です")
	}
	if data.CompanyName == "" {
		t.Error("CompanyNameが空です")
	}
	if data.DocType == "" {
		t.Error("DocTypeが空です")
	}
	if data.Period == "" {
		t.Error("Periodが空です")
	}

	// 財務データの存在確認
	if data.NetSales == "" {
		t.Error("NetSalesが空です")
	}
	if data.ProfitLoss == "" {
		t.Error("ProfitLossが空です")
	}
	if data.TotalAssets == "" {
		t.Error("TotalAssetsが空です")
	}
} 