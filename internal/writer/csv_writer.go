package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"edinet-api-test/internal/config"
	"edinet-api-test/internal/models"
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

// ExtractFinancialValues 財務タグから値を抽出
func (c *CSVWriter) ExtractFinancialValues(values map[string]string) []string {
	var result []string
	
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
	
	return result
} 