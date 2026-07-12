import { anytokenBrand } from './anytoken'
import { ikunBrand } from './ikun'
import type { BrandConfig, BrandKey } from '../types'

export const DEFAULT_BRAND: BrandKey = 'ikun'

export const brandConfigs = {
  ikun: ikunBrand,
  anytoken: anytokenBrand,
} satisfies Record<BrandKey, BrandConfig>

export function resolveBrandConfig(value?: string): BrandConfig {
  const normalized = (value || DEFAULT_BRAND).trim().toLowerCase()

  if (normalized in brandConfigs) {
    return brandConfigs[normalized as BrandKey]
  }

  const supportedBrands = Object.keys(brandConfigs).join(', ')
  throw new Error(`Unsupported VITE_BRAND "${value}". Supported brands: ${supportedBrands}`)
}
