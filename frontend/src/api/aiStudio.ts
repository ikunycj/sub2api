import { apiClient } from './client'

export type AIStudioMode = 'chat' | 'image'

export interface AIStudioConversation {
  id: number
  mode: AIStudioMode
  title: string
  pinned: boolean
  last_active_at: string
  created_at: string
  updated_at: string
}

export interface AIStudioAttachment {
  id: number
  conversation_id?: number | null
  message_id?: number | null
  kind: 'upload' | 'generated_image'
  filename: string
  content_type: string
  byte_size: number
  expires_at: string
  created_at: string
  url: string
}

export interface AIStudioMessage {
  id: number
  conversation_id: number
  role: 'assistant' | 'user' | 'system'
  kind: AIStudioMode
  content: string
  model: string
  status: string
  sequence: number
  metadata?: Record<string, unknown>
  created_at: string
  updated_at: string
  attachments?: AIStudioAttachment[]
}

export interface AIStudioConversationDetail {
  conversation: AIStudioConversation
  messages: AIStudioMessage[]
}

export interface AIStudioKeyOption {
  id: number
  name: string
  group_id: number
  group_name: string
  platform: string
  status: string
  available: boolean
  unavailable_reason: string
  models: string[]
  chat_models: string[]
  image_models: string[]
  allow_image_generation: boolean
  last_used_at?: string | null
}

export interface AIStudioRunPayload {
  mode: AIStudioMode
  conversation_id?: number
  prompt: string
  model: string
  attachment_ids?: number[]
  api_key_id?: number
  store_history: boolean
  resend_from_message_id?: number
  system_prompt?: string
  temperature?: number
  max_tokens?: number
  thinking_enabled?: boolean
  reasoning_effort?: string
  image_size?: string
  image_quality?: string
  local_context?: Array<{
    role: 'assistant' | 'user' | 'system'
    kind: AIStudioMode
    content: string
  }>
}

export interface PaginatedAIStudioConversations {
  items: AIStudioConversation[]
  total: number
  page: number
  page_size: number
  pages: number
}

export const aiStudioApi = {
  async listKeyOptions(): Promise<AIStudioKeyOption[]> {
    const response = await apiClient.get<{ items: AIStudioKeyOption[] }>('/ai-studio/key-options')
    return response.data.items
  },

  async listConversations(): Promise<PaginatedAIStudioConversations> {
    const response = await apiClient.get('/ai-studio/conversations', {
      params: { page: 1, page_size: 100 }
    })
    return response.data
  },

  async getConversation(id: number): Promise<AIStudioConversationDetail> {
    const response = await apiClient.get(`/ai-studio/conversations/${id}`)
    return response.data
  },

  async updateConversation(id: number, payload: { title?: string; pinned?: boolean }): Promise<AIStudioConversation> {
    const response = await apiClient.patch(`/ai-studio/conversations/${id}`, payload)
    return response.data
  },

  async deleteConversation(id: number): Promise<void> {
    await apiClient.delete(`/ai-studio/conversations/${id}`)
  },

  async uploadAttachment(file: File): Promise<AIStudioAttachment> {
    const form = new FormData()
    form.append('file', file)
    const response = await apiClient.post('/ai-studio/attachments', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 120000,
    })
    return response.data
  },

  runURL(): string {
    const base = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    return `${base}/ai-studio/runs`
  },
}
