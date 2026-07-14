import { createHash } from 'node:crypto'
import { existsSync, readFileSync, readdirSync } from 'node:fs'
import { dirname, extname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const testDir = dirname(fileURLToPath(import.meta.url))
const srcDir = resolve(testDir, '../..')
const frontendDir = resolve(srcDir, '..')

function parseEnv(path: string): Record<string, string> {
  return Object.fromEntries(
    readFileSync(path, 'utf8')
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter(Boolean)
      .map((line) => {
        const separator = line.indexOf('=')
        return [line.slice(0, separator), line.slice(separator + 1)]
      }),
  )
}

function collectSourceFiles(directory: string): string[] {
  const files: string[] = []
  for (const entry of readdirSync(directory, { withFileTypes: true })) {
    const path = resolve(directory, entry.name)
    if (entry.isDirectory()) {
      if (entry.name !== '__tests__') files.push(...collectSourceFiles(path))
      continue
    }
    if (['.css', '.ts', '.vue'].includes(extname(entry.name))) files.push(path)
  }
  return files
}

describe('site style isolation', () => {
  it('allows exactly three site-specific environment variables', () => {
    const expectedKeys = ['VITE_PRIMARY_SCHEME', 'VITE_SITE_LOGO', 'VITE_SITE_NAME']
    const ikunEnv = parseEnv(resolve(frontendDir, '.env.ikun'))
    const anytokenEnv = parseEnv(resolve(frontendDir, '.env.anytoken'))

    expect(Object.keys(ikunEnv).sort()).toEqual(expectedKeys)
    expect(Object.keys(anytokenEnv).sort()).toEqual(expectedKeys)
    expect(ikunEnv).toEqual({
      VITE_SITE_NAME: 'ikun.love',
      VITE_PRIMARY_SCHEME: 'orange',
      VITE_SITE_LOGO: '/logo.png',
    })
    expect(anytokenEnv).toEqual({
      VITE_SITE_NAME: 'AnyToken',
      VITE_PRIMARY_SCHEME: 'blue',
      VITE_SITE_LOGO: '/brands/anytoken/logo.svg',
    })
  })

  it('has no page-level site switch or duplicate site content config', () => {
    const source = collectSourceFiles(srcDir)
      .map((path) => readFileSync(path, 'utf8'))
      .join('\n')

    expect(source).not.toContain('activeBrand.key')
    expect(source).not.toContain('VITE_BRAND')
    expect(existsSync(resolve(srcDir, 'brand/configs/ikun.ts'))).toBe(false)
    expect(existsSync(resolve(srcDir, 'brand/configs/anytoken.ts'))).toBe(false)
    expect(existsSync(resolve(srcDir, 'brand/configs/shared.ts'))).toBe(true)
  })

  it('keeps the fixed cyan secondary palette outside runtime site config', () => {
    const tailwindSource = readFileSync(resolve(frontendDir, 'tailwind.config.js'), 'utf8')
    const themeSource = readFileSync(resolve(srcDir, 'brand/theme.ts'), 'utf8')
    const viteSource = readFileSync(resolve(frontendDir, 'vite.config.ts'), 'utf8')

    expect(tailwindSource).toContain("500: '#06b6d4'")
    expect(tailwindSource).not.toContain('--color-secondary-')
    expect(themeSource).not.toContain('--color-secondary-')
    expect(viteSource).not.toContain('--color-secondary-')
  })

  it('reads first-paint identity from the same three Vite variables', () => {
    const viteSource = readFileSync(resolve(frontendDir, 'vite.config.ts'), 'utf8')
    const packageJson = JSON.parse(readFileSync(resolve(frontendDir, 'package.json'), 'utf8')) as {
      scripts: Record<string, string>
    }

    expect(viteSource).toContain('env.VITE_SITE_NAME')
    expect(viteSource).toContain('env.VITE_PRIMARY_SCHEME')
    expect(viteSource).toContain('env.VITE_SITE_LOGO')
    expect(viteSource).not.toContain('siteUrl:')
    expect(viteSource).not.toContain("siteName: 'AnyToken'")
    expect(existsSync(resolve(frontendDir, 'vite.config.js'))).toBe(false)
    for (const script of ['dev', 'dev:ikun', 'dev:anytoken', 'build', 'build:ikun', 'build:anytoken']) {
      expect(packageJson.scripts[script]).toContain('--config vite.config.ts')
    }
  })

  it('keeps the original custom-branch ikun logo and removes the duplicate fake logo', () => {
    const logo = readFileSync(resolve(frontendDir, 'public/logo.png'))
    expect(createHash('sha256').update(logo).digest('hex')).toBe(
      'e6e4ff64988dd2bbc661709d56e1a7a761c6771390b1737dae3df487e0dae023',
    )
    expect(existsSync(resolve(frontendDir, 'public/brands/ikun/logo.svg'))).toBe(false)
  })

  it('keeps the console brand link pointed at the shared home route', () => {
    const sidebarSource = readFileSync(
      resolve(srcDir, 'components/layout/AppSidebar.vue'),
      'utf8',
    )

    expect(sidebarSource).toMatch(/<router-link\s+to="\/home"\s+class="sidebar-header"/)
  })
})
