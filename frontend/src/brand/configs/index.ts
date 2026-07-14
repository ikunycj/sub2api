import { DEFAULT_BRAND_ENV, PRIMARY_COLOR_SCHEMES } from '../palettes'
import type { BrandPrimaryScheme } from '../palettes'
import type { BrandConfig, BrandEnvironmentConfig } from '../types'

export { sharedBrandContent } from './shared'
export { DEFAULT_BRAND_ENV, PRIMARY_COLOR_SCHEMES } from '../palettes'

export function resolveBrandConfig(input: BrandEnvironmentConfig = {}): BrandConfig {
  const siteName = input.siteName?.trim() || DEFAULT_BRAND_ENV.siteName
  const logo = input.logo?.trim() || DEFAULT_BRAND_ENV.logo
  const primaryScheme = (input.primaryScheme?.trim().toLowerCase() || DEFAULT_BRAND_ENV.primaryScheme) as BrandPrimaryScheme
  const primary = PRIMARY_COLOR_SCHEMES[primaryScheme]

  if (!primary) {
    throw new Error(`Unsupported VITE_PRIMARY_SCHEME "${input.primaryScheme}". Supported schemes: ${Object.keys(PRIMARY_COLOR_SCHEMES).join(', ')}`)
  }

  return {
    siteName,
    logo,
    primaryScheme,
    theme: { primary },
  }
}
