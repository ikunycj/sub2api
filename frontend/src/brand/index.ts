import { resolveBrandConfig } from './configs'
import { applyBrandTheme } from './theme'

export { brandConfigs, DEFAULT_BRAND, resolveBrandConfig } from './configs'
export { replaceBrandTokens, withBrandTokens } from './text'
export { applyBrandTheme, getBrandCssText, getBrandCssVariables } from './theme'
export type {
  BrandColorScale,
  BrandColorStep,
  BrandConfig,
  BrandKey,
  HomeCopy,
  LocaleKey,
  PublicLayoutCopy,
} from './types'

export const activeBrand = resolveBrandConfig(import.meta.env.VITE_BRAND)
export const applyActiveBrandTheme = () => applyBrandTheme(activeBrand)
