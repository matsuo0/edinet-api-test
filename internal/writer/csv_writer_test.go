package writer

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"

	"edinet-api-test/internal/models"
)

func TestNewCSVWriter(t *testing.T) {
	// 一時ファイル名
	tmpFile := "test_output.csv"
	defer os.Remove(tmpFile)

	writer, err := NewCSVWriter(tmpFile)
	if err != nil {
		t.Fatalf("CSVWriter作成エラー: %v", err)
	}
	defer writer.Close()

	if writer == nil {
		t.Error("CSVWriterがnilです")
	}
	if writer.file == nil {
		t.Error("ファイルがnilです")
	}
	if writer.writer == nil {
		t.Error("CSV writerがnilです")
	}
	if len(writer.headers) == 0 {
		t.Error("ヘッダーが空です")
	}
	if len(writer.financialTags) == 0 {
		t.Error("財務タグが空です")
	}
}

func TestNewCSVWriter_InvalidPath(t *testing.T) {
	// 無効なパス（権限がないディレクトリなど）
	invalidPath := "/root/invalid/test.csv"

	_, err := NewCSVWriter(invalidPath)
	if err == nil {
		t.Error("無効なパスの場合、エラーが発生すべきです")
	}
}

func TestCSVWriter_WriteHeader(t *testing.T) {
	tmpFile := "test_header.csv"
	defer os.Remove(tmpFile)

	writer, err := NewCSVWriter(tmpFile)
	if err != nil {
		t.Fatalf("CSVWriter作成エラー: %v", err)
	}
	defer writer.Close()

	err = writer.WriteHeader()
	if err != nil {
		t.Fatalf("ヘッダー書き込みエラー: %v", err)
	}

	// ファイルを閉じてから内容を確認
	writer.Close()

	// ファイル内容を読み込み
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ファイル読み込みエラー: %v", err)
	}

	// CSVとして解析
	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("CSV解析エラー: %v", err)
	}

	if len(records) != 1 {
		t.Errorf("レコード数不一致: 期待=1, 実際=%d", len(records))
	}

	// 最初の5つのヘッダーを確認
	expectedHeaders := []string{"日付", "証券コード", "会社名", "文書タイプ", "会計期間"}
	for i, expected := range expectedHeaders {
		if records[0][i] != expected {
			t.Errorf("ヘッダー[%d]不一致: 期待=%s, 実際=%s", i, expected, records[0][i])
		}
	}
}

func TestCSVWriter_WriteRow(t *testing.T) {
	tmpFile := "test_row.csv"
	defer os.Remove(tmpFile)

	writer, err := NewCSVWriter(tmpFile)
	if err != nil {
		t.Fatalf("CSVWriter作成エラー: %v", err)
	}
	defer writer.Close()

	// ヘッダーを書き込み
	err = writer.WriteHeader()
	if err != nil {
		t.Fatalf("ヘッダー書き込みエラー: %v", err)
	}

	// テストデータを書き込み（ヘッダーと同じ長さのデータ）
	testRow := make([]string, len(writer.headers))
	testRow[0] = "2025-01-01"
	testRow[1] = "12345"
	testRow[2] = "テスト株式会社"
	testRow[3] = "有価証券報告書"
	testRow[4] = "2024年度"
	testRow[5] = "1000000000"
	
	err = writer.WriteRow(testRow)
	if err != nil {
		t.Fatalf("行書き込みエラー: %v", err)
	}

	// ファイルを閉じてから内容を確認
	writer.Close()

	// ファイル内容を読み込み
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ファイル読み込みエラー: %v", err)
	}

	// CSVとして解析
	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("CSV解析エラー: %v", err)
	}

	if len(records) != 2 { // ヘッダー + データ行
		t.Errorf("レコード数不一致: 期待=2, 実際=%d", len(records))
	}

	// データ行を確認
	dataRow := records[1]
	for i, expected := range testRow[:5] { // 最初の5つだけ確認
		if dataRow[i] != expected {
			t.Errorf("データ[%d]不一致: 期待=%s, 実際=%s", i, expected, dataRow[i])
		}
	}
}

func TestCSVWriter_WriteFinancialData(t *testing.T) {
	tmpFile := "test_financial.csv"
	defer os.Remove(tmpFile)

	writer, err := NewCSVWriter(tmpFile)
	if err != nil {
		t.Fatalf("CSVWriter作成エラー: %v", err)
	}
	defer writer.Close()

	// テスト用の財務データを作成
	financialData := &models.FinancialData{
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

	err = writer.WriteFinancialData(financialData)
	if err != nil {
		t.Fatalf("財務データ書き込みエラー: %v", err)
	}

	// ファイルを閉じてから内容を確認
	writer.Close()

	// ファイル内容を読み込み
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ファイル読み込みエラー: %v", err)
	}

	// CSVとして解析
	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("CSV解析エラー: %v", err)
	}

	if len(records) != 1 {
		t.Errorf("レコード数不一致: 期待=1, 実際=%d", len(records))
	}

	// 最初の5つのフィールドを確認
	expectedFields := []string{"2025-01-01", "12345", "テスト株式会社", "有価証券報告書", "2024年度"}
	dataRow := records[0]
	for i, expected := range expectedFields {
		if dataRow[i] != expected {
			t.Errorf("フィールド[%d]不一致: 期待=%s, 実際=%s", i, expected, dataRow[i])
		}
	}
}

func TestCSVWriter_ExtractFinancialValues(t *testing.T) {
	writer, err := NewCSVWriter("test_extract.csv")
	if err != nil {
		t.Fatalf("CSVWriter作成エラー: %v", err)
	}
	defer writer.Close()
	defer os.Remove("test_extract.csv")

	// テスト用のXBRL値を設定
	values := map[string]string{
		"jppfs_cor:NetSales|contextRef=CurrentYearDuration|unitRef=JPY": "1000000000",
		"jppfs_cor:ProfitLoss|contextRef=CurrentYearDuration|unitRef=JPY": "100000000",
		"jppfs_cor:TotalAssets|contextRef=CurrentYearDuration|unitRef=JPY": "2000000000",
		"jppfs_cor:NetSalesTextBlock|contextRef=CurrentYearDuration": "売上高の説明",
	}

	// 財務値を抽出
	result := writer.ExtractFinancialValues(values)

	// 結果の長さを確認
	expectedLength := len(writer.financialTags)
	if len(result) != expectedLength {
		t.Errorf("結果の長さ不一致: 期待=%d, 実際=%d", expectedLength, len(result))
	}

	// 最初のタグ（NetSales）が正しく抽出されているか確認
	if len(result) > 0 && result[0] != "1000000000" {
		t.Errorf("NetSalesの抽出値不一致: 期待=1000000000, 実際=%s", result[0])
	}
}

func TestCSVWriter_ExtractFinancialValues_Empty(t *testing.T) {
	writer, err := NewCSVWriter("test_empty.csv")
	if err != nil {
		t.Fatalf("CSVWriter作成エラー: %v", err)
	}
	defer writer.Close()
	defer os.Remove("test_empty.csv")

	// 空の値を設定
	values := map[string]string{}

	// 財務値を抽出
	result := writer.ExtractFinancialValues(values)

	// 結果の長さを確認
	expectedLength := len(writer.financialTags)
	if len(result) != expectedLength {
		t.Errorf("結果の長さ不一致: 期待=%d, 実際=%d", expectedLength, len(result))
	}

	// 全ての値が空文字列であることを確認
	for i, value := range result {
		if value != "" {
			t.Errorf("値[%d]が空でない: %s", i, value)
		}
	}
} 