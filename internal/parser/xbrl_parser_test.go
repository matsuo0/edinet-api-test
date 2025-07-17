package parser

import (
	"archive/zip"
	"os"
	"strings"
	"testing"
)

func TestNewXBRLParser(t *testing.T) {
	parser := NewXBRLParser()
	if parser == nil {
		t.Error("XBRLParserがnilです")
	}
}

func TestXBRLParser_GetDocTypeName(t *testing.T) {
	parser := NewXBRLParser()

	testCases := []struct {
		docTypeCode string
		expected    string
	}{
		{"120", "有価証券報告書"},
		{"130", "四半期報告書"},
		{"140", "有価証券届出書"},
		{"150", "発行登録書"},
		{"160", "発行登録追補書類"},
		{"170", "訂正有価証券報告書"},
		{"180", "訂正四半期報告書"},
		{"190", "訂正有価証券届出書"},
		{"200", "有価証券報告書（外国会社）"},
		{"210", "四半期報告書（外国会社）"},
		{"999", "その他"},
		{"", "その他"},
	}

	for _, tc := range testCases {
		result := parser.GetDocTypeName(tc.docTypeCode)
		if result != tc.expected {
			t.Errorf("DocTypeCode=%s: 期待=%s, 実際=%s", tc.docTypeCode, tc.expected, result)
		}
	}
}

func TestXBRLParser_GetQuarterInfo(t *testing.T) {
	parser := NewXBRLParser()

	testCases := []struct {
		startDate string
		endDate   string
		expected  string
	}{
		// 年次報告書（期間が300日以上）
		{"2024-04-01", "2025-03-31", "2024年度"},
		{"2023-04-01", "2024-03-31", "2023年度"},
		{"2024-01-01", "2024-12-31", "2024年度"},
		
		// 四半期報告書（期間が300日未満）
		{"2024-04-01", "2024-06-30", "2024Q1"},
		{"2024-07-01", "2024-09-30", "2024Q2"},
		{"2024-10-01", "2024-12-31", "2024Q3"},
		{"2024-01-01", "2024-03-31", "2023Q4"},
		
		// エラーケース
		{"", "2024-03-31", "不明"},
		{"2024-04-01", "", "不明"},
		{"", "", "不明"},
		{"invalid", "2024-03-31", "不明"},
		{"2024-04-01", "invalid", "不明"},
	}

	for _, tc := range testCases {
		result := parser.GetQuarterInfo(tc.startDate, tc.endDate)
		if result != tc.expected {
			t.Errorf("startDate=%s, endDate=%s: 期待=%s, 実際=%s", 
				tc.startDate, tc.endDate, tc.expected, result)
		}
	}
}

func TestXBRLParser_ExtractAccountingPeriod(t *testing.T) {
	parser := NewXBRLParser()

	// テスト用のXBRLファイルを作成
	testXBRL := `<?xml version="1.0" encoding="UTF-8"?>
<xbrli:xbrl xmlns:xbrli="http://www.xbrl.org/2003/instance">
  <xbrli:context id="CurrentYearDuration">
    <xbrli:entity>
      <xbrli:identifier scheme="http://disclosure.edinet-fsa.go.jp">E00763-000</xbrli:identifier>
    </xbrli:entity>
    <xbrli:period>
      <xbrli:startDate>2024-05-01</xbrli:startDate>
      <xbrli:endDate>2025-04-30</xbrli:endDate>
    </xbrli:period>
  </xbrli:context>
</xbrli:xbrl>`

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "test_*.xbrl")
	if err != nil {
		t.Fatalf("一時ファイル作成エラー: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// テストデータを書き込み
	if _, err := tmpFile.WriteString(testXBRL); err != nil {
		t.Fatalf("テストデータ書き込みエラー: %v", err)
	}
	tmpFile.Close()

	// テスト実行
	startDate, endDate := parser.ExtractAccountingPeriod(tmpFile.Name())

	expectedStartDate := "2024-05-01"
	expectedEndDate := "2025-04-30"

	if startDate != expectedStartDate {
		t.Errorf("開始日不一致: 期待=%s, 実際=%s", expectedStartDate, startDate)
	}
	if endDate != expectedEndDate {
		t.Errorf("終了日不一致: 期待=%s, 実際=%s", expectedEndDate, endDate)
	}
}

func TestXBRLParser_ExtractAccountingPeriod_FileNotFound(t *testing.T) {
	parser := NewXBRLParser()

	startDate, endDate := parser.ExtractAccountingPeriod("nonexistent_file.xbrl")

	if startDate != "" {
		t.Errorf("開始日が空でない: %s", startDate)
	}
	if endDate != "" {
		t.Errorf("終了日が空でない: %s", endDate)
	}
}

func TestXBRLParser_ParseAllXBRL(t *testing.T) {
	parser := NewXBRLParser()

	// テスト用のXBRLファイルを作成
	testXBRL := `<?xml version="1.0" encoding="UTF-8"?>
<xbrli:xbrl xmlns:jppfs_cor="http://disclosure.edinet-fsa.go.jp/taxonomy/jppfs/2024-11-01/jppfs_cor">
  <jppfs_cor:NetSales contextRef="CurrentYearDuration" unitRef="JPY">1000000000</jppfs_cor:NetSales>
  <jppfs_cor:ProfitLoss contextRef="CurrentYearDuration" unitRef="JPY">100000000</jppfs_cor:ProfitLoss>
  <jppfs_cor:TotalAssets contextRef="CurrentYearDuration" unitRef="JPY">2000000000</jppfs_cor:TotalAssets>
</xbrli:xbrl>`

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "test_*.xbrl")
	if err != nil {
		t.Fatalf("一時ファイル作成エラー: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// テストデータを書き込み
	if _, err := tmpFile.WriteString(testXBRL); err != nil {
		t.Fatalf("テストデータ書き込みエラー: %v", err)
	}
	tmpFile.Close()

	// テスト実行
	values, err := parser.ParseAllXBRL(tmpFile.Name())
	if err != nil {
		t.Fatalf("XBRL解析エラー: %v", err)
	}

	// 結果を検証（実際に抽出されたキーを使用）
	foundKeys := 0
	for k, v := range values {
		if strings.Contains(k, "NetSales") && v == "1000000000" {
			foundKeys++
		}
		if strings.Contains(k, "ProfitLoss") && v == "100000000" {
			foundKeys++
		}
		if strings.Contains(k, "TotalAssets") && v == "2000000000" {
			foundKeys++
		}
	}

	if foundKeys < 3 {
		t.Errorf("期待される値が見つかりません。見つかったキー数: %d", foundKeys)
		t.Logf("実際に抽出されたキー:")
		for k, v := range values {
			t.Logf("  %s = %s", k, v)
		}
	}
}

func TestXBRLParser_ParseAllXBRL_FileNotFound(t *testing.T) {
	parser := NewXBRLParser()

	_, err := parser.ParseAllXBRL("nonexistent_file.xbrl")
	if err == nil {
		t.Error("ファイルが存在しない場合、エラーが発生すべきです")
	}
}

func TestXBRLParser_ExtractPublicDocXBRL(t *testing.T) {
	parser := NewXBRLParser()

	// テスト用のZIPファイルを作成
	zipFile := "test_publicdoc.zip"
	defer os.Remove(zipFile)

	// 実際のZIPファイルを作成するのは複雑なので、
	// 存在しないファイルでテスト
	_, err := parser.ExtractPublicDocXBRL("nonexistent.zip")
	if err == nil {
		t.Error("存在しないZIPファイルの場合、エラーが発生すべきです")
	}
}

func TestXBRLParser_ExtractPublicDocXBRL_WithRealZip(t *testing.T) {
	parser := NewXBRLParser()

	// テスト用のZIPファイルを作成
	zipFile := "test_publicdoc.zip"
	defer os.Remove(zipFile)

	// 実際のZIPファイルを作成
	zipFileWriter, err := os.Create(zipFile)
	if err != nil {
		t.Fatalf("ZIPファイル作成エラー: %v", err)
	}

	zipWriter := zip.NewWriter(zipFileWriter)

	// PublicDocのXBRLファイルを追加
	xbrlWriter, err := zipWriter.Create("PublicDoc_Test.xbrl")
	if err != nil {
		t.Fatalf("XBRLファイル作成エラー: %v", err)
	}

	testXBRL := `<?xml version="1.0" encoding="UTF-8"?>
<xbrli:xbrl xmlns:jppfs_cor="http://disclosure.edinet-fsa.go.jp/taxonomy/jppfs/2024-11-01/jppfs_cor">
  <jppfs_cor:NetSales contextRef="CurrentYearDuration" unitRef="JPY">1000000000</jppfs_cor:NetSales>
</xbrli:xbrl>`

	_, err = xbrlWriter.Write([]byte(testXBRL))
	if err != nil {
		t.Fatalf("XBRLデータ書き込みエラー: %v", err)
	}

	// ZIPファイルを閉じる
	err = zipWriter.Close()
	if err != nil {
		t.Fatalf("ZIPファイルクローズエラー: %v", err)
	}
	zipFileWriter.Close()

	// テスト実行
	xbrlPath, err := parser.ExtractPublicDocXBRL(zipFile)
	if err != nil {
		t.Fatalf("XBRL抽出エラー: %v", err)
	}
	defer os.Remove(xbrlPath)

	if !strings.Contains(xbrlPath, "PublicDoc_Test.xbrl") {
		t.Errorf("XBRLファイル名不一致: 期待=PublicDoc_Test.xbrlを含む, 実際=%s", xbrlPath)
	}
}

func TestXBRLParser_ParseAllXBRL_Complex(t *testing.T) {
	parser := NewXBRLParser()

	// より複雑なXBRLファイルを作成
	testXBRL := `<?xml version="1.0" encoding="UTF-8"?>
<xbrli:xbrl xmlns:jppfs_cor="http://disclosure.edinet-fsa.go.jp/taxonomy/jppfs/2024-11-01/jppfs_cor"
            xmlns:xbrli="http://www.xbrl.org/2003/instance">
  <xbrli:context id="CurrentYearDuration">
    <xbrli:entity>
      <xbrli:identifier scheme="http://disclosure.edinet-fsa.go.jp">E00763-000</xbrli:identifier>
    </xbrli:entity>
    <xbrli:period>
      <xbrli:startDate>2024-04-01</xbrli:startDate>
      <xbrli:endDate>2025-03-31</xbrli:endDate>
    </xbrli:period>
  </xbrli:context>
  <jppfs_cor:NetSales contextRef="CurrentYearDuration" unitRef="JPY">1000000000</jppfs_cor:NetSales>
  <jppfs_cor:GrossProfit contextRef="CurrentYearDuration" unitRef="JPY">500000000</jppfs_cor:GrossProfit>
  <jppfs_cor:OperatingIncome contextRef="CurrentYearDuration" unitRef="JPY">200000000</jppfs_cor:OperatingIncome>
  <jppfs_cor:ProfitLoss contextRef="CurrentYearDuration" unitRef="JPY">100000000</jppfs_cor:ProfitLoss>
  <jppfs_cor:TotalAssets contextRef="CurrentYearDuration" unitRef="JPY">2000000000</jppfs_cor:TotalAssets>
  <jppfs_cor:NetSalesTextBlock contextRef="CurrentYearDuration">売上高の説明文</jppfs_cor:NetSalesTextBlock>
</xbrli:xbrl>`

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "test_complex_*.xbrl")
	if err != nil {
		t.Fatalf("一時ファイル作成エラー: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// テストデータを書き込み
	if _, err := tmpFile.WriteString(testXBRL); err != nil {
		t.Fatalf("テストデータ書き込みエラー: %v", err)
	}
	tmpFile.Close()

	// テスト実行
	values, err := parser.ParseAllXBRL(tmpFile.Name())
	if err != nil {
		t.Fatalf("XBRL解析エラー: %v", err)
	}

	// 結果を検証
	expectedCount := 6 // NetSales, GrossProfit, OperatingIncome, ProfitLoss, TotalAssets, NetSalesTextBlock
	if len(values) < expectedCount {
		t.Errorf("抽出された値の数が少なすぎます: 期待>=%d, 実際=%d", expectedCount, len(values))
	}

	// 特定の値が含まれているか確認
	foundValues := 0
	for k, v := range values {
		if strings.Contains(k, "NetSales") && !strings.Contains(k, "TextBlock") && v == "1000000000" {
			foundValues++
		}
		if strings.Contains(k, "GrossProfit") && v == "500000000" {
			foundValues++
		}
		if strings.Contains(k, "OperatingIncome") && v == "200000000" {
			foundValues++
		}
		if strings.Contains(k, "ProfitLoss") && v == "100000000" {
			foundValues++
		}
		if strings.Contains(k, "TotalAssets") && v == "2000000000" {
			foundValues++
		}
	}

	if foundValues < 4 {
		t.Errorf("期待される値が見つかりません。見つかった値数: %d", foundValues)
		t.Logf("実際に抽出されたキー:")
		for k, v := range values {
			t.Logf("  %s = %s", k, v)
		}
	}
}

func TestXBRLParser_ExtractAccountingPeriod_Complex(t *testing.T) {
	parser := NewXBRLParser()

	// より複雑なXBRLファイルを作成（複数のcontextを含む）
	testXBRL := `<?xml version="1.0" encoding="UTF-8"?>
<xbrli:xbrl xmlns:xbrli="http://www.xbrl.org/2003/instance">
  <xbrli:context id="CurrentYearDuration">
    <xbrli:entity>
      <xbrli:identifier scheme="http://disclosure.edinet-fsa.go.jp">E00763-000</xbrli:identifier>
    </xbrli:entity>
    <xbrli:period>
      <xbrli:startDate>2024-04-01</xbrli:startDate>
      <xbrli:endDate>2025-03-31</xbrli:endDate>
    </xbrli:period>
  </xbrli:context>
  <xbrli:context id="PreviousYearDuration">
    <xbrli:entity>
      <xbrli:identifier scheme="http://disclosure.edinet-fsa.go.jp">E00763-000</xbrli:identifier>
    </xbrli:entity>
    <xbrli:period>
      <xbrli:startDate>2023-04-01</xbrli:startDate>
      <xbrli:endDate>2024-03-31</xbrli:endDate>
    </xbrli:period>
  </xbrli:context>
</xbrli:xbrl>`

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "test_accounting_*.xbrl")
	if err != nil {
		t.Fatalf("一時ファイル作成エラー: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// テストデータを書き込み
	if _, err := tmpFile.WriteString(testXBRL); err != nil {
		t.Fatalf("テストデータ書き込みエラー: %v", err)
	}
	tmpFile.Close()

	// テスト実行
	startDate, endDate := parser.ExtractAccountingPeriod(tmpFile.Name())

	// 最初に見つかった期間を確認
	if startDate != "2024-04-01" {
		t.Errorf("開始日不一致: 期待=2024-04-01, 実際=%s", startDate)
	}
	if endDate != "2025-03-31" {
		t.Errorf("終了日不一致: 期待=2025-03-31, 実際=%s", endDate)
	}
} 