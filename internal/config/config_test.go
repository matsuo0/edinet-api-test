package config

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Success(t *testing.T) {
	// テスト用の環境変数を設定
	os.Setenv("EDINET_API_KEY", "test-api-key")
	defer os.Unsetenv("EDINET_API_KEY")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("設定読み込みエラー: %v", err)
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey不一致: 期待=test-api-key, 実際=%s", cfg.APIKey)
	}
	if cfg.StartDate != "2025-07-10" {
		t.Errorf("StartDate不一致: 期待=2025-07-10, 実際=%s", cfg.StartDate)
	}
	if cfg.EndDate != "2025-07-16" {
		t.Errorf("EndDate不一致: 期待=2025-07-16, 実際=%s", cfg.EndDate)
	}
	if cfg.TargetSecCode != "" {
		t.Errorf("TargetSecCode不一致: 期待=空文字列, 実際=%s", cfg.TargetSecCode)
	}
	if cfg.OutputFile != "xbrl_financial_items.csv" {
		t.Errorf("OutputFile不一致: 期待=xbrl_financial_items.csv, 実際=%s", cfg.OutputFile)
	}
}

func TestLoadConfig_NoAPIKey(t *testing.T) {
	// 環境変数をクリア
	os.Unsetenv("EDINET_API_KEY")

	_, err := LoadConfig()
	if err == nil {
		t.Error("APIキーが設定されていない場合、エラーが発生すべきです")
	}

	configErr, ok := err.(*ConfigError)
	if !ok {
		t.Error("ConfigError型のエラーが返されるべきです")
	}

	expectedMessage := "EDINET_API_KEYが設定されていません。.envファイルを確認してください。"
	if configErr.Message != expectedMessage {
		t.Errorf("エラーメッセージ不一致: 期待=%s, 実際=%s", expectedMessage, configErr.Message)
	}
}

func TestConfig_GetDateRange_Success(t *testing.T) {
	cfg := &Config{
		StartDate: "2025-01-01",
		EndDate:   "2025-01-31",
	}

	start, end, err := cfg.GetDateRange()
	if err != nil {
		t.Fatalf("日付範囲取得エラー: %v", err)
	}

	expectedStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	if !start.Equal(expectedStart) {
		t.Errorf("開始日不一致: 期待=%v, 実際=%v", expectedStart, start)
	}
	if !end.Equal(expectedEnd) {
		t.Errorf("終了日不一致: 期待=%v, 実際=%v", expectedEnd, end)
	}
}

func TestConfig_GetDateRange_InvalidStartDate(t *testing.T) {
	cfg := &Config{
		StartDate: "invalid-date",
		EndDate:   "2025-01-31",
	}

	_, _, err := cfg.GetDateRange()
	if err == nil {
		t.Error("無効な開始日の場合、エラーが発生すべきです")
	}
}

func TestConfig_GetDateRange_InvalidEndDate(t *testing.T) {
	cfg := &Config{
		StartDate: "2025-01-01",
		EndDate:   "invalid-date",
	}

	_, _, err := cfg.GetDateRange()
	if err == nil {
		t.Error("無効な終了日の場合、エラーが発生すべきです")
	}
}

func TestConfigError_Error(t *testing.T) {
	configErr := &ConfigError{
		Message: "テストエラーメッセージ",
	}

	errorMessage := configErr.Error()
	if errorMessage != "テストエラーメッセージ" {
		t.Errorf("エラーメッセージ不一致: 期待=テストエラーメッセージ, 実際=%s", errorMessage)
	}
}

func TestJapaneseHeaders_Length(t *testing.T) {
	expectedLength := 28 // 実際のヘッダー数
	if len(JapaneseHeaders) != expectedLength {
		t.Errorf("日本語ヘッダーの長さ不一致: 期待=%d, 実際=%d", expectedLength, len(JapaneseHeaders))
	}

	// 最初の5つのヘッダーを確認
	expectedFirstHeaders := []string{"日付", "証券コード", "会社名", "文書タイプ", "会計期間"}
	for i, expected := range expectedFirstHeaders {
		if JapaneseHeaders[i] != expected {
			t.Errorf("ヘッダー[%d]不一致: 期待=%s, 実際=%s", i, expected, JapaneseHeaders[i])
		}
	}
}

func TestFinancialTags_Length(t *testing.T) {
	expectedLength := 23 // 実際の財務タグの数
	if len(FinancialTags) != expectedLength {
		t.Errorf("財務タグの長さ不一致: 期待=%d, 実際=%d", expectedLength, len(FinancialTags))
	}

	// 最初のタグを確認
	expectedFirstTag := "jppfs_cor:NetSales"
	if FinancialTags[0] != expectedFirstTag {
		t.Errorf("最初の財務タグ不一致: 期待=%s, 実際=%s", expectedFirstTag, FinancialTags[0])
	}

	// 最後のタグを確認
	expectedLastTag := "jppfs_cor:DividendsFromSurplus"
	if FinancialTags[len(FinancialTags)-1] != expectedLastTag {
		t.Errorf("最後の財務タグ不一致: 期待=%s, 実際=%s", expectedLastTag, FinancialTags[len(FinancialTags)-1])
	}
}

func TestLoadConfig_WithCommandLineArgs(t *testing.T) {
	// テスト用の環境変数を設定
	os.Setenv("EDINET_API_KEY", "test-api-key")
	defer os.Unsetenv("EDINET_API_KEY")

	// コマンドライン引数を設定
	os.Args = []string{"test", "-start", "2025-01-01", "-end", "2025-01-31", "-code", "6758", "-output", "test_output.csv"}
	
	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("設定読み込みエラー: %v", err)
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey不一致: 期待=test-api-key, 実際=%s", cfg.APIKey)
	}
	if cfg.StartDate != "2025-01-01" {
		t.Errorf("StartDate不一致: 期待=2025-01-01, 実際=%s", cfg.StartDate)
	}
	if cfg.EndDate != "2025-01-31" {
		t.Errorf("EndDate不一致: 期待=2025-01-31, 実際=%s", cfg.EndDate)
	}
	if cfg.TargetSecCode != "67580" {
		t.Errorf("TargetSecCode不一致: 期待=67580, 実際=%s", cfg.TargetSecCode)
	}
	if cfg.OutputFile != "test_output.csv" {
		t.Errorf("OutputFile不一致: 期待=test_output.csv, 実際=%s", cfg.OutputFile)
	}
	if cfg.QuarterOnly != false {
		t.Errorf("QuarterOnly不一致: 期待=false, 実際=%t", cfg.QuarterOnly)
	}
}

func TestLoadConfig_WithQuarterOnly(t *testing.T) {
	// テスト用の環境変数を設定
	os.Setenv("EDINET_API_KEY", "test-api-key")
	defer os.Unsetenv("EDINET_API_KEY")

	// コマンドライン引数を設定
	os.Args = []string{"test", "-quarter"}
	
	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("設定読み込みエラー: %v", err)
	}

	if cfg.QuarterOnly != true {
		t.Errorf("QuarterOnly不一致: 期待=true, 実際=%t", cfg.QuarterOnly)
	}
} 