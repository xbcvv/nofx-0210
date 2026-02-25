#!/usr/bin/env python3
"""
NOFX 每小时交易数据分析脚本
生成分析报告，检测异常，检索全球基因
"""

import os
import sys
import json
import sqlite3
import glob
from datetime import datetime, timedelta
from pathlib import Path

# 添加项目根目录到路径
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from lib.data_collector import DataCollector
from lib.analyzer import DecisionAnalyzer, ProfitAnalyzer, AnomalyDetector
from lib.evo_client import EvoMapClient
from lib.report_generator import HourlyReportGenerator

# 配置路径
CONFIG_PATH = Path(__file__).parent.parent / "config" / "evolver.yaml"
OUTPUT_DIR = Path("docs/openclaw")
STRATEGY_DIR = Path("docs/openclaw/strategies")

class HourlyAnalyzer:
    def __init__(self):
        self.collector = DataCollector()
        self.evo_client = EvoMapClient()
        self.report_gen = HourlyReportGenerator()
        self.now = datetime.now()
        self.hour_str = self.now.strftime("%Y%m%d_%H")
        
    def run(self):
        """执行每小时分析"""
        print(f"🚀 开始每小时分析 - {self.now.strftime('%Y-%m-%d %H:%M')}")
        
        # Phase 1: 收集数据
        print("📊 Phase 1: 收集数据...")
        raw_decisions = self._collect_raw_decisions()
        db_data = self._query_database()
        log_data = self._parse_logs()
        
        # Phase 2: 本地分析
        print("🔍 Phase 2: 本地分析...")
        analysis = {
            "decision_quality": self._analyze_decision_quality(raw_decisions, db_data),
            "profit_loss": self._analyze_profit_loss(db_data),
            "anomalies": self._detect_anomalies(raw_decisions, db_data),
            "symbol_performance": self._analyze_symbols(db_data),
        }
        
        # Phase 3: 全球基因检索 (新增)
        print("🌐 Phase 3: 检索全球基因...")
        gene_recommendations = self._search_global_genes(analysis)
        
        # Phase 4: 融合分析
        print("🧬 Phase 4: 融合分析...")
        fusion_suggestions = self._generate_fusion_suggestions(
            analysis, gene_recommendations
        )
        
        # Phase 5: 生成报告
        print("📝 Phase 5: 生成报告...")
        report_path = self._generate_report(analysis, gene_recommendations, fusion_suggestions)
        
        # Phase 6: 保存当前策略快照
        self._save_strategy_snapshot()
        
        print(f"✅ 完成！报告: {report_path}")
        return report_path
    
    def _collect_raw_decisions(self):
        """收集原始决策数据"""
        decisions = []
        # 读取最近2小时的决策文件
        for i in range(2):
            hour = (self.now - timedelta(hours=i)).strftime("%Y%m%d_%H")
            file_path = Path(f"data/raw/decisions_{hour}.json")
            if file_path.exists():
                with open(file_path, 'r') as f:
                    decisions.extend(json.load(f))
        return decisions
    
    def _query_database(self):
        """查询交易数据库"""
        conn = sqlite3.connect("data/data.db")
        cursor = conn.cursor()
        
        # 最近1小时的订单
        cursor.execute("""
            SELECT symbol, side, status, avg_fill_price, 
                   quantity, created_at
            FROM trader_orders
            WHERE created_at >= datetime('now', '-1 hour')
        """)
        orders = cursor.fetchall()
        
        # 最近1小时的持仓
        cursor.execute("""
            SELECT symbol, side, realized_pnl, unrealized_pnl,
                   entry_price, current_price
            FROM trader_positions
            WHERE created_at >= datetime('now', '-1 hour')
        """)
        positions = cursor.fetchall()
        
        conn.close()
        return {"orders": orders, "positions": positions}
    
    def _parse_logs(self):
        """解析系统日志"""
        log_file = f"data/nofx_{self.now.strftime('%Y-%m-%d')}.log"
        errors = []
        warnings = []
        
        if os.path.exists(log_file):
            with open(log_file, 'r') as f:
                for line in f:
                    if self.now.strftime('%H') in line:
                        if 'ERROR' in line or 'FATA' in line:
                            errors.append(line.strip())
                        elif 'WARN' in line:
                            warnings.append(line.strip())
        
        return {"errors": errors, "warnings": warnings}
    
    def _analyze_decision_quality(self, decisions, db_data):
        """分析决策质量"""
        results = []
        for decision in decisions:
            # 对比信心值与实际结果
            confidence = decision.get('confidence', 0)
            result = decision.get('result', 'unknown')
            expected_pnl = decision.get('expected_pnl', 0)
            actual_pnl = decision.get('actual_pnl', 0)
            
            match = self._check_prediction_match(
                confidence, expected_pnl, actual_pnl
            )
            
            results.append({
                "timestamp": decision.get('timestamp'),
                "symbol": decision.get('symbol'),
                "confidence": confidence,
                "match": match,
                "actual_pnl": actual_pnl
            })
        
        # 统计
        high_conf = [r for r in results if r['confidence'] > 80]
        high_conf_accuracy = len([r for r in high_conf if r['match']]) / len(high_conf) if high_conf else 0
        
        return {
            "total_decisions": len(results),
            "high_confidence_decisions": len(high_conf),
            "high_confidence_accuracy": high_conf_accuracy,
            "details": results
        }
    
    def _check_prediction_match(self, confidence, expected_pnl, actual_pnl):
        """检查预测与实际是否匹配"""
        if expected_pnl > 0 and actual_pnl > 0:
            return True
        if expected_pnl < 0 and actual_pnl < 0:
            return True
        return False
    
    def _analyze_profit_loss(self, db_data):
        """分析盈亏"""
        positions = db_data.get('positions', [])
        
        wins = [p for p in positions if p[2] > 0]  # realized_pnl > 0
        losses = [p for p in positions if p[2] < 0]
        
        total_pnl = sum(p[2] for p in positions)
        win_count = len(wins)
        loss_count = len(losses)
        
        return {
            "total_pnl": round(total_pnl, 2),
            "win_count": win_count,
            "loss_count": loss_count,
            "win_rate": round(win_count / (win_count + loss_count) * 100, 1) if (win_count + loss_count) > 0 else 0,
            "avg_win": round(sum(p[2] for p in wins) / len(wins), 2) if wins else 0,
            "avg_loss": round(sum(p[2] for p in losses) / len(losses), 2) if losses else 0,
        }
    
    def _detect_anomalies(self, decisions, db_data):
        """检测异常"""
        anomalies = []
        
        # 检查连续亏损
        positions = db_data.get('positions', [])
        symbol_losses = {}
        for p in positions:
            symbol = p[0]
            pnl = p[2]
            if pnl < 0:
                symbol_losses[symbol] = symbol_losses.get(symbol, 0) + 1
        
        for symbol, count in symbol_losses.items():
            if count >= 3:
                anomalies.append({
                    "type": "consecutive_losses",
                    "symbol": symbol,
                    "count": count,
                    "severity": "high"
                })
        
        # 检查高信心亏损
        for decision in decisions:
            if decision.get('confidence', 0) > 80 and decision.get('actual_pnl', 0) < -10:
                anomalies.append({
                    "type": "high_confidence_loss",
                    "symbol": decision.get('symbol'),
                    "confidence": decision.get('confidence'),
                    "loss": decision.get('actual_pnl'),
                    "severity": "medium"
                })
        
        return anomalies
    
    def _analyze_symbols(self, db_data):
        """分析各币种表现"""
        positions = db_data.get('positions', [])
        symbol_stats = {}
        
        for p in positions:
            symbol = p[0]
            pnl = p[2]
            
            if symbol not in symbol_stats:
                symbol_stats[symbol] = {"trades": 0, "wins": 0, "total_pnl": 0}
            
            symbol_stats[symbol]["trades"] += 1
            symbol_stats[symbol]["total_pnl"] += pnl
            if pnl > 0:
                symbol_stats[symbol]["wins"] += 1
        
        # 计算胜率
        for symbol, stats in symbol_stats.items():
            stats["win_rate"] = round(stats["wins"] / stats["trades"] * 100, 1)
            stats["total_pnl"] = round(stats["total_pnl"], 2)
        
        return symbol_stats
    
    def _search_global_genes(self, analysis):
        """检索全球基因胶囊"""
        recommendations = []
        
        # 基于异常搜索相关基因
        for anomaly in analysis.get("anomalies", []):
            if anomaly["type"] == "consecutive_losses":
                genes = self.evo_client.search_genes(
                    signals=[f"consecutive_losses:{anomaly['count']}+", 
                            f"symbol:{anomaly['symbol']}"],
                    gene_types=["PromptGene", "RiskGene"]
                )
                recommendations.extend(genes)
        
        # 基于低胜率搜索
        pl = analysis.get("profit_loss", {})
        if pl.get("win_rate", 100) < 50:
            genes = self.evo_client.search_genes(
                signals=["low_win_rate", "confidence_calibration"],
                gene_types=["PromptGene"]
            )
            recommendations.extend(genes)
        
        return recommendations[:5]  # 最多5个推荐
    
    def _generate_fusion_suggestions(self, analysis, gene_recommendations):
        """生成本地优化与基因融合建议"""
        suggestions = []
        
        # 本地优化建议
        if analysis["anomalies"]:
            suggestions.append({
                "type": "local",
                "priority": "high",
                "description": f"修复 {len(analysis['anomalies'])} 个检测到的异常",
                "risk": "low"
            })
        
        # 基因融合建议
        for gene in gene_recommendations:
            suggestions.append({
                "type": "fusion",
                "priority": "medium" if gene.get("rating", 0) > 4.5 else "low",
                "gene_id": gene.get("id"),
                "description": gene.get("description", ""),
                "risk": self._assess_gene_risk(gene),
                "source": gene.get("author", "unknown")
            })
        
        return suggestions
    
    def _assess_gene_risk(self, gene):
        """评估基因风险等级"""
        # 基于基因类型和修改范围评估
        gene_type = gene.get("gene_type", "")
        
        if gene_type == "PromptGene":
            return "high"  # 修改Prompt需要谨慎
        elif gene_type == "RiskGene":
            return "medium"  # 风控参数
        else:
            return "low"  # 指标权重等低风险
    
    def _generate_report(self, analysis, gene_recommendations, fusion_suggestions):
        """生成报告"""
        OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
        
        report_path = OUTPUT_DIR / f"hourly_analysis_{self.hour_str}.md"
        
        report = self.report_gen.generate(
            timestamp=self.now,
            analysis=analysis,
            genes=gene_recommendations,
            suggestions=fusion_suggestions
        )
        
        with open(report_path, 'w') as f:
            f.write(report)
        
        return report_path
    
    def _save_strategy_snapshot(self):
        """保存当前策略快照"""
        STRATEGY_DIR.mkdir(parents=True, exist_ok=True)
        
        # 从数据库导出当前策略
        conn = sqlite3.connect("data/data.db")
        cursor = conn.cursor()
        
        cursor.execute("""
            SELECT id, name, config FROM strategies 
            WHERE is_active = 1 LIMIT 1
        """)
        row = cursor.fetchone()
        conn.close()
        
        if row:
            strategy = {
                "id": row[0],
                "name": row[1],
                "config": json.loads(row[2]),
                "exported_at": self.now.isoformat()
            }
            
            snapshot_path = STRATEGY_DIR / "current" / f"strategy_{self.hour_str}.json"
            snapshot_path.parent.mkdir(parents=True, exist_ok=True)
            
            with open(snapshot_path, 'w') as f:
                json.dump(strategy, f, indent=2)

if __name__ == "__main__":
    analyzer = HourlyAnalyzer()
    analyzer.run()