package filter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// DynamicOverride 指令集结构 (V3.0 精准对账版)
type DynamicOverride struct {
	Header struct {
		Timestamp string `json:"timestamp"`
		Source    string `json:"source"`
		Version   string `json:"version"`
	} `json:"header"`
	MarketRegimeOverride string `json:"market_regime_override"`
	RiskParameters       struct {
		Mode           string  `json:"mode"`
		MaxMarginRatio float64 `json:"max_margin_ratio"`
	} `json:"risk_parameters"`
	SymbolRecommendations []struct {
		Symbol         string   `json:"symbol"`
		ActionOverride string   `json:"action_override"`
		BurstScore     int      `json:"burst_score"`
		Flags          []string `json:"flags"`
	} `json:"symbol_recommendations"`
	Reflection string `json:"reflection"`
}

// GetDynamicOverrideText 读取并翻译 openclaw_rec.json 为 Prompt 文本
func GetDynamicOverrideText() string {
	// 获取可执行文件路径或直接定位到 /root/nofx/data
	// 这里直接硬编码路径以保证在 Docker 挂载下 100% 准确
	recPath := "/app/data/openclaw_rec.json"
	
	// 如果文件不存在，返回空
	if _, err := os.Stat(recPath); os.IsNotExist(err) {
		return ""
	}

	data, err := os.ReadFile(recPath)
	if err != nil {
		log.Printf("[Commander] Error reading rec file: %v", err)
		return ""
	}

	var rec DynamicOverride
	if err := json.Unmarshal(data, &rec); err != nil {
		log.Printf("[Commander] Error unmarshaling rec file: %v", err)
		return ""
	}

	// 检查新鲜度 (必须是 1 小时内的指令)
	t, err := time.Parse(time.RFC3339, rec.Header.Timestamp)
	if err == nil {
		if time.Since(t) > 1*time.Hour {
			log.Printf("[Commander] Override file is too old: %v", rec.Header.Timestamp)
			return ""
		}
	}

	// 开始组装 Prompt 文本
	var sb strings.Builder
	sb.WriteString("[COMMANDER STRATEGIC OVERRIDE]\n")
	sb.WriteString(fmt.Sprintf("Global Regime: %s\n", rec.MarketRegimeOverride))
	sb.WriteString(fmt.Sprintf("Risk Strategy: %s (Max Margin: %.1f%%)\n", 
		rec.RiskParameters.Mode, rec.RiskParameters.MaxMarginRatio*100))
	
	if len(rec.SymbolRecommendations) > 0 {
		sb.WriteString("Specific Directives:\n")
		for _, item := range rec.SymbolRecommendations {
			sb.WriteString(fmt.Sprintf("- %s: %s [Priority: %d]", 
				item.Symbol, item.ActionOverride, item.BurstScore))
			if len(item.Flags) > 0 {
				sb.WriteString(fmt.Sprintf(" (Tags: %s)", strings.Join(item.Flags, ", ")))
			}
			sb.WriteString("\n")
		}
	}
	
	if rec.Reflection != "" {
		sb.WriteString(fmt.Sprintf("\nStrategic Intelligence: %s\n", rec.Reflection))
	}

	return sb.String()
}
