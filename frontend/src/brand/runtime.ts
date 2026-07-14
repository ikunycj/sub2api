import { resolveBrandConfig } from './configs'

export const activeBrand = resolveBrandConfig({
  siteName: import.meta.env.VITE_SITE_NAME,
  logo: import.meta.env.VITE_SITE_LOGO,
  primaryScheme: import.meta.env.VITE_PRIMARY_SCHEME,
})
