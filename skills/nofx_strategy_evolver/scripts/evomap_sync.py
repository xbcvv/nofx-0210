#!/usr/bin/env python3
"""
EvoMap 同步工具
发布本地基因、搜索全球基因、同步声誉数据
"""

import os
import sys
import json
import argparse
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))
from lib.evo_client import EvoMapClient

def main():
    parser = argparse.ArgumentParser(description="EvoMap 同步工具")
    parser.add_argument("--publish", action="store_true", help="发布本地基因")
    parser.add_argument("--search", action="store_true", help="搜索全球基因")
    parser.add_argument("--sync-reputation", action="store_true", help="同步声誉数据")
    parser.add_argument("--gene-file", type=str, help="要发布的基因文件路径")
    parser.add_argument("--search-signals", nargs="+", help="搜索信号")
    parser.add_argument("--min-rating", type=float, default=4.0, help="最低评分")
    
    args = parser.parse_args()
    
    client = EvoMapClient()
    
    if args.publish and args.gene_file:
        print(f"📤 发布基因: {args.gene_file}")
        with open(args.gene_file, 'r') as f:
            gene = json.load(f)
        result = client.publish_gene(gene)
        print(f"✅ 已发布，ID: {result}")
    
    elif args.search:
        signals = args.search_signals or ["consecutive_losses", "market_regime"]
        print(f"🔍 搜索基因，信号: {signals}")
        genes = client.search_genes(
            signals=signals,
            min_rating=args.min_rating
        )
        print(f"📊 找到 {len(genes)} 个基因:")
        for gene in genes:
            print(f"  - {gene.get('id')}: {gene.get('rating')}⭐ ({gene.get('author')})")
    
    elif args.sync_reputation:
        print("🔄 同步声誉数据...")
        reputation = client.get_reputation()
        print(f"💎 当前声誉: {reputation}")
    
    else:
        parser.print_help()

if __name__ == "__main__":
    main()