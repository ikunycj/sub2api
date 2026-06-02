<template>
  <AppLayout>
    <div class="ai-studio-page mx-auto w-full max-w-[1540px]">
      <section class="ai-studio-shell overflow-hidden rounded-3xl border border-amber-100/80 bg-white/86 shadow-sm shadow-primary-500/10 dark:border-dark-700/70 dark:bg-dark-900/82">
        <aside class="ai-studio-side flex flex-col border-b border-amber-100/80 bg-white/76 dark:border-dark-700/70 dark:bg-dark-900/72 lg:border-b-0 lg:border-r">
          <div class="p-4">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <p class="inline-flex items-center gap-1.5 rounded-full border border-primary-200/80 bg-primary-50 px-2.5 py-1 text-xs font-semibold text-primary-700 dark:border-primary-800/60 dark:bg-primary-950/30 dark:text-primary-200">
                  <Icon name="sparkles" size="xs" />
                  {{ t('aiStudio.entry.badge') }}
                </p>
                <h2 class="mt-3 truncate text-xl font-bold tracking-normal text-gray-950 dark:text-white">
                  {{ t('aiStudio.title') }}
                </h2>
                <p class="mt-1 line-clamp-2 text-xs leading-5 text-gray-500 dark:text-dark-400">
                  {{ t('aiStudio.description') }}
                </p>
              </div>

              <button type="button" class="btn btn-secondary btn-icon lg:hidden" @click="showSettingsModal = true">
                <Icon name="cog" size="sm" />
              </button>
            </div>

            <button type="button" class="btn btn-primary mt-4 w-full justify-center" @click="startNewConversation">
              <Icon name="plus" size="sm" />
              {{ t('aiStudio.side.newChat') }}
            </button>
          </div>

          <div class="px-4">
            <div class="grid grid-cols-2 rounded-2xl border border-amber-100/90 bg-primary-50/45 p-1 dark:border-dark-700 dark:bg-dark-800/65">
              <button
                v-for="mode in modeOptions"
                :key="mode.mode"
                type="button"
                class="inline-flex items-center justify-center gap-2 rounded-xl px-3 py-2 text-xs font-semibold transition-all active:scale-[0.98]"
                :class="activeMode === mode.mode
                  ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-950 dark:text-primary-200'
                  : 'text-gray-500 hover:text-primary-700 dark:text-dark-300 dark:hover:text-primary-200'"
                @click="switchMode(mode.mode)"
              >
                <Icon :name="mode.icon" size="xs" />
                {{ mode.label }}
              </button>
            </div>
          </div>

          <div class="mt-4 flex-1 overflow-y-auto px-3 pb-4">
            <label class="flex items-center gap-2 rounded-2xl border border-amber-100/80 bg-white px-3 py-2 text-gray-400 shadow-sm shadow-primary-500/5 dark:border-dark-700 dark:bg-dark-950/45">
              <Icon name="search" size="xs" />
              <input
                v-model="historyQuery"
                class="min-w-0 flex-1 border-0 bg-transparent text-xs text-gray-700 outline-none placeholder:text-gray-400 focus:ring-0 dark:text-dark-100"
                :placeholder="t('aiStudio.side.searchPlaceholder')"
              />
            </label>

            <div class="mt-4 flex items-center justify-between px-1">
              <div>
                <p class="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-dark-500">
                  {{ t('aiStudio.side.history') }}
                </p>
                <p class="mt-1 text-[11px] leading-4 text-gray-400 dark:text-dark-500">
                  {{ t('aiStudio.side.serverRetention') }}
                </p>
              </div>
              <span class="text-xs text-gray-400 dark:text-dark-500">{{ filteredConversations.length }}</span>
            </div>

            <div class="mt-2 space-y-1.5">
              <div
                v-for="conversation in filteredConversations"
                :key="conversation.id"
                class="group rounded-2xl px-3 py-2.5 transition-all"
                :class="activeHistoryId === conversation.id
                  ? 'bg-primary-50 text-primary-700 shadow-sm dark:bg-primary-950/35 dark:text-primary-200'
                  : 'text-gray-600 hover:bg-white hover:text-gray-950 hover:shadow-sm dark:text-dark-300 dark:hover:bg-dark-950/45 dark:hover:text-white'"
              >
                <div class="flex items-start gap-2.5">
                  <button
                    type="button"
                    class="mt-0.5 flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-xl"
                    :class="conversation.mode === 'image'
                      ? 'bg-amber-100 text-primary-600 dark:bg-primary-950/45 dark:text-primary-200'
                      : 'bg-white text-gray-500 ring-1 ring-amber-100/80 dark:bg-dark-900 dark:text-dark-300 dark:ring-dark-700'"
                    @click="openHistory(conversation)"
                  >
                    <Icon :name="conversation.mode === 'image' ? 'sparkles' : 'chat'" size="xs" />
                  </button>
                  <button type="button" class="min-w-0 flex-1 text-left" @click="openHistory(conversation)">
                    <span class="flex min-w-0 items-center gap-1.5">
                      <Icon v-if="conversation.pinned" name="badge" size="xs" class="flex-shrink-0 text-primary-500" />
                      <span class="block truncate text-sm font-semibold">{{ conversation.title }}</span>
                    </span>
                    <span class="mt-0.5 flex flex-wrap items-center gap-1.5 text-xs text-gray-400 dark:text-dark-500">
                      <Icon name="clock" size="xs" />
                      {{ conversation.time }}
                      <span v-if="conversation.source === 'local'" class="rounded-full bg-primary-50 px-1.5 py-0.5 text-[10px] font-semibold text-primary-600 dark:bg-primary-950/45 dark:text-primary-200">
                        {{ t('aiStudio.side.localBadge') }}
                      </span>
                      <span v-if="conversation.autoSaved" class="rounded-full bg-emerald-50 px-1.5 py-0.5 text-[10px] font-semibold text-emerald-600 dark:bg-emerald-950/30 dark:text-emerald-300">
                        {{ t('aiStudio.side.autoBadge') }}
                      </span>
                    </span>
                  </button>
                  <div class="flex flex-shrink-0 items-center gap-0.5 opacity-100 md:opacity-0 md:transition-opacity md:group-hover:opacity-100">
                    <button
                      type="button"
                      class="rounded-lg p-1.5 text-gray-400 hover:bg-primary-50 hover:text-primary-600 dark:hover:bg-primary-950/35 dark:hover:text-primary-200"
                      :title="conversation.pinned ? t('aiStudio.actions.unpin') : t('aiStudio.actions.pin')"
                      @click="toggleConversationPin(conversation)"
                    >
                      <Icon name="badge" size="xs" />
                    </button>
                    <button
                      type="button"
                      class="rounded-lg p-1.5 text-gray-400 hover:bg-primary-50 hover:text-primary-600 dark:hover:bg-primary-950/35 dark:hover:text-primary-200"
                      :title="t('aiStudio.actions.rename')"
                      @click="renameConversation(conversation)"
                    >
                      <Icon name="edit" size="xs" />
                    </button>
                    <button
                      type="button"
                      class="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950/35 dark:hover:text-red-300"
                      :title="t('aiStudio.actions.delete')"
                      @click="deleteConversation(conversation)"
                    >
                      <Icon name="trash" size="xs" />
                    </button>
                  </div>
                </div>
              </div>

              <div v-if="filteredConversations.length === 0" class="rounded-2xl border border-dashed border-amber-200/80 px-3 py-6 text-center text-xs text-gray-400 dark:border-dark-700 dark:text-dark-500">
                {{ t('aiStudio.side.noHistory') }}
              </div>
            </div>
          </div>

          <div class="border-t border-amber-100/80 p-4 dark:border-dark-700/70">
            <button
              type="button"
              class="flex w-full items-center gap-3 rounded-2xl bg-white px-3 py-3 text-left shadow-sm shadow-primary-500/5 ring-1 ring-amber-100/80 transition-colors hover:text-primary-700 dark:bg-dark-950/45 dark:ring-dark-700 dark:hover:text-primary-200"
              @click="showSettingsModal = true"
            >
              <span class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-primary-950/45 dark:text-primary-200">
                <Icon name="cog" size="sm" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="block text-sm font-semibold text-gray-900 dark:text-white">{{ t('aiStudio.settings.title') }}</span>
                <span class="block truncate text-xs text-gray-400 dark:text-dark-500">{{ settingsSummary }}</span>
              </span>
              <Icon name="chevronRight" size="xs" />
            </button>
          </div>
        </aside>

        <main class="ai-studio-main flex min-h-0 min-w-0 flex-col overflow-hidden">
          <header class="flex min-h-[68px] flex-shrink-0 items-center justify-between gap-3 border-b border-amber-100/80 bg-white/72 px-4 py-3 backdrop-blur dark:border-dark-700/70 dark:bg-dark-900/65 md:px-5">
            <div class="flex min-w-0 items-center gap-3">
              <span class="hidden h-9 w-9 flex-shrink-0 items-center justify-center rounded-2xl bg-primary-50 text-primary-600 dark:bg-primary-950/40 dark:text-primary-200 sm:flex">
                <Icon :name="activeMode === 'image' ? 'sparkles' : 'chat'" size="sm" />
              </span>

              <div class="min-w-0">
                <h3 class="truncate text-sm font-semibold text-gray-950 dark:text-white md:text-base">
                  {{ currentTitle }}
                </h3>
                <p class="truncate text-xs text-gray-500 dark:text-dark-400">
                  {{ currentSubtitle }}
                </p>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <span class="hidden rounded-full border border-amber-100/80 bg-white px-3 py-1.5 text-xs font-semibold text-gray-700 shadow-sm dark:border-dark-700 dark:bg-dark-950/45 dark:text-dark-200 md:inline-flex">
                {{ currentModel }}
              </span>
              <button type="button" class="btn btn-secondary btn-icon" @click="showSettingsModal = true">
                <Icon name="cog" size="sm" />
              </button>
            </div>
          </header>

          <div class="ai-studio-messages min-h-0 flex-1 overflow-y-auto px-4 py-5 dark:bg-dark-900/30 md:px-6">
            <div class="mx-auto max-w-[860px] space-y-5">
              <div v-if="messages.length === 0" class="mx-auto mb-8 max-w-xl pt-14 text-center">
                <div class="mx-auto flex h-14 w-14 items-center justify-center rounded-3xl bg-primary-100 text-primary-600 shadow-sm dark:bg-primary-950/50 dark:text-primary-300">
                  <Icon :name="activeMode === 'image' ? 'sparkles' : 'chat'" size="lg" />
                </div>
                <h2 class="mt-4 text-2xl font-semibold tracking-normal text-gray-950 dark:text-white">
                  {{ activeMode === 'image' ? t('aiStudio.conversation.imageGreeting') : t('aiStudio.conversation.chatGreeting') }}
                </h2>
                <p class="mt-2 text-sm leading-6 text-gray-500 dark:text-dark-400">
                  {{ activeMode === 'image' ? t('aiStudio.conversation.imageHint') : t('aiStudio.conversation.chatHint') }}
                </p>
              </div>

              <div
                v-for="message in messages"
                :key="message.id"
                class="flex gap-3"
                :class="message.role === 'user' ? 'justify-end' : 'justify-start'"
              >
                <div
                  v-if="message.role === 'assistant'"
                  class="mt-1 flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-2xl bg-white text-primary-600 shadow-sm ring-1 ring-amber-100 dark:bg-dark-950 dark:text-primary-300 dark:ring-dark-700"
                >
                  <Icon :name="message.kind === 'image' ? 'sparkles' : 'chat'" size="sm" />
                </div>

                <article
                  class="max-w-[86%] rounded-3xl px-4 py-3 shadow-sm"
                  :class="message.role === 'user'
                    ? 'rounded-br-md bg-gradient-to-r from-primary-500 to-primary-600 text-white shadow-primary-500/20'
                    : 'rounded-bl-md border border-amber-100/80 bg-white text-gray-700 dark:border-dark-700 dark:bg-dark-950/60 dark:text-dark-200'"
                >
                  <template v-if="editingMessageId === message.id">
                    <textarea
                      v-model="editingMessageContent"
                      rows="4"
                      class="min-h-[104px] w-full resize-none rounded-2xl border border-white/40 bg-white/95 px-3 py-2 text-sm leading-6 text-gray-900 outline-none transition focus:border-white focus:ring-2 focus:ring-white/35 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-100"
                      :placeholder="t('aiStudio.chat.editPlaceholder')"
                      @keydown.enter.exact.prevent="submitEditedMessage(message)"
                    ></textarea>
                    <div class="mt-3 flex flex-wrap justify-end gap-2">
                      <button
                        type="button"
                        class="rounded-full px-3 py-1.5 text-xs font-semibold transition-colors"
                        :class="message.role === 'user'
                          ? 'bg-white/10 text-white hover:bg-white/20'
                          : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-800 dark:text-dark-200'"
                        @click="cancelEditingMessage"
                      >
                        {{ t('common.cancel') }}
                      </button>
                      <button
                        type="button"
                        class="inline-flex items-center gap-1 rounded-full bg-white px-3 py-1.5 text-xs font-semibold text-primary-700 transition-colors hover:bg-primary-50 disabled:cursor-not-allowed disabled:opacity-60"
                        :disabled="!canSubmitEditedMessage(message)"
                        @click="submitEditedMessage(message)"
                      >
                        <Icon name="refresh" size="xs" />
                        {{ t('aiStudio.chat.resend') }}
                      </button>
                    </div>
                  </template>
                  <div v-else>
                    <div
                      v-if="message.role === 'assistant' && !message.content.trim() && !message.reasoning?.trim() && message.isStreaming"
                      class="inline-flex items-center gap-2 text-sm text-gray-500 dark:text-dark-300"
                    >
                      <Icon name="refresh" size="sm" class="animate-spin text-primary-500 dark:text-primary-300" />
                      <span>{{ t('aiStudio.chat.waiting') }}</span>
                    </div>

                    <details
                      v-if="message.role === 'assistant' && message.reasoning?.trim()"
                      class="mb-3 rounded-2xl border border-primary-100 bg-primary-50/65 px-3 py-2 text-xs text-gray-600 dark:border-primary-900/50 dark:bg-primary-950/25 dark:text-dark-200"
                      open
                    >
                      <summary class="flex cursor-pointer list-none items-center gap-2 font-semibold text-primary-700 dark:text-primary-200">
                        <Icon name="brain" size="xs" />
                        <span>{{ message.isStreaming && !message.content.trim() ? t('aiStudio.chat.thinking') : t('aiStudio.chat.thought') }}</span>
                        <Icon v-if="message.isStreaming && !message.content.trim()" name="refresh" size="xs" class="ml-auto animate-spin opacity-70" />
                      </summary>
                      <p class="mt-2 whitespace-pre-line leading-5">{{ message.reasoning }}</p>
                    </details>

                    <p v-if="message.content.trim()" class="whitespace-pre-line text-sm leading-6">{{ message.content }}</p>
                  </div>

                  <div v-if="message.attachments?.length" class="mt-3 space-y-2">
                    <div
                      v-for="attachment in message.attachments"
                      :key="attachment.id"
                      class="flex items-center gap-2 rounded-2xl px-3 py-2 text-xs"
                      :class="message.role === 'user'
                        ? 'bg-white/15 text-white'
                        : 'bg-primary-50 text-gray-600 dark:bg-dark-800 dark:text-dark-300'"
                    >
                      <Icon name="document" size="xs" />
                      <span class="min-w-0 flex-1 truncate">{{ attachment.name }}</span>
                      <span class="flex-shrink-0 opacity-75">{{ attachment.sizeLabel }}</span>
                    </div>
                  </div>

                  <div v-if="message.kind === 'image' && message.role === 'assistant' && message.images?.length" class="mt-3 grid grid-cols-2 gap-2">
                    <a
                      v-for="item in message.images"
                      :key="item.id"
                      :href="item.url"
                      target="_blank"
                      rel="noreferrer"
                      class="relative aspect-square overflow-hidden rounded-2xl border border-white/70 bg-primary-50 shadow-sm dark:border-dark-700 dark:bg-dark-900"
                    >
                      <img :src="item.url" :alt="item.name" class="h-full w-full object-cover" />
                    </a>
                  </div>

                  <div
                    class="mt-3 flex items-center gap-2 text-xs"
                    :class="message.role === 'user' ? 'text-primary-100' : 'text-gray-400 dark:text-dark-400'"
                  >
                    <span>{{ message.time }}</span>
                    <span class="inline-flex items-center gap-1">
                      <Icon :name="message.kind === 'image' ? 'sparkles' : 'chat'" size="xs" />
                      {{ message.kind === 'image' ? t('aiStudio.image.title') : t('aiStudio.chat.title') }}
                    </span>
                    <button
                      v-if="message.role === 'assistant'"
                      type="button"
                      class="inline-flex items-center gap-1 rounded-lg px-2 py-1 transition-colors hover:bg-primary-50 hover:text-primary-600 dark:hover:bg-primary-950/30 dark:hover:text-primary-300"
                      @click="copyToClipboard(message.content)"
                    >
                      <Icon name="copy" size="xs" />
                      {{ t('aiStudio.chat.copy') }}
                    </button>
                    <button
                      v-if="message.role === 'user' && editingMessageId !== message.id"
                      type="button"
                      class="inline-flex items-center gap-1 rounded-lg px-2 py-1 transition-colors hover:bg-white/15 hover:text-white"
                      @click="startEditingMessage(message)"
                    >
                      <Icon name="edit" size="xs" />
                      {{ t('aiStudio.chat.edit') }}
                    </button>
                  </div>
                </article>
              </div>
            </div>
          </div>

          <form class="flex-shrink-0 border-t border-amber-100/80 bg-white/82 p-3 backdrop-blur dark:border-dark-700 dark:bg-dark-900/85 md:p-4" @submit.prevent="handleSubmit">
            <input ref="fileInputRef" type="file" multiple class="hidden" @change="handleFilesSelected" />
            <div class="mx-auto max-w-[860px] rounded-3xl border border-amber-100/90 bg-white p-2 shadow-sm shadow-primary-500/10 dark:border-dark-700 dark:bg-dark-950/70">
              <div class="flex flex-wrap items-center gap-2 border-b border-amber-100/70 px-2 pb-2 dark:border-dark-700">
                <button
                  v-for="mode in modeOptions"
                  :key="mode.mode"
                  type="button"
                  class="inline-flex items-center gap-2 rounded-full px-3 py-1.5 text-xs font-semibold transition-all active:scale-[0.98]"
                  :class="activeMode === mode.mode
                    ? 'bg-primary-100 text-primary-700 dark:bg-primary-950/55 dark:text-primary-200'
                    : 'text-gray-500 hover:bg-primary-50 hover:text-primary-700 dark:text-dark-300 dark:hover:bg-primary-950/30 dark:hover:text-primary-200'"
                  @click="switchMode(mode.mode)"
                >
                  <Icon :name="mode.icon" size="xs" />
                  {{ mode.label }}
                </button>

                <div class="ml-auto flex items-center gap-2">
                  <select v-model="currentModel" class="max-w-[180px] rounded-full border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
                    <option v-if="availableModels.length === 0" value="" disabled>{{ t('aiStudio.settings.modelUnavailable') }}</option>
                    <option v-for="model in availableModels" :key="model" :value="model">{{ model }}</option>
                  </select>
                </div>
              </div>

              <div v-if="pendingAttachments.length > 0" class="flex flex-wrap gap-2 px-2 pt-2">
                <span
                  v-for="attachment in pendingAttachments"
                  :key="attachment.id"
                  class="inline-flex max-w-full items-center gap-2 rounded-full bg-primary-50 px-3 py-1.5 text-xs font-medium text-primary-700 dark:bg-primary-950/35 dark:text-primary-200"
                >
                  <Icon name="document" size="xs" />
                  <span class="max-w-[180px] truncate">{{ attachment.name }}</span>
                  <span class="text-primary-500/70 dark:text-primary-200/70">{{ attachment.sizeLabel }}</span>
                  <button type="button" class="rounded-full p-0.5 hover:bg-primary-100 dark:hover:bg-primary-900" @click="removeAttachment(attachment.id)">
                    <Icon name="x" size="xs" />
                  </button>
                </span>
              </div>

              <div class="flex items-end gap-2 pt-2">
                <button
                  type="button"
                  class="mb-1 flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-amber-100 bg-white text-gray-500 transition-colors hover:border-primary-200 hover:bg-primary-50 hover:text-primary-600 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300 dark:hover:text-primary-200"
                  :title="t('aiStudio.input.upload')"
                  @click="openFilePicker"
                >
                  <Icon name="upload" size="sm" />
                </button>
                <textarea
                  v-model="prompt"
                  rows="2"
                  class="min-h-[56px] flex-1 resize-none border-0 bg-transparent px-3 py-2 text-sm leading-6 text-gray-900 outline-none placeholder:text-gray-400 focus:ring-0 dark:text-gray-100 dark:placeholder:text-dark-400"
                  :placeholder="activeMode === 'image' ? t('aiStudio.input.imagePlaceholder') : t('aiStudio.input.chatPlaceholder')"
                  @keydown.enter.exact.prevent="handleSubmit"
                ></textarea>
                <button type="submit" class="btn btn-primary btn-icon mb-1" :disabled="!canSubmit">
                  <Icon name="arrowUp" size="sm" />
                </button>
              </div>
            </div>
            <p class="mt-2 text-center text-xs text-gray-400 dark:text-dark-500">{{ t('aiStudio.conversation.footerHint') }}</p>
          </form>
        </main>
      </section>
    </div>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showSettingsModal"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-sm"
          @click.self="showSettingsModal = false"
        >
          <section class="w-full max-w-2xl overflow-hidden rounded-3xl border border-amber-100 bg-white shadow-2xl shadow-black/15 dark:border-dark-700 dark:bg-dark-900">
            <header class="flex items-center justify-between gap-3 border-b border-amber-100/80 px-5 py-4 dark:border-dark-700">
              <div>
                <h3 class="text-base font-semibold text-gray-950 dark:text-white">{{ t('aiStudio.settings.title') }}</h3>
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.description') }}</p>
              </div>
              <button type="button" class="btn btn-secondary btn-icon" @click="showSettingsModal = false">
                <Icon name="x" size="sm" />
              </button>
            </header>

            <div class="max-h-[70dvh] space-y-5 overflow-y-auto px-5 py-5">
              <div class="rounded-2xl border border-amber-100/80 bg-primary-50/35 p-4 dark:border-dark-700 dark:bg-dark-950/35">
                <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <h4 class="text-sm font-semibold text-gray-950 dark:text-white">{{ t('aiStudio.settings.historyStorageTitle') }}</h4>
                    <p class="mt-1 text-xs leading-5 text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.historyStorageHint') }}</p>
                  </div>
                  <button type="button" class="btn btn-secondary flex-shrink-0" @click="saveCurrentSessionToLocalHistory">
                    <Icon name="download" size="sm" />
                    {{ t('aiStudio.settings.saveLocalHistory') }}
                  </button>
                </div>
                <label class="mt-4 flex cursor-pointer items-start gap-3 rounded-2xl border border-amber-100/80 bg-white/70 p-3 dark:border-dark-700 dark:bg-dark-900/55">
                  <input v-model="autoSaveLocalHistory" type="checkbox" class="mt-1 h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                  <span>
                    <span class="block text-sm font-semibold text-gray-900 dark:text-white">{{ t('aiStudio.settings.autoSaveLocalHistory') }}</span>
                    <span class="mt-0.5 block text-xs leading-5 text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.autoSaveLocalHistoryHint') }}</span>
                  </span>
                </label>
                <label class="mt-3 block rounded-2xl border border-amber-100/80 bg-white/70 p-3 dark:border-dark-700 dark:bg-dark-900/55">
                  <span class="block text-sm font-semibold text-gray-900 dark:text-white">{{ t('aiStudio.settings.localHistoryRetention') }}</span>
                  <span class="mt-0.5 block text-xs leading-5 text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.localHistoryRetentionHint') }}</span>
                  <div class="mt-3 flex items-center gap-2">
                    <input
                      v-model.number="localHistoryRetentionDays"
                      type="number"
                      min="1"
                      step="1"
                      class="input max-w-[140px]"
                      :placeholder="t('aiStudio.settings.noRetentionLimit')"
                      @blur="normalizeLocalHistoryRetention"
                    />
                    <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.daysAfter') }}</span>
                  </div>
                </label>
                <p class="mt-3 text-xs text-gray-400 dark:text-dark-500">
                  {{ t('aiStudio.settings.localHistoryCount', { count: localHistoryCount }) }}
                </p>
              </div>

              <label class="space-y-2">
                <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.apiKey') }}</span>
                <select v-model.number="selectedAPIKeyId" class="input" :disabled="keyOptionsLoading">
                  <option v-if="keyOptions.length === 0" :value="null" disabled>{{ t('aiStudio.settings.apiKeyPlaceholder') }}</option>
                  <option v-else-if="availableKeyOptions.length === 0" :value="null" disabled>{{ t('aiStudio.settings.noUsableAPIKey') }}</option>
                  <option v-for="key in keyOptions" :key="key.id" :value="key.id" :disabled="!key.available">
                    {{ formatAPIKeyOptionLabel(key) }}
                  </option>
                </select>
                <p v-if="selectedKeyOption && !selectedKeyOption.available" class="text-xs leading-5 text-amber-600 dark:text-amber-300">
                  {{ t('aiStudio.settings.selectedKeyUnavailable', { reason: apiKeyUnavailableReasonLabel(selectedKeyOption.unavailable_reason) }) }}
                </p>
                <p v-else-if="keyOptions.length > 0 && availableKeyOptions.length === 0" class="text-xs leading-5 text-amber-600 dark:text-amber-300">
                  {{ t('aiStudio.settings.noUsableAPIKeyHint', { reason: firstUnavailableReasonLabel() }) }}
                </p>
                <button
                  v-if="keyOptions.length === 0 || availableKeyOptions.length === 0"
                  type="button"
                  class="mt-2 inline-flex items-center gap-1 text-xs font-semibold text-primary-600 hover:text-primary-700 dark:text-primary-300"
                  @click="handleAPIKeyHelpClick"
                >
                  <Icon name="plus" size="xs" />
                  {{ t('aiStudio.settings.goManageAPIKey') }}
                </button>
              </label>

              <div class="grid gap-4 md:grid-cols-2">
                <label class="space-y-2">
                  <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.chatModel') }}</span>
                  <select v-model="chatModel" class="input">
                    <option v-if="chatModels.length === 0" value="" disabled>{{ t('aiStudio.settings.modelUnavailable') }}</option>
                    <option v-for="model in chatModels" :key="model" :value="model">{{ model }}</option>
                  </select>
                </label>

                <label class="space-y-2">
                  <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.imageModel') }}</span>
                  <select v-model="imageModel" class="input">
                    <option v-if="imageModels.length === 0" value="" disabled>{{ t('aiStudio.settings.modelUnavailable') }}</option>
                    <option v-for="model in imageModels" :key="model" :value="model">{{ model }}</option>
                  </select>
                </label>
              </div>

              <div class="grid gap-4 md:grid-cols-3">
                <label class="space-y-2">
                  <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.maxTokens') }}</span>
                  <input v-model.number="maxTokens" type="number" min="256" step="256" class="input" />
                </label>

                <label class="space-y-2">
                  <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.imageSize') }}</span>
                  <select v-model="imageSize" class="input">
                    <option value="1024x1024">1024x1024</option>
                    <option value="1024x1536">1024x1536</option>
                    <option value="1536x1024">1536x1024</option>
                  </select>
                </label>

                <label class="space-y-2">
                  <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.quality') }}</span>
                  <select v-model="imageQuality" class="input">
                    <option value="standard">{{ t('aiStudio.settings.qualityStandard') }}</option>
                    <option value="hd">{{ t('aiStudio.settings.qualityHd') }}</option>
                  </select>
                </label>
              </div>

              <div class="rounded-2xl border border-amber-100/80 bg-white/70 p-3 dark:border-dark-700 dark:bg-dark-900/55">
                <label class="flex cursor-pointer items-start gap-3">
                  <input v-model="thinkingEnabled" type="checkbox" class="mt-1 h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                  <span>
                    <span class="block text-sm font-semibold text-gray-900 dark:text-white">{{ t('aiStudio.settings.thinkingMode') }}</span>
                    <span class="mt-0.5 block text-xs leading-5 text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.thinkingModeHint') }}</span>
                  </span>
                </label>
                <label v-if="thinkingEnabled" class="mt-3 block space-y-2">
                  <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.reasoningEffort') }}</span>
                  <select v-model="reasoningEffort" class="input">
                    <option value="low">{{ t('aiStudio.settings.reasoningEffortLow') }}</option>
                    <option value="medium">{{ t('aiStudio.settings.reasoningEffortMedium') }}</option>
                    <option value="high">{{ t('aiStudio.settings.reasoningEffortHigh') }}</option>
                    <option value="xhigh">{{ t('aiStudio.settings.reasoningEffortXHigh') }}</option>
                  </select>
                </label>
              </div>

              <label class="space-y-2">
                <span class="text-xs font-semibold text-gray-500 dark:text-dark-400">{{ t('aiStudio.settings.systemPrompt') }}</span>
                <textarea
                  v-model="systemPrompt"
                  rows="4"
                  class="input min-h-[104px] resize-none"
                  :placeholder="t('aiStudio.settings.systemPromptPlaceholder')"
                ></textarea>
              </label>
            </div>

            <footer class="flex justify-end gap-2 border-t border-amber-100/80 px-5 py-4 dark:border-dark-700">
              <button type="button" class="btn btn-secondary" @click="showSettingsModal = false">
                {{ t('common.cancel') }}
              </button>
              <button type="button" class="btn btn-primary" @click="saveSettings">
                {{ t('common.save') }}
              </button>
            </footer>
          </section>
        </div>
      </Transition>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores/app'
import { useClipboard } from '@/composables/useClipboard'
import { aiStudioApi, type AIStudioAttachment, type AIStudioConversation, type AIStudioKeyOption, type AIStudioMessage, type AIStudioRunPayload } from '@/api/aiStudio'

type StudioMode = 'chat' | 'image'
type ConversationSource = 'server' | 'local'

type ChatAttachment = {
  id: string
  serverId?: number
  name: string
  size: number
  sizeLabel: string
  type: string
  url?: string
  file?: File
}

type ChatImage = {
  id: string
  name: string
  url: string
}

type ChatMessage = {
  id: number | string
  serverId?: number
  role: 'assistant' | 'user' | 'system'
  kind: StudioMode
  content: string
  reasoning?: string
  isStreaming?: boolean
  time: string
  attachments?: ChatAttachment[]
  images?: ChatImage[]
}

type ConversationItem = {
  id: string
  mode: StudioMode
  title: string
  time: string
  source: ConversationSource
  pinned?: boolean
  autoSaved?: boolean
}

type StoredSession = {
  id: string
  activeHistoryId: string
  mode: StudioMode
  prompt: string
  title: string
  messages: ChatMessage[]
  updatedAt: number
  pinned?: boolean
  autoSaved?: boolean
}

const ACTIVE_SESSION_KEY = 'sub2api_ai_studio_active_session_v1'
const LOCAL_HISTORY_KEY = 'sub2api_ai_studio_local_history_v1'
const AUTO_SAVE_KEY = 'sub2api_ai_studio_auto_save_local_v1'
const LOCAL_HISTORY_RETENTION_KEY = 'sub2api_ai_studio_local_history_retention_days_v1'
const SELECTED_API_KEY_KEY = 'sub2api_ai_studio_selected_api_key_v1'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const { copyToClipboard } = useClipboard()

const restoredSession = readActiveSession()
const initialMode = isRouteMode(route.params.mode) ? normalizeMode(route.params.mode) : restoredSession?.mode ?? 'chat'

const activeMode = ref<StudioMode>(initialMode)
const prompt = ref(restoredSession?.prompt ?? '')
const messages = ref<ChatMessage[]>(restoredSession?.messages ?? [])
const currentSessionId = ref(restoredSession?.id ?? createSessionId())
const activeHistoryId = ref(restoredSession?.activeHistoryId ?? currentSessionId.value)
const keyOptions = ref<AIStudioKeyOption[]>([])
const keyOptionsLoading = ref(false)
const selectedAPIKeyId = ref<number | null>(readSelectedAPIKeyId())
const chatModel = ref('')
const imageModel = ref('')
const maxTokens = ref(2048)
const thinkingEnabled = ref(false)
const reasoningEffort = ref('medium')
const imageSize = ref('1024x1024')
const imageQuality = ref('standard')
const systemPrompt = ref('')
const historyQuery = ref('')
const showSettingsModal = ref(false)
const pendingAttachments = ref<ChatAttachment[]>([])
const serverConversations = ref<ConversationItem[]>([])
const isSubmitting = ref(false)
const editingMessageId = ref<number | string | null>(null)
const editingMessageContent = ref('')
const localHistory = ref<StoredSession[]>(readLocalHistory())
const autoSaveLocalHistory = ref(readAutoSavePreference())
const localHistoryRetentionDays = ref<number | null>(readLocalHistoryRetentionDays())
const fileInputRef = ref<HTMLInputElement | null>(null)
let suppressNextAutoSave = false
let noAPIKeyPromptShown = false
let noUsableAPIKeyPromptShown = false

const availableKeyOptions = computed(() => keyOptions.value.filter((item) => item.available))
const selectedKeyOption = computed(() => keyOptions.value.find((item) => item.id === selectedAPIKeyId.value) ?? null)
const selectedAvailableKeyOption = computed(() => selectedKeyOption.value?.available ? selectedKeyOption.value : null)
const chatModels = computed(() => selectedAvailableKeyOption.value?.chat_models ?? [])
const imageModels = computed(() => selectedAvailableKeyOption.value?.image_models ?? [])
const currentModel = computed({
  get: () => activeMode.value === 'image' ? imageModel.value : chatModel.value,
  set: (value: string) => {
    if (activeMode.value === 'image') {
      imageModel.value = value
      return
    }
    chatModel.value = value
  },
})
const availableModels = computed(() => activeMode.value === 'image' ? imageModels.value : chatModels.value)
const currentTitle = computed(() => activeMode.value === 'image' ? t('aiStudio.image.title') : t('aiStudio.chat.title'))
const currentSubtitle = computed(() => activeMode.value === 'image' ? t('aiStudio.image.subtitle') : t('aiStudio.chat.subtitle'))
const canSubmit = computed(() => (
  !isSubmitting.value
  && Boolean(selectedAvailableKeyOption.value)
  && availableModels.value.length > 0
  && Boolean(currentModel.value)
  && (prompt.value.trim().length > 0 || pendingAttachments.value.length > 0)
))
const localHistoryCount = computed(() => localHistory.value.length)

const modeOptions = computed(() => [
  { mode: 'chat' as StudioMode, label: t('aiStudio.chat.title'), icon: 'chat' as const },
  { mode: 'image' as StudioMode, label: t('aiStudio.image.title'), icon: 'sparkles' as const },
])

const settingsSummary = computed(() => {
  const keyName = selectedAvailableKeyOption.value?.name
    ?? (selectedKeyOption.value ? t('aiStudio.settings.unavailableKeySummary', { name: selectedKeyOption.value.name }) : t('aiStudio.settings.noAPIKey'))
  const modelName = currentModel.value || t('aiStudio.settings.modelUnavailable')
  return `${keyName} · ${modelName}`
})

const localConversations = computed<ConversationItem[]>(() => localHistory.value.map((session) => ({
  id: session.id,
  mode: session.mode,
  title: session.title,
  time: formatHistoryTime(session.updatedAt),
  source: 'local',
  pinned: session.pinned,
  autoSaved: session.autoSaved,
})))

const filteredConversations = computed(() => {
  const query = historyQuery.value.trim().toLowerCase()
  const all = [...localConversations.value, ...serverConversations.value]
  const visible = query ? all.filter((item) => item.title.toLowerCase().includes(query)) : all
  return visible.slice().sort((a, b) => Number(Boolean(b.pinned)) - Number(Boolean(a.pinned)))
})

watch(
  () => route.params.mode,
  (mode) => {
    if (isRouteMode(mode)) {
      activeMode.value = normalizeMode(mode)
    }
  },
)

watch(
  [activeMode, prompt, messages],
  () => {
    persistActiveSession()
    if (autoSaveLocalHistory.value) {
      if (suppressNextAutoSave) {
        suppressNextAutoSave = false
        return
      }
      upsertCurrentSessionToLocalHistory(true)
    }
  },
  { deep: true },
)

watch(autoSaveLocalHistory, (enabled) => {
  if (typeof window !== 'undefined') {
    window.localStorage.setItem(AUTO_SAVE_KEY, enabled ? 'true' : 'false')
  }
  if (enabled) {
    upsertCurrentSessionToLocalHistory(true)
    appStore.showSuccess(t('aiStudio.settings.autoSaveEnabled'))
  }
})

watch(localHistoryRetentionDays, (days) => {
  if (typeof window === 'undefined') return
  if (typeof days === 'number' && Number.isFinite(days) && days > 0) {
    const normalized = Math.floor(days)
    window.localStorage.setItem(LOCAL_HISTORY_RETENTION_KEY, String(normalized))
    if (normalized !== days) {
      localHistoryRetentionDays.value = normalized
      return
    }
    pruneLocalHistoryByRetention()
    return
  }
  window.localStorage.removeItem(LOCAL_HISTORY_RETENTION_KEY)
})

watch(selectedAPIKeyId, (id) => {
  if (typeof window !== 'undefined') {
    if (typeof id === 'number' && Number.isFinite(id)) {
      window.localStorage.setItem(SELECTED_API_KEY_KEY, String(id))
    } else {
      window.localStorage.removeItem(SELECTED_API_KEY_KEY)
    }
  }
  ensureSelectedModels()
})

watch([selectedKeyOption, activeMode], () => {
  ensureSelectedModels()
})

onBeforeUnmount(() => {
  persistActiveSession()
})

onMounted(() => {
  void loadKeyOptions()
  void loadServerConversations()
})

function isRouteMode(mode: unknown): mode is StudioMode {
  return mode === 'image' || mode === 'chat'
}

function normalizeMode(mode: unknown): StudioMode {
  return mode === 'image' ? 'image' : 'chat'
}

function createSessionId() {
  return `session-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function readActiveSession(): StoredSession | null {
  if (typeof window === 'undefined') return null
  const raw = window.sessionStorage.getItem(ACTIVE_SESSION_KEY)
  if (!raw) return null
  try {
    const session = normalizeStoredSession(JSON.parse(raw))
    if (!session) {
      window.sessionStorage.removeItem(ACTIVE_SESSION_KEY)
    }
    return session
  } catch {
    window.sessionStorage.removeItem(ACTIVE_SESSION_KEY)
    return null
  }
}

function readLocalHistory(): StoredSession[] {
  if (typeof window === 'undefined') return []
  const raw = window.localStorage.getItem(LOCAL_HISTORY_KEY)
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    const items = Array.isArray(parsed) ? parsed.slice(0, 50).map(normalizeStoredSession).filter((item): item is StoredSession => Boolean(item)) : []
    return applyLocalHistoryRetention(items)
  } catch {
    window.localStorage.removeItem(LOCAL_HISTORY_KEY)
    return []
  }
}

function normalizeStoredSession(value: unknown): StoredSession | null {
  if (!value || typeof value !== 'object') return null
  const raw = value as Partial<StoredSession>
  const mode = isRouteMode(raw.mode) ? raw.mode : 'chat'
  const id = typeof raw.id === 'string' && raw.id ? raw.id : createSessionId()
  const messages = Array.isArray(raw.messages) ? raw.messages.map(normalizeStoredMessage).filter((item): item is ChatMessage => Boolean(item)) : []
  return {
    id,
    activeHistoryId: typeof raw.activeHistoryId === 'string' && raw.activeHistoryId ? raw.activeHistoryId : id,
    mode,
    prompt: typeof raw.prompt === 'string' ? raw.prompt : '',
    title: typeof raw.title === 'string' && raw.title ? raw.title : (mode === 'image' ? t('aiStudio.history.newImage') : t('aiStudio.history.newChat')),
    messages,
    updatedAt: typeof raw.updatedAt === 'number' && Number.isFinite(raw.updatedAt) ? raw.updatedAt : Date.now(),
    pinned: Boolean(raw.pinned),
    autoSaved: Boolean(raw.autoSaved),
  }
}

function normalizeStoredMessage(value: unknown): ChatMessage | null {
  if (!value || typeof value !== 'object') return null
  const raw = value as Partial<ChatMessage>
  const role = raw.role === 'user' || raw.role === 'assistant' || raw.role === 'system' ? raw.role : null
  if (!role) return null
  return {
    id: typeof raw.id === 'number' || typeof raw.id === 'string' ? raw.id : Date.now(),
    serverId: typeof raw.serverId === 'number' ? raw.serverId : undefined,
    role,
    kind: isRouteMode(raw.kind) ? raw.kind : 'chat',
    content: typeof raw.content === 'string' ? raw.content : '',
    reasoning: typeof raw.reasoning === 'string' ? raw.reasoning : '',
    isStreaming: false,
    time: typeof raw.time === 'string' ? raw.time : '',
    attachments: Array.isArray(raw.attachments) ? raw.attachments.filter(isChatAttachmentLike) as ChatAttachment[] : undefined,
    images: Array.isArray(raw.images) ? raw.images.filter(isChatImageLike) as ChatImage[] : undefined,
  }
}

function isChatAttachmentLike(value: unknown): value is ChatAttachment {
  return Boolean(value && typeof value === 'object' && typeof (value as ChatAttachment).id === 'string' && typeof (value as ChatAttachment).name === 'string')
}

function isChatImageLike(value: unknown): value is ChatImage {
  return Boolean(value && typeof value === 'object' && typeof (value as ChatImage).id === 'string' && typeof (value as ChatImage).url === 'string')
}

function readLocalHistoryRetentionDays(): number | null {
  if (typeof window === 'undefined') return null
  const raw = window.localStorage.getItem(LOCAL_HISTORY_RETENTION_KEY)
  if (!raw) return null
  const value = Number(raw)
  return Number.isFinite(value) && value > 0 ? Math.floor(value) : null
}

function applyLocalHistoryRetention(items: StoredSession[]) {
  const days = readLocalHistoryRetentionDays()
  if (!days) return items
  const cutoff = Date.now() - days * 24 * 60 * 60 * 1000
  return items.filter((item) => item.updatedAt >= cutoff)
}

function readAutoSavePreference() {
  if (typeof window === 'undefined') return false
  return window.localStorage.getItem(AUTO_SAVE_KEY) === 'true'
}

function readSelectedAPIKeyId(): number | null {
  if (typeof window === 'undefined') return null
  const raw = window.localStorage.getItem(SELECTED_API_KEY_KEY)
  if (!raw) return null
  const id = Number(raw)
  return Number.isFinite(id) && id > 0 ? id : null
}

function writeLocalHistory(items: StoredSession[]) {
  if (typeof window === 'undefined') return
  const retained = applyLocalHistoryRetention(items).slice(0, 50)
  window.localStorage.setItem(LOCAL_HISTORY_KEY, JSON.stringify(retained))
  localHistory.value = retained
}

function pruneLocalHistoryByRetention() {
  writeLocalHistory(localHistory.value)
}

function normalizeLocalHistoryRetention() {
  if (typeof localHistoryRetentionDays.value !== 'number' || !Number.isFinite(localHistoryRetentionDays.value) || localHistoryRetentionDays.value <= 0) {
    localHistoryRetentionDays.value = null
    return
  }
  localHistoryRetentionDays.value = Math.floor(localHistoryRetentionDays.value)
}

function buildStoredSession(): StoredSession {
  const existing = localHistory.value.find((item) => item.id === currentSessionId.value)
  const storedMessages = messages.value.map((message) => ({ ...message, isStreaming: false }))
  return {
    id: currentSessionId.value,
    activeHistoryId: activeHistoryId.value,
    mode: activeMode.value,
    prompt: prompt.value,
    title: existing?.title || deriveSessionTitle(),
    messages: storedMessages,
    updatedAt: Date.now(),
    pinned: existing?.pinned,
    autoSaved: existing?.autoSaved || autoSaveLocalHistory.value,
  }
}

function persistActiveSession() {
  if (typeof window === 'undefined') return
  window.sessionStorage.setItem(ACTIVE_SESSION_KEY, JSON.stringify(buildStoredSession()))
}

function deriveSessionTitle() {
  const firstUserMessage = messages.value.find((message) => message.role === 'user')
  const firstContent = firstUserMessage?.content.trim()
  if (firstContent) return firstContent.slice(0, 24)
  if (pendingAttachments.value.length > 0) return pendingAttachments.value[0].name
  return activeMode.value === 'image' ? t('aiStudio.history.newImage') : t('aiStudio.history.newChat')
}

function formatHistoryTime(value: number) {
  return new Date(value).toLocaleDateString([], { month: '2-digit', day: '2-digit' })
}

async function loadKeyOptions() {
  keyOptionsLoading.value = true
  try {
    keyOptions.value = await aiStudioApi.listKeyOptions()
    const available = keyOptions.value.filter((item) => item.available)
    if (!available.some((item) => item.id === selectedAPIKeyId.value)) {
      selectedAPIKeyId.value = available[0]?.id ?? null
    } else {
      ensureSelectedModels()
    }
    if (keyOptions.value.length === 0) {
      promptCreateAPIKey()
    } else if (available.length === 0) {
      promptNoUsableAPIKey()
    }
  } catch (error) {
    appStore.showError(extractErrorMessage(error, t('aiStudio.errors.loadKeyOptionsFailed')))
  } finally {
    keyOptionsLoading.value = false
  }
}

function ensureSelectedModels() {
  if (!chatModels.value.includes(chatModel.value)) {
    chatModel.value = chatModels.value[0] ?? ''
  }
  if (!imageModels.value.includes(imageModel.value)) {
    imageModel.value = imageModels.value[0] ?? ''
  }
}

function formatAPIKeyOptionLabel(key: AIStudioKeyOption) {
  const groupName = key.group_name || key.platform || t('aiStudio.settings.noGroup')
  const base = `${key.name} · ${groupName}`
  if (key.available) return base
  return `${base} (${t('aiStudio.settings.unavailableKey')}: ${apiKeyUnavailableReasonLabel(key.unavailable_reason)})`
}

function firstUnavailableReasonLabel() {
  const first = keyOptions.value.find((item) => !item.available)
  return apiKeyUnavailableReasonLabel(first?.unavailable_reason)
}

function apiKeyUnavailableReasonLabel(reason?: string) {
  switch (reason) {
    case 'disabled':
      return t('aiStudio.settings.apiKeyUnavailableReasons.disabled')
    case 'expired':
      return t('aiStudio.settings.apiKeyUnavailableReasons.expired')
    case 'quota_exhausted':
      return t('aiStudio.settings.apiKeyUnavailableReasons.quotaExhausted')
    case 'missing_group':
      return t('aiStudio.settings.apiKeyUnavailableReasons.missingGroup')
    case 'non_openai_group':
      return t('aiStudio.settings.apiKeyUnavailableReasons.nonOpenAIGroup')
    case 'group_disabled':
      return t('aiStudio.settings.apiKeyUnavailableReasons.groupDisabled')
    default:
      return t('aiStudio.settings.apiKeyUnavailableReasons.unknown')
  }
}

function promptCreateAPIKey(force = false) {
  if (!force && noAPIKeyPromptShown) return
  noAPIKeyPromptShown = true
  const shouldGo = window.confirm(t('aiStudio.settings.noAPIKeyMessage'))
  if (shouldGo) {
    router.push({ name: 'Keys' })
  }
}

function promptNoUsableAPIKey(force = false) {
  if (!force && noUsableAPIKeyPromptShown) return
  noUsableAPIKeyPromptShown = true
  const shouldGo = window.confirm(t('aiStudio.settings.noUsableAPIKeyMessage', { reason: firstUnavailableReasonLabel() }))
  if (shouldGo) {
    router.push({ name: 'Keys' })
  }
}

function handleAPIKeyHelpClick() {
  if (keyOptions.value.length === 0) {
    promptCreateAPIKey(true)
    return
  }
  promptNoUsableAPIKey(true)
}

function switchMode(mode: StudioMode) {
  activeMode.value = mode
  persistActiveSession()
  if (route.name === 'AIStudioConversation') {
    router.replace({ name: 'AIStudioConversation', params: { mode } })
  }
}

async function openHistory(conversation: ConversationItem) {
  persistActiveSession()
  activeHistoryId.value = conversation.id
  activeMode.value = conversation.mode
  currentSessionId.value = conversation.id
  prompt.value = ''
  pendingAttachments.value = []

  if (conversation.source === 'local') {
    const stored = localHistory.value.find((item) => item.id === conversation.id)
    messages.value = stored?.messages ?? []
  } else {
    try {
      const detail = await aiStudioApi.getConversation(Number(conversation.id))
      messages.value = detail.messages.map(serverMessageToChatMessage).filter((message) => message.role !== 'system')
      currentSessionId.value = String(detail.conversation.id)
      activeHistoryId.value = String(detail.conversation.id)
      activeMode.value = detail.conversation.mode
    } catch (error) {
      appStore.showError(extractErrorMessage(error, t('aiStudio.errors.loadConversationFailed')))
    }
  }
  persistActiveSession()
}

function createMessage(role: ChatMessage['role'], kind: StudioMode, content: string, attachments?: ChatAttachment[]): ChatMessage {
  return {
    id: Date.now() + Math.floor(Math.random() * 1000),
    role,
    kind,
    content,
    reasoning: '',
    isStreaming: role === 'assistant' && content.trim() === '',
    attachments,
    time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
  }
}

function startNewConversation() {
  persistActiveSession()
  currentSessionId.value = createSessionId()
  activeHistoryId.value = currentSessionId.value
  prompt.value = ''
  messages.value = []
  pendingAttachments.value = []
  persistActiveSession()
}

function saveSettings() {
  showSettingsModal.value = false
  appStore.showSuccess(t('aiStudio.settings.saved'))
}

function saveCurrentSessionToLocalHistory() {
  const session = upsertCurrentSessionToLocalHistory(false)
  activeHistoryId.value = session.id
  appStore.showSuccess(t('aiStudio.settings.localHistorySaved'))
}

function upsertCurrentSessionToLocalHistory(autoSaved: boolean) {
  const session = { ...buildStoredSession(), autoSaved }
  const next = [session, ...localHistory.value.filter((item) => item.id !== session.id)]
  localHistory.value = next.slice(0, 50)
  writeLocalHistory(localHistory.value)
  return session
}

async function renameConversation(conversation: ConversationItem) {
  const nextTitle = window.prompt(t('aiStudio.actions.renamePrompt'), conversation.title)?.trim()
  if (!nextTitle) return
  if (conversation.source === 'local') {
    localHistory.value = localHistory.value.map((item) => item.id === conversation.id ? { ...item, title: nextTitle } : item)
    writeLocalHistory(localHistory.value)
    if (currentSessionId.value === conversation.id) persistActiveSession()
  } else {
    try {
      const updated = await aiStudioApi.updateConversation(Number(conversation.id), { title: nextTitle })
      upsertServerConversation(updated)
    } catch (error) {
      appStore.showError(extractErrorMessage(error, t('aiStudio.errors.renameFailed')))
    }
  }
}

async function deleteConversation(conversation: ConversationItem) {
  if (!window.confirm(t('aiStudio.actions.deleteConfirm'))) return
  if (conversation.source === 'local') {
    localHistory.value = localHistory.value.filter((item) => item.id !== conversation.id)
    writeLocalHistory(localHistory.value)
  } else {
    try {
      await aiStudioApi.deleteConversation(Number(conversation.id))
      serverConversations.value = serverConversations.value.filter((item) => item.id !== conversation.id)
    } catch (error) {
      appStore.showError(extractErrorMessage(error, t('aiStudio.errors.deleteFailed')))
      return
    }
  }
  if (activeHistoryId.value === conversation.id) {
    suppressNextAutoSave = true
    startNewConversation()
  }
}

async function toggleConversationPin(conversation: ConversationItem) {
  if (conversation.source === 'local') {
    localHistory.value = localHistory.value.map((item) => (
      item.id === conversation.id ? { ...item, pinned: !item.pinned } : item
    ))
    writeLocalHistory(localHistory.value)
    return
  }
  try {
    const updated = await aiStudioApi.updateConversation(Number(conversation.id), { pinned: !conversation.pinned })
    upsertServerConversation(updated)
  } catch (error) {
    appStore.showError(extractErrorMessage(error, t('aiStudio.errors.updateFailed')))
  }
}

function openFilePicker() {
  fileInputRef.value?.click()
}

function handleFilesSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  if (files.length === 0) return
  pendingAttachments.value.push(...files.map(fileToAttachment))
  input.value = ''
  persistActiveSession()
}

function fileToAttachment(file: File): ChatAttachment {
  return {
    id: `file-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    name: file.name,
    size: file.size,
    sizeLabel: formatFileSize(file.size),
    type: file.type || 'application/octet-stream',
    file,
  }
}

function removeAttachment(id: string) {
  pendingAttachments.value = pendingAttachments.value.filter((item) => item.id !== id)
}

function formatFileSize(size: number) {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

async function handleSubmit() {
  if (keyOptions.value.length === 0) {
    promptCreateAPIKey(true)
    return
  }
  if (!selectedAvailableKeyOption.value) {
    promptNoUsableAPIKey(true)
    return
  }
  if (availableModels.value.length === 0 || !currentModel.value) {
    appStore.showWarning(t('aiStudio.settings.modelUnavailable'))
    return
  }
  if (!canSubmit.value) return
  isSubmitting.value = true
  const value = prompt.value.trim()
  const attachments = [...pendingAttachments.value]
  const userContent = value || t('aiStudio.input.attachmentOnlyMessage')
  const userMessage = createMessage('user', activeMode.value, userContent, attachments.length > 0 ? attachments : undefined)
  const assistantMessage = createMessage('assistant', activeMode.value, '')
  messages.value.push(userMessage)
  messages.value.push(assistantMessage)
  prompt.value = ''
  pendingAttachments.value = []
  persistActiveSession()
  try {
    const uploaded = await uploadPendingAttachments(attachments)
    userMessage.attachments = uploaded.map(serverAttachmentToChatAttachment)
    await runAIStudioRequest(assistantMessage, {
      prompt: userContent,
      attachmentIds: uploaded.map((item) => item.id),
      excludeMessageIds: [userMessage.id, assistantMessage.id],
    })
  } catch (error) {
    assistantMessage.isStreaming = false
    assistantMessage.content = extractErrorMessage(error, t('aiStudio.errors.runFailed'))
    appStore.showError(assistantMessage.content)
  } finally {
    isSubmitting.value = false
    persistActiveSession()
  }
}

function startEditingMessage(message: ChatMessage) {
  if (message.role !== 'user') return
  editingMessageId.value = message.id
  editingMessageContent.value = message.content
}

function cancelEditingMessage() {
  editingMessageId.value = null
  editingMessageContent.value = ''
}

function canSubmitEditedMessage(message: ChatMessage) {
  return editingMessageContent.value.trim().length > 0 || Boolean(message.attachments?.length)
}

async function submitEditedMessage(message: ChatMessage) {
  if (message.role !== 'user' || !canSubmitEditedMessage(message)) return
  const messageIndex = messages.value.findIndex((item) => item.id === message.id)
  if (messageIndex === -1) return

  const nextContent = editingMessageContent.value.trim() || t('aiStudio.input.attachmentOnlyMessage')
  const updatedMessage: ChatMessage = {
    ...message,
    content: nextContent,
    time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
  }

  messages.value = [
    ...messages.value.slice(0, messageIndex),
    updatedMessage,
  ]
  const assistantMessage = createMessage('assistant', updatedMessage.kind, '')
  messages.value.push(assistantMessage)
  cancelEditingMessage()
  persistActiveSession()
  isSubmitting.value = true
  try {
    await runAIStudioRequest(assistantMessage, {
      prompt: nextContent,
      attachmentIds: (updatedMessage.attachments ?? []).map((item) => item.serverId).filter((id): id is number => typeof id === 'number'),
      resendFromMessageId: updatedMessage.serverId,
      excludeMessageIds: [updatedMessage.id, assistantMessage.id],
    })
  } catch (error) {
    assistantMessage.isStreaming = false
    assistantMessage.content = extractErrorMessage(error, t('aiStudio.errors.runFailed'))
    appStore.showError(assistantMessage.content)
  } finally {
    isSubmitting.value = false
    persistActiveSession()
  }
}

async function loadServerConversations() {
  try {
    const result = await aiStudioApi.listConversations()
    serverConversations.value = result.items.map(serverConversationToItem)
  } catch (error) {
    appStore.showError(extractErrorMessage(error, t('aiStudio.errors.loadHistoryFailed')))
  }
}

function upsertServerConversation(conversation: AIStudioConversation) {
  if (!conversation || !conversation.id) return
  const item = serverConversationToItem(conversation)
  serverConversations.value = [
    item,
    ...serverConversations.value.filter((existing) => existing.id !== item.id),
  ]
}

function serverConversationToItem(item: AIStudioConversation): ConversationItem {
  return {
    id: String(item.id),
    mode: item.mode,
    title: item.title,
    time: formatHistoryTime(new Date(item.last_active_at).getTime()),
    source: 'server',
    pinned: item.pinned,
  }
}

function serverMessageToChatMessage(item: AIStudioMessage): ChatMessage {
  const images = (item.attachments ?? [])
    .filter((attachment) => attachment.kind === 'generated_image')
    .map((attachment) => ({
      id: String(attachment.id),
      name: attachment.filename,
      url: attachment.url,
    }))
  return {
    id: item.id,
    serverId: item.id,
    role: item.role,
    kind: item.kind,
    content: item.content,
    reasoning: '',
    isStreaming: false,
    time: new Date(item.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
    attachments: (item.attachments ?? []).filter((attachment) => attachment.kind === 'upload').map(serverAttachmentToChatAttachment),
    images,
  }
}

function serverAttachmentToChatAttachment(item: AIStudioAttachment): ChatAttachment {
  return {
    id: String(item.id),
    serverId: item.id,
    name: item.filename,
    size: item.byte_size,
    sizeLabel: formatFileSize(item.byte_size),
    type: item.content_type,
    url: item.url,
  }
}

async function uploadPendingAttachments(attachments: ChatAttachment[]) {
  const uploaded: AIStudioAttachment[] = []
  for (const attachment of attachments) {
    if (!attachment.file) continue
    uploaded.push(await aiStudioApi.uploadAttachment(attachment.file))
  }
  return uploaded
}

async function runAIStudioRequest(
  assistantMessage: ChatMessage,
  options: { prompt: string; attachmentIds: number[]; resendFromMessageId?: number; excludeMessageIds?: Array<number | string> },
) {
  if (keyOptions.value.length === 0) {
    promptCreateAPIKey(true)
    throw new Error(t('aiStudio.settings.noAPIKey'))
  }
  if (!selectedAvailableKeyOption.value) {
    promptNoUsableAPIKey(true)
    throw new Error(t('aiStudio.settings.noUsableAPIKey'))
  }
  if (!currentModel.value || availableModels.value.length === 0) {
    throw new Error(t('aiStudio.settings.modelUnavailable'))
  }
  const excluded = new Set(options.excludeMessageIds ?? [assistantMessage.id])
  const localContext = autoSaveLocalHistory.value
    ? messages.value
      .filter((item) => item.role === 'user' || item.role === 'assistant')
      .slice(-10)
      .filter((item) => !excluded.has(item.id))
      .map((item) => ({ role: item.role as 'assistant' | 'user', kind: item.kind, content: item.content }))
    : undefined
  const payload: AIStudioRunPayload = {
    mode: activeMode.value,
    conversation_id: numericConversationId(),
    prompt: options.prompt,
    model: currentModel.value,
    attachment_ids: options.attachmentIds,
    api_key_id: selectedAvailableKeyOption.value.id,
    store_history: !autoSaveLocalHistory.value,
    resend_from_message_id: options.resendFromMessageId,
    system_prompt: systemPrompt.value,
    max_tokens: maxTokens.value,
    thinking_enabled: activeMode.value === 'chat' && thinkingEnabled.value,
    reasoning_effort: reasoningEffort.value,
    image_size: imageSize.value,
    image_quality: imageQuality.value,
    local_context: localContext,
  }
  await consumeRunStream(payload, assistantMessage)
}

async function consumeRunStream(payload: AIStudioRunPayload, assistantMessage: ChatMessage) {
  const token = window.localStorage.getItem('auth_token')
  const runResponse = await fetch(aiStudioApi.runURL(), {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(payload),
  })
  if (!runResponse.ok || !runResponse.body) {
    throw new Error(parseRunError(await runResponse.text(), runResponse.statusText))
  }
  const reader = runResponse.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  while (true) {
    const { value, done } = await reader.read()
    if (done) break
    buffer += decoder.decode(value, { stream: true }).replace(/\r\n/g, '\n')
    const events = buffer.split(/\n{2,}/)
    buffer = events.pop() ?? ''
    for (const event of events) {
      handleSSEEvent(event, assistantMessage)
    }
  }
  if (buffer.trim()) {
    handleSSEEvent(buffer, assistantMessage)
  }
  assistantMessage.isStreaming = false
  if (!assistantMessage.content.trim() && activeMode.value === 'image' && assistantMessage.images?.length) {
    assistantMessage.content = t('aiStudio.image.completed')
  }
  if (!assistantMessage.content.trim() && !assistantMessage.reasoning?.trim() && !assistantMessage.images?.length) {
    assistantMessage.content = t('aiStudio.errors.emptyResponse')
  }
  await loadServerConversations()
}

function handleSSEEvent(raw: string, assistantMessage: ChatMessage) {
  const lines = raw.split('\n')
  const eventName = lines.find((line) => line.startsWith('event:'))?.slice(6).trim()
  const data = lines
    .filter((line) => line.startsWith('data:'))
    .map((line) => line.slice(5).trimStart())
    .join('\n')
  if (!data || data === '[DONE]') return
  if (eventName === 'conversation') {
    const parsed = safeJSON(data) as AIStudioConversation | null
    if (parsed?.id) {
      currentSessionId.value = String(parsed.id)
      activeHistoryId.value = String(parsed.id)
      upsertServerConversation(parsed)
    }
    return
  }
  if (eventName === 'image') {
    const parsed = safeJSON(data) as { attachments?: AIStudioAttachment[]; conversation?: AIStudioConversation } | null
    assistantMessage.images = (parsed?.attachments ?? [])
      .filter((attachment) => attachment?.url)
      .map((attachment) => ({
      id: String(attachment.id),
      name: attachment.filename,
      url: attachment.url,
    }))
    assistantMessage.content = t('aiStudio.image.completed')
    assistantMessage.isStreaming = false
    if (parsed?.conversation) {
      currentSessionId.value = String(parsed.conversation.id)
      activeHistoryId.value = String(parsed.conversation.id)
      upsertServerConversation(parsed.conversation)
    }
    return
  }
  const parsed = safeJSON(data)
  if (!parsed || typeof parsed !== 'object') return
  if (appendErrorEvent(parsed, assistantMessage)) {
    return
  }
  if (appendDirectProviderDelta(parsed, assistantMessage)) {
    return
  }
  const reasoningDelta = firstPathString(parsed, [
    ['choices', 0, 'delta', 'reasoning_content'],
    ['choices', 0, 'delta', 'reasoning'],
    ['choices', 0, 'delta', 'thinking'],
    ['delta', 'reasoning_content'],
    ['delta', 'reasoning'],
    ['delta', 'thinking'],
  ])
  if (reasoningDelta) {
    assistantMessage.reasoning = `${assistantMessage.reasoning ?? ''}${reasoningDelta}`
  }
  const contentDelta = firstPathString(parsed, [
    ['choices', 0, 'delta', 'content'],
    ['choices', 0, 'message', 'content'],
    ['delta', 'content'],
    ['message', 'content'],
    ['content'],
  ])
  if (contentDelta) {
    assistantMessage.content += contentDelta
  }
  const finishReason = readPathString(parsed, ['choices', 0, 'finish_reason'])
  if (finishReason) {
    assistantMessage.isStreaming = false
  }
}

function appendErrorEvent(parsed: unknown, assistantMessage: ChatMessage) {
  const message = firstPathString(parsed, [
    ['error', 'message'],
    ['message'],
    ['detail'],
  ])
  if (!message) return false
  assistantMessage.content = message
  assistantMessage.isStreaming = false
  return true
}

function parseRunError(body: string, fallback: string) {
  const parsed = safeJSON(body)
  if (parsed && typeof parsed === 'object') {
    const message = firstPathString(parsed, [
      ['error', 'message'],
      ['message'],
      ['detail'],
    ])
    if (message) return message
  }
  return body || fallback
}

function appendDirectProviderDelta(parsed: unknown, assistantMessage: ChatMessage) {
  const eventType = readPathString(parsed, ['type'])
  if (eventType === 'response.output_text.delta') {
    appendAssistantContent(assistantMessage, readPathString(parsed, ['delta']))
    return true
  }
  if (eventType === 'response.output_text.done') {
    appendAssistantContentIfEmpty(assistantMessage, readPathString(parsed, ['text']))
    return true
  }
  if (eventType === 'response.reasoning_summary_text.delta' || eventType === 'response.reasoning_text.delta') {
    appendAssistantReasoning(assistantMessage, readPathString(parsed, ['delta']))
    return true
  }
  if (eventType === 'response.reasoning_summary_text.done' || eventType === 'response.reasoning_text.done') {
    appendAssistantReasoningIfEmpty(assistantMessage, readPathString(parsed, ['text']))
    return true
  }
  if (eventType === 'response.failed' || eventType === 'response.incomplete') {
    const message = firstPathString(parsed, [
      ['response', 'error', 'message'],
      ['response', 'incomplete_details', 'reason'],
      ['error', 'message'],
    ])
    if (message) assistantMessage.content = message
    assistantMessage.isStreaming = false
    return true
  }
  if (eventType === 'content_block_delta') {
    const deltaType = readPathString(parsed, ['delta', 'type'])
    if (deltaType === 'thinking_delta') {
      appendAssistantReasoning(assistantMessage, readPathString(parsed, ['delta', 'thinking']))
      return true
    }
    if (deltaType === 'text_delta') {
      appendAssistantContent(assistantMessage, readPathString(parsed, ['delta', 'text']))
      return true
    }
  }
  if (eventType === 'message_delta') {
    assistantMessage.isStreaming = false
    return true
  }
  if (eventType === 'response.completed') {
    appendCompletedResponseTextIfEmpty(parsed, assistantMessage)
    assistantMessage.isStreaming = false
    return true
  }
  return false
}

function appendAssistantReasoning(assistantMessage: ChatMessage, delta: string) {
  if (!delta) return
  assistantMessage.reasoning = `${assistantMessage.reasoning ?? ''}${delta}`
}

function appendAssistantContent(assistantMessage: ChatMessage, delta: string) {
  if (!delta) return
  assistantMessage.content += delta
}

function appendAssistantContentIfEmpty(assistantMessage: ChatMessage, text: string) {
  if (!text || assistantMessage.content.trim()) return
  assistantMessage.content = text
}

function appendAssistantReasoningIfEmpty(assistantMessage: ChatMessage, text: string) {
  if (!text || assistantMessage.reasoning?.trim()) return
  assistantMessage.reasoning = text
}

function appendCompletedResponseTextIfEmpty(parsed: unknown, assistantMessage: ChatMessage) {
  if (assistantMessage.content.trim()) return
  const response = readPath(parsed, ['response'])
  if (!response || typeof response !== 'object') return
  const output = (response as Record<string, unknown>).output
  if (!Array.isArray(output)) return
  const parts: string[] = []
  for (const item of output) {
    if (!item || typeof item !== 'object') continue
    const record = item as Record<string, unknown>
    if (record.type !== 'message' || !Array.isArray(record.content)) continue
    for (const part of record.content) {
      if (!part || typeof part !== 'object') continue
      const partRecord = part as Record<string, unknown>
      if (partRecord.type === 'output_text' && typeof partRecord.text === 'string') {
        parts.push(partRecord.text)
      }
    }
  }
  if (parts.length) {
    assistantMessage.content = parts.join('')
  }
}

function numericConversationId() {
  if (autoSaveLocalHistory.value) return undefined
  const id = Number(currentSessionId.value)
  return Number.isFinite(id) && id > 0 ? id : undefined
}

function safeJSON(value: string): unknown {
  try {
    return JSON.parse(value)
  } catch {
    return null
  }
}

function readPathString(value: unknown, path: Array<string | number>): string {
  const current = readPath(value, path)
  return typeof current === 'string' ? current : ''
}

function readPath(value: unknown, path: Array<string | number>): unknown {
  let current: unknown = value
  for (const key of path) {
    if (typeof key === 'number') {
      if (!Array.isArray(current)) return ''
      current = current[key]
      continue
    }
    if (!current || typeof current !== 'object') return ''
    current = (current as Record<string, unknown>)[key]
  }
  return current
}

function firstPathString(value: unknown, paths: Array<Array<string | number>>) {
  for (const path of paths) {
    const found = readPathString(value, path)
    if (found) return found
  }
  return ''
}

function extractErrorMessage(error: unknown, fallback: string) {
  if (error instanceof Error && error.message) return error.message
  if (error && typeof error === 'object' && 'message' in error) {
    const message = String((error as { message?: unknown }).message || '')
    if (message) return message
  }
  return fallback
}
</script>

<style scoped>
.ai-studio-shell {
  display: grid;
  grid-template-columns: 1fr;
  background:
    radial-gradient(circle at 88% 8%, rgba(249, 115, 22, 0.13), transparent 19rem),
    linear-gradient(135deg, rgba(255, 255, 255, 0.94), rgba(255, 248, 237, 0.82));
}

.ai-studio-side {
  min-height: auto;
}

.ai-studio-main {
  height: min(760px, calc(100dvh - 7rem));
  min-height: 0;
  background:
    radial-gradient(circle at 94% 12%, rgba(249, 115, 22, 0.1), transparent 18rem),
    linear-gradient(180deg, rgba(255, 255, 255, 0.76), rgba(255, 248, 237, 0.4));
}

.ai-studio-messages {
  overscroll-behavior: contain;
}

:global(.dark) .ai-studio-shell {
  background:
    radial-gradient(circle at 88% 8%, rgba(249, 115, 22, 0.12), transparent 19rem),
    linear-gradient(135deg, rgba(30, 41, 59, 0.86), rgba(15, 23, 42, 0.82));
}

:global(.dark) .ai-studio-main {
  background:
    radial-gradient(circle at 94% 12%, rgba(249, 115, 22, 0.1), transparent 18rem),
    linear-gradient(180deg, rgba(15, 23, 42, 0.5), rgba(2, 6, 23, 0.36));
}

.ai-image-tile::after {
  content: '';
  position: absolute;
  inset: 12%;
  border: 1px solid rgba(255, 255, 255, 0.55);
  border-radius: 9999px;
  transform: rotate(-18deg);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.18s ease;
}

.modal-enter-active > section,
.modal-leave-active > section {
  transition: transform 0.18s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > section,
.modal-leave-to > section {
  transform: translateY(8px) scale(0.98);
}

@media (min-width: 1024px) {
  .ai-studio-shell {
    grid-template-columns: 292px minmax(0, 1fr);
    height: calc(100dvh - 9rem);
    max-height: calc(100dvh - 9rem);
    min-height: 0;
  }

  .ai-studio-side {
    height: 100%;
    min-height: 0;
    max-height: calc(100dvh - 9rem);
  }

  .ai-studio-main {
    height: 100%;
    max-height: calc(100dvh - 9rem);
  }
}
</style>
