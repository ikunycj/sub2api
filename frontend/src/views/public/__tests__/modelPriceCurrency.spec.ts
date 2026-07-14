import { describe, expect, it } from 'vitest'
import {
  calculateModelRechargePriceCny,
  calculateOfficialReferencePriceCny,
  formatModelCnyAmount,
  formatModelRateMultiplier,
} from '../modelPriceCurrency'

describe('public model price currency formatting', () => {
  it('calculates the RMB recharge needed after the group multiplier', () => {
    const rechargePrice = calculateModelRechargePriceCny(1.75 / 1_000_000, 0.1, 1_000_000)

    expect(rechargePrice).toBeCloseTo(0.175)
    expect(formatModelCnyAmount(rechargePrice)).toBe('¥0.175')
  })

  it('converts the official USD reference price to RMB separately', () => {
    const officialPrice = calculateOfficialReferencePriceCny(1.75 / 1_000_000, 1_000_000)

    expect(officialPrice).toBeCloseTo(12.6875)
    expect(formatModelCnyAmount(officialPrice)).toBe('¥12.6875')
  })

  it('formats backend group multipliers consistently', () => {
    expect(formatModelRateMultiplier(2)).toBe('×2')
    expect(formatModelRateMultiplier(0.5)).toBe('×0.5')
  })
})
