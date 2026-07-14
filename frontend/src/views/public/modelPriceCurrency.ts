export const BALANCE_USD_PER_CNY = 1
export const OFFICIAL_USD_TO_CNY_REFERENCE_RATE = 7.25

export function formatCompactModelAmount(value: number): string {
  return value.toPrecision(10).replace(/\.?0+$/, '')
}

export function calculateModelRechargePriceCny(unitPriceUsd: number, groupRateMultiplier: number, scale = 1): number {
  return unitPriceUsd * groupRateMultiplier * scale / BALANCE_USD_PER_CNY
}

export function calculateOfficialReferencePriceCny(unitPriceUsd: number, scale = 1): number {
  return unitPriceUsd * scale * OFFICIAL_USD_TO_CNY_REFERENCE_RATE
}

export function formatModelCnyAmount(amountCny: number): string {
  return `¥${formatCompactModelAmount(amountCny)}`
}

export function formatModelRateMultiplier(rateMultiplier: number): string {
  return `×${formatCompactModelAmount(rateMultiplier)}`
}
