package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

const (
	AIStudioModeChat  = "chat"
	AIStudioModeImage = "image"

	AIStudioRoleSystem    = "system"
	AIStudioRoleUser      = "user"
	AIStudioRoleAssistant = "assistant"

	AIStudioAttachmentUpload         = "upload"
	AIStudioAttachmentGeneratedImage = "generated_image"

	AIStudioDefaultHistoryRetentionDays = 14
	AIStudioDefaultAttachmentTTL        = 24 * time.Hour
	AIStudioRecentRoundLimit            = 5
)

var (
	ErrAIStudioConversationNotFound = infraerrors.NotFound("AI_STUDIO_CONVERSATION_NOT_FOUND", "conversation not found")
	ErrAIStudioMessageNotFound      = infraerrors.NotFound("AI_STUDIO_MESSAGE_NOT_FOUND", "message not found")
	ErrAIStudioAttachmentNotFound   = infraerrors.NotFound("AI_STUDIO_ATTACHMENT_NOT_FOUND", "attachment not found")
	ErrAIStudioInvalidMode          = infraerrors.BadRequest("AI_STUDIO_INVALID_MODE", "mode must be chat or image")
	ErrAIStudioOpenAIKeyRequired    = infraerrors.BadRequest("AI_STUDIO_OPENAI_KEY_REQUIRED", "please create or bind an active OpenAI group API key first")
	ErrAIStudioOpenAIKeyOnly        = infraerrors.BadRequest("AI_STUDIO_OPENAI_KEY_ONLY", "AI Studio currently only supports OpenAI group API keys")
	ErrAIStudioAttachmentExpired    = infraerrors.NotFound("AI_STUDIO_ATTACHMENT_EXPIRED", "attachment has expired")
)

type AIStudioConversation struct {
	ID               int64
	UserID           int64
	Mode             string
	Title            string
	Pinned           bool
	Summary          string
	SummaryMessageID *int64
	LastActiveAt     time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type AIStudioMessage struct {
	ID             int64
	ConversationID int64
	UserID         int64
	Role           string
	Kind           string
	Content        string
	Model          string
	Status         string
	Sequence       int64
	Metadata       map[string]any
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Attachments    []AIStudioAttachment
}

type AIStudioAttachment struct {
	ID             int64
	UserID         int64
	ConversationID *int64
	MessageID      *int64
	Kind           string
	Filename       string
	StoragePath    string
	ContentType    string
	ByteSize       int64
	SHA256         string
	ExpiresAt      time.Time
	CreatedAt      time.Time
}

type AIStudioKeyOption struct {
	ID                   int64
	Name                 string
	GroupID              int64
	GroupName            string
	Platform             string
	Status               string
	Available            bool
	UnavailableReason    string
	AllowImageGeneration bool
	LastUsedAt           *time.Time
}

type AIStudioExpiredAttachment struct {
	ID          int64
	StoragePath string
}

type AIStudioRepository interface {
	ListConversations(ctx context.Context, userID int64, limit, offset int) ([]AIStudioConversation, int64, error)
	GetConversation(ctx context.Context, userID, conversationID int64) (*AIStudioConversation, error)
	CreateConversation(ctx context.Context, conversation *AIStudioConversation) error
	UpdateConversation(ctx context.Context, conversation *AIStudioConversation) error
	DeleteConversation(ctx context.Context, userID, conversationID int64) error
	TouchConversation(ctx context.Context, userID, conversationID int64, at time.Time) error

	ListMessages(ctx context.Context, userID, conversationID int64) ([]AIStudioMessage, error)
	ListRecentMessages(ctx context.Context, userID, conversationID int64, limit int) ([]AIStudioMessage, error)
	CountMessages(ctx context.Context, userID, conversationID int64) (int64, error)
	CreateMessage(ctx context.Context, message *AIStudioMessage) error
	UpdateUserMessageAndDeleteAfter(ctx context.Context, userID, messageID int64, content string, metadata map[string]any) (*AIStudioMessage, error)
	AttachToMessage(ctx context.Context, userID int64, attachmentIDs []int64, conversationID, messageID int64) error
	UpdateSummary(ctx context.Context, userID, conversationID int64, summary string, summaryMessageID *int64) error

	CreateAttachment(ctx context.Context, attachment *AIStudioAttachment) error
	GetAttachment(ctx context.Context, userID, attachmentID int64) (*AIStudioAttachment, error)
	ListAttachmentsByIDs(ctx context.Context, userID int64, ids []int64) ([]AIStudioAttachment, error)
	ListAttachmentsByMessageIDs(ctx context.Context, userID int64, messageIDs []int64) (map[int64][]AIStudioAttachment, error)

	DeleteExpiredConversations(ctx context.Context, cutoff time.Time, batchSize int) (int64, error)
	ListExpiredAttachments(ctx context.Context, now time.Time, limit int) ([]AIStudioExpiredAttachment, error)
	DeleteAttachmentsByIDs(ctx context.Context, ids []int64) (int64, error)
}

type AIStudioService struct {
	repo          AIStudioRepository
	apiKeyService *APIKeyService
	cfg           *config.Config
}

func NewAIStudioService(repo AIStudioRepository, apiKeyService *APIKeyService, cfg *config.Config) *AIStudioService {
	return &AIStudioService{repo: repo, apiKeyService: apiKeyService, cfg: cfg}
}

func (s *AIStudioService) ListConversations(ctx context.Context, userID int64, params pagination.PaginationParams) ([]AIStudioConversation, *pagination.PaginationResult, error) {
	if s == nil || s.repo == nil {
		return nil, nil, errors.New("ai studio service is not configured")
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 30
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	items, total, err := s.repo.ListConversations(ctx, userID, params.Limit(), params.Offset())
	if err != nil {
		return nil, nil, fmt.Errorf("list ai studio conversations: %w", err)
	}
	return items, paginationResult(total, params), nil
}

func (s *AIStudioService) GetConversationWithMessages(ctx context.Context, userID, conversationID int64) (*AIStudioConversation, []AIStudioMessage, error) {
	if s == nil || s.repo == nil {
		return nil, nil, errors.New("ai studio service is not configured")
	}
	conversation, err := s.repo.GetConversation(ctx, userID, conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("get ai studio conversation: %w", err)
	}
	messages, err := s.repo.ListMessages(ctx, userID, conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("list ai studio messages: %w", err)
	}
	if err := s.hydrateMessageAttachments(ctx, userID, messages); err != nil {
		return nil, nil, err
	}
	return conversation, messages, nil
}

func (s *AIStudioService) UpdateConversation(ctx context.Context, userID, conversationID int64, title *string, pinned *bool) (*AIStudioConversation, error) {
	conversation, err := s.repo.GetConversation(ctx, userID, conversationID)
	if err != nil {
		return nil, fmt.Errorf("get ai studio conversation: %w", err)
	}
	if title != nil {
		conversation.Title = strings.TrimSpace(*title)
		if conversation.Title == "" {
			conversation.Title = defaultAIStudioConversationTitle(conversation.Mode)
		}
	}
	if pinned != nil {
		conversation.Pinned = *pinned
	}
	if err := s.repo.UpdateConversation(ctx, conversation); err != nil {
		return nil, fmt.Errorf("update ai studio conversation: %w", err)
	}
	return conversation, nil
}

func (s *AIStudioService) DeleteConversation(ctx context.Context, userID, conversationID int64) error {
	if err := s.repo.DeleteConversation(ctx, userID, conversationID); err != nil {
		return fmt.Errorf("delete ai studio conversation: %w", err)
	}
	return nil
}

func (s *AIStudioService) SaveUploadAttachment(ctx context.Context, userID int64, filename, contentType string, data []byte) (*AIStudioAttachment, error) {
	if len(data) == 0 {
		return nil, infraerrors.BadRequest("AI_STUDIO_EMPTY_ATTACHMENT", "attachment is empty")
	}
	storagePath, err := s.writeAttachmentFile(userID, filename, contentType, data)
	if err != nil {
		return nil, fmt.Errorf("save ai studio attachment: %w", err)
	}
	attachment := &AIStudioAttachment{
		UserID:      userID,
		Kind:        AIStudioAttachmentUpload,
		Filename:    sanitizeAttachmentFilename(filename),
		StoragePath: storagePath,
		ContentType: normalizeContentType(contentType, data),
		ByteSize:    int64(len(data)),
		SHA256:      sha256Hex(data),
		ExpiresAt:   time.Now().UTC().Add(AIStudioDefaultAttachmentTTL),
	}
	if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
		_ = os.Remove(storagePath)
		return nil, fmt.Errorf("create ai studio attachment: %w", err)
	}
	return attachment, nil
}

func (s *AIStudioService) GetAttachmentContent(ctx context.Context, userID, attachmentID int64) (*AIStudioAttachment, error) {
	attachment, err := s.repo.GetAttachment(ctx, userID, attachmentID)
	if err != nil {
		return nil, fmt.Errorf("get ai studio attachment: %w", err)
	}
	if time.Now().UTC().After(attachment.ExpiresAt) {
		return nil, ErrAIStudioAttachmentExpired
	}
	return attachment, nil
}

func (s *AIStudioService) ListOpenAIKeyOptions(ctx context.Context, userID int64) ([]AIStudioKeyOption, error) {
	if s == nil || s.apiKeyService == nil {
		return nil, errors.New("api key service is not configured")
	}
	keys, _, err := s.apiKeyService.List(ctx, userID, pagination.PaginationParams{
		Page:      1,
		PageSize:  100,
		SortBy:    "last_used_at",
		SortOrder: pagination.SortOrderDesc,
	}, APIKeyListFilters{})
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		if keys[i].LastUsedAt == nil && keys[j].LastUsedAt != nil {
			return false
		}
		if keys[i].LastUsedAt != nil && keys[j].LastUsedAt == nil {
			return true
		}
		if keys[i].LastUsedAt != nil && keys[j].LastUsedAt != nil && !keys[i].LastUsedAt.Equal(*keys[j].LastUsedAt) {
			return keys[i].LastUsedAt.After(*keys[j].LastUsedAt)
		}
		return keys[i].ID > keys[j].ID
	})
	out := make([]AIStudioKeyOption, 0, len(keys))
	for i := range keys {
		key := &keys[i]
		available, unavailableReason := classifyAIStudioOpenAIKey(key)
		groupID := int64(0)
		groupName := ""
		platform := ""
		allowImageGeneration := false
		if key.GroupID != nil {
			groupID = *key.GroupID
		}
		if key.Group != nil {
			groupName = key.Group.Name
			platform = key.Group.Platform
			allowImageGeneration = key.Group.AllowImageGeneration
		}
		out = append(out, AIStudioKeyOption{
			ID:                   key.ID,
			Name:                 key.Name,
			GroupID:              groupID,
			GroupName:            groupName,
			Platform:             platform,
			Status:               key.Status,
			Available:            available,
			UnavailableReason:    unavailableReason,
			AllowImageGeneration: allowImageGeneration,
			LastUsedAt:           key.LastUsedAt,
		})
	}
	return out, nil
}

func (s *AIStudioService) SelectOpenAIAPIKey(ctx context.Context, userID int64, apiKeyID *int64) (*APIKey, error) {
	if s == nil || s.apiKeyService == nil {
		return nil, errors.New("api key service is not configured")
	}
	if apiKeyID != nil {
		apiKey, err := s.apiKeyService.GetByID(ctx, *apiKeyID)
		if err != nil {
			return nil, fmt.Errorf("get api key: %w", err)
		}
		if apiKey.UserID != userID {
			return nil, ErrAIStudioOpenAIKeyRequired
		}
		if err := validateAIStudioOpenAIKey(apiKey); err != nil {
			return nil, err
		}
		return apiKey, nil
	}

	keys, _, err := s.apiKeyService.List(ctx, userID, pagination.PaginationParams{
		Page:      1,
		PageSize:  100,
		SortBy:    "last_used_at",
		SortOrder: pagination.SortOrderDesc,
	}, APIKeyListFilters{Status: StatusAPIKeyActive})
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		if keys[i].LastUsedAt == nil && keys[j].LastUsedAt != nil {
			return false
		}
		if keys[i].LastUsedAt != nil && keys[j].LastUsedAt == nil {
			return true
		}
		if keys[i].LastUsedAt != nil && keys[j].LastUsedAt != nil && !keys[i].LastUsedAt.Equal(*keys[j].LastUsedAt) {
			return keys[i].LastUsedAt.After(*keys[j].LastUsedAt)
		}
		return keys[i].ID > keys[j].ID
	})
	for i := range keys {
		key := &keys[i]
		if isAIStudioOpenAIKeyCandidate(key) {
			full, err := s.apiKeyService.GetByID(ctx, key.ID)
			if err != nil {
				continue
			}
			if err := validateAIStudioOpenAIKey(full); err == nil {
				return full, nil
			}
		}
	}
	return nil, ErrAIStudioOpenAIKeyRequired
}

func (s *AIStudioService) TouchSelectedAPIKey(ctx context.Context, apiKeyID int64) {
	if s == nil || s.apiKeyService == nil {
		return
	}
	_ = s.apiKeyService.TouchLastUsed(ctx, apiKeyID)
}

type AIStudioPrepareRunInput struct {
	UserID              int64
	Mode                string
	ConversationID      *int64
	Prompt              string
	Model               string
	AttachmentIDs       []int64
	StoreHistory        bool
	ResendFromMessageID *int64
	SystemPrompt        string
	MaxTokens           int
	Temperature         *float64
	LocalContext        []AIStudioMessage
}

type AIStudioPreparedRun struct {
	Conversation *AIStudioConversation
	UserMessage  *AIStudioMessage
	Messages     []map[string]any
	Attachments  []AIStudioAttachment
}

func (s *AIStudioService) PrepareRun(ctx context.Context, input AIStudioPrepareRunInput) (*AIStudioPreparedRun, error) {
	if err := validateAIStudioMode(input.Mode); err != nil {
		return nil, err
	}
	prompt := strings.TrimSpace(input.Prompt)
	if prompt == "" && len(input.AttachmentIDs) == 0 {
		return nil, infraerrors.BadRequest("AI_STUDIO_PROMPT_REQUIRED", "prompt or attachment is required")
	}

	attachments, err := s.repo.ListAttachmentsByIDs(ctx, input.UserID, input.AttachmentIDs)
	if err != nil {
		return nil, fmt.Errorf("list ai studio attachments: %w", err)
	}
	if len(attachments) != len(uniqueInt64s(input.AttachmentIDs)) {
		return nil, ErrAIStudioAttachmentNotFound
	}
	now := time.Now().UTC()
	for _, attachment := range attachments {
		if now.After(attachment.ExpiresAt) {
			return nil, ErrAIStudioAttachmentExpired
		}
	}

	var conversation *AIStudioConversation
	var userMessage *AIStudioMessage
	var history []AIStudioMessage
	if input.StoreHistory {
		if input.ResendFromMessageID != nil {
			userMessage, err = s.repo.UpdateUserMessageAndDeleteAfter(ctx, input.UserID, *input.ResendFromMessageID, prompt, buildAttachmentMetadata(input.AttachmentIDs))
			if err != nil {
				return nil, fmt.Errorf("update ai studio message: %w", err)
			}
			conversation, err = s.repo.GetConversation(ctx, input.UserID, userMessage.ConversationID)
			if err != nil {
				return nil, fmt.Errorf("get ai studio conversation: %w", err)
			}
			conversation.Mode = input.Mode
			conversation.Title = deriveAIStudioTitle(conversation.Title, prompt, input.Mode)
			conversation.LastActiveAt = now
			if err := s.repo.UpdateConversation(ctx, conversation); err != nil {
				return nil, fmt.Errorf("update ai studio conversation: %w", err)
			}
		} else {
			if input.ConversationID != nil && *input.ConversationID > 0 {
				conversation, err = s.repo.GetConversation(ctx, input.UserID, *input.ConversationID)
				if err != nil {
					return nil, fmt.Errorf("get ai studio conversation: %w", err)
				}
				conversation.Mode = input.Mode
				conversation.Title = deriveAIStudioTitle(conversation.Title, prompt, input.Mode)
				conversation.LastActiveAt = now
				if err := s.repo.UpdateConversation(ctx, conversation); err != nil {
					return nil, fmt.Errorf("update ai studio conversation: %w", err)
				}
			} else {
				conversation = &AIStudioConversation{
					UserID:       input.UserID,
					Mode:         input.Mode,
					Title:        deriveAIStudioTitle("", prompt, input.Mode),
					LastActiveAt: now,
				}
				if err := s.repo.CreateConversation(ctx, conversation); err != nil {
					return nil, fmt.Errorf("create ai studio conversation: %w", err)
				}
			}
			userMessage = &AIStudioMessage{
				ConversationID: conversation.ID,
				UserID:         input.UserID,
				Role:           AIStudioRoleUser,
				Kind:           input.Mode,
				Content:        prompt,
				Model:          input.Model,
				Status:         "completed",
				Metadata:       buildAttachmentMetadata(input.AttachmentIDs),
			}
			if err := s.repo.CreateMessage(ctx, userMessage); err != nil {
				return nil, fmt.Errorf("create ai studio user message: %w", err)
			}
		}
		if len(input.AttachmentIDs) > 0 {
			if err := s.repo.AttachToMessage(ctx, input.UserID, input.AttachmentIDs, conversation.ID, userMessage.ID); err != nil {
				return nil, fmt.Errorf("attach ai studio files: %w", err)
			}
		}
		history, err = s.repo.ListRecentMessages(ctx, input.UserID, conversation.ID, AIStudioRecentRoundLimit*2)
		if err != nil {
			return nil, fmt.Errorf("list recent ai studio messages: %w", err)
		}
		history = removeMessageByID(history, userMessage.ID)
	} else {
		history = input.LocalContext
	}

	messages, err := s.BuildOpenAIChatMessages(ctx, input, conversation, history, prompt, attachments)
	if err != nil {
		return nil, err
	}
	return &AIStudioPreparedRun{
		Conversation: conversation,
		UserMessage:  userMessage,
		Messages:     messages,
		Attachments:  attachments,
	}, nil
}

func (s *AIStudioService) BuildOpenAIChatMessages(
	ctx context.Context,
	input AIStudioPrepareRunInput,
	conversation *AIStudioConversation,
	history []AIStudioMessage,
	prompt string,
	attachments []AIStudioAttachment,
) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(history)+3)
	if systemPrompt := strings.TrimSpace(input.SystemPrompt); systemPrompt != "" {
		out = append(out, map[string]any{"role": "system", "content": systemPrompt})
	}
	if conversation != nil && strings.TrimSpace(conversation.Summary) != "" {
		out = append(out, map[string]any{
			"role":    "system",
			"content": "Conversation summary so far:\n" + strings.TrimSpace(conversation.Summary),
		})
	}
	for _, msg := range recentRounds(history, AIStudioRecentRoundLimit) {
		if msg.Role != AIStudioRoleUser && msg.Role != AIStudioRoleAssistant && msg.Role != AIStudioRoleSystem {
			continue
		}
		content := strings.TrimSpace(msg.Content)
		if content == "" {
			continue
		}
		out = append(out, map[string]any{"role": msg.Role, "content": content})
	}
	currentContent, err := s.buildUserContent(ctx, prompt, attachments, input.Mode)
	if err != nil {
		return nil, err
	}
	out = append(out, map[string]any{"role": "user", "content": currentContent})
	return out, nil
}

func (s *AIStudioService) CompleteAssistantMessage(ctx context.Context, userID int64, prepared *AIStudioPreparedRun, mode, model, content string, metadata map[string]any) (*AIStudioMessage, error) {
	if prepared == nil || prepared.Conversation == nil {
		return nil, nil
	}
	message := &AIStudioMessage{
		ConversationID: prepared.Conversation.ID,
		UserID:         userID,
		Role:           AIStudioRoleAssistant,
		Kind:           mode,
		Content:        content,
		Model:          model,
		Status:         "completed",
		Metadata:       metadata,
	}
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("create ai studio assistant message: %w", err)
	}
	_ = s.repo.TouchConversation(ctx, userID, prepared.Conversation.ID, time.Now().UTC())
	return message, nil
}

func (s *AIStudioService) CompleteGeneratedImage(ctx context.Context, userID int64, prepared *AIStudioPreparedRun, filename, contentType string, data []byte) (*AIStudioAttachment, error) {
	if len(data) == 0 {
		return nil, nil
	}
	storagePath, err := s.writeAttachmentFile(userID, filename, contentType, data)
	if err != nil {
		return nil, fmt.Errorf("save ai studio generated image: %w", err)
	}
	var conversationID *int64
	if prepared != nil && prepared.Conversation != nil {
		id := prepared.Conversation.ID
		conversationID = &id
	}
	attachment := &AIStudioAttachment{
		UserID:         userID,
		ConversationID: conversationID,
		Kind:           AIStudioAttachmentGeneratedImage,
		Filename:       sanitizeAttachmentFilename(filename),
		StoragePath:    storagePath,
		ContentType:    normalizeContentType(contentType, data),
		ByteSize:       int64(len(data)),
		SHA256:         sha256Hex(data),
		ExpiresAt:      time.Now().UTC().Add(AIStudioDefaultAttachmentTTL),
	}
	if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
		_ = os.Remove(storagePath)
		return nil, fmt.Errorf("create ai studio generated image: %w", err)
	}
	return attachment, nil
}

func (s *AIStudioService) AttachGeneratedAttachmentsToMessage(ctx context.Context, userID int64, attachmentIDs []int64, conversationID, messageID int64) error {
	if s == nil || s.repo == nil || len(attachmentIDs) == 0 || conversationID <= 0 || messageID <= 0 {
		return nil
	}
	if err := s.repo.AttachToMessage(ctx, userID, attachmentIDs, conversationID, messageID); err != nil {
		return fmt.Errorf("attach ai studio generated images: %w", err)
	}
	return nil
}

func (s *AIStudioService) MaybeCompactSummary(ctx context.Context, userID int64, conversationID int64) {
	if s == nil || s.repo == nil || conversationID <= 0 {
		return
	}
	count, err := s.repo.CountMessages(ctx, userID, conversationID)
	if err != nil || count <= int64(AIStudioRecentRoundLimit*2+4) {
		return
	}
	messages, err := s.repo.ListMessages(ctx, userID, conversationID)
	if err != nil || len(messages) == 0 {
		return
	}
	summary, coveredID := summarizeMessagesDeterministically(messages)
	if strings.TrimSpace(summary) == "" {
		return
	}
	_ = s.repo.UpdateSummary(ctx, userID, conversationID, summary, coveredID)
}

func (s *AIStudioService) ReadAttachmentBytes(attachment *AIStudioAttachment, limit int64) ([]byte, error) {
	if attachment == nil {
		return nil, ErrAIStudioAttachmentNotFound
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(attachment.StoragePath)), "http://") ||
		strings.HasPrefix(strings.ToLower(strings.TrimSpace(attachment.StoragePath)), "https://") {
		return nil, infraerrors.BadRequest("AI_STUDIO_REMOTE_ATTACHMENT", "remote attachment content is not available")
	}
	if limit <= 0 {
		limit = 25 << 20
	}
	f, err := os.Open(attachment.StoragePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(io.LimitReader(f, limit))
}

func (s *AIStudioService) hydrateMessageAttachments(ctx context.Context, userID int64, messages []AIStudioMessage) error {
	if len(messages) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(messages))
	for i := range messages {
		ids = append(ids, messages[i].ID)
	}
	grouped, err := s.repo.ListAttachmentsByMessageIDs(ctx, userID, ids)
	if err != nil {
		return fmt.Errorf("list ai studio message attachments: %w", err)
	}
	for i := range messages {
		messages[i].Attachments = grouped[messages[i].ID]
	}
	return nil
}

func (s *AIStudioService) buildUserContent(ctx context.Context, prompt string, attachments []AIStudioAttachment, mode string) (any, error) {
	if len(attachments) == 0 {
		return prompt, nil
	}
	if mode == AIStudioModeImage {
		return prompt, nil
	}
	parts := make([]map[string]any, 0, len(attachments)+1)
	if strings.TrimSpace(prompt) != "" {
		parts = append(parts, map[string]any{"type": "text", "text": prompt})
	}
	textSnippets := make([]string, 0)
	for _, attachment := range attachments {
		if isTextContentType(attachment.ContentType) {
			text, err := readTextAttachment(attachment.StoragePath)
			if err != nil {
				return nil, err
			}
			textSnippets = append(textSnippets, fmt.Sprintf("File: %s\n%s", attachment.Filename, text))
			continue
		}
		if strings.HasPrefix(strings.ToLower(attachment.ContentType), "image/") {
			dataURL, err := fileDataURL(attachment.StoragePath, attachment.ContentType)
			if err != nil {
				return nil, err
			}
			parts = append(parts, map[string]any{
				"type": "image_url",
				"image_url": map[string]any{
					"url": dataURL,
				},
			})
		}
	}
	if len(textSnippets) > 0 {
		parts = append(parts, map[string]any{
			"type": "text",
			"text": "Attached text files:\n\n" + strings.Join(textSnippets, "\n\n---\n\n"),
		})
	}
	if len(parts) == 0 {
		return prompt, nil
	}
	return parts, nil
}

func (s *AIStudioService) writeAttachmentFile(userID int64, filename, contentType string, data []byte) (string, error) {
	base := "./data"
	if s != nil && s.cfg != nil && strings.TrimSpace(s.cfg.Pricing.DataDir) != "" {
		base = strings.TrimSpace(s.cfg.Pricing.DataDir)
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		exts, _ := mime.ExtensionsByType(normalizeContentType(contentType, data))
		if len(exts) > 0 {
			ext = exts[0]
		}
	}
	if ext == "" || len(ext) > 16 {
		ext = ".bin"
	}
	dir := filepath.Join(base, "ai-studio", fmt.Sprintf("%d", userID), time.Now().UTC().Format("20060102"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randomSuffix(), ext))
	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", err
	}
	return path, nil
}

func validateAIStudioMode(mode string) error {
	switch mode {
	case AIStudioModeChat, AIStudioModeImage:
		return nil
	default:
		return ErrAIStudioInvalidMode
	}
}

func validateAIStudioOpenAIKey(apiKey *APIKey) error {
	if apiKey.User == nil {
		return ErrAIStudioOpenAIKeyRequired
	}
	if apiKey.UserID <= 0 {
		apiKey.UserID = apiKey.User.ID
	}
	available, reason := classifyAIStudioOpenAIKey(apiKey)
	if !available {
		if reason == AIStudioKeyUnavailableNonOpenAIGroup {
			return ErrAIStudioOpenAIKeyOnly
		}
		return ErrAIStudioOpenAIKeyRequired
	}
	return nil
}

func isAIStudioOpenAIKeyCandidate(apiKey *APIKey) bool {
	available, _ := classifyAIStudioOpenAIKey(apiKey)
	return available
}

const (
	AIStudioKeyUnavailableMissingKey     = "missing_key"
	AIStudioKeyUnavailableDisabled       = "disabled"
	AIStudioKeyUnavailableExpired        = "expired"
	AIStudioKeyUnavailableQuotaExhausted = "quota_exhausted"
	AIStudioKeyUnavailableMissingGroup   = "missing_group"
	AIStudioKeyUnavailableNonOpenAIGroup = "non_openai_group"
	AIStudioKeyUnavailableGroupDisabled  = "group_disabled"
)

func classifyAIStudioOpenAIKey(apiKey *APIKey) (bool, string) {
	if apiKey == nil {
		return false, AIStudioKeyUnavailableMissingKey
	}
	if apiKey.Status == StatusAPIKeyExpired || apiKey.IsExpired() {
		return false, AIStudioKeyUnavailableExpired
	}
	if apiKey.Status == StatusAPIKeyQuotaExhausted || apiKey.IsQuotaExhausted() {
		return false, AIStudioKeyUnavailableQuotaExhausted
	}
	if apiKey.Status != StatusAPIKeyActive {
		return false, AIStudioKeyUnavailableDisabled
	}
	if apiKey.Group == nil || apiKey.GroupID == nil {
		return false, AIStudioKeyUnavailableMissingGroup
	}
	if apiKey.Group.Platform != PlatformOpenAI {
		return false, AIStudioKeyUnavailableNonOpenAIGroup
	}
	if !apiKey.Group.IsActive() {
		return false, AIStudioKeyUnavailableGroupDisabled
	}
	return true, ""
}

func deriveAIStudioTitle(current, prompt, mode string) string {
	current = strings.TrimSpace(current)
	if current != "" && current != defaultAIStudioConversationTitle(mode) {
		return current
	}
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return defaultAIStudioConversationTitle(mode)
	}
	return truncateRunes(prompt, 32)
}

func defaultAIStudioConversationTitle(mode string) string {
	if mode == AIStudioModeImage {
		return "新生图"
	}
	return "新对话"
}

func recentRounds(messages []AIStudioMessage, rounds int) []AIStudioMessage {
	if rounds <= 0 || len(messages) == 0 {
		return nil
	}
	limit := rounds * 2
	if len(messages) <= limit {
		return messages
	}
	return messages[len(messages)-limit:]
}

func removeMessageByID(messages []AIStudioMessage, id int64) []AIStudioMessage {
	out := messages[:0]
	for _, msg := range messages {
		if msg.ID != id {
			out = append(out, msg)
		}
	}
	return out
}

func buildAttachmentMetadata(ids []int64) map[string]any {
	if len(ids) == 0 {
		return map[string]any{}
	}
	return map[string]any{"attachment_ids": uniqueInt64s(ids)}
}

func summarizeMessagesDeterministically(messages []AIStudioMessage) (string, *int64) {
	if len(messages) <= AIStudioRecentRoundLimit*2 {
		return "", nil
	}
	cutoff := len(messages) - AIStudioRecentRoundLimit*2
	covered := messages[:cutoff]
	lines := make([]string, 0, len(covered))
	for _, msg := range covered {
		content := strings.TrimSpace(msg.Content)
		if content == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s: %s", msg.Role, truncateRunes(content, 180)))
	}
	if len(lines) == 0 {
		return "", nil
	}
	summary := "Earlier conversation highlights:\n" + strings.Join(lines, "\n")
	summary = truncateRunes(summary, 3000)
	last := covered[len(covered)-1].ID
	return summary, &last
}

func truncateRunes(s string, max int) string {
	if max <= 0 || utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}

func paginationResult(total int64, params pagination.PaginationParams) *pagination.PaginationResult {
	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pages < 1 {
		pages = 1
	}
	return &pagination.PaginationResult{Total: total, Page: page, PageSize: pageSize, Pages: pages}
}

func EncodeAIStudioMetadata(metadata map[string]any) string {
	if len(metadata) == 0 {
		return "{}"
	}
	data, err := json.Marshal(metadata)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func DecodeAIStudioMetadata(raw string) map[string]any {
	if strings.TrimSpace(raw) == "" {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil || out == nil {
		return map[string]any{}
	}
	return out
}

func sanitizeAttachmentFilename(filename string) string {
	filename = strings.TrimSpace(filepath.Base(filename))
	if filename == "." || filename == "/" || filename == "" {
		return "attachment"
	}
	return truncateRunes(filename, 160)
}

func normalizeContentType(contentType string, data []byte) string {
	contentType = strings.TrimSpace(contentType)
	if contentType == "" && len(data) > 0 {
		contentType = http.DetectContentType(data)
	}
	if contentType == "" {
		return "application/octet-stream"
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil && mediaType != "" {
		return mediaType
	}
	return contentType
}

func uniqueInt64s(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:])
}

func randomSuffix() string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b[:])
}

func isTextContentType(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	return strings.HasPrefix(contentType, "text/") ||
		strings.Contains(contentType, "json") ||
		strings.Contains(contentType, "xml") ||
		strings.Contains(contentType, "yaml") ||
		strings.Contains(contentType, "javascript") ||
		strings.Contains(contentType, "typescript")
}

func readTextAttachment(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 64<<10))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fileDataURL(path, contentType string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	contentType = normalizeContentType(contentType, data)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data)), nil
}
