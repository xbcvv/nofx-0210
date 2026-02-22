import { useState } from 'react'
import { Plus, X, Database, TrendingUp, TrendingDown, List, Ban, Zap, Shuffle } from 'lucide-react'
import type { CoinSourceConfig } from '../../types'

interface CoinSourceEditorProps {
  config: CoinSourceConfig
  onChange: (config: CoinSourceConfig) => void
  disabled?: boolean
  language: string
}

export function CoinSourceEditor({
  config,
  onChange,
  disabled,
  language,
}: CoinSourceEditorProps) {
  const [newCoin, setNewCoin] = useState('')
  const [newExcludedCoin, setNewExcludedCoin] = useState('')

  const t = (key: string) => {
    const translations: Record<string, Record<string, string>> = {
      sourceType: { zh: '数据来源类型', en: 'Source Type' },
      static: { zh: '静态列表', en: 'Static List' },
      ai500: { zh: 'AI500 数据源', en: 'AI500 Data Provider' },
      oi_top: { zh: 'OI 持仓增加', en: 'OI Increase' },
      oi_low: { zh: 'OI 持仓减少', en: 'OI Decrease' },
      mixed: { zh: '混合模式', en: 'Mixed Mode' },
      staticCoins: { zh: '自定义币种', en: 'Custom Coins' },
      addCoin: { zh: '添加币种', en: 'Add Coin' },
      useAI500: { zh: '启用 AI500 数据源', en: 'Enable AI500 Data Provider' },
      ai500Limit: { zh: '数量上限', en: 'Limit' },
      useOITop: { zh: '启用 OI 持仓增加榜', en: 'Enable OI Increase' },
      oiTopLimit: { zh: '数量上限', en: 'Limit' },
      useOILow: { zh: '启用 OI 持仓减少榜', en: 'Enable OI Decrease' },
      oiLowLimit: { zh: '数量上限', en: 'Limit' },
      staticDesc: { zh: '手动指定交易币种列表', en: 'Manually specify trading coins' },
      ai500Desc: {
        zh: '使用 AI500 智能筛选的热门币种',
        en: 'Use AI500 smart-filtered popular coins',
      },
      oiTopDesc: {
        zh: '持仓增加榜，适合做多',
        en: 'OI increase ranking, for long',
      },
      oi_lowDesc: {
        zh: '持仓减少榜，适合做空',
        en: 'OI decrease ranking, for short',
      },
      mixedDesc: {
        zh: '组合多种数据源',
        en: 'Combine multiple sources',
      },
      mixedConfig: { zh: '组合数据源配置', en: 'Combined Sources Configuration' },
      mixedSummary: { zh: '已选组合', en: 'Selected Sources' },
      maxCoins: { zh: '最多', en: 'Up to' },
      coins: { zh: '个币种', en: 'coins' },
      dataSourceConfig: { zh: '数据源配置', en: 'Data Source Configuration' },
      excludedCoins: { zh: '排除币种', en: 'Excluded Coins' },
      excludedCoinsDesc: { zh: '这些币种将从所有数据源中排除，不会被交易', en: 'These coins will be excluded from all sources and will not be traded' },
      addExcludedCoin: { zh: '添加排除', en: 'Add Excluded' },
      nofxosNote: { zh: '使用 NofxOS API Key（在指标配置中设置）', en: 'Uses NofxOS API Key (set in Indicators config)' },
      useBinanceSource: { zh: '启用币安全网海选', en: 'Enable Binance Top Volume' },
      binanceLimit: { zh: '预抓数量', en: 'Fetch Limit' },
      binanceInterval: { zh: '缓存间隔(分)', en: 'Filter Interval (m)' },
      fetchAllData: { zh: '获取全部名单', en: 'Fetch All Data' },
      blackboxSource: { zh: '核心风控黑盒', en: 'Core Blackbox Filter' },
      blackboxLimit: { zh: '终极名额(交付AI)', en: 'Final Output Limit (To AI)' },
      pipelineDesc: { zh: '外部输入源经过风控管线洗选后，统一由黑盒限制最大输出名额。', en: 'Sources are filtered by risk control pipeline, then truncated by blackbox.' },
    }
    return translations[key]?.[language] || key
  }

  const sourceTypes = [
    { value: 'static', icon: List, color: '#848E9C' },
    { value: 'ai500', icon: Database, color: '#F0B90B' },
    { value: 'oi_top', icon: TrendingUp, color: '#0ECB81' },
    { value: 'oi_low', icon: TrendingDown, color: '#F6465D' },
    { value: 'mixed', icon: Shuffle, color: '#60a5fa' },
  ] as const

  // xyz dex assets (stocks, forex, commodities) - should NOT get USDT suffix
  const xyzDexAssets = new Set([
    // Stocks
    'TSLA', 'NVDA', 'AAPL', 'MSFT', 'META', 'AMZN', 'GOOGL', 'AMD', 'COIN', 'NFLX',
    'PLTR', 'HOOD', 'INTC', 'MSTR', 'TSM', 'ORCL', 'MU', 'RIVN', 'COST', 'LLY',
    'CRCL', 'SKHX', 'SNDK',
    // Forex
    'EUR', 'JPY',
    // Commodities
    'GOLD', 'SILVER',
    // Index
    'XYZ100',
  ])

  const isXyzDexAsset = (symbol: string): boolean => {
    const base = symbol.toUpperCase().replace(/^XYZ:/, '').replace(/USDT$|USD$|-USDC$/, '')
    return xyzDexAssets.has(base)
  }

  const handleAddCoin = () => {
    if (!newCoin.trim()) return
    const symbol = newCoin.toUpperCase().trim()

    // For xyz dex assets (stocks, forex, commodities), use xyz: prefix without USDT
    let formattedSymbol: string
    if (isXyzDexAsset(symbol)) {
      // Remove xyz: prefix (case-insensitive) and any USD suffixes
      const base = symbol.replace(/^xyz:/i, '').replace(/USDT$|USD$|-USDC$/i, '')
      formattedSymbol = `xyz:${base}`
    } else {
      formattedSymbol = symbol.endsWith('USDT') ? symbol : `${symbol}USDT`
    }

    const currentCoins = config.static_coins || []
    if (!currentCoins.includes(formattedSymbol)) {
      onChange({
        ...config,
        static_coins: [...currentCoins, formattedSymbol],
      })
    }
    setNewCoin('')
  }

  const handleRemoveCoin = (coin: string) => {
    onChange({
      ...config,
      static_coins: (config.static_coins || []).filter((c) => c !== coin),
    })
  }

  const handleAddExcludedCoin = () => {
    if (!newExcludedCoin.trim()) return
    const symbol = newExcludedCoin.toUpperCase().trim()

    // For xyz dex assets, use xyz: prefix without USDT
    let formattedSymbol: string
    if (isXyzDexAsset(symbol)) {
      const base = symbol.replace(/^xyz:/i, '').replace(/USDT$|USD$|-USDC$/i, '')
      formattedSymbol = `xyz:${base}`
    } else {
      formattedSymbol = symbol.endsWith('USDT') ? symbol : `${symbol}USDT`
    }

    const currentExcluded = config.excluded_coins || []
    if (!currentExcluded.includes(formattedSymbol)) {
      onChange({
        ...config,
        excluded_coins: [...currentExcluded, formattedSymbol],
      })
    }
    setNewExcludedCoin('')
  }

  const handleRemoveExcludedCoin = (coin: string) => {
    onChange({
      ...config,
      excluded_coins: (config.excluded_coins || []).filter((c) => c !== coin),
    })
  }

  // NofxOS badge component
  const NofxOSBadge = () => (
    <span
      className="text-[9px] px-1.5 py-0.5 rounded font-medium bg-purple-500/20 text-purple-400 border border-purple-500/30"
    >
      NofxOS
    </span>
  )

  return (
    <div className="space-y-6">
      {/* Source Type Selector */}
      <div>
        <label className="block text-sm font-medium mb-3 text-nofx-text">
          {t('sourceType')}
        </label>
        <div className="grid grid-cols-5 gap-2">
          {sourceTypes.map(({ value, icon: Icon, color }) => (
            <button
              key={value}
              onClick={() =>
                !disabled &&
                onChange({ ...config, source_type: value as CoinSourceConfig['source_type'] })
              }
              disabled={disabled}
              className={`p-4 rounded-lg border transition-all ${config.source_type === value
                ? 'ring-2 ring-nofx-gold bg-nofx-gold/10'
                : 'hover:bg-white/5 bg-nofx-bg'
                } border-nofx-gold/20`}
            >
              <Icon className="w-6 h-6 mx-auto mb-2" style={{ color }} />
              <div className="text-sm font-medium text-nofx-text">
                {t(value)}
              </div>
              <div className="text-xs mt-1 text-nofx-text-muted">
                {t(`${value}Desc`)}
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Static Coins - only for static mode */}
      {config.source_type === 'static' && (
        <div>
          <label className="block text-sm font-medium mb-3 text-nofx-text">
            {t('staticCoins')}
          </label>
          <div className="flex flex-wrap gap-2 mb-3">
            {(config.static_coins || []).map((coin) => (
              <span
                key={coin}
                className="flex items-center gap-1 px-3 py-1.5 rounded-full text-sm bg-nofx-bg-lighter text-nofx-text"
              >
                {coin}
                {!disabled && (
                  <button
                    onClick={() => handleRemoveCoin(coin)}
                    className="ml-1 hover:text-red-400 transition-colors"
                  >
                    <X className="w-3 h-3" />
                  </button>
                )}
              </span>
            ))}
          </div>
          {!disabled && (
            <div className="flex gap-2">
              <input
                type="text"
                value={newCoin}
                onChange={(e) => setNewCoin(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleAddCoin()}
                placeholder="BTC, ETH, SOL..."
                className="flex-1 px-4 py-2 rounded-lg bg-nofx-bg border border-nofx-gold/20 text-nofx-text"
              />
              <button
                onClick={handleAddCoin}
                className="px-4 py-2 rounded-lg flex items-center gap-2 transition-colors bg-nofx-gold text-black hover:bg-yellow-500"
              >
                <Plus className="w-4 h-4" />
                {t('addCoin')}
              </button>
            </div>
          )}
        </div>
      )}

      {/* Excluded Coins */}
      <div>
        <div className="flex items-center gap-2 mb-3">
          <Ban className="w-4 h-4 text-nofx-danger" />
          <label className="text-sm font-medium text-nofx-text">
            {t('excludedCoins')}
          </label>
        </div>
        <p className="text-xs mb-3 text-nofx-text-muted">
          {t('excludedCoinsDesc')}
        </p>
        <div className="flex flex-wrap gap-2 mb-3">
          {(config.excluded_coins || []).map((coin) => (
            <span
              key={coin}
              className="flex items-center gap-1 px-3 py-1.5 rounded-full text-sm bg-nofx-danger/15 text-nofx-danger"
            >
              {coin}
              {!disabled && (
                <button
                  onClick={() => handleRemoveExcludedCoin(coin)}
                  className="ml-1 hover:text-white transition-colors"
                >
                  <X className="w-3 h-3" />
                </button>
              )}
            </span>
          ))}
          {(config.excluded_coins || []).length === 0 && (
            <span className="text-xs italic text-nofx-text-muted">
              {language === 'zh' ? '无' : 'None'}
            </span>
          )}
        </div>
        {!disabled && (
          <div className="flex gap-2">
            <input
              type="text"
              value={newExcludedCoin}
              onChange={(e) => setNewExcludedCoin(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleAddExcludedCoin()}
              placeholder="BTC, ETH, DOGE..."
              className="flex-1 px-4 py-2 rounded-lg text-sm bg-nofx-bg border border-nofx-gold/20 text-nofx-text"
            />
            <button
              onClick={handleAddExcludedCoin}
              className="px-4 py-2 rounded-lg flex items-center gap-2 transition-colors text-sm bg-nofx-danger text-white hover:bg-red-600"
            >
              <Ban className="w-4 h-4" />
              {t('addExcludedCoin')}
            </button>
          </div>
        )}
      </div>

      {/* AI500 Options - only for ai500 mode */}
      {config.source_type === 'ai500' && (
        <div
          className="p-4 rounded-lg bg-nofx-gold/5 border border-nofx-gold/20"
        >
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <Zap className="w-4 h-4 text-nofx-gold" />
              <span className="text-sm font-medium text-nofx-text">
                AI500 {t('dataSourceConfig')}
              </span>
              <NofxOSBadge />
            </div>
          </div>

          <div className="space-y-3">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.use_ai500}
                onChange={(e) =>
                  !disabled && onChange({ ...config, use_ai500: e.target.checked })
                }
                disabled={disabled}
                className="w-5 h-5 rounded accent-nofx-gold"
              />
              <span className="text-nofx-text">{t('useAI500')}</span>
            </label>

            {config.use_ai500 && (
              <div className="flex items-center gap-3 pl-8">
                <span className="text-sm text-nofx-text-muted">
                  {t('ai500Limit')}:
                </span>
                <select
                  value={config.ai500_limit || 10}
                  onChange={(e) =>
                    !disabled &&
                    onChange({ ...config, ai500_limit: parseInt(e.target.value) || 10 })
                  }
                  disabled={disabled}
                  className="px-3 py-1.5 rounded bg-nofx-bg border border-nofx-gold/20 text-nofx-text"
                >
                  {[5, 10, 15, 20, 30, 50].map(n => (
                    <option key={n} value={n}>{n}</option>
                  ))}
                </select>
              </div>
            )}

            <p className="text-xs pl-8 text-nofx-text-muted">
              {t('nofxosNote')}
            </p>
          </div>
        </div>
      )}

      {/* OI Top Options - only for oi_top mode */}
      {config.source_type === 'oi_top' && (
        <div
          className="p-4 rounded-lg bg-nofx-success/5 border border-nofx-success/20"
        >
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <TrendingUp className="w-4 h-4 text-nofx-success" />
              <span className="text-sm font-medium text-nofx-text">
                OI {language === 'zh' ? '持仓增加榜' : 'Increase'} {t('dataSourceConfig')}
              </span>
              <NofxOSBadge />
            </div>
          </div>

          <div className="space-y-3">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.use_oi_top}
                onChange={(e) =>
                  !disabled && onChange({ ...config, use_oi_top: e.target.checked })
                }
                disabled={disabled}
                className="w-5 h-5 rounded accent-nofx-success"
              />
              <span className="text-nofx-text">{t('useOITop')}</span>
            </label>

            {config.use_oi_top && (
              <div className="flex items-center gap-3 pl-8">
                <span className="text-sm text-nofx-text-muted">
                  {t('oiTopLimit')}:
                </span>
                <select
                  value={config.oi_top_limit || 10}
                  onChange={(e) =>
                    !disabled &&
                    onChange({ ...config, oi_top_limit: parseInt(e.target.value) || 10 })
                  }
                  disabled={disabled}
                  className="px-3 py-1.5 rounded bg-nofx-bg border border-nofx-gold/20 text-nofx-text"
                >
                  {[5, 10, 15, 20, 30, 50].map(n => (
                    <option key={n} value={n}>{n}</option>
                  ))}
                </select>
              </div>
            )}

            <p className="text-xs pl-8 text-nofx-text-muted">
              {t('nofxosNote')}
            </p>
          </div>
        </div>
      )}

      {/* OI Low Options - only for oi_low mode */}
      {config.source_type === 'oi_low' && (
        <div
          className="p-4 rounded-lg bg-nofx-danger/5 border border-nofx-danger/20"
        >
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <TrendingDown className="w-4 h-4 text-nofx-danger" />
              <span className="text-sm font-medium text-nofx-text">
                OI {language === 'zh' ? '持仓减少榜' : 'Decrease'} {t('dataSourceConfig')}
              </span>
              <NofxOSBadge />
            </div>
          </div>

          <div className="space-y-3">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.use_oi_low}
                onChange={(e) =>
                  !disabled && onChange({ ...config, use_oi_low: e.target.checked })
                }
                disabled={disabled}
                className="w-5 h-5 rounded accent-red-500"
              />
              <span className="text-nofx-text">{t('useOILow')}</span>
            </label>

            {config.use_oi_low && (
              <div className="flex items-center gap-3 pl-8">
                <span className="text-sm text-nofx-text-muted">
                  {t('oiLowLimit')}:
                </span>
                <select
                  value={config.oi_low_limit || 10}
                  onChange={(e) =>
                    !disabled &&
                    onChange({ ...config, oi_low_limit: parseInt(e.target.value) || 10 })
                  }
                  disabled={disabled}
                  className="px-3 py-1.5 rounded bg-nofx-bg border border-nofx-gold/20 text-nofx-text"
                >
                  {[5, 10, 15, 20, 30, 50].map(n => (
                    <option key={n} value={n}>{n}</option>
                  ))}
                </select>
              </div>
            )}

            <p className="text-xs pl-8 text-nofx-text-muted">
              {t('nofxosNote')}
            </p>
          </div>
        </div>
      )}

      {/* Mixed Mode - Unified Data Pipeline Selector */}
      {config.source_type === 'mixed' && (
        <div className="p-4 rounded-lg bg-blue-500/5 border border-blue-500/20">
          <div className="flex items-center gap-2 mb-4">
            <Shuffle className="w-4 h-4 text-blue-400" />
            <span className="text-sm font-medium text-nofx-text">
              {t('mixedConfig')}
            </span>
          </div>
          <p className="text-xs text-nofx-text-muted mb-4">
            {t('pipelineDesc')}
          </p>

          <div className="flex flex-col md:flex-row gap-4">
            {/* Column 1: Sources (AI500 & Binance) */}
            <div className="flex-1 space-y-3 p-3 rounded-lg border border-nofx-border bg-nofx-bg">
              <div className="text-xs font-semibold text-nofx-text-muted mb-2 uppercase tracking-wider">
                1. {language === 'zh' ? '主输入源' : 'Primary Sources'}
              </div>

              {/* AI500 */}
              <div
                className={`p-3 rounded-lg border transition-all cursor-pointer ${config.use_ai500
                    ? 'bg-nofx-gold/10 border-nofx-gold/50'
                    : 'bg-nofx-bg border-nofx-border hover:border-nofx-gold/30'
                  }`}
                onClick={() => !disabled && onChange({ ...config, use_ai500: !config.use_ai500, ai500_fetch_all: !config.use_ai500 ? true : config.ai500_fetch_all })}
              >
                <div className="flex items-center gap-2 mb-2">
                  <input
                    type="checkbox"
                    checked={config.use_ai500}
                    onChange={(e) => !disabled && onChange({ ...config, use_ai500: e.target.checked, ai500_fetch_all: e.target.checked ? true : config.ai500_fetch_all })}
                    disabled={disabled}
                    className="w-4 h-4 rounded accent-nofx-gold"
                    onClick={(e) => e.stopPropagation()}
                  />
                  <Database className="w-4 h-4 text-nofx-gold" />
                  <span className="text-sm font-medium text-nofx-text">AI500</span>
                  <NofxOSBadge />
                </div>
                {config.use_ai500 && (
                  <div className="mt-2 space-y-2 pl-6" onClick={(e) => e.stopPropagation()}>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={config.ai500_fetch_all}
                        onChange={(e) => !disabled && onChange({ ...config, ai500_fetch_all: e.target.checked })}
                        disabled={disabled}
                        className="w-3.5 h-3.5 rounded accent-nofx-gold"
                      />
                      <span className="text-xs text-nofx-text cursor-pointer">{t('fetchAllData')}</span>
                    </label>
                    {!config.ai500_fetch_all && (
                      <div className="flex items-center gap-2">
                        <span className="text-xs text-nofx-text-muted w-14">Limit:</span>
                        <input
                          type="number"
                          value={config.ai500_limit || 10}
                          onChange={(e) => !disabled && onChange({ ...config, ai500_limit: Math.max(1, parseInt(e.target.value) || 10) })}
                          disabled={disabled}
                          min={1}
                          className="w-20 px-2 py-1 rounded text-xs bg-nofx-bg border border-nofx-gold/20 text-nofx-text"
                        />
                      </div>
                    )}
                  </div>
                )}
              </div>

              {/* Binance Top Volume */}
              <div
                className={`p-3 rounded-lg border transition-all cursor-pointer ${config.use_binance_top_vol
                    ? 'bg-nofx-success/10 border-nofx-success/50'
                    : 'bg-nofx-bg border-nofx-border hover:border-nofx-success/30'
                  }`}
                onClick={() => !disabled && onChange({ ...config, use_binance_top_vol: !config.use_binance_top_vol, binance_top_vol_limit: !config.use_binance_top_vol ? 100 : config.binance_top_vol_limit, binance_filter_interval: config.binance_filter_interval || 30 })}
              >
                <div className="flex items-center gap-2 mb-2">
                  <input
                    type="checkbox"
                    checked={config.use_binance_top_vol}
                    onChange={(e) => !disabled && onChange({ ...config, use_binance_top_vol: e.target.checked, binance_top_vol_limit: e.target.checked ? 100 : config.binance_top_vol_limit, binance_filter_interval: config.binance_filter_interval || 30 })}
                    disabled={disabled}
                    className="w-4 h-4 rounded accent-nofx-success"
                    onClick={(e) => e.stopPropagation()}
                  />
                  <TrendingUp className="w-4 h-4 text-nofx-success" />
                  <span className="text-sm font-medium text-nofx-text">
                    {t('useBinanceSource')}
                  </span>
                </div>
                {config.use_binance_top_vol && (
                  <div className="mt-2 space-y-2 pl-6" onClick={(e) => e.stopPropagation()}>
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-nofx-text-muted w-14">{t('binanceLimit')}:</span>
                      <input
                        type="number"
                        value={config.binance_top_vol_limit || 100}
                        onChange={(e) => !disabled && onChange({ ...config, binance_top_vol_limit: Math.max(1, parseInt(e.target.value) || 100) })}
                        disabled={disabled}
                        min={1}
                        className="w-20 px-2 py-1 rounded text-xs bg-nofx-bg border border-nofx-success/20 text-nofx-text"
                      />
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-nofx-text-muted w-14">{t('binanceInterval')}:</span>
                      <select
                        value={config.binance_filter_interval || 30}
                        onChange={(e) => !disabled && onChange({ ...config, binance_filter_interval: parseInt(e.target.value) || 30 })}
                        disabled={disabled}
                        className="w-20 px-2 py-1 rounded text-xs bg-nofx-bg border border-nofx-success/20 text-nofx-text"
                      >
                        <option value={15}>15 m</option>
                        <option value={30}>30 m</option>
                        <option value={60}>60 m</option>
                        <option value={120}>120 m</option>
                      </select>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Middle visual arrow for desktop */}
            <div className="hidden md:flex flex-col items-center justify-center relative">
              <div className="text-xs text-nofx-text-muted mb-1 absolute -top-2 whitespace-nowrap bg-nofx-bg px-1 hidden lg:block">
                {language === 'zh' ? '风控洗选' : 'Filtering'}
              </div>
              <div className="h-px w-8 bg-nofx-border"></div>
              <TrendingUp className="w-4 h-4 text-nofx-text-muted rotate-90 md:rotate-0 absolute text-nofx-border" />
            </div>

            {/* Column 2: Secondary Sources (OI & Custom) */}
            <div className="flex-1 space-y-3 p-3 rounded-lg border border-nofx-border bg-nofx-bg">
              <div className="text-xs font-semibold text-nofx-text-muted mb-2 uppercase tracking-wider">
                2. {language === 'zh' ? '补充源' : 'Secondary'}
              </div>

              {/* OI Top Card */}
              <div className="flex items-center justify-between">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={config.use_oi_top}
                    onChange={(e) => !disabled && onChange({ ...config, use_oi_top: e.target.checked })}
                    disabled={disabled}
                    className="w-3.5 h-3.5 rounded accent-nofx-success"
                  />
                  <span className="text-xs font-medium text-nofx-text">{language === 'zh' ? 'OI 增加' : 'OI Increase'}</span>
                </label>
                {config.use_oi_top && (
                  <div className="flex items-center gap-1">
                    <span className="text-[10px] text-nofx-text-muted">Lmt:</span>
                    <input
                      type="number"
                      value={config.oi_top_limit || 10}
                      onChange={(e) => !disabled && onChange({ ...config, oi_top_limit: Math.max(1, parseInt(e.target.value) || 10) })}
                      disabled={disabled}
                      min={1}
                      className="w-12 px-1 py-0.5 rounded text-[10px] bg-nofx-bg border border-nofx-border text-nofx-text"
                    />
                  </div>
                )}
              </div>

              {/* OI Low Card */}
              <div className="flex items-center justify-between">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={config.use_oi_low}
                    onChange={(e) => !disabled && onChange({ ...config, use_oi_low: e.target.checked })}
                    disabled={disabled}
                    className="w-3.5 h-3.5 rounded accent-red-500"
                  />
                  <span className="text-xs font-medium text-nofx-text">{language === 'zh' ? 'OI 减少' : 'OI Decrease'}</span>
                </label>
                {config.use_oi_low && (
                  <div className="flex items-center gap-1">
                    <span className="text-[10px] text-nofx-text-muted">Lmt:</span>
                    <input
                      type="number"
                      value={config.oi_low_limit || 10}
                      onChange={(e) => !disabled && onChange({ ...config, oi_low_limit: Math.max(1, parseInt(e.target.value) || 10) })}
                      disabled={disabled}
                      min={1}
                      className="w-12 px-1 py-0.5 rounded text-[10px] bg-nofx-bg border border-nofx-border text-nofx-text"
                    />
                  </div>
                )}
              </div>

              {/* Static/Custom Card */}
              <div className="mt-4 pt-3 border-t border-nofx-border">
                <div className="flex items-center gap-2 mb-2">
                  <List className="w-3.5 h-3.5 text-gray-400" />
                  <span className="text-xs font-medium text-nofx-text">
                    {language === 'zh' ? '自定义自选' : 'Custom'}
                  </span>
                  {(config.static_coins || []).length > 0 && (
                    <span className="text-[10px] px-1 py-0.5 rounded bg-gray-500/20 text-gray-400">
                      {config.static_coins?.length}
                    </span>
                  )}
                </div>
                <div className="flex flex-wrap gap-1 mt-1 mb-2">
                  {(config.static_coins || []).slice(0, 3).map((coin) => (
                    <span key={coin} className="flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] bg-nofx-bg-lighter text-nofx-text">
                      {coin}
                      {!disabled && (
                        <button onClick={(e) => { e.stopPropagation(); handleRemoveCoin(coin); }} className="hover:text-red-400 transition-colors">
                          <X className="w-2 h-2" />
                        </button>
                      )}
                    </span>
                  ))}
                  {(config.static_coins || []).length > 3 && (
                    <span className="text-[10px] text-nofx-text-muted">+{(config.static_coins?.length || 0) - 3}</span>
                  )}
                </div>
                {!disabled && (
                  <div className="flex gap-1">
                    <input
                      type="text"
                      value={newCoin}
                      onChange={(e) => setNewCoin(e.target.value)}
                      onKeyDown={(e) => { e.stopPropagation(); if (e.key === 'Enter') handleAddCoin(); }}
                      placeholder="BTC, ETH..."
                      className="flex-1 px-2 py-1 rounded text-[10px] bg-nofx-bg border border-nofx-border text-nofx-text"
                    />
                    <button onClick={(e) => { e.stopPropagation(); handleAddCoin(); }} className="px-2 py-1 rounded text-[10px] bg-nofx-gold text-black hover:bg-yellow-500">
                      <Plus className="w-2.5 h-2.5" />
                    </button>
                  </div>
                )}
              </div>
            </div>

            {/* Middle visual arrow for desktop */}
            <div className="hidden md:flex flex-col items-center justify-center relative">
              <div className="text-xs text-nofx-text-muted mb-1 absolute -top-2 whitespace-nowrap bg-nofx-bg px-1 hidden lg:block">
                {language === 'zh' ? '汇流截断' : 'Truncation'}
              </div>
              <div className="h-px w-8 bg-nofx-border"></div>
              <TrendingUp className="w-4 h-4 text-nofx-text-muted rotate-90 md:rotate-0 absolute text-nofx-border" />
            </div>

            {/* Column 3: Blackbox Filter Cutoff */}
            <div className="flex-1 space-y-3 p-3 rounded-lg border border-red-500/20 bg-red-500/5">
              <div className="text-xs font-semibold text-red-400 mb-2 uppercase tracking-wider flex items-center gap-1">
                <Ban className="w-3.5 h-3.5" />
                3. {t('blackboxSource')}
              </div>

              <div className="p-3 bg-nofx-bg rounded border border-nofx-border relative overflow-hidden group">
                <div className="absolute top-0 right-0 w-16 h-16 bg-red-500/5 rounded-bl-full -mr-8 -mt-8 transition-transform group-hover:scale-110"></div>
                <div className="text-sm font-medium text-nofx-text mb-2 relative z-10">
                  {t('blackboxLimit')}
                </div>
                <div className="flex items-center gap-2 relative z-10">
                  <input
                    type="number"
                    value={config.blackbox_cutoff_limit || 5}
                    onChange={(e) => !disabled && onChange({ ...config, blackbox_cutoff_limit: Math.max(1, parseInt(e.target.value) || 5) })}
                    disabled={disabled}
                    min={1}
                    className="w-full px-3 py-2 rounded text-base font-bold bg-nofx-bg-lighter border border-red-500/30 text-nofx-text text-center focus:border-red-500/50 outline-none transition-colors"
                  />
                </div>
                <p className="text-[10px] text-nofx-text-muted mt-3 leading-relaxed relative z-10">
                  {language === 'zh' ? '无论左侧提取多少数据，经过严苛的本金与存活期清洗后，最多只允许此数量的候选者见 AI。' : 'Max number of survivors allowed to pass to AI after strict pipeline filtration.'}
                </p>
              </div>
            </div>

          </div>

          {/* Summary */}
          {(() => {
            const sources = []
            if (config.use_ai500) sources.push('AI500')
            if (config.use_binance_top_vol) sources.push('BinanceVol')
            if (config.use_oi_top) sources.push('OI↑')
            if (config.use_oi_low) sources.push('OI↓')
            if ((config.static_coins || []).length > 0) sources.push('Custom')

            if (sources.length === 0) return null
            return (
              <div className="mt-4 p-2 rounded bg-nofx-bg border border-nofx-border">
                <div className="flex items-center justify-between text-xs">
                  <span className="text-nofx-text-muted">{t('mixedSummary')}:</span>
                  <span className="text-nofx-text font-medium">
                    {sources.join(' + ')} <span className="text-nofx-text-muted mx-1">➜</span> {t('blackboxLimit')}: {config.blackbox_cutoff_limit || 5}
                  </span>
                </div>
              </div>
            )
          })()}
        </div>
      )}
    </div>
  )
}
