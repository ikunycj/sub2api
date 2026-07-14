import { activeBrand } from './runtime'
import type { BrandConfig } from './types'

export function getCurrentOrigin(): string {
  return typeof window === 'undefined' ? '' : window.location.origin
}

export function getCurrentHost(): string {
  return typeof window === 'undefined' ? '' : window.location.host
}

export function getCurrentApiBaseUrl(): string {
  const origin = getCurrentOrigin()
  return origin ? `${origin}/v1` : '/v1'
}

export function getBrandSlug(brand: BrandConfig = activeBrand): string {
  return brand.siteName
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '_')
    .replace(/^_+|_+$/g, '') || 'site'
}

export function replaceBrandTokens(text: string, brand: BrandConfig = activeBrand): string {
  const origin = getCurrentOrigin()
  const host = getCurrentHost()
  const apiBaseUrl = getCurrentApiBaseUrl()
  const siteSlug = getBrandSlug(brand)
  const replacements: Array<[RegExp, string]> = [
    [/\{\{siteName\}\}/g, brand.siteName],
    [/\{\{siteSlug\}\}/g, siteSlug],
    [/\{\{origin\}\}/g, origin || '/'],
    [/\{\{apiBaseUrl\}\}/g, apiBaseUrl],
    [/\{\{host\}\}/g, host || brand.siteName],
    [/https:\/\/api\.anytoken\.com\/v1/g, apiBaseUrl],
    [/https:\/\/api\.anytoken\.com/g, origin || '/'],
    [/api\.anytoken\.com\/v1/g, host ? `${host}/v1` : '/v1'],
    [/anytoken\.com/g, host || brand.siteName],
    [/\[model_providers\.anytoken\]/g, `[model_providers.${siteSlug}]`],
    [/(model_provider\s*=\s*["'])anytoken(["'])/g, `$1${siteSlug}$2`],
    [/("provider"\s*:\s*\{\s*\n\s*")anytoken("\s*:\s*\{)/g, `$1${siteSlug}$2`],
    [/ANYTOKEN_API_KEY/g, 'SITE_API_KEY'],
    [/AnyToken/g, brand.siteName],
    [/\banytoken\b/g, brand.siteName],
  ]

  return replacements.reduce((current, [pattern, replacement]) => {
    return current.replace(pattern, replacement)
  }, text)
}

export function withBrandTokens<T>(value: T, brand: BrandConfig = activeBrand): T {
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
