import { apiClient } from './client'
import type { BillingMode } from '@/constants/channel'

export interface PublicPricingInterval {
  min_tokens: number
  max_tokens: number | null
  tier_label?: string
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  per_request_price: number | null
}

export interface PublicModelPricing {
  billing_mode: BillingMode
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  image_output_price: number | null
  per_request_price: number | null
  intervals: PublicPricingInterval[]
}

export interface PublicCatalogModel {
  name: string
  platform: string
  pricing: PublicModelPricing | null
  official_pricing?: PublicModelPricing | null
}

export interface PublicCatalogGroup {
  id: number
  name: string
  platform: string
  rate_multiplier: number
  models: PublicCatalogModel[]
}

export interface PublicModelCatalog {
  groups: PublicCatalogGroup[]
}

export async function getPublicModelCatalog(options?: { signal?: AbortSignal }): Promise<PublicModelCatalog> {
  const { data } = await apiClient.get<PublicModelCatalog>('/models/catalog', {
    signal: options?.signal,
  })
  return data
}

export const publicModelsAPI = {
  getPublicModelCatalog,
}

export default publicModelsAPI
