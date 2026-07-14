import type { BrandConfig } from './types'

type BrandIdentitySettings = {
  site_name?: string
  site_logo?: string
}

export function withBrandIdentity<T extends BrandIdentitySettings>(settings: T, brand: BrandConfig): T {
  return {
    ...settings,
    site_name: brand.siteName,
    site_logo: brand.logo,
  }
}
