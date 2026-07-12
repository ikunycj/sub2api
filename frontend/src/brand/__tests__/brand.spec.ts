import { describe, expect, it } from 'vitest'
import { brandConfigs, resolveBrandConfig, withBrandTokens } from '@/brand'
import { getBrandCssVariables } from '@/brand/theme'

describe('brand config', () => {
  it('defaults to ikun when VITE_BRAND is not provided', () => {
    expect(resolveBrandConfig().key).toBe('ikun')
  })

  it('resolves supported brands and rejects unsupported values', () => {
    expect(resolveBrandConfig('ikun').siteName).toBe('ikun.love')
    expect(resolveBrandConfig('anytoken').siteName).toBe('AnyToken')
    expect(() => resolveBrandConfig('unknown')).toThrow(/Unsupported VITE_BRAND/)
  })

  it('keeps every brand self-contained', () => {
    Object.values(brandConfigs).forEach((brand) => {
      expect(brand.logo).toBe(`/brands/${brand.key}/logo.svg`)
      expect(brand.favicon).toBe(`/brands/${brand.key}/logo.svg`)
      expect(brand.publicLayout.zh.start).toBeTruthy()
      expect(brand.publicLayout.en.start).toBeTruthy()
      expect(brand.home.zh.hero.title).toBeTruthy()
      expect(brand.home.en.hero.title).toBeTruthy()
      expect(Object.keys(brand.theme.primary)).toHaveLength(11)
    })
  })

  it('exposes primary colors as CSS variable channel values', () => {
    expect(getBrandCssVariables(brandConfigs.ikun)['--color-primary-500']).toBe('249 115 22')
    expect(getBrandCssVariables(brandConfigs.anytoken)['--color-primary-500']).toBe('37 99 235')
  })

  it('keeps shared model/docs pages brand-aware through token replacement', () => {
    const text = 'AnyToken uses https://api.anytoken.com/v1 with ANYTOKEN_API_KEY on anytoken.com.'
    expect(withBrandTokens(text, brandConfigs.ikun)).toBe(
      'ikun.love uses https://ikun.love/v1 with IKUN_API_KEY on ikun.love.'
    )
    expect(withBrandTokens(text, brandConfigs.anytoken)).toBe(
      'AnyToken uses https://api.anytoken.com/v1 with ANYTOKEN_API_KEY on anytoken.com.'
    )
  })

  it('lets only brand copy differ while the home component stays shared', () => {
    expect(brandConfigs.ikun.home.zh.hero.title).not.toBe(brandConfigs.anytoken.home.zh.hero.title)
    expect(brandConfigs.ikun.home.zh.hero.secondaryCta).toBeTruthy()
    expect(brandConfigs.anytoken.home.zh.hero.secondaryCta).toBeTruthy()
  })
})
