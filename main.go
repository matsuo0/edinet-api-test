package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"

	"edinet-api-test/internal/api"
	"edinet-api-test/internal/config"
	"edinet-api-test/internal/parser"
	"edinet-api-test/internal/writer"
	"edinet-api-test/internal/models"
)

func main() {
	// ヘルプメッセージを設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "EDINET API XBRL財務データ抽出ツール\n\n")
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s [オプション]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "オプション:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n例:\n")
		fmt.Fprintf(os.Stderr, "  %s -start 2025-01-01 -end 2025-01-31 -code 40260\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -start 2024-12-01 -end 2024-12-31 -code 6758 -output toshiba_data.csv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n注意: EDINET_API_KEY環境変数が設定されている必要があります。\n")
	}

	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .envファイルを読み込めませんでした: %v", err)
	}

	// 設定を読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 設定情報を表示
	fmt.Printf("設定情報:\n")
	fmt.Printf("  開始日: %s\n", cfg.StartDate)
	fmt.Printf("  終了日: %s\n", cfg.EndDate)
	fmt.Printf("  対象証券コード: %s\n", cfg.TargetSecCode)
	fmt.Printf("  出力ファイル: %s\n", cfg.OutputFile)
	if cfg.QuarterOnly {
		fmt.Printf("  対象文書: 四半期報告書のみ\n")
	} else {
		fmt.Printf("  対象文書: 有価証券報告書・四半期報告書\n")
	}
	fmt.Printf("\n")

	// 日付範囲を取得
	start, end, err := cfg.GetDateRange()
	if err != nil {
		log.Fatalf("日付範囲の取得エラー: %v", err)
	}

	// 各コンポーネントを初期化
	edinetAPI := api.NewEdinetAPI(cfg.APIKey)
	xbrlParser := parser.NewXBRLParser()
	csvWriter, err := writer.NewCSVWriter(cfg.OutputFile)
	if err != nil {
		log.Fatalf("CSV出力器の初期化エラー: %v", err)
	}
	defer csvWriter.Close()

	// ヘッダーを書き込み
	if err := csvWriter.WriteHeader(); err != nil {
		log.Fatalf("ヘッダー書き込みエラー: %v", err)
	}

	// 処理件数をカウント
	processedCount := 0

	// 日付範囲でループ
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		fmt.Printf("処理中: %s\n", dateStr)
		
		// 文書一覧を取得
		docList, err := edinetAPI.GetDocuments(d)
		if err != nil {
			log.Printf("文書一覧取得エラー (%s): %v", dateStr, err)
			continue
		}

		// 文書をフィルタリング
		filteredDocs := api.FilterDocuments(docList.Results, cfg.TargetSecCode, cfg.QuarterOnly)

		// 各文書を処理
		for _, doc := range filteredDocs {
			if err := processDocument(doc, dateStr, edinetAPI, xbrlParser, csvWriter); err != nil {
				log.Printf("文書処理エラー (%s): %v", doc.DocID, err)
				continue
			}
			processedCount++
		}
	}

	fmt.Printf("\n処理完了: %d件の文書を処理し、%s に主要財務項目を出力しました。\n", processedCount, cfg.OutputFile)
}

// processDocument 個別文書を処理
func processDocument(doc models.DocInfo, dateStr string, edinetAPI *api.EdinetAPI, xbrlParser *parser.XBRLParser, csvWriter *writer.CSVWriter) error {
	// XBRL ZIPをダウンロード
	zipData, err := edinetAPI.DownloadXBRLZip(doc.DocID)
	if err != nil {
		return fmt.Errorf("ZIPダウンロード失敗: %v", err)
	}

	// 一時ZIPファイルを作成
	zipFile := doc.DocID + ".zip"
	if err := ioutil.WriteFile(zipFile, zipData, 0644); err != nil {
		return fmt.Errorf("ZIPファイル保存失敗: %v", err)
	}
	defer os.Remove(zipFile)

	// XBRLファイルを抽出
	xbrlPath, err := xbrlParser.ExtractPublicDocXBRL(zipFile)
	if err != nil {
		return fmt.Errorf("XBRL抽出失敗: %v", err)
	}
	defer os.Remove(xbrlPath)

	// XBRLファイルを解析
	values, err := xbrlParser.ParseAllXBRL(xbrlPath)
	if err != nil {
		return fmt.Errorf("XBRLパース失敗: %v", err)
	}

	// 会計期間を抽出
	startDate, endDate := xbrlParser.ExtractAccountingPeriod(xbrlPath)
	quarterInfo := xbrlParser.GetQuarterInfo(startDate, endDate)

	// 文書タイプ名を取得
	docTypeName := xbrlParser.GetDocTypeName(doc.DocTypeCode)

	// 財務値を抽出
	financialValues := csvWriter.ExtractFinancialValues(values)

	// 行データを作成
	row := []string{
		dateStr,           // 日付
		doc.SecCode,       // 証券コード
		doc.FilerName,     // 会社名
		docTypeName,       // 文書タイプ
		quarterInfo,       // 会計期間
	}
	row = append(row, financialValues...)

	// CSVに書き込み
	return csvWriter.WriteRow(row)
}