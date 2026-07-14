import { applyBrandTheme } from './theme'
import { activeBrand } from './runtime'

export { DEFAULT_BRAND_ENV, PRIMARY_COLOR_SCHEMES, resolveBrandConfig, sharedBrandContent } from './configs'
export { activeBrand } from './runtime'
export {
  getBrandSlug,
  getCurrentApiBaseUrl,
  getCurrentHost,
  getCurrentOrigin,
  replaceBrandTokens,
  withBrandTokens,
} from './text'
export { withBrandIdentity } from './settings'
export { applyBrandTheme, getBrandCssText, getBrandCssVariables } from './theme'
export { getBrandChartPalette, hexToRgba } from './chart'
export type {
  BrandEnvironmentConfig,
  BrandPrimaryScheme,
  BrandColorScale,
  BrandColorStep,
  BrandConfig,
  HomeCopy,
  LocaleKey,
  PublicLayoutCopy,
  SharedBrandContent,
} from './types'

export const applyActiveBrandTheme = () => applyBrandTheme(activeBrand)
