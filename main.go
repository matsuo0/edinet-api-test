package main

import (
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
	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .envファイルを読み込めませんでした: %v", err)
	}

	// 設定を読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	// 日付範囲でループ
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		
		// 文書一覧を取得
		docList, err := edinetAPI.GetDocuments(d)
		if err != nil {
			log.Printf("文書一覧取得エラー (%s): %v", dateStr, err)
			continue
		}

		// 文書をフィルタリング
		filteredDocs := api.FilterDocuments(docList.Results, cfg.TargetSecCode)

		// 各文書を処理
		for _, doc := range filteredDocs {
			if err := processDocument(doc, dateStr, edinetAPI, xbrlParser, csvWriter); err != nil {
				log.Printf("文書処理エラー (%s): %v", doc.DocID, err)
				continue
			}
		}
	}

	fmt.Printf("%s に主要財務項目を出力しました。\n", cfg.OutputFile)
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