package market

import (
	"math"
)

// calculateADX calculates the Average Directional Index (ADX)
// Returns ADX, DI+, DI- (all last values)
func calculateADX(klines []Kline, period int) (float64, float64, float64) {
	if len(klines) < period*2 {
		return 0, 0, 0
	}

	// 1. Calculate TR, +DM, -DM
	trs := make([]float64, len(klines))
	plusDMs := make([]float64, len(klines))
	minusDMs := make([]float64, len(klines))

	for i := 1; i < len(klines); i++ {
		high := klines[i].High
		low := klines[i].Low
		prevHigh := klines[i-1].High
		prevLow := klines[i-1].Low
		prevClose := klines[i-1].Close

		// True Range
		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)
		trs[i] = math.Max(tr1, math.Max(tr2, tr3))

		// Directional Movement
		upMove := high - prevHigh
		downMove := prevLow - low

		if upMove > downMove && upMove > 0 {
			plusDMs[i] = upMove
		} else {
			plusDMs[i] = 0
		}

		if downMove > upMove && downMove > 0 {
			minusDMs[i] = downMove
		} else {
			minusDMs[i] = 0
		}
	}

	// 2. Smoothed TR, +DM, -DM (Wilder's Smoothing)
	// First value is simple sum of period
	smtTr := 0.0
	smtPlusDM := 0.0
	smtMinusDM := 0.0

	for i := 1; i <= period; i++ {
		smtTr += trs[i]
		smtPlusDM += plusDMs[i]
		smtMinusDM += minusDMs[i]
	}

	// Calculate initial DI
	// Avoid division by zero
	if smtTr == 0 {
		return 0, 0, 0
	}

	// Result arrays for DX calculation
	dxValues := make([]float64, len(klines))

	// Iterate from period + 1
	for i := period + 1; i < len(klines); i++ {
		// Wilder's Smoothing: Previous - (Previous/N) + Current
		smtTr = smtTr - (smtTr / float64(period)) + trs[i]
		smtPlusDM = smtPlusDM - (smtPlusDM / float64(period)) + plusDMs[i]
		smtMinusDM = smtMinusDM - (smtMinusDM / float64(period)) + minusDMs[i]

		if smtTr == 0 {
			dxValues[i] = 0
			continue
		}

		diPlus := (smtPlusDM / smtTr) * 100
		diMinus := (smtMinusDM / smtTr) * 100

		sumDI := diPlus + diMinus
		if sumDI == 0 {
			dxValues[i] = 0
		} else {
			dxValues[i] = (math.Abs(diPlus-diMinus) / sumDI) * 100
		}
	}

	// 3. Calculate ADX (Smoothed DX)
	// First ADX is average of period DX values
	startIdx := period * 2 // We need period samples for DX, so first ADX is at 2*period
	if startIdx >= len(klines) {
		return 0, 0, 0 // Not enough data
	}

	adxSum := 0.0
	for i := period + 1; i <= startIdx; i++ {
		adxSum += dxValues[i]
	}
	currentADX := adxSum / float64(period)

	// Calculate subsequent ADX
	for i := startIdx + 1; i < len(klines); i++ {
		currentADX = ((currentADX * float64(period-1)) + dxValues[i]) / float64(period)
	}

	// Calculate final DI values for return
	finalDIPlus := (smtPlusDM / smtTr) * 100
	finalDIMinus := (smtMinusDM / smtTr) * 100

	return currentADX, finalDIPlus, finalDIMinus
}

// calculateADXSeries calculates ADX, DI+, DI- series
func calculateADXSeries(klines []Kline, period int) ([]float64, []float64, []float64) {
	if len(klines) < period*2 {
		return nil, nil, nil
	}

	adxValues := make([]float64, len(klines))
	diPlusValues := make([]float64, len(klines))
	diMinusValues := make([]float64, len(klines))

	// 1. Calculate TR, +DM, -DM
	trs := make([]float64, len(klines))
	plusDMs := make([]float64, len(klines))
	minusDMs := make([]float64, len(klines))

	for i := 1; i < len(klines); i++ {
		high := klines[i].High
		low := klines[i].Low
		prevHigh := klines[i-1].High
		prevLow := klines[i-1].Low
		prevClose := klines[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)
		trs[i] = math.Max(tr1, math.Max(tr2, tr3))

		upMove := high - prevHigh
		downMove := prevLow - low

		if upMove > downMove && upMove > 0 {
			plusDMs[i] = upMove
		} else {
			plusDMs[i] = 0
		}

		if downMove > upMove && downMove > 0 {
			minusDMs[i] = downMove
		} else {
			minusDMs[i] = 0
		}
	}

	// 2. Smoothed TR, +DM, -DM (Wilder's Smoothing)
	smtTr := 0.0
	smtPlusDM := 0.0
	smtMinusDM := 0.0

	for i := 1; i <= period; i++ {
		smtTr += trs[i]
		smtPlusDM += plusDMs[i]
		smtMinusDM += minusDMs[i]
	}

	// Result arrays for DX calculation
	dxValues := make([]float64, len(klines))

	for i := period + 1; i < len(klines); i++ {
		smtTr = smtTr - (smtTr / float64(period)) + trs[i]
		smtPlusDM = smtPlusDM - (smtPlusDM / float64(period)) + plusDMs[i]
		smtMinusDM = smtMinusDM - (smtMinusDM / float64(period)) + minusDMs[i]

		if smtTr == 0 {
			dxValues[i] = 0
			diPlusValues[i] = 0
			diMinusValues[i] = 0
			continue
		}

		diPlus := (smtPlusDM / smtTr) * 100
		diMinus := (smtMinusDM / smtTr) * 100

		diPlusValues[i] = diPlus
		diMinusValues[i] = diMinus

		sumDI := diPlus + diMinus
		if sumDI == 0 {
			dxValues[i] = 0
		} else {
			dxValues[i] = (math.Abs(diPlus-diMinus) / sumDI) * 100
		}
	}

	// 3. Calculate ADX (Smoothed DX)
	startIdx := period * 2
	if startIdx >= len(klines) {
		return nil, nil, nil
	}

	adxSum := 0.0
	for i := period + 1; i <= startIdx; i++ {
		adxSum += dxValues[i]
	}
	// First ADX
	adxValues[startIdx] = adxSum / float64(period)

	for i := startIdx + 1; i < len(klines); i++ {
		adxValues[i] = ((adxValues[i-1] * float64(period-1)) + dxValues[i]) / float64(period)
	}

	return adxValues, diPlusValues, diMinusValues
}
