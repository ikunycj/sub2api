export function isDiscountedPrice(current: number | null | undefined, official: number | null | undefined): boolean {
  if (current == null || official == null || !Number.isFinite(current) || !Number.isFinite(official) || official <= 0) {
    return false
  }
  return current < official * (1 - 1e-9)
}

export function hasOfficialReferencePrice(official: number | null | undefined): boolean {
  return official != null && Number.isFinite(official) && official >= 0
}

export function calculateSavingsPercent(current: number | null | undefined, official: number | null | undefined): number | null {
  if (!isDiscountedPrice(current, official) || current == null || official == null) {
    return null
  }
  return Math.max(1, Math.min(99, Math.round((1 - current / official) * 100)))
}
