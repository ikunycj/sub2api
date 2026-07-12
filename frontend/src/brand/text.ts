import { resolveBrandConfig } from './configs'
import type { BrandConfig } from './types'

const currentBrand = resolveBrandConfig(import.meta.env.VITE_BRAND)

function apiRoot(apiBaseUrl: string): string {
  return apiBaseUrl.replace(/\/v1\/?$/, '')
}

function apiBaseHostPath(apiBaseUrl: string): string {
  return apiBaseUrl.replace(/^https?:\/\//, '')
}

export function replaceBrandTokens(text: string, brand: BrandConfig = currentBrand): string {
  const replacements: Array<[RegExp, string]> = [
    [/https:\/\/api\.anytoken\.com\/v1/g, brand.apiBaseUrl],
    [/https:\/\/api\.anytoken\.com/g, apiRoot(brand.apiBaseUrl)],
    [/api\.anytoken\.com\/v1/g, apiBaseHostPath(brand.apiBaseUrl)],
    [/anytoken\.com/g, brand.domain],
    [/ANYTOKEN_API_KEY/g, brand.apiKeyEnvVar],
    [/AnyToken/g, brand.siteName],
    [/\banytoken\b/g, brand.shortName],
  ]

  return replacements.reduce((current, [pattern, replacement]) => {
    return current.replace(pattern, replacement)
  }, text)
}

export function withBrandTokens<T>(value: T, brand: BrandConfig = currentBrand): T {
  if (typeof value === 'string') {
    return replaceBrandTokens(value, brand) as T
  }

  if (Array.isArray(value)) {
    return value.map((item) => withBrandTokens(item, brand)) as T
  }

  if (value && typeof value === 'object') {
    return Object.fromEntries(
      Object.entries(value).map(([key, item]) => [key, withBrandTokens(item, brand)])
    ) as T
  }

  return value
}
