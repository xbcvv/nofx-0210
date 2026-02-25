# iflow CLI 使用指南 - NOFX 项目

## 项目关联
- **GitHub 仓库**: https://github.com/xbcvv/nofx-0210
- **本地路径**: /root/nofx-0210
- **iflow 配置**: /root/nofx-0210/.iflow/

## 快速开始

### 1. 进入项目目录启动 iflow
```bash
cd ~/nofx-0210
iflow
```

### 2. 常用命令

#### 代码分析
```
分析 trader 目录的交易策略实现
解释 provider 目录中的交易所接入逻辑
优化 backtest 目录的回测系统
```

#### GitHub 集成 (需要设置 GITHUB_TOKEN)
```
查看最近的 GitHub Issues
创建一个新的 Pull Request
提交代码变更
```

#### 代码修改
```
/init
# 分析项目结构后开始工作

请帮我重构 trader/position.go 文件
优化策略执行性能
添加新的技术指标支持
```

### 3. 非交互式使用
```bash
# 分析代码并生成文档
iflow -p "分析 trader 目录的代码结构并生成文档"

# 检查代码问题
iflow -p "检查 trader 目录中的潜在 bug 和安全问题"
```

## 环境变量设置 (如需 GitHub 集成)
```bash
export GITHUB_TOKEN="your_github_personal_access_token"
iflow
```

## 项目特定的 iflow 能力

### 交易策略分析
- 分析和优化 trader/ 目录下的策略代码
- 解释 provider/ 中的交易所 API 接入
- 优化 backtest/ 回测逻辑

### 配置管理
- 修改 config/ 中的配置文件
- 更新 docker-compose.yml
- 调整 .env 环境变量

### 数据库操作
- 查询 SQLite 数据库
- 分析交易记录
- 生成报表

## 注意事项
1. 本项目是 Fork 自 NoFxAiOS/nofx 的修改版本
2. 主要修改在 trader/ 和 strategy/ 目录
3. 使用 Go 语言开发，需要 Go 1.21+
4. Docker 部署为主，生产环境使用 docker-compose.prod.yml

## 相关代理
- **coder_agent** (vpscli/gemini-3.1-pro-high) - 可调用 iflow 处理编程任务
- **nofx_trading_agent** (aipro/claude-opus-4-6) - 交易分析和策略优化