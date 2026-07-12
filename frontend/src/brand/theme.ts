import type { BrandConfig, BrandColorStep } from './types'

export const BRAND_COLOR_STEPS: BrandColorStep[] = [
  '50',
  '100',
  '200',
  '300',
  '400',
  '500',
  '600',
  '700',
  '800',
  '900',
  '950',
]

function hexToRgbChannels(hex: string): string {
  const normalized = hex.replace('#', '').trim()
  const value = normalized.length === 3
    ? normalized.split('').map((char) => `${char}${char}`).join('')
    : normalized

  const red = Number.parseInt(value.slice(0, 2), 16)
  const green = Number.parseInt(value.slice(2, 4), 16)
  const blue = Number.parseInt(value.slice(4, 6), 16)

  if ([red, green, blue].some((channel) => Number.isNaN(channel))) {
    throw new Error(`Invalid brand color: ${hex}`)
  }

  return `${red} ${green} ${blue}`
}

export function getBrandCssVariables(brand: BrandConfig): Record<string, string> {
  return BRAND_COLOR_STEPS.reduce<Record<string, string>>((variables, step) => {
    variables[`--color-primary-${step}`] = hexToRgbChannels(brand.theme.primary[step])
    return variables
  }, {})
}

export function getBrandCssText(brand: BrandConfig): string {
  const variables = getBrandCssVariables(brand)
  return Object.entries(variables)
    .map(([name, value]) => `${name}: ${value};`)
    .join(' ')
}

export function applyBrandTheme(brand: BrandConfig, root: HTMLElement = document.documentElement): void {
  const variables = getBrandCssVariables(brand)
  Object.entries(variables).forEach(([name, value]) => {
    root.style.setProperty(name, value)
  })
  root.dataset.brand = brand.key
}
