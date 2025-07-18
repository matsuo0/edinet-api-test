# EDINET API XBRL財務データ抽出ツール

EDINET APIを使用して有価証券報告書のXBRLデータを取得し、主要財務項目をCSVに出力するGo言語製ツールです。

## 機能

- EDINET APIから指定期間の文書一覧を取得
- 有価証券報告書・四半期報告書のXBRLファイルを自動ダウンロード
- 四半期報告書のみを対象とした四半期データ取得
- XBRLファイルから主要財務項目を抽出
- 日本語ヘッダー付きCSVファイルに出力
- 期間指定による複数日分の一括処理
- 証券コード指定による特定企業のデータ取得

## 必要な環境

- Go 1.19以上
- EDINET APIキー

## セットアップ

1. リポジトリをクローン
```bash
git clone <repository-url>
cd edinet-api-test
```

2. 依存関係をインストール
```bash
go mod tidy
```

3. 環境変数を設定
```bash
# .envファイルを作成
echo "EDINET_API_KEY=your-api-key-here" > .env
```

## 使用方法

### 基本的な使用方法

```bash
# デフォルト設定で実行（2025-07-10から2025-07-16、証券コード40260）
go run main.go
```

### 期間指定での実行

```bash
# 2025年1月のデータを取得
go run main.go -start 2025-01-01 -end 2025-01-31

# 特定の証券コードを指定
go run main.go -start 2025-01-01 -end 2025-01-31 -code 6758

# 出力ファイル名を指定
go run main.go -start 2025-01-01 -end 2025-01-31 -code 6758 -output toshiba_data.csv
```

### コマンドラインオプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `-start` | 開始日 (YYYY-MM-DD形式) | 2025-07-10 |
| `-end` | 終了日 (YYYY-MM-DD形式) | 2025-07-16 |
| `-code` | 対象証券コード | 40260 |
| `-output` | 出力ファイル名 | xbrl_financial_items.csv |
| `-quarter` | 四半期報告書のみを対象にする | false |

### 主要企業の証券コード例

| 企業名 | 4桁証券コード | EDINET証券コード | 使用例 |
|--------|-------------|----------------|--------|
| ソニー | 6758 | 67580 | `-code 6758` または `-code 67580` |
| トヨタ自動車 | 7203 | 72030 | `-code 7203` または `-code 72030` |
| 日産自動車 | 7201 | 72010 | `-code 7201` または `-code 72010` |
| 本田技研工業 | 7267 | 72670 | `-code 7267` または `-code 72670` |
| 任天堂 | 7974 | 79740 | `-code 7974` または `-code 79740` |
| ソフトバンクグループ | 9984 | 99840 | `-code 9984` または `-code 99840` |
| NTT | 9432 | 94320 | `-code 9432` または `-code 94320` |
| 三菱UFJフィナンシャル・グループ | 8306 | 83060 | `-code 8306` または `-code 83060` |
| みずほフィナンシャルグループ | 8411 | 84110 | `-code 8411` または `-code 84110` |

### 証券コードについて

- **4桁証券コード**: 一般的な証券コード（例：6758）
- **EDINET証券コード**: EDINETで使用される5桁の証券コード（例：67580）
- **自動変換**: 4桁の証券コードを指定すると、自動的に末尾に0を付けて5桁に変換されます
- **両方対応**: 4桁または5桁のどちらでも指定可能です

### 使用例

```bash
# トヨタ自動車の2024年10月のデータを取得
go run main.go -start 2024-10-01 -end 2024-10-31 -code 7203 -output toyota_202410.csv

# 任天堂の2024年9月のデータを取得
go run main.go -start 2024-09-01 -end 2024-09-30 -code 7974 -output nintendo_202409.csv

# 四半期報告書のみを取得（年度報告書は除外）
go run main.go -start 2024-01-01 -end 2024-12-31 -code 6758 -quarter -output toshiba_quarterly_2024.csv

# 証券コードのみ指定（期間はデフォルト）
go run main.go -code 6758

# ヘルプを表示
go run main.go -h
```

## 出力される財務項目

CSVファイルには以下の財務項目が含まれます：

1. **基本情報**: 日付、証券コード、会社名、文書タイプ、会計期間
2. **損益計算書項目**: 売上高、売上総利益、営業利益、経常利益、税引前当期純利益、当期純利益、1株当たり当期純利益
3. **貸借対照表項目**: 総資産、流動資産、固定資産、総負債、流動負債、固定負債、純資産、資本金、利益剰余金
4. **キャッシュフロー項目**: 営業CF、投資CF、財務CF、現金及び現金同等物
5. **その他**: 1株当たり純資産、自己資本比率、配当金

## アーキテクチャ

```
edinet-api-test/
├── main.go                 # メインエントリーポイント
├── internal/
│   ├── models/            # データ構造定義
│   ├── config/            # 設定管理
│   ├── api/               # EDINET APIクライアント
│   ├── parser/            # XBRLファイル解析
│   └── writer/            # CSV出力
├── go.mod
├── go.sum
└── README.md
```

## テスト

```bash
# 全テストを実行
go test ./...

# カバレッジを確認
go test ./... -cover

# 詳細なカバレッジレポート
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 注意事項

- EDINET APIキーが必要です（[EDINET API](https://disclosure.edinet-fsa.go.jp/guide/guide.html)で取得）
- EDINET APIキーの取得はこちらから [EDINET API KEY](https://api.edinet-fsa.go.jp/api/auth/index.aspx?mode=1)
- APIの利用制限にご注意ください
- 大量のデータを取得する場合は、適切な間隔を空けて実行してください
- 出力されるCSVファイルは一時ファイルと一緒に作成されます

## トラブルシューティング

### よくあるエラー

1. **APIキーエラー**
   ```
   EDINET_API_KEYが設定されていません。.envファイルを確認してください。
   ```
   → `.env`ファイルに正しいAPIキーを設定してください

2. **日付形式エラー**
   ```
   日付範囲の取得エラー: parsing time "invalid-date"
   ```
   → 日付は`YYYY-MM-DD`形式で指定してください

3. **ファイル出力エラー**
   ```
   CSV作成エラー: open /invalid/path: permission denied
   ```
   → 出力先ディレクトリの権限を確認してください

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
