#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
任天堂株式会社の財務モデリング
EDINETから取得した財務データを使用した財務分析と予測
"""

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from sklearn.linear_model import LinearRegression
from sklearn.metrics import r2_score, mean_squared_error
import warnings
warnings.filterwarnings('ignore')

# 日本語フォント設定
import matplotlib
matplotlib.rcParams['font.family'] = ['Hiragino Sans', 'Yu Gothic', 'Meiryo', 'Takao', 'IPAexGothic', 'IPAPGothic', 'VL PGothic', 'Noto Sans CJK JP']
matplotlib.rcParams['axes.unicode_minus'] = False

class NintendoFinancialModeling:
    def __init__(self, csv_file):
        """財務モデリングクラスの初期化"""
        self.df = pd.read_csv(csv_file)
        self.prepare_data()
        
    def prepare_data(self):
        """データの前処理"""
        # 数値列を適切な型に変換
        numeric_columns = ['売上高', '売上総利益', '営業利益', '経常利益', 
                          '税引前当期純利益', '当期純利益', '総資産', 
                          '流動資産', '固定資産', '総負債', '流動負債', 
                          '固定負債', '純資産', '資本金', '利益剰余金',
                          '営業CF', '投資CF', '財務CF', '現金及び現金同等物',
                          '1株当たり純資産', '自己資本比率', '配当金']
        
        for col in numeric_columns:
            if col in self.df.columns:
                self.df[col] = pd.to_numeric(self.df[col], errors='coerce')
        
        # 日付をdatetime型に変換
        self.df['日付'] = pd.to_datetime(self.df['日付'])
        
        # 会計期間の修正
        self.correct_fiscal_period()
        
        # 年度を抽出（会計期間から）
        self.df['年度'] = self.df['会計期間'].str.extract(r'(\d{4})').astype(int)
        
        # 有価証券報告書のみを抽出（年度データとして使用）
        self.annual_data = self.df[self.df['文書タイプ'] == '有価証券報告書'].copy()
        
        print(f"データ準備完了: {len(self.df)}件の全データ, {len(self.annual_data)}件の年度データ")
        print("会計期間を修正しました:")
        for _, row in self.annual_data.iterrows():
            print(f"  提出日: {row['日付'].strftime('%Y-%m-%d')} → 会計期間: {row['会計期間']}")
    
    def correct_fiscal_period(self):
        """会計期間を提出日から正しく計算"""
        print("会計期間の修正中...")
        
        for idx, row in self.df.iterrows():
            submit_date = row['日付']
            doc_type = row['文書タイプ']
            
            # 有価証券報告書の場合：提出年の前年度が会計期間
            if doc_type == '有価証券報告書':
                fiscal_year = submit_date.year - 1
                self.df.at[idx, '会計期間'] = f"{fiscal_year}年度"
            
            # 四半期報告書の場合：提出年の前年度の四半期
            elif doc_type == '四半期報告書':
                fiscal_year = submit_date.year - 1
                # 提出月から四半期を推定
                month = submit_date.month
                if month <= 3:  # 1-3月提出 → 前年度第3四半期
                    quarter = 3
                elif month <= 6:  # 4-6月提出 → 前年度第4四半期
                    quarter = 4
                elif month <= 9:  # 7-9月提出 → 当年度第1四半期
                    quarter = 1
                    fiscal_year = submit_date.year
                else:  # 10-12月提出 → 当年度第2四半期
                    quarter = 2
                    fiscal_year = submit_date.year
                
                self.df.at[idx, '会計期間'] = f"{fiscal_year}年度第{quarter}四半期"
    
    def calculate_financial_ratios(self):
        """財務比率の計算"""
        print("\n=== 財務比率の計算 ===")
        
        # 年度データのみで財務比率を計算
        ratios = self.annual_data.copy()
        
        # データの検証と修正（直接実行）
        print("データ検証中...")
        for idx, row in ratios.iterrows():
            print(f"{row['年度']}年度: 総資産={row['総資産']:,.0f}, 純資産={row['純資産']:,.0f}")
            
            # 総資産が欠損または異常値の場合、流動資産+固定資産で推定
            if pd.isna(row['総資産']) or row['総資産'] <= 0:
                if not pd.isna(row['流動資産']) and not pd.isna(row['固定資産']):
                    estimated_total_assets = row['流動資産'] + row['固定資産']
                    ratios.at[idx, '総資産'] = estimated_total_assets
                    print(f"  総資産を修正: {estimated_total_assets:,.0f}")
            
            # 純資産が異常に小さい場合（総資産の10%未満）も修正
            if not pd.isna(row['総資産']) and not pd.isna(row['総負債']) and row['純資産'] < row['総資産'] * 0.1:
                print(f"  純資産が異常に小さい: {row['純資産']:,.0f} (総資産の{(row['純資産']/row['総資産']*100):.1f}%)")
                estimated_equity = row['総資産'] - abs(row['総負債'])
                ratios.at[idx, '純資産'] = estimated_equity
                print(f"  純資産を推定値に修正: {estimated_equity:,.0f}")
        
        # 収益性比率
        ratios['売上高利益率'] = (ratios['営業利益'] / ratios['売上高']) * 100
        ratios['当期純利益率'] = (ratios['当期純利益'] / ratios['売上高']) * 100
        
        # ROAの計算（総資産が有効な場合のみ）
        ratios['ROA'] = np.where(
            (ratios['総資産'] > 0) & (ratios['総資産'].notna()),
            (ratios['当期純利益'] / ratios['総資産']) * 100,
            np.nan
        )
        
        # ROEの計算（純資産が有効な場合のみ、異常値を制限）
        ratios['ROE'] = np.where(
            (ratios['純資産'] > 10000000000) & (ratios['純資産'].notna()),  # 100億円以上
            np.minimum((ratios['当期純利益'] / ratios['純資産']) * 100, 50),  # 最大50%に制限
            np.nan
        )
        
        # 効率性比率
        ratios['総資産回転率'] = np.where(
            (ratios['総資産'] > 0) & (ratios['総資産'].notna()),
            ratios['売上高'] / ratios['総資産'],
            np.nan
        )
        ratios['固定資産回転率'] = np.where(
            (ratios['固定資産'] > 0) & (ratios['固定資産'].notna()),
            ratios['売上高'] / ratios['固定資産'],
            np.nan
        )
        
        # 安全性比率
        ratios['流動比率'] = np.where(
            (ratios['流動負債'] > 0) & (ratios['流動負債'].notna()),
            (ratios['流動資産'] / ratios['流動負債']) * 100,
            np.nan
        )
        ratios['固定比率'] = np.where(
            (ratios['純資産'] > 0) & (ratios['純資産'].notna()),
            (ratios['固定資産'] / ratios['純資産']) * 100,
            np.nan
        )
        
        # 成長率
        ratios['売上高成長率'] = ratios['売上高'].pct_change() * 100
        ratios['営業利益成長率'] = ratios['営業利益'].pct_change() * 100
        ratios['当期純利益成長率'] = ratios['当期純利益'].pct_change() * 100
        
        self.ratios = ratios
        return ratios
    
    def validate_and_fix_financial_data(self, df):
        """財務データの検証と修正"""
        print("財務データの検証中...")
        
        for idx, row in df.iterrows():
            print(f"\n{row['年度']}年度のデータ検証:")
            print(f"  元データ - 総資産: {row['総資産']:,.0f}, 純資産: {row['純資産']:,.0f}, 総負債: {row['総負債']:,.0f}")
            
            # 総資産の検証
            if pd.isna(row['総資産']) or row['総資産'] <= 0:
                print(f"  総資産が無効: {row['総資産']}")
                # 流動資産 + 固定資産で推定
                if not pd.isna(row['流動資産']) and not pd.isna(row['固定資産']):
                    estimated_total_assets = row['流動資産'] + row['固定資産']
                    df.at[idx, '総資産'] = estimated_total_assets
                    print(f"  総資産を推定値に修正: {estimated_total_assets:,.0f}")
            
            # 純資産の検証
            if pd.isna(row['純資産']) or row['純資産'] <= 0:
                print(f"  純資産が無効: {row['純資産']}")
                # 総資産 - 総負債で推定
                if not pd.isna(row['総資産']) and not pd.isna(row['総負債']):
                    estimated_equity = row['総資産'] - abs(row['総負債'])
                    df.at[idx, '純資産'] = estimated_equity
                    print(f"  純資産を推定値に修正: {estimated_equity:,.0f}")
            
            # 純資産が異常に小さい場合（総資産の10%未満）も修正
            elif not pd.isna(row['総資産']) and not pd.isna(row['総負債']) and row['純資産'] < row['総資産'] * 0.1:
                print(f"  純資産が異常に小さい: {row['純資産']:,.0f} (総資産の{(row['純資産']/row['総資産']*100):.1f}%)")
                estimated_equity = row['総資産'] - abs(row['総負債'])
                df.at[idx, '純資産'] = estimated_equity
                print(f"  純資産を推定値に修正: {estimated_equity:,.0f}")
            
            # 負債の検証（負の値を正に修正）
            if not pd.isna(row['総負債']) and row['総負債'] < 0:
                print(f"  総負債が負の値: {row['総負債']}")
                df.at[idx, '総負債'] = abs(row['総負債'])
                print(f"  総負債を正の値に修正: {abs(row['総負債']):,.0f}")
            
            if not pd.isna(row['流動負債']) and row['流動負債'] < 0:
                print(f"  流動負債が負の値: {row['流動負債']}")
                df.at[idx, '流動負債'] = abs(row['流動負債'])
                print(f"  流動負債を正の値に修正: {abs(row['流動負債']):,.0f}")
            
            if not pd.isna(row['固定負債']) and row['固定負債'] < 0:
                print(f"  固定負債が負の値: {row['固定負債']}")
                df.at[idx, '固定負債'] = abs(row['固定負債'])
                print(f"  固定負債を正の値に修正: {abs(row['固定負債']):,.0f}")
            
            # 修正後の値を表示
            print(f"  修正後 - 総資産: {df.at[idx, '総資産']:,.0f}, 純資産: {df.at[idx, '純資産']:,.0f}")
            
            # ROEの計算例を表示
            if not pd.isna(df.at[idx, '当期純利益']) and not pd.isna(df.at[idx, '純資産']) and df.at[idx, '純資産'] > 0:
                roe = (df.at[idx, '当期純利益'] / df.at[idx, '純資産']) * 100
                print(f"  ROE計算例: {df.at[idx, '当期純利益']:,.0f} / {df.at[idx, '純資産']:,.0f} * 100 = {roe:.2f}%")
    
    def plot_financial_trends(self):
        """財務指標の推移をプロット"""
        print("\n=== 財務指標の推移分析 ===")
        
        fig, axes = plt.subplots(2, 2, figsize=(15, 12))
        fig.suptitle('任天堂株式会社 財務指標の推移 (2020-2024)', fontsize=16)
        
        # 売上高と利益の推移
        ax1 = axes[0, 0]
        ax1.plot(self.annual_data['年度'], self.annual_data['売上高'] / 1e12, 
                marker='o', label='売上高', linewidth=2)
        ax1.plot(self.annual_data['年度'], self.annual_data['営業利益'] / 1e12, 
                marker='s', label='営業利益', linewidth=2)
        ax1.plot(self.annual_data['年度'], self.annual_data['当期純利益'] / 1e12, 
                marker='^', label='当期純利益', linewidth=2)
        ax1.set_title('売上高・利益の推移')
        ax1.set_ylabel('金額 (兆円)')
        ax1.legend()
        ax1.grid(True, alpha=0.3)
        
        # 収益性比率の推移
        ax2 = axes[0, 1]
        ax2.plot(self.ratios['年度'], self.ratios['売上高利益率'], 
                marker='o', label='売上高利益率', linewidth=2)
        ax2.plot(self.ratios['年度'], self.ratios['当期純利益率'], 
                marker='s', label='当期純利益率', linewidth=2)
        ax2.set_title('収益性比率の推移')
        ax2.set_ylabel('比率 (%)')
        ax2.legend()
        ax2.grid(True, alpha=0.3)
        
        # ROA・ROEの推移
        ax3 = axes[1, 0]
        ax3.plot(self.ratios['年度'], self.ratios['ROA'], 
                marker='o', label='ROA', linewidth=2)
        ax3.plot(self.ratios['年度'], self.ratios['ROE'], 
                marker='s', label='ROE', linewidth=2)
        ax3.set_title('ROA・ROEの推移')
        ax3.set_ylabel('比率 (%)')
        ax3.legend()
        ax3.grid(True, alpha=0.3)
        
        # 自己資本比率の推移
        ax4 = axes[1, 1]
        ax4.plot(self.annual_data['年度'], self.annual_data['自己資本比率'] * 100, 
                marker='o', color='green', linewidth=2)
        ax4.set_title('自己資本比率の推移')
        ax4.set_ylabel('比率 (%)')
        ax4.grid(True, alpha=0.3)
        
        plt.tight_layout()
        plt.savefig('nintendo_financial_trends.png', dpi=300, bbox_inches='tight')
        plt.close()
    
    def build_forecast_model(self):
        """売上高予測モデルの構築"""
        print("\n=== 売上高予測モデルの構築 ===")
        
        # 時系列データの準備
        X = self.annual_data['年度'].values.reshape(-1, 1)
        y = self.annual_data['売上高'].values
        
        # 線形回帰モデルの構築
        model = LinearRegression()
        model.fit(X, y)
        
        # 予測
        y_pred = model.predict(X)
        
        # モデル評価
        r2 = r2_score(y, y_pred)
        rmse = np.sqrt(mean_squared_error(y, y_pred))
        
        print(f"決定係数 (R²): {r2:.4f}")
        print(f"RMSE: {rmse:,.0f} 百万円")
        print(f"傾き (年間成長率): {model.coef_[0]:,.0f} 百万円/年")
        
        # 将来予測
        future_years = np.array([2025, 2026, 2027, 2028, 2029]).reshape(-1, 1)
        future_sales = model.predict(future_years)
        
        print("\n将来予測 (売上高):")
        for year, sales in zip(future_years.flatten(), future_sales):
            print(f"{year}年: {sales:,.0f} 百万円")
        
        # 予測結果のプロット
        plt.figure(figsize=(12, 8))
        plt.plot(X, y, 'o-', label='実績値', linewidth=2, markersize=8)
        plt.plot(X, y_pred, 'r--', label='予測値', linewidth=2)
        plt.plot(future_years, future_sales, 'g--', label='将来予測', linewidth=2)
        
        plt.title('任天堂株式会社 売上高予測モデル', fontsize=16)
        plt.xlabel('年度')
        plt.ylabel('売上高 (百万円)')
        plt.legend()
        plt.grid(True, alpha=0.3)
        plt.savefig('nintendo_sales_forecast.png', dpi=300, bbox_inches='tight')
        plt.close()
        
        return model, future_sales
    
    def calculate_valuation_metrics(self):
        """企業価値評価指標の計算"""
        print("\n=== 企業価値評価指標 ===")
        
        latest_data = self.annual_data.iloc[-1]
        
        # 基本的な評価指標
        print(f"最新年度: {latest_data['年度']}")
        print(f"売上高: {latest_data['売上高']:,.0f} 百万円")
        print(f"営業利益: {latest_data['営業利益']:,.0f} 百万円")
        print(f"当期純利益: {latest_data['当期純利益']:,.0f} 百万円")
        print(f"総資産: {latest_data['総資産']:,.0f} 百万円")
        print(f"純資産: {latest_data['純資産']:,.0f} 百万円")
        
        # 収益性指標
        print(f"\n収益性指標:")
        print(f"売上高利益率: {(latest_data['営業利益'] / latest_data['売上高'] * 100):.2f}%")
        print(f"ROA: {(latest_data['当期純利益'] / latest_data['総資産'] * 100):.2f}%")
        print(f"ROE: {(latest_data['当期純利益'] / latest_data['純資産'] * 100):.2f}%")
        
        # 安全性指標
        print(f"\n安全性指標:")
        print(f"自己資本比率: {latest_data['自己資本比率'] * 100:.2f}%")
        print(f"流動比率: {(latest_data['流動資産'] / latest_data['流動負債'] * 100):.2f}%")
        
        # 成長性指標
        if len(self.annual_data) > 1:
            prev_data = self.annual_data.iloc[-2]
            sales_growth = ((latest_data['売上高'] - prev_data['売上高']) / prev_data['売上高'] * 100)
            profit_growth = ((latest_data['当期純利益'] - prev_data['当期純利益']) / prev_data['当期純利益'] * 100)
            
            print(f"\n成長性指標:")
            print(f"売上高成長率: {sales_growth:.2f}%")
            print(f"当期純利益成長率: {profit_growth:.2f}%")
    
    def generate_financial_report(self):
        """財務分析レポートの生成"""
        print("\n=== 財務分析レポート ===")
        
        # 財務比率の計算
        self.calculate_financial_ratios()
        
        # 財務指標の推移プロット
        self.plot_financial_trends()
        
        # 予測モデルの構築
        self.build_forecast_model()
        
        # 企業価値評価指標の計算
        self.calculate_valuation_metrics()
        
        # レポートの保存
        self.save_financial_report()
    
    def save_financial_report(self):
        """財務分析レポートをCSVに保存"""
        # 財務比率データをCSVに保存
        self.ratios.to_csv('nintendo_financial_ratios.csv', index=False, encoding='utf-8-sig')
        print("\n財務分析レポートを保存しました:")
        print("- nintendo_financial_ratios.csv: 財務比率データ")
        print("- nintendo_financial_trends.png: 財務指標推移グラフ")
        print("- nintendo_sales_forecast.png: 売上高予測グラフ")

def main():
    """メイン実行関数"""
    print("任天堂株式会社 財務モデリング開始")
    print("=" * 50)
    
    # 財務モデリングの実行
    model = NintendoFinancialModeling('nintendo_.csv')
    model.generate_financial_report()
    
    print("\n財務モデリング完了！")

if __name__ == "__main__":
    main() 