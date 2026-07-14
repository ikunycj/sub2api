/**
 * Application State Store
 * Manages global UI state including sidebar, loading indicators, and toast notifications
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Toast, ToastType, PublicSettings } from '@/types'
import { i18n } from '@/i18n'
import { getPublicSettings as fetchPublicSettingsAPI } from '@/api/auth'
import { activeBrand, withBrandIdentity } from '@/brand'

const PUBLIC_SETTINGS_CACHE_KEY = 'sub2api_public_settings_cache_v3'
const PUBLIC_SETTINGS_CACHE_TTL_MS = 5 * 60 * 1000
const PUBLIC_SETTINGS_STALE_TTL_MS = 60 * 60 * 1000

type PublicSettingsCacheRecord = {
  savedAt: number
  data: PublicSettings
}

function readPublicSettingsBrowserCache(maxAgeMs: number): PublicSettings | null {
  if (typeof window === 'undefined') return null

  try {
    const raw = window.localStorage.getItem(PUBLIC_SETTINGS_CACHE_KEY)
    if (!raw) return null

    const parsed = JSON.parse(raw) as Partial<PublicSettingsCacheRecord>
    if (!parsed.data || typeof parsed.savedAt !== 'number') return null
    if (Date.now() - parsed.savedAt > maxAgeMs) return null

    return parsed.data
  } catch {
    window.localStorage.removeItem(PUBLIC_SETTINGS_CACHE_KEY)
    return null
  }
}

function writePublicSettingsBrowserCache(data: PublicSettings): void {
  if (typeof window === 'undefined') return

  try {
    window.localStorage.setItem(PUBLIC_SETTINGS_CACHE_KEY, JSON.stringify({
      savedAt: Date.now(),
      data,
    } satisfies PublicSettingsCacheRecord))
  } catch {
    // Storage may be unavailable or full; in-memory cache still works.
  }
}

function clearPublicSettingsBrowserCache(): void {
  if (typeof window === 'undefined') return
  window.localStorage.removeItem(PUBLIC_SETTINGS_CACHE_KEY)
}

export const useAppStore = defineStore('app', () => {
  // ==================== State ====================

  const sidebarCollapsed = ref<boolean>(false)
  const mobileOpen = ref<boolean>(false)
  const loading = ref<boolean>(false)
  const toasts = ref<Toast[]>([])

  // Public settings cache state
  const publicSettingsLoaded = ref<boolean>(false)
  const publicSettingsLoading = ref<boolean>(false)
  const siteName = ref<string>(activeBrand.siteName)
  const siteLogo = ref<string>(activeBrand.logo)
  const siteVersion = ref<string>('')
  const contactInfo = ref<string>('')
  const apiBaseUrl = ref<string>('')
  const docUrl = ref<string>('')
  const cachedPublicSettings = ref<PublicSettings | null>(null)

  // Auto-incrementing ID for toasts
  let toastIdCounter = 0

  // ==================== Computed ====================

  const hasActiveToasts = computed(() => toasts.value.length > 0)
  const backendModeEnabled = computed(() => cachedPublicSettings.value?.backend_mode_enabled ?? false)

  const loadingCount = ref<number>(0)

  // ==================== Actions ====================

  /**
   * Toggle sidebar collapsed state
   */
  function toggleSidebar(): void {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  /**
   * Set sidebar collapsed state explicitly
   * @param collapsed - Whether sidebar should be collapsed
   */
  function setSidebarCollapsed(collapsed: boolean): void {
    sidebarCollapsed.value = collapsed
  }

  /**
   * Toggle mobile sidebar open state
   */
  function toggleMobileSidebar(): void {
    mobileOpen.value = !mobileOpen.value
  }

  /**
   * Set mobile sidebar open state explicitly
   * @param open - Whether mobile sidebar should be open
   */
  function setMobileOpen(open: boolean): void {
    mobileOpen.value = open
  }

  /**
   * Set global loading state
   * @param isLoading - Whether app is in loading state
   */
  function setLoading(isLoading: boolean): void {
    if (isLoading) {
      loadingCount.value++
    } else {
      loadingCount.value = Math.max(0, loadingCount.value - 1)
    }
    loading.value = loadingCount.value > 0
  }

  /**
   * Show a toast notification
   * @param type - Type of toast (success, error, info, warning)
   * @param message - Toast message content
   * @param duration - Auto-dismiss duration in ms (undefined = no auto-dismiss)
   * @returns Toast ID for manual dismissal
   */
  function showToast(type: ToastType, message: string, duration?: number): string {
    const id = `toast-${++toastIdCounter}`
    const toast: Toast = {
      id,
      type,
      message,
      duration,
      startTime: duration !== undefined ? Date.now() : undefined
    }

    toasts.value.push(toast)

    // Auto-dismiss if duration is specified
    if (duration !== undefined) {
      setTimeout(() => {
        hideToast(id)
      }, duration)
    }

    return id
  }

  /**
   * Show a success toast
   * @param message - Success message
   * @param duration - Auto-dismiss duration in ms (default: 3000)
   */
  function showSuccess(message: string, duration: number = 3000): string {
    return showToast('success', message, duration)
  }

  /**
   * Show an error toast
   * @param message - Error message
   * @param duration - Auto-dismiss duration in ms (default: 5000)
   */
  function showError(message: string, duration: number = 5000): string {
    return showToast('error', message, duration)
  }

  /**
   * Show an info toast
   * @param message - Info message
   * @param duration - Auto-dismiss duration in ms (default: 3000)
   */
  function showInfo(message: string, duration: number = 3000): string {
    return showToast('info', message, duration)
  }

  /**
   * Show a warning toast
   * @param message - Warning message
   * @param duration - Auto-dismiss duration in ms (default: 4000)
   */
  function showWarning(message: string, duration: number = 4000): string {
    return showToast('warning', message, duration)
  }

  /**
   * Hide a specific toast by ID
   * @param id - Toast ID to hide
   */
  function hideToast(id: string): void {
    const index = toasts.value.findIndex((t) => t.id === id)
    if (index !== -1) {
      toasts.value.splice(index, 1)
    }
  }

  /**
   * Clear all toasts
   */
  function clearAllToasts(): void {
    toasts.value = []
  }

  /**
   * Execute an async operation with loading state
   * Automatically manages loading indicator
   * @param operation - Async operation to execute
   * @returns Promise resolving to operation result
   */
  async function withLoading<T>(operation: () => Promise<T>): Promise<T> {
    setLoading(true)
    try {
      return await operation()
    } finally {
      setLoading(false)
    }
  }

  /**
   * Execute an async operation with loading and error handling
   * Shows error toast on failure
   * @param operation - Async operation to execute
   * @param errorMessage - Custom error message (optional)
   * @returns Promise resolving to operation result or null on error
   */
  async function withLoadingAndError<T>(
    operation: () => Promise<T>,
    errorMessage?: string
  ): Promise<T | null> {
    setLoading(true)
    try {
      return await operation()
    } catch (error) {
      const message =
        errorMessage ||
        (error as { message?: string }).message ||
        i18n.global.t('common.unknownError')
      showError(message)
      return null
    } finally {
      setLoading(false)
    }
  }

  /**
   * Reset app state to defaults
   * Useful for cleanup or testing
   */
  function reset(): void {
    sidebarCollapsed.value = false
    loading.value = false
    loadingCount.value = 0
    toasts.value = []
  }

  // ==================== Public Settings Management ====================

  /**
   * Apply settings to store state (internal helper to avoid code duplication)
   */
  function applySettings(config: PublicSettings): PublicSettings {
    const brandedConfig = withBrandIdentity(config, activeBrand)
    if (typeof window !== 'undefined') {
      window.__APP_CONFIG__ = { ...brandedConfig }
    }
    cachedPublicSettings.value = brandedConfig
    siteName.value = brandedConfig.site_name
    siteLogo.value = brandedConfig.site_logo
    siteVersion.value = brandedConfig.version || ''
    contactInfo.value = brandedConfig.contact_info || ''
    apiBaseUrl.value = brandedConfig.api_base_url || ''
    docUrl.value = brandedConfig.doc_url || ''
    publicSettingsLoaded.value = true
    return brandedConfig
  }

  /**
   * Fetch public settings (uses cache unless force=true)
   * @param force - Force refresh from API
   */
  async function fetchPublicSettings(force = false): Promise<PublicSettings | null> {
    // Check for injected config from server (eliminates flash)
    if (!publicSettingsLoaded.value && !force && window.__APP_CONFIG__) {
      const config = applySettings(window.__APP_CONFIG__)
      writePublicSettingsBrowserCache(config)
      return config
    }

    // Return cached data if available and not forcing refresh
    if (publicSettingsLoaded.value && !force) {
      if (cachedPublicSettings.value) {
        return { ...cachedPublicSettings.value }
      }
      return {
        registration_enabled: false,
        email_verify_enabled: false,
        force_email_on_third_party_signup: false,
        registration_email_suffix_whitelist: [],
        promo_code_enabled: true,
        password_reset_enabled: false,
        invitation_code_enabled: false,
        turnstile_enabled: false,
        turnstile_site_key: '',
        site_name: siteName.value,
        site_logo: siteLogo.value,
        site_subtitle: '',
        api_base_url: apiBaseUrl.value,
        contact_info: contactInfo.value,
        doc_url: docUrl.value,
        home_content: '',
        hide_ccs_import_button: false,
        payment_enabled: false,
        payment_display_mode: 'off',
        table_default_page_size: 20,
        table_page_size_options: [10, 20, 50, 100],
        custom_menu_items: [],
        custom_endpoints: [],
        linuxdo_oauth_enabled: false,
        wechat_oauth_enabled: false,
        wechat_oauth_open_enabled: false,
        wechat_oauth_mp_enabled: false,
        wechat_oauth_mobile_enabled: false,
        oidc_oauth_enabled: false,
        oidc_oauth_provider_name: 'OIDC',
        github_oauth_enabled: false,
        google_oauth_enabled: false,
        backend_mode_enabled: false,
        version: siteVersion.value,
        balance_low_notify_enabled: false,
        account_quota_notify_enabled: false,
        balance_low_notify_threshold: 0,
        channel_monitor_enabled: true,
        channel_monitor_default_interval_seconds: 60,
        available_channels_enabled: false,
        risk_control_enabled: false,
        service_quota_enabled: false,
        affiliate_enabled: false,
        allow_user_view_error_requests: false,
      }
    }

    if (!force) {
      const freshCachedSettings = readPublicSettingsBrowserCache(PUBLIC_SETTINGS_CACHE_TTL_MS)
      if (freshCachedSettings) {
        return applySettings(freshCachedSettings)
      }

      const staleCachedSettings = readPublicSettingsBrowserCache(PUBLIC_SETTINGS_STALE_TTL_MS)
      if (staleCachedSettings) {
        const config = applySettings(staleCachedSettings)
        if (!publicSettingsLoading.value) {
          void fetchPublicSettings(true)
        }
        return config
      }
    }

    // Prevent duplicate requests
    if (publicSettingsLoading.value) {
      return null
    }

    publicSettingsLoading.value = true
    try {
      const data = await fetchPublicSettingsAPI()
      const config = applySettings(data)
      writePublicSettingsBrowserCache(config)
      return config
    } catch (error) {
      console.error('Failed to fetch public settings:', error)
      return null
    } finally {
      publicSettingsLoading.value = false
    }
  }

  /**
   * Clear public settings cache
   */
  function clearPublicSettingsCache(): void {
    publicSettingsLoaded.value = false
    cachedPublicSettings.value = null
    clearPublicSettingsBrowserCache()
  }

  /**
   * Initialize settings from injected config (window.__APP_CONFIG__)
   * This is called synchronously before Vue app mounts to prevent flash
   * @returns true if config was found and applied, false otherwise
   */
  function initFromInjectedConfig(): boolean {
    if (window.__APP_CONFIG__) {
      const config = applySettings(window.__APP_CONFIG__)
      writePublicSettingsBrowserCache(config)
      return true
    }
    return false
  }

  // ==================== Return Store API ====================

  return {
    // State
    sidebarCollapsed,
    mobileOpen,
    loading,
    toasts,

    // Public settings state
    publicSettingsLoaded,
    siteName,
    siteLogo,
    siteVersion,
    contactInfo,
    apiBaseUrl,
    docUrl,
    cachedPublicSettings,

    // Computed
    hasActiveToasts,
    backendModeEnabled,

    // Actions
    toggleSidebar,
    setSidebarCollapsed,
    toggleMobileSidebar,
    setMobileOpen,
    setLoading,
    showToast,
    showSuccess,
    showError,
    showInfo,
    showWarning,
    hideToast,
    clearAllToasts,
    withLoading,
    withLoadingAndError,
    reset,

    // Public settings actions
    fetchPublicSettings,
    clearPublicSettingsCache,
    initFromInjectedConfig
  }
})
