<template>
  <div class="brand-console-shell min-h-screen bg-primary-50 dark:bg-dark-950">
    <!-- Background Decoration -->
    <div class="pointer-events-none fixed inset-0 overflow-hidden">
      <div class="absolute inset-0 bg-mesh-gradient"></div>
      <div class="brand-console-grid absolute inset-0"></div>
      <div class="brand-console-arc brand-console-arc-main"></div>
      <div class="brand-console-arc brand-console-arc-small"></div>
    </div>

    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main Content Area -->
    <div
      class="relative min-h-screen transition-all duration-300"
      :class="[sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64']"
    >
      <!-- Header -->
      <AppHeader />

      <!-- Main Content -->
      <main class="relative p-4 md:p-6 lg:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.user?.role === 'admin')

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
})

defineExpose({ replayTour })
</script>

<style scoped>
.brand-console-shell {
  background:
    radial-gradient(circle at 12% 0%, rgb(var(--color-primary-400) / 0.16), transparent 30rem),
    radial-gradient(circle at 92% 8%, rgb(6 182 212 / 0.1), transparent 30rem),
    linear-gradient(135deg, rgb(var(--color-primary-50)) 0%, #ffffff 52%, #f8fafc 100%);
}

:global(.dark) .brand-console-shell {
  background:
    radial-gradient(circle at 12% 0%, rgb(var(--color-primary-500) / 0.1), transparent 30rem),
    radial-gradient(circle at 92% 8%, rgb(13 148 136 / 0.08), transparent 28rem),
    linear-gradient(135deg, #020617 0%, #0f172a 58%, #111827 100%);
}

.brand-console-grid {
  background-image:
    linear-gradient(rgb(var(--color-primary-800) / 0.055) 1px, transparent 1px),
    linear-gradient(90deg, rgb(var(--color-primary-800) / 0.055) 1px, transparent 1px);
  background-size: 72px 72px;
  mask-image: linear-gradient(to bottom, black 0%, transparent 76%);
}

:global(.dark) .brand-console-grid {
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.035) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.035) 1px, transparent 1px);
}

.brand-console-arc {
  position: absolute;
  border: 1px solid rgb(var(--color-primary-400) / 0.28);
  border-bottom-color: transparent;
  border-left-color: transparent;
  border-radius: 9999px;
}

.brand-console-arc-main {
  right: -12rem;
  top: 5rem;
  width: 34rem;
  height: 34rem;
  transform: rotate(138deg);
}

.brand-console-arc-small {
  left: 10rem;
  bottom: -12rem;
  width: 28rem;
  height: 28rem;
  transform: rotate(28deg);
}

:global(.dark) .brand-console-arc {
  border-color: rgb(var(--color-primary-500) / 0.18);
  border-bottom-color: transparent;
  border-left-color: transparent;
}
</style>
