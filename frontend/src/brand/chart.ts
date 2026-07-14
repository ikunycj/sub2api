import type { BrandConfig } from './types'

export function getBrandChartPalette(brand: BrandConfig): string[] {
  return [
    brand.theme.primary['500'],
    '#2563eb',
    '#059669',
    '#dc2626',
    '#0891b2',
    '#ca8a04',
    '#65a30d',
    '#0d9488',
    '#0284c7',
    brand.theme.primary['600'],
    '#16a34a',
    '#475569',
  ]
}

export function hexToRgba(hex: string, alpha: number): string {
  const normalized = hex.replace('#', '')
  const red = Number.parseInt(normalized.slice(0, 2), 16)
  const green = Number.parseInt(normalized.slice(2, 4), 16)
  const blue = Number.parseInt(normalized.slice(4, 6), 16)
  return `rgba(${red}, ${green}, ${blue}, ${alpha})`
}
