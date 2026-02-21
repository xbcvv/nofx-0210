# 量化系统选币过滤机制 (Coin Selection Pipeline)
> **文档属性:** 核心架构设计与防护指南 / 适用于 NoFx 交易引擎  
> **更新时间:** 本地缓存化与解耦重构后 (RightSide_Flow_Hunter_V4+)

---

## 1. 核心目标 (Why we need this)
原有的 AI500 和 OI_Top 等外部信号源虽然能精准捕捉全网“热度”，但无法完全阻挡**劣币入侵**。为了保护 AI 分析大脑不受扭曲盘面的污染（如上市不足10天的新币导致均线失真、深度不足的微盘币导致的画门插针），我们引入了物理切割的**双时间轴防封禁体系**（The Dual-Loop Architecture）。

本机制的宗旨是：**宁缺毋滥。绝不让一颗脏币流入 AI 的推演引擎。**

---

## 2. 架构原理 (Architecture Design)

系统内部设立了一个独立微服务单例：`CoinFilterManager`（位于 `filter/coin_filter_manager.go`）。
此设计将过滤与引擎微观轮询（3~15 分钟）完全隔离开来，从根本上杜绝了对 Binance API 的高频 DDoS 风险，保证了主交易集群的绝对高可用。

### ⏳ 慢循环层（后台洗刷层 - 默认 30 分钟）
这是一个在后台常驻的定时守护进程 (`daemonLoop`)：
1. **聚合拉取**: 获取币安的 `/fapi/v1/ticker/24hr` （全球统一接口，权重极低）。
2. **本地入库**: 通过上述 1 次接口，就在内存中刷平了几百个合约的交易量、市值计算等。
3. **初指补充**: 每当市面上出现“陌生币”时，请求其首根周线/日线来记录出世时间，并将该“出生纸”永久落盘到 `data/listing_cache.json` 中，终身仅查 1 次。
4. **OI并流**: 对那些“交易量达标”的初步幸存者，单独请求实时 OI。

### ⚡ 快循环层（前台交易层 - 依赖前端配置）
1. AI 交易引擎被触发（如 每 3 分钟）。
2. 调用 `nofxosClient.GetTopRatedCoins(200)` 从外部抓取 200 个热度排名。
3. 把这 200 个名字送进 `CoinFilterManager.GetCleanCoins(raw, limit)` 的黑盒中。
4. **黑盒裁决**: 依据后台刚算好的那张“生死簿”，淘汰掉那些脏币。
5. **截断交付**: 截取存活下来的、排名最靠前的 `limit` (如 10) 个发送给 AI 处理室。

---

## 3. 过滤网细则 (The Hard Filters)

这些门槛极尽严苛，任何一条未达标，都会在入围 AI 分析前被秒杀淘汰（通过常量 `FilterConfig` 设定）。

1. **历史纵深防失真 (`MinListingDays >= 90`)** 
   - 过滤上市未满 3 个月的次新资产，防止其早期无序波动和 EMA(50) 物理层缺失导致的大模型幻觉。
2. **深度红线 (`MinQuoteVolume24h >= 50,000,000 USDT`)**
   - 过去 24 小时内的真实成交额低于 5000 万美金的，直接判为死水，防止强庄插针爆破。
3. **持仓体量红线 (`MinOpenInterest >= 10,000,000 USDT`)**
   - 多空双方的博弈盘子如果不到一千万美金，没有大资金进场，就不具备“量化吃鱼头”的肥膘基础。

*(即将或可拓展接入：)*
4. **极端资金费率与逼空拦截**：通过内存中的 Ticker 处理异常数值。
5. **流通盘安全预估**：防止市值极小的妖币混淆视听。

---

## 4. 开发者维护与注意事项 (Operation & Maintenance)

### 🚨 注意事项（Caution!）
- **不要在交易主循环 `GetFullDecision` 中直接写死过滤 API**：绝不允许为了获取某个币的某项过滤指标，在快循环框架内并发发起百次的 HTTP GET，必会招致 `429 Too Many Requests`。
- **配置的修改**：若需放宽/紧缩阈值，请直接修改 `filter/coin_filter_manager.go` 中的 `DefaultFilterConfig` 结构体：
  ```go
  var DefaultFilterConfig = FilterConfig{
  	MinListingDays:    90,           // 天
  	MinQuoteVolume24h: 50_000_000,   // 5000w美金
  	MinOpenInterest:   10_000_000,   // 1000w美金
  }
  ```
- **配置与前端关联度**：无论前端用户将 `ai500_limit` 设为 3 个 还是 100 个，这都是在**完成以上全部物理清洗之后**，进行的存活名额分配。即使配置要 100 个，如果该时段通过过滤的只有 40 个绝对安全资产，系统也**绝对只会将这通过了审查的 40 个提交给 AI**，绝不进行劣币凑数。

### 📁 核心关联文件
- `filter/coin_filter_manager.go`：独立解耦的核心裁决中心。
- `market/api_client.go`：底层封装的 `GetAll24hrTickers` 大全量数据流。
- `data/listing_cache.json`：由系统自动生成的币种“出生纸”永久保存库。
