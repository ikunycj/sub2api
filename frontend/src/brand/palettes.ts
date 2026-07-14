export type BrandColorStep =
  | '50'
  | '100'
  | '200'
  | '300'
  | '400'
  | '500'
  | '600'
  | '700'
  | '800'
  | '900'
  | '950'

export type BrandColorScale = Record<BrandColorStep, string>
export type BrandPrimaryScheme = 'orange' | 'blue'

export const DEFAULT_BRAND_ENV = {
  siteName: 'ikun.love',
  logo: '/logo.png',
  primaryScheme: 'orange',
} as const

export const PRIMARY_COLOR_SCHEMES = {
  orange: {
    50: '#fff8ed', 100: '#ffefd0', 200: '#ffdaa1', 300: '#ffbe6b',
    400: '#ff9832', 500: '#f97316', 600: '#ea580c', 700: '#c2410c',
    800: '#9a3412', 900: '#7c2d12', 950: '#431407',
  },
  blue: {
    50: '#eff6ff', 100: '#dbeafe', 200: '#bfdbfe', 300: '#93c5fd',
    400: '#60a5fa', 500: '#2563eb', 600: '#1d4ed8', 700: '#1e40af',
    800: '#1e3a8a', 900: '#172554', 950: '#0b1120',
  },
} satisfies Record<BrandPrimaryScheme, BrandColorScale>
