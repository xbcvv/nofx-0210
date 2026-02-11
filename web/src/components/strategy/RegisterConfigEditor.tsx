import { BookOpen } from 'lucide-react'

interface RegisterConfig {
  enabled: boolean
  max_records: number
  include_decisions: boolean
  include_market_data: boolean
}

interface RegisterConfigEditorProps {
  config: RegisterConfig
  onChange: (config: RegisterConfig) => void
  disabled?: boolean
  language: string
}

export function RegisterConfigEditor({
  config,
  onChange,
  disabled,
  language,
}: RegisterConfigEditorProps) {
  const t = (key: string) => {
    const translations: Record<string, Record<string, string>> = {
      registerConfig: { zh: '寄存器配置', en: 'Register Config' },
      registerEnabled: { zh: '启用寄存器', en: 'Enable Register' },
      registerEnabledDesc: { zh: '记录历史决策并在下次轮询时作为参考', en: 'Record historical decisions and use as reference in next polling' },
      maxRecords: { zh: '最大记录数', en: 'Max Records' },
      maxRecordsDesc: { zh: '保存的最大决策记录数量', en: 'Maximum number of decision records to save' },
      includeDecisions: { zh: '包含完整决策', en: 'Include Full Decisions' },
      includeDecisionsDesc: { zh: '在寄存器中包含完整的决策详情', en: 'Include complete decision details in register' },
      includeMarketData: { zh: '包含市场数据', en: 'Include Market Data' },
      includeMarketDataDesc: { zh: '在寄存器中包含市场数据', en: 'Include market data in register' },
      tokenUsage: { zh: 'Token 使用提示', en: 'Token Usage Tip' },
      tokenUsageDesc: { zh: '更多记录会增加 AI 提示词长度，建议保持在 5-10 条记录', en: 'More records increase AI prompt length, recommended to keep 5-10 records' },
    }
    return translations[key]?.[language] || key
  }

  const updateField = <K extends keyof RegisterConfig>(
    key: K,
    value: RegisterConfig[K]
  ) => {
    if (!disabled) {
      onChange({ ...config, [key]: value })
    }
  }

  return (
    <div className="space-y-6">
      {/* Register Config */}
      <div>
        <div className="flex items-center gap-2 mb-4">
          <BookOpen className="w-5 h-5" style={{ color: '#6366F1' }} />
          <h3 className="font-medium" style={{ color: '#EAECEF' }}>
            {t('registerConfig')}
          </h3>
        </div>

        {/* Enable Register */}
        <div
          className="p-4 rounded-lg"
          style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
        >
          <div className="flex items-center justify-between">
            <div>
              <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
                {t('registerEnabled')}
              </label>
              <p className="text-xs" style={{ color: '#848E9C' }}>
                {t('registerEnabledDesc')}
              </p>
            </div>
            <button
              type="button"
              onClick={() => updateField('enabled', !config.enabled)}
              disabled={disabled}
              className={`w-12 h-6 rounded-full transition-colors ${
                config.enabled
                  ? 'bg-indigo-600'
                  : 'bg-gray-700'
              }`}
            >
              <div
                className={`w-5 h-5 rounded-full bg-white transition-transform ${
                  config.enabled
                    ? 'transform translate-x-6'
                    : 'transform translate-x-1'
                }`}
              />
            </button>
          </div>
        </div>

        {/* Register Settings */}
        {config.enabled && (
          <div className="space-y-4 mt-4">
            {/* Max Records */}
            <div
              className="p-4 rounded-lg"
              style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
            >
              <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
                {t('maxRecords')}
              </label>
              <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
                {t('maxRecordsDesc')}
              </p>
              <div className="flex items-center gap-2">
                <input
                  type="range"
                  value={config.max_records || 5}
                  onChange={(e) =>
                    updateField('max_records', parseInt(e.target.value))
                  }
                  disabled={disabled}
                  min={1}
                  max={20}
                  className="flex-1 accent-indigo-500"
                />
                <span
                  className="w-12 text-center font-mono"
                  style={{ color: '#6366F1' }}
                >
                  {config.max_records || 5}
                </span>
              </div>
            </div>

            {/* Include Decisions */}
            <div
              className="p-4 rounded-lg"
              style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
            >
              <div className="flex items-center justify-between">
                <div>
                  <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
                    {t('includeDecisions')}
                  </label>
                  <p className="text-xs" style={{ color: '#848E9C' }}>
                    {t('includeDecisionsDesc')}
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => updateField('include_decisions', !config.include_decisions)}
                  disabled={disabled}
                  className={`w-12 h-6 rounded-full transition-colors ${
                    config.include_decisions
                      ? 'bg-indigo-600'
                      : 'bg-gray-700'
                  }`}
                >
                  <div
                    className={`w-5 h-5 rounded-full bg-white transition-transform ${
                      config.include_decisions
                        ? 'transform translate-x-6'
                        : 'transform translate-x-1'
                    }`}
                  />
                </button>
              </div>
            </div>

            {/* Include Market Data */}
            <div
              className="p-4 rounded-lg"
              style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
            >
              <div className="flex items-center justify-between">
                <div>
                  <label className="block text-sm mb-1" style={{ color: '#EAECEF' }}>
                    {t('includeMarketData')}
                  </label>
                  <p className="text-xs" style={{ color: '#848E9C' }}>
                    {t('includeMarketDataDesc')}
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => updateField('include_market_data', !config.include_market_data)}
                  disabled={disabled}
                  className={`w-12 h-6 rounded-full transition-colors ${
                    config.include_market_data
                      ? 'bg-indigo-600'
                      : 'bg-gray-700'
                  }`}
                >
                  <div
                    className={`w-5 h-5 rounded-full bg-white transition-transform ${
                      config.include_market_data
                        ? 'transform translate-x-6'
                        : 'transform translate-x-1'
                    }`}
                  />
                </button>
              </div>
            </div>

            {/* Token Usage Tip */}
            <div
              className="p-4 rounded-lg"
              style={{ background: '#1A1F29', border: '1px solid #2D3748' }}
            >
              <div className="flex items-start gap-2">
                <div className="flex-shrink-0 mt-0.5">
                  <div className="w-5 h-5 rounded-full bg-indigo-100 flex items-center justify-center">
                    <span className="text-indigo-800 text-xs font-bold">!</span>
                  </div>
                </div>
                <div>
                  <h4 className="text-sm font-medium" style={{ color: '#E2E8F0' }}>
                    {t('tokenUsage')}
                  </h4>
                  <p className="text-xs mt-1" style={{ color: '#A0AEC0' }}>
                    {t('tokenUsageDesc')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
