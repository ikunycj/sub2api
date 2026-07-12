export type BrandKey = 'ikun' | 'anytoken'
export type LocaleKey = 'zh' | 'en'
export type BrandColorStep = '50' | '100' | '200' | '300' | '400' | '500' | '600' | '700' | '800' | '900' | '950'
export type BrandColorScale = Record<BrandColorStep, string>

export interface BrandTheme {
  primary: BrandColorScale
}

export interface PublicLayoutCopy {
  home: string
  models: string
  docs: string
  start: string
  product: string
  endpoint: string
  footerIntro: string
  lightMode: string
  darkMode: string
}

export interface HomeFeatureItem {
  mark: string
  title: string
  description: string
}

export interface HomeStepItem {
  title: string
  description: string
  meta: string
}

export interface HomeOpsItem {
  number: string
  title: string
  description: string
  cta: string
  to: string
}

export interface HomePricingRow {
  name: string
  modality: string
  input: string
  output: string
  cache: string
  context: string
  policy: string
}

export interface HomeTestimonialItem {
  role: string
  quote: string
}

export interface HomeFaqItem {
  question: string
  answer: string
}

export interface HomeCopy {
  metaTitle: string
  hero: {
    badge: string
    title: string
    description: string
    primaryCta: string
    secondaryCta: string
  }
  heroPanel: {
    title: string
    status: string
  }
  providersTitle: string
  modalities: {
    kicker: string
    title: string
    description: string
    items: HomeFeatureItem[]
  }
  steps: {
    kicker: string
    title: string
    description: string
    items: HomeStepItem[]
  }
  ops: {
    kicker: string
    title: string
    description: string
    items: HomeOpsItem[]
  }
  sdk: {
    title: string
    description: string
    primaryCta: string
    secondaryCta: string
    code: string
  }
  pricing: {
    kicker: string
    title: string
    description: string
    cta: string
    columns: string[]
    rows: HomePricingRow[]
  }
  testimonials: {
    kicker: string
    title: string
    items: HomeTestimonialItem[]
  }
  faq: {
    kicker: string
    title: string
    items: HomeFaqItem[]
  }
}

export interface BrandConfig {
  key: BrandKey
  siteName: string
  shortName: string
  domain: string
  apiBaseUrl: string
  apiKeyEnvVar: string
  logo: string
  favicon: string
  logoInitial: string
  titleSuffix: string
  theme: BrandTheme
  publicLayout: Record<LocaleKey, PublicLayoutCopy>
  home: Record<LocaleKey, HomeCopy>
}
