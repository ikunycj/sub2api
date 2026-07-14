export interface CatalogModelWithName {
  name: string
}

const DATE_SUFFIX_PATTERNS = [
  /-\d{4}-\d{2}-\d{2}$/,
  /-\d{8}$/,
]

function normalizeModelId(value: string): string {
  return value.trim().toLowerCase()
}

export function createCatalogModelKey(groupId: string | number, platform: string, modelId: string): string {
  return `${String(groupId).trim()}:${platform.trim().toLowerCase()}:${normalizeModelId(modelId)}`
}

function findAliasTarget(modelId: string, availableIds: ReadonlySet<string>): string | null {
  for (const pattern of DATE_SUFFIX_PATTERNS) {
    const candidate = modelId.replace(pattern, '')
    if (candidate !== modelId && availableIds.has(candidate)) {
      return candidate
    }
  }

  if (modelId.endsWith('-chat-latest')) {
    const candidate = modelId.slice(0, -'-chat-latest'.length)
    if (availableIds.has(candidate)) {
      return candidate
    }
  }

  return null
}

export function selectCanonicalCatalogModels<T extends CatalogModelWithName>(models: readonly T[]): T[] {
  const uniqueModels = new Map<string, T>()
  for (const model of models) {
    const modelId = normalizeModelId(model.name)
    if (modelId && !uniqueModels.has(modelId)) {
      uniqueModels.set(modelId, model)
    }
  }

  const availableIds = new Set(uniqueModels.keys())
  return Array.from(uniqueModels, ([modelId, model]) => ({ modelId, model }))
    .filter(({ modelId }) => findAliasTarget(modelId, availableIds) == null)
    .map(({ model }) => model)
}
