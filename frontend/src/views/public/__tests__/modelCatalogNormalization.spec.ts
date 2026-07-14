import { describe, expect, it } from 'vitest'
import { createCatalogModelKey, selectCanonicalCatalogModels } from '../modelCatalogNormalization'

describe('public model catalog normalization', () => {
  it('collapses dated and chat-latest aliases when canonical models exist', () => {
    const models = [
      { name: 'gpt-5.2' },
      { name: 'gpt-5.2-2025-12-11' },
      { name: 'gpt-5.2-chat-latest' },
      { name: 'gpt-5.2-pro' },
      { name: 'gpt-5.2-pro-2025-12-11' },
      { name: 'gpt-5.3-codex' },
    ]

    expect(selectCanonicalCatalogModels(models).map((model) => model.name)).toEqual([
      'gpt-5.2',
      'gpt-5.2-pro',
      'gpt-5.3-codex',
    ])
  })

  it('keeps dated models when no canonical alias is configured', () => {
    const models = [
      { name: 'claude-3-5-haiku-20241022' },
      { name: 'claude-3-7-sonnet-20250219' },
    ]

    expect(selectCanonicalCatalogModels(models)).toEqual(models)
  })

  it('deduplicates exact model ids case-insensitively', () => {
    const models = [
      { name: ' gpt-5.4 ' },
      { name: 'GPT-5.4' },
    ]

    expect(selectCanonicalCatalogModels(models)).toEqual([{ name: ' gpt-5.4 ' }])
  })

  it('keeps identical model ids in different groups as distinct Vue rows', () => {
    expect(createCatalogModelKey(22, 'openai', 'gpt-5.2')).toBe('22:openai:gpt-5.2')
    expect(createCatalogModelKey(26, 'openai', 'gpt-5.2')).toBe('26:openai:gpt-5.2')
    expect(createCatalogModelKey(22, 'openai', 'gpt-5.2')).not.toBe(createCatalogModelKey(26, 'openai', 'gpt-5.2'))
  })
})
