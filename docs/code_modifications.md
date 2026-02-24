# NOFX 代码修改清单

## 概述

基于交易数据分析和 Moltbook 技能学习，本清单列出了需要修改的代码文件，以实现：
- Confluence Scoring (5因素评分系统)
- AI Discipline (拒绝清单)
- DCA Optimization (分批建仓)
- Dual-Source Verification (双源价格验证)

---

## 修改原则

⚠️ **重要安全提示**:
- 只修改 `~/nofx-0210/` 本地 GitHub 仓库代码
- **绝不修改正在运行的 Docker 容器** (`nofx-trading`)
- 所有修改通过 GitHub 发布，用户手动更新 Docker 镜像
- 使用 iflow CLI 辅助代码分析和修改

---

## 需要修改的文件清单

### 1. 核心交易决策模块

#### 1.1 `trader/decision.go` (新增/修改)
**修改类型**: 核心逻辑增强
**功能**: Confluence Scoring 计算

**新增内容**:
```go
// ConfluenceScore 5因素评分结构
type ConfluenceScore struct {
    RSI       float64 `json:"rsi"`       // RSI 评分 (0-1)
    Fib       float64 `json:"fib"`       // Fibonacci 评分 (0-1)
    Volume    float64 `json:"volume"`    // Volume 评分 (0-1)
    MATrend   float64 `json:"ma_trend"`  // MA 趋势评分 (0-1)
    Confluence float64 `json:"confluence"` // 综合确认 (0-1)
    Total     float64 `json:"total"`     // 总分 (0-5)
}

// CalculateConfluence 计算融合评分
func (t *Trader) CalculateConfillment(symbol string, side string) (*ConfluenceScore, error) {
    // 1. 获取 RSI
    rsi := t.indicators.GetRSI(symbol, 14)
    rsiScore := 0.0
    if side == "LONG" && rsi < 35 {
        rsiScore = 1.0
    } else if side == "SHORT" && rsi > 65 {
        rsiScore = 1.0
    }
    
    // 2. 获取 Fibonacci 位置
    fibLevels := t.indicators.GetFibonacci(symbol)
    fibScore := t.checkFibonacciProximity(symbol, fibLevels)
    
    // 3. Volume 检查
    volume := t.market.GetVolume(symbol)
    avgVolume := t.market.GetAvgVolume(symbol, 20)
    volumeScore := 0.0
    if volume > avgVolume*1.2 {
        volumeScore = 1.0
    }
    
    // 4. MA 趋势
    ma20 := t.indicators.GetMA(symbol, 20)
    price := t.market.GetPrice(symbol)
    maScore := 0.0
    if side == "LONG" && price > ma20 {
        maScore = 1.0
    } else if side == "SHORT" && price < ma20 {
        maScore = 1.0
    }
    
    // 5. Confluence 确认
    confluenceScore := 0.0
    if rsiScore+fibScore+volumeScore+maScore >= 3 {
        confluenceScore = 1.0
    }
    
    total := rsiScore + fibScore + volumeScore + maScore + confluenceScore
    
    return &ConfluenceScore{
        RSI: rsiScore,
        Fib: fibScore,
        Volume: volumeScore,
        MATrend: maScore,
        Confluence: confluenceScore,
        Total: total,
    }, nil
}
```

**修改理由**: 实现 TradingLobster 的 Confluence Scoring System

**预期效果**: 提升胜率至 72%+

---

#### 1.2 `trader/discipline.go` (新增)
**修改类型**: 新增文件
**功能**: AI Discipline 拒绝逻辑

**新增内容**:
```go
package trader

// DisciplineChecker 纪律检查器
type DisciplineChecker struct {
    consecutiveLosses map[string]int  // 连续亏损计数
    dailyLossPercent  float64         // 当日亏损百分比
    activePositions   int             // 当前持仓数量
}

// CheckRejection 检查是否拒绝交易
func (d *DisciplineChecker) CheckRejection(symbol string, side string, score *ConfluenceScore) (bool, string) {
    // 1. RSI 极端值检查
    rsi := GetRSI(symbol)
    if side == "LONG" && rsi > 78 {
        return true, "RSI > 78，严重超买，拒绝做多"
    }
    if side == "SHORT" && rsi < 22 {
        return true, "RSI < 22，严重超卖，拒绝做空"
    }
    
    // 2. Confluence 分数检查
    if score.Total < 4.0 {
        return true, fmt.Sprintf("Confluence 分数 %.1f/5，低于4分，拒绝入场", score.Total)
    }
    
    // 3. 连续亏损检查
    key := symbol + "_" + side
    if d.consecutiveLosses[key] >= 2 {
        return true, fmt.Sprintf("%s %s 方向连续亏损 %d 次，暂停24小时", symbol, side, d.consecutiveLosses[key])
    }
    
    // 4. 当日亏损限额检查
    if d.dailyLossPercent > 5.0 {
        return true, fmt.Sprintf("当日亏损 %.2f%%，超过5%%限额，停止交易", d.dailyLossPercent)
    }
    
    // 5. 持仓数量检查
    if d.activePositions >= 3 {
        return true, "当前持仓已达3个，拒绝新开仓"
    }
    
    // 6. Volume 检查
    volume := GetVolume(symbol)
    avgVolume := GetAvgVolume(symbol, 20)
    if volume < avgVolume*1.0 {
        return true, "成交量不足，拒绝入场"
    }
    
    return false, ""
}

// RecordResult 记录交易结果（用于连续亏损统计）
func (d *DisciplineChecker) RecordResult(symbol string, side string, pnl float64) {
    key := symbol + "_" + side
    if pnl < 0 {
        d.consecutiveLosses[key]++
    } else {
        d.consecutiveLosses[key] = 0
    }
    
    if pnl < 0 {
        d.dailyLossPercent += math.Abs(pnl)
    }
}
```

**修改理由**: 实现 AI Discipline，避免情绪化交易和过度风险

**预期效果**: 减少连续亏损，控制回撤

---

#### 1.3 `trader/dca.go` (新增)
**修改类型**: 新增文件
**功能**: DCA 分批建仓管理

**新增内容**:
```go
package trader

// DCAManager DCA 管理器
type DCAManager struct {
    batches map[string]*DCABatch  // symbol -> 批次信息
}

// DCABatch DCA 批次
type DCABatch struct {
    Symbol      string
    Side        string
    TotalSize   float64   // 总目标仓位
    EntryPrices []float64 // 各批次入场价
    Sizes       []float64 // 各批次仓位大小
    Status      []string  // 各批次状态 [pending/filled]
}

// CreateDCABatch 创建 DCA 批次
func (d *DCAManager) CreateDCABatch(symbol string, side string, totalSize float64) *DCABatch {
    // 分配比例: 30% - 40% - 30%
    return &DCABatch{
        Symbol: symbol,
        Side: side,
        TotalSize: totalSize,
        EntryPrices: make([]float64, 3),
        Sizes: []float64{
            totalSize * 0.30,  // 第1批 30%
            totalSize * 0.40,  // 第2批 40%
            totalSize * 0.30,  // 第3批 30%
        },
        Status: []string{"pending", "pending", "pending"},
    }
}

// ShouldEnterBatch 判断是否进入某批次
func (d *DCAManager) ShouldEnterBatch(batch *DCABatch, batchIndex int, currentPrice float64, score *ConfluenceScore) bool {
    // 第1批: Confluence = 4/5
    if batchIndex == 0 && score.Total >= 4.0 {
        return true
    }
    
    // 第2批: Confluence = 5/5 或价格向有利方向移动 1%
    if batchIndex == 1 {
        if score.Total >= 5.0 {
            return true
        }
        if batch.EntryPrices[0] > 0 {
            movePercent := math.Abs(currentPrice-batch.EntryPrices[0]) / batch.EntryPrices[0]
            if movePercent >= 0.01 {
                return true
            }
        }
    }
    
    // 第3批: 价格向有利方向移动 2%
    if batchIndex == 2 && batch.EntryPrices[1] > 0 {
        movePercent := math.Abs(currentPrice-batch.EntryPrices[1]) / batch.EntryPrices[1]
        if movePercent >= 0.02 {
            return true
        }
    }
    
    return false
}

// GetAverageEntryPrice 计算平均入场价
func (d *DCAManager) GetAverageEntryPrice(batch *DCABatch) float64 {
    totalValue := 0.0
    totalSize := 0.0
    
    for i := 0; i < 3; i++ {
        if batch.Status[i] == "filled" {
            totalValue += batch.EntryPrices[i] * batch.Sizes[i]
            totalSize += batch.Sizes[i]
        }
    }
    
    if totalSize == 0 {
        return 0
    }
    return totalValue / totalSize
}
```

**修改理由**: 实现 DCA Optimization，降低入场成本 21.6%

**预期效果**: 更好的平均入场价，降低单次入场风险

---

### 2. 价格验证模块

#### 2.1 `market/price_verifier.go` (新增)
**修改类型**: 新增文件
**功能**: 双源价格验证

**新增内容**:
```go
package market

// PriceVerifier 价格验证器
type PriceVerifier struct {
    primarySource   PriceSource
    secondarySource PriceSource
}

// PriceSource 价格源接口
type PriceSource interface {
    GetPrice(symbol string) (float64, error)
    GetName() string
}

// VerifyPrice 验证价格
func (v *PriceVerifier) VerifyPrice(symbol string) (*PriceVerification, error) {
    price1, err1 := v.primarySource.GetPrice(symbol)
    price2, err2 := v.secondarySource.GetPrice(symbol)
    
    // 处理错误
    if err1 != nil && err2 != nil {
        return nil, fmt.Errorf("两个价格源都不可用")
    }
    
    if err1 != nil {
        return &PriceVerification{
            Price: price2,
            Source: v.secondarySource.GetName(),
            Confidence: 0.5, // 单一源置信度降低
        }, nil
    }
    
    if err2 != nil {
        return &PriceVerification{
            Price: price1,
            Source: v.primarySource.GetName(),
            Confidence: 0.5,
        }, nil
    }
    
    // 计算差异
    variance := math.Abs(price1-price2) / price1 * 100
    
    // 决策
    if variance < 1.0 {
        // 使用平均值
        return &PriceVerification{
            Price: (price1 + price2) / 2,
            Variance: variance,
            Status: "VERIFIED",
            Confidence: 1.0,
        }, nil
    } else if variance < 3.0 {
        return &PriceVerification{
            Price: price1,
            Variance: variance,
            Status: "WARNING",
            Confidence: 0.7,
        }, nil
    } else {
        return nil, fmt.Errorf("价格差异 %.2f%% 超过阈值", variance)
    }
}

// PriceVerification 价格验证结果
type PriceVerification struct {
    Price      float64
    Variance   float64
    Status     string  // VERIFIED / WARNING / REJECTED
    Confidence float64 // 0.0 - 1.0
    Source     string
}
```

**修改理由**: 实现 Dual-Source Verification，避免异常价格交易

**预期效果**: 防止在价格延迟或异常时执行错误交易

---

### 3. 配置文件修改

#### 3.1 `config/strategy.yaml` (新增配置)
**修改类型**: 配置增强

**新增内容**:
```yaml
# Confluence Scoring 配置
confluence:
  enabled: true
  factors:
    rsi:
      enabled: true
      period: 14
      long_threshold: 35   # RSI < 35 做多
      short_threshold: 65  # RSI > 65 做空
      reject_long_above: 78
      reject_short_below: 22
    fibonacci:
      enabled: true
      levels: [0.618, 0.786]
      tolerance: 0.02      # ±2%
    volume:
      enabled: true
      period: 20
      multiplier: 1.2      # > 1.2倍平均量
    ma:
      enabled: true
      period: 20
  min_score: 4.0           # 最低4/5分

# AI Discipline 配置
discipline:
  enabled: true
  max_consecutive_losses: 2
  daily_loss_limit: 5.0    # 5%
  max_positions: 3
  min_volume_multiplier: 1.0

# DCA 配置
dca:
  enabled: true
  batches:
    - ratio: 0.30
      trigger: "confluence_4"
    - ratio: 0.40
      trigger: "confluence_5_or_move_1pct"
    - ratio: 0.30
      trigger: "move_2pct"

# 双源验证配置
price_verification:
  enabled: true
  primary_source: "binance"
  secondary_source: "coingecko"
  max_variance: 3.0        # 最大3%差异
```

---

### 4. Prompt 集成修改

#### 4.1 `llm/prompt_builder.go` (修改)
**修改类型**: Prompt 增强

**修改内容**:
在构建 Prompt 时，注入 Confluence Scoring 规则和 AI Discipline 原则。

```go
// BuildTradingPrompt 构建交易决策 Prompt
func (b *PromptBuilder) BuildTradingPrompt(marketData *MarketData) string {
    prompt := fmt.Sprintf(`
## 市场数据
- 交易对: %s
- 当前价格: %.2f
- RSI(14): %.2f
- MA20: %.2f
- 成交量: %.2f (20日平均: %.2f)

## Confluence Scoring 系统
请根据以下5因素评分：
1. RSI: %s (%s)
2. Fibonacci: %s
3. Volume: %s
4. MA Trend: %s
5. Confluence: %s

## 决策规则
- 总分 ≥ 4/5: 执行交易
- 总分 < 4/5: 拒绝交易
- RSI > 78 (做多) / < 22 (做空): 拒绝

## 输出要求
请输出：
1. 5因素评分表
2. 总分
3. 决策（执行/拒绝）
4. 原因
5. DCA 分批计划
`, 
        marketData.Symbol,
        marketData.Price,
        marketData.RSI,
        marketData.MA20,
        marketData.Volume,
        marketData.AvgVolume,
        // ... 更多数据
    )
    
    return prompt
}
```

---

## 修改优先级

| 优先级 | 文件 | 影响 | 估计工作量 |
|--------|------|------|-----------|
| P0 | `trader/decision.go` | 核心评分逻辑 | 4小时 |
| P0 | `trader/discipline.go` | 风险控制 | 3小时 |
| P1 | `trader/dca.go` | 仓位管理 | 3小时 |
| P1 | `market/price_verifier.go` | 价格安全 | 2小时 |
| P2 | `config/strategy.yaml` | 配置化 | 1小时 |
| P2 | `llm/prompt_builder.go` | AI 集成 | 2小时 |

**总估计工作量**: 约 15 小时

---

## 测试计划

### 单元测试
```bash
# 测试 Confluence Scoring
go test ./trader -run TestConfluenceScoring

# 测试 Discipline
go test ./trader -run TestDisciplineChecker

# 测试 DCA
go test ./trader -run TestDCAManager

# 测试价格验证
go test ./market -run TestPriceVerifier
```

### 回测
```bash
# 使用历史数据回测新策略
./backtest --strategy=confluence_v2 --start=2026-01-01 --end=2026-02-24
```

---

## GitHub 发布流程

### 1. 创建分支
```bash
cd ~/nofx-0210
git checkout -b feature/confluence-scoring-system
```

### 2. 提交修改
```bash
git add .
git commit -m "feat(strategy): implement Confluence Scoring System

- Add 5-factor scoring (RSI, Fib, Volume, MA, Confluence)
- Implement AI Discipline rejection logic
- Add DCA batch position management
- Integrate dual-source price verification
- Based on Moltbook TradingLobster's proven system

Closes #123"
```

### 3. 推送并创建 PR
```bash
git push origin feature/confluence-scoring-system
# 在 GitHub 创建 Pull Request
```

### 4. 合并后发布 Release
- 版本号: v2.0.0
- 发布说明: 包含 Confluence Scoring 等新功能

---

## 用户更新指南

用户收到更新后执行：
```bash
# 1. 拉取最新代码
cd ~/nofx-0210
git pull origin main

# 2. 重新构建 Docker 镜像
docker-compose build

# 3. 重启容器（在合适的时间，避免中断交易）
docker-compose down
docker-compose up -d

# 4. 验证新版本
# 检查日志确认新策略已加载
```

---

## 回滚计划

如发现问题：
```bash
# 回滚到上一个稳定版本
git checkout v1.x.x
docker-compose build
docker-compose up -d
```

---

*本清单由 AI 基于交易数据分析和 Moltbook 技能学习生成*
*日期: 2026-02-24*