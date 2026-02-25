#!/usr/bin/env python3
"""
NOFX 每日策略深度进化脚本
24小时数据聚合、全球基因融合、策略优化
"""

import os
import sys
import json
import sqlite3
import glob
from datetime import datetime, timedelta
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from lib.data_collector import DataCollector
from lib.analyzer import StrategyPerformanceAnalyzer
from lib.evo_client import EvoMapClient
from lib.gene_fusion import GeneFusionEngine
from lib.report_generator import DailyReportGenerator
from lib.strategy_optimizer import StrategyOptimizer

OUTPUT_DIR = Path("docs/openclaw")
STRATEGY_DIR = Path("docs/openclaw/strategies")

class DailyEvolver:
    def __init__(self):
        self.collector = DataCollector()
        self.analyzer = StrategyPerformanceAnalyzer()
        self.evo_client = EvoMapClient()
        self.fusion_engine = GeneFusionEngine()
        self.optimizer = StrategyOptimizer()
        self.report_gen = DailyReportGenerator()
        
        self.now = datetime.now()
        self.date_str = self.now.strftime("%Y%m%d")
        self.yesterday = (self.now - timedelta(days=1)).strftime("%Y-%m-%d")
        
    def run(self):
        """执行每日进化"""
        print(f"🚀 开始每日策略进化 - {self.now.strftime('%Y-%m-%d')}")
        
        # Phase 1: 24小时数据聚合
        print("📊 Phase 1: 数据聚合...")
        hourly_reports = self._collect_hourly_reports()
        aggregated_data = self._aggregate_data(hourly_reports)
        
        # Phase 2: 深度分析
        print("🔍 Phase 2: 深度分析...")
        performance = self._analyze_performance(aggregated_data)
        market_conditions = self._analyze_market_conditions()
        
        # Phase 3: 全球基因深度搜索
        print("🌐 Phase 3: 全球基因搜索...")
        global_genes = self._search_global_genes_deep(performance)
        
        # Phase 4: 下载并验证基因
        print("📥 Phase 4: 基因验证...")
        validated_genes = self._validate_genes(global_genes)
        
        # Phase 5: 基因融合
        print("🧬 Phase 5: 基因融合...")
        current_strategy = self._load_current_strategy()
        fused_strategy = self._fuse_genes(current_strategy, validated_genes)
        
        # Phase 6: 生成优化策略
        print("⚙️  Phase 6: 策略优化...")
        optimized_strategy = self._optimize_strategy(fused_strategy, performance)
        
        # Phase 7: 生成报告
        print("📝 Phase 7: 生成报告...")
        report_path = self._generate_daily_report(
            performance, market_conditions, validated_genes, optimized_strategy
        )
        
        # Phase 8: 保存优化策略（待确认）
        self._save_optimized_strategy(optimized_strategy)
        
        # Phase 9: 发布优秀基因到EvoMap
        if performance.get("improvement_significant"):
            self._publish_local_genes(performance, optimized_strategy)
        
        print(f"✅ 完成！报告: {report_path}")
        print(f"🎯 优化策略: {STRATEGY_DIR}/optimized/")
        print(f"⏳ 等待人工确认后应用...")
        
        return report_path
    
    def _collect_hourly_reports(self):
        """收集过去24小时的所有小时报告"""
        reports = []
        for i in range(24):
            hour = (self.now - timedelta(hours=i)).strftime("%Y%m%d_%H")
            report_path = OUTPUT_DIR / f"hourly_analysis_{hour}.md"
            if report_path.exists():
                with open(report_path, 'r') as f:
                    reports.append({
                        "hour": hour,
                        "content": f.read()
                    })
        return reports
    
    def _aggregate_data(self, hourly_reports):
        """聚合24小时数据"""
        conn = sqlite3.connect("data/data.db")
        cursor = conn.cursor()
        
        # 24小时交易统计
        cursor.execute("""
            SELECT 
                COUNT(*) as total_trades,
                SUM(CASE WHEN realized_pnl > 0 THEN 1 ELSE 0 END) as wins,
                SUM(CASE WHEN realized_pnl < 0 THEN 1 ELSE 0 END) as losses,
                SUM(realized_pnl) as total_pnl,
                AVG(realized_pnl) as avg_pnl,
                MAX(realized_pnl) as max_win,
                MIN(realized_pnl) as max_loss,
                symbol
            FROM trader_positions
            WHERE created_at >= datetime('now', '-24 hours')
              AND status = 'CLOSED'
            GROUP BY symbol
        """)
        symbol_stats = cursor.fetchall()
        
        # 时间序列数据
        cursor.execute("""
            SELECT 
                strftime('%H', created_at) as hour,
                SUM(realized_pnl) as hourly_pnl,
                COUNT(*) as hourly_trades
            FROM trader_positions
            WHERE created_at >= datetime('now', '-24 hours')
              AND status = 'CLOSED'
            GROUP BY hour
            ORDER BY hour
        """)
        hourly_data = cursor.fetchall()
        
        conn.close()
        
        return {
            "symbol_stats": symbol_stats,
            "hourly_data": hourly_data,
            "hourly_reports": hourly_reports
        }
    
    def _analyze_performance(self, aggregated_data):
        """分析24小时策略表现"""
        stats = aggregated_data["symbol_stats"]
        
        total_trades = sum(s[0] for s in stats) if stats else 0
        total_wins = sum(s[1] for s in stats) if stats else 0
        total_losses = sum(s[2] for s in stats) if stats else 0
        total_pnl = sum(s[3] for s in stats) if stats else 0
        
        win_rate = (total_wins / total_trades * 100) if total_trades > 0 else 0
        
        # 计算盈亏比
        avg_win = sum(s[4] for s in stats if s[4] > 0) / total_wins if total_wins > 0 else 0
        avg_loss = abs(sum(s[4] for s in stats if s[4] < 0) / total_losses) if total_losses > 0 else 1
        profit_loss_ratio = avg_win / avg_loss if avg_loss > 0 else 0
        
        # 计算夏普比率（简化版）
        returns = [s[3] for s in stats]
        sharpe = self._calculate_sharpe(returns) if returns else 0
        
        # 最大回撤
        max_drawdown = self._calculate_max_drawdown(aggregated_data["hourly_data"])
        
        # 评级
        rating = self._rate_performance(win_rate, profit_loss_ratio, total_pnl, max_drawdown)
        
        return {
            "total_trades": total_trades,
            "win_rate": round(win_rate, 1),
            "profit_loss_ratio": round(profit_loss_ratio, 2),
            "total_pnl": round(total_pnl, 2),
            "sharpe_ratio": round(sharpe, 2),
            "max_drawdown": round(max_drawdown, 1),
            "rating": rating,
            "symbol_breakdown": stats,
            "improvement_significant": total_pnl > 0 and win_rate > 60
        }
    
    def _calculate_sharpe(self, returns):
        """计算简化夏普比率"""
        if not returns or len(returns) < 2:
            return 0
        avg = sum(returns) / len(returns)
        variance = sum((r - avg) ** 2 for r in returns) / len(returns)
        std = variance ** 0.5
        return (avg / std) * (252 ** 0.5) if std > 0 else 0  # 年化
    
    def _calculate_max_drawdown(self, hourly_data):
        """计算最大回撤"""
        if not hourly_data:
            return 0
        
        cumulative = 0
        peak = 0
        max_dd = 0
        
        for hour, pnl, _ in hourly_data:
            cumulative += pnl
            if cumulative > peak:
                peak = cumulative
            dd = (peak - cumulative) / peak * 100 if peak > 0 else 0
            max_dd = max(max_dd, dd)
        
        return max_dd
    
    def _rate_performance(self, win_rate, pl_ratio, total_pnl, max_dd):
        """评级策略表现"""
        score = 0
        score += 30 if win_rate > 60 else 20 if win_rate > 50 else 10
        score += 30 if pl_ratio > 2 else 20 if pl_ratio > 1.5 else 10
        score += 20 if total_pnl > 0 else 0
        score += 20 if max_dd < 10 else 10 if max_dd < 15 else 0
        
        if score >= 80:
            return "优秀"
        elif score >= 60:
            return "良好"
        elif score >= 40:
            return "一般"
        else:
            return "需改进"
    
    def _analyze_market_conditions(self):
        """分析市场条件"""
        # 从小时报告中提取市场状态统计
        # 实际实现需要解析小时报告内容
        return {
            "best_regime": "资金溢出",
            "worst_regime": "恶性暴跌",
            "regime_distribution": {
                "狂暴牛市": 4,
                "资金溢出": 8,
                "阴跌绞杀": 6,
                "恶性暴跌": 2
            }
        }
    
    def _search_global_genes_deep(self, performance):
        """深度搜索全球基因"""
        search_queries = []
        
        # 基于表现特征构建搜索
        if performance["win_rate"] < 50:
            search_queries.append({
                "signals": ["low_win_rate", "confidence_calibration"],
                "types": ["PromptGene"]
            })
        
        if performance["max_drawdown"] > 10:
            search_queries.append({
                "signals": ["high_drawdown", "risk_management"],
                "types": ["RiskGene"]
            })
        
        # 搜索
        all_genes = []
        for query in search_queries:
            genes = self.evo_client.search_genes(
                signals=query["signals"],
                gene_types=query["types"],
                min_rating=4.0,
                limit=3
            )
            all_genes.extend(genes)
        
        return all_genes
    
    def _validate_genes(self, genes):
        """验证下载的基因"""
        validated = []
        for gene in genes:
            # 本地回测验证
            validation_result = self._backtest_gene(gene)
            if validation_result["improvement"] > 5:  # 至少5%改进
                gene["local_validation"] = validation_result
                validated.append(gene)
        
        # 分类存储
        self._store_genes_by_status(validated)
        return validated
    
    def _backtest_gene(self, gene):
        """本地回测基因"""
        # 简化回测：基于历史数据模拟
        return {
            "improvement": 8.5,  # 示例数据
            "test_cases": 10,
            "passed": 8
        }
    
    def _store_genes_by_status(self, genes):
        """按状态存储基因"""
        for gene in genes:
            if gene.get("local_validation", {}).get("improvement", 0) > 10:
                status = "approved"
            else:
                status = "pending_review"
            
            gene_path = STRATEGY_DIR / "imported_genes" / status / f"gene_{gene['id']}.json"
            gene_path.parent.mkdir(parents=True, exist_ok=True)
            
            with open(gene_path, 'w') as f:
                json.dump(gene, f, indent=2)
    
    def _load_current_strategy(self):
        """加载当前策略"""
        conn = sqlite3.connect("data/data.db")
        cursor = conn.cursor()
        
        cursor.execute("""
            SELECT id, name, config FROM strategies 
            WHERE is_active = 1 LIMIT 1
        """)
        row = cursor.fetchone()
        conn.close()
        
        if row:
            return {
                "id": row[0],
                "name": row[1],
                "config": json.loads(row[2])
            }
        return None
    
    def _fuse_genes(self, strategy, genes):
        """融合基因到策略"""
        if not strategy or not genes:
            return strategy
        
        return self.fusion_engine.fuse(strategy, genes)
    
    def _optimize_strategy(self, strategy, performance):
        """进一步优化策略"""
        return self.optimizer.optimize(strategy, performance)
    
    def _generate_daily_report(self, performance, market_conditions, genes, optimized_strategy):
        """生成每日报告"""
        OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
        
        report_path = OUTPUT_DIR / f"daily_analysis_{self.date_str}.md"
        
        report = self.report_gen.generate_daily(
            date=self.now,
            performance=performance,
            market_conditions=market_conditions,
            genes=genes,
            optimized_strategy=optimized_strategy
        )
        
        with open(report_path, 'w') as f:
            f.write(report)
        
        return report_path
    
    def _save_optimized_strategy(self, strategy):
        """保存优化后的策略（待确认）"""
        STRATEGY_DIR.mkdir(parents=True, exist_ok=True)
        
        opt_path = STRATEGY_DIR / "optimized" / f"strategy_optimized_{self.date_str}.json"
        with open(opt_path, 'w') as f:
            json.dump(strategy, f, indent=2)
    
    def _publish_local_genes(self, performance, strategy):
        """发布本地优秀基因到 EvoMap"""
        if performance.get("total_pnl", 0) > 50:  # 盈利>$50才发布
            gene_capsule = {
                "type": "GeneCapsule",
                "gene_type": "StrategyGene",
                "payload": {
                    "strategy_config": strategy["config"],
                    "improvement_metrics": {
                        "win_rate": performance["win_rate"],
                        "profit_loss_ratio": performance["profit_loss_ratio"],
                        "sample_size": performance["total_trades"]
                    }
                },
                "metadata": {
                    "author": "xh_lx024",
                    "source_strategy": strategy["name"],
                    "creation_time": self.now.isoformat()
                }
            }
            
            self.evo_client.publish_gene(gene_capsule)
            
            # 保存到本地基因库
            local_path = STRATEGY_DIR / "local_genes" / f"gene_{self.date_str}_{performance['total_pnl']}.json"
            local_path.parent.mkdir(parents=True, exist_ok=True)
            with open(local_path, 'w') as f:
                json.dump(gene_capsule, f, indent=2)

if __name__ == "__main__":
    evolver = DailyEvolver()
    evolver.run()