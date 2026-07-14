import { describe, expect, it } from 'vitest'
import { calculateSavingsPercent, hasOfficialReferencePrice, isDiscountedPrice } from '../modelPriceComparison'

describe('model price comparison', () => {
  it('calculates the displayed saving for a 1/50 price', () => {
    expect(isDiscountedPrice(0.02, 1)).toBe(true)
    expect(calculateSavingsPercent(0.02, 1)).toBe(98)
  })

  it('does not present equal or higher prices as discounts', () => {
    expect(hasOfficialReferencePrice(1)).toBe(true)
    expect(isDiscountedPrice(1, 1)).toBe(false)
    expect(isDiscountedPrice(1.1, 1)).toBe(false)
    expect(calculateSavingsPercent(1, 1)).toBeNull()
  })

  it('ignores missing and invalid official prices', () => {
    expect(hasOfficialReferencePrice(null)).toBe(false)
    expect(hasOfficialReferencePrice(Number.NaN)).toBe(false)
    expect(hasOfficialReferencePrice(-1)).toBe(false)
    expect(isDiscountedPrice(null, 1)).toBe(false)
    expect(isDiscountedPrice(0.5, null)).toBe(false)
    expect(isDiscountedPrice(0.5, 0)).toBe(false)
  })
})
