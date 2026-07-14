import { describe, expect, it } from 'vitest'
import {
  getBrandChartPalette,
  getBrandSlug,
  resolveBrandConfig,
  sharedBrandContent,
  withBrandIdentity,
  withBrandTokens,
} from '@/brand'
import { getBrandCssVariables } from '@/brand/theme'

const ikun = resolveBrandConfig({
  siteName: 'ikun.love',
  logo: '/logo.png',
  primaryScheme: 'orange',
})

const anytoken = resolveBrandConfig({
  siteName: 'AnyToken',
  logo: '/brands/anytoken/logo.svg',
  primaryScheme: 'blue',
})

describe('site brand config', () => {
  it('uses the original ikun identity and palette by default', () => {
    expect(resolveBrandConfig()).toEqual(ikun)
    expect(ikun.theme.primary['500']).toBe('#f97316')
  })

  it('accepts only a supported primary scheme', () => {
    expect(anytoken.theme.primary['500']).toBe('#2563eb')
    expect(() => resolveBrandConfig({ primaryScheme: 'unknown' })).toThrow(
      /Unsupported VITE_PRIMARY_SCHEME/,
    )
  })

  it('contains no per-site page, domain, secondary color, or feature config', () => {
    expect(Object.keys(ikun).sort()).toEqual(['logo', 'primaryScheme', 'siteName', 'theme'])
    expect(Object.keys(anytoken).sort()).toEqual(['logo', 'primaryScheme', 'siteName', 'theme'])
    expect(Object.keys(ikun.theme)).toEqual(['primary'])
    expect(Object.keys(anytoken.theme)).toEqual(['primary'])
  })

  it('exposes only the selected primary palette as runtime CSS variables', () => {
    expect(getBrandCssVariables(ikun)['--color-primary-500']).toBe('249 115 22')
    expect(getBrandCssVariables(anytoken)['--color-primary-500']).toBe('37 99 235')
    expect(Object.keys(getBrandCssVariables(ikun)).some((key) => key.includes('secondary'))).toBe(false)
  })

  it('keeps backend settings from overriding the configured name and logo', () => {
    expect(withBrandIdentity({ site_name: 'backend', site_logo: '/backend.png' }, ikun)).toEqual({
      site_name: 'ikun.love',
      site_logo: '/logo.png',
    })
  })

  it('changes only the two primary shades in the shared chart palette', () => {
    const ikunPalette = getBrandChartPalette(ikun)
    const anytokenPalette = getBrandChartPalette(anytoken)
    expect(ikunPalette[0]).toBe('#f97316')
    expect(anytokenPalette[0]).toBe('#2563eb')
    expect(ikunPalette[9]).toBe('#ea580c')
    expect(anytokenPalette[9]).toBe('#1d4ed8')
    expect(ikunPalette.filter((_, index) => ![0, 9].includes(index))).toEqual(
      anytokenPalette.filter((_, index) => ![0, 9].includes(index)),
    )
  })

  it('renders the same shared copy with only runtime identity and endpoint tokens changed', () => {
    const template = 'AnyToken uses https://api.anytoken.com/v1 with ANYTOKEN_API_KEY.'
    const expectedEndpoint = `${window.location.origin}/v1`

    expect(withBrandTokens(template, ikun)).toBe(
      `ikun.love uses ${expectedEndpoint} with SITE_API_KEY.`,
    )
    expect(withBrandTokens(template, anytoken)).toBe(
      `AnyToken uses ${expectedEndpoint} with SITE_API_KEY.`,
    )
    expect(getBrandSlug(ikun)).toBe('ikun_love')
    expect(getBrandSlug(anytoken)).toBe('anytoken')
  })

  it('keeps one shared homepage and public-layout content tree', () => {
    const copy = JSON.stringify(sharedBrandContent)
    expect(copy).toContain('OpenAI')
    expect(copy).toContain('Claude')
    expect(copy).toContain('1/50')
    expect(copy).not.toMatch(/Gemini|DeepSeek|Qwen|Bytedance|ElevenLabs|Minimax|Kling|Vidu|Grok|Runway/)
  })
})
