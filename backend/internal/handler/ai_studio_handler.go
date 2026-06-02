package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type AIStudioHandler struct {
	studio              *service.AIStudioService
	gatewayService      *service.GatewayService
	openaiGateway       *OpenAIGatewayHandler
	subscriptionService *service.SubscriptionService
	cfg                 *config.Config
}

func NewAIStudioHandler(
	studio *service.AIStudioService,
	gatewayService *service.GatewayService,
	openaiGateway *OpenAIGatewayHandler,
	subscriptionService *service.SubscriptionService,
	cfg *config.Config,
) *AIStudioHandler {
	return &AIStudioHandler{
		studio:              studio,
		gatewayService:      gatewayService,
		openaiGateway:       openaiGateway,
		subscriptionService: subscriptionService,
		cfg:                 cfg,
	}
}

type aiStudioKeyOption struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	GroupID              int64      `json:"group_id"`
	GroupName            string     `json:"group_name"`
	Platform             string     `json:"platform"`
	Status               string     `json:"status"`
	Available            bool       `json:"available"`
	UnavailableReason    string     `json:"unavailable_reason"`
	Models               []string   `json:"models"`
	ChatModels           []string   `json:"chat_models"`
	ImageModels          []string   `json:"image_models"`
	AllowImageGeneration bool       `json:"allow_image_generation"`
	LastUsedAt           *time.Time `json:"last_used_at"`
}

func (h *AIStudioHandler) ListKeyOptions(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	options, err := h.studio.ListOpenAIKeyOptions(c.Request.Context(), userID)
	if response.ErrorFrom(c, err) {
		return
	}
	items := make([]aiStudioKeyOption, 0, len(options))
	for i := range options {
		item := options[i]
		models := []string{}
		if item.Available {
			models = normalizeAIStudioModels(h.availableModels(c.Request.Context(), item.GroupID))
		}
		items = append(items, aiStudioKeyOption{
			ID:                   item.ID,
			Name:                 item.Name,
			GroupID:              item.GroupID,
			GroupName:            item.GroupName,
			Platform:             item.Platform,
			Status:               item.Status,
			Available:            item.Available,
			UnavailableReason:    item.UnavailableReason,
			Models:               models,
			ChatModels:           chatCapableModels(models),
			ImageModels:          imageCapableModels(models, item.AllowImageGeneration),
			AllowImageGeneration: item.AllowImageGeneration,
			LastUsedAt:           item.LastUsedAt,
		})
	}
	response.Success(c, gin.H{"items": items})
}

func (h *AIStudioHandler) ListConversations(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	page, pageSize := response.ParsePagination(c)
	items, pag, err := h.studio.ListConversations(c.Request.Context(), userID, pagination.PaginationParams{Page: page, PageSize: pageSize})
	if response.ErrorFrom(c, err) {
		return
	}
	response.Paginated(c, mapConversations(items), pag.Total, pag.Page, pag.PageSize)
}

func (h *AIStudioHandler) GetConversation(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid conversation id")
		return
	}
	conversation, messages, err := h.studio.GetConversationWithMessages(c.Request.Context(), userID, id)
	if response.ErrorFrom(c, err) {
		return
	}
	response.Success(c, gin.H{
		"conversation": mapConversation(*conversation),
		"messages":     mapMessages(messages),
	})
}

func (h *AIStudioHandler) UpdateConversation(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid conversation id")
		return
	}
	var req struct {
		Title  *string `json:"title"`
		Pinned *bool   `json:"pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	conversation, err := h.studio.UpdateConversation(c.Request.Context(), userID, id, req.Title, req.Pinned)
	if response.ErrorFrom(c, err) {
		return
	}
	response.Success(c, mapConversation(*conversation))
}

func (h *AIStudioHandler) DeleteConversation(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid conversation id")
		return
	}
	if response.ErrorFrom(c, h.studio.DeleteConversation(c.Request.Context(), userID, id)) {
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *AIStudioHandler) UploadAttachment(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 25<<20))
	if err != nil {
		response.BadRequest(c, "failed to read file")
		return
	}
	contentType := header.Header.Get("Content-Type")
	attachment, err := h.studio.SaveUploadAttachment(c.Request.Context(), userID, header.Filename, contentType, data)
	if response.ErrorFrom(c, err) {
		return
	}
	response.Created(c, mapAttachment(*attachment))
}

func (h *AIStudioHandler) GetAttachmentContent(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid attachment id")
		return
	}
	attachment, err := h.studio.GetAttachmentContent(c.Request.Context(), userID, id)
	if response.ErrorFrom(c, err) {
		return
	}
	if strings.HasPrefix(attachment.StoragePath, "http://") || strings.HasPrefix(attachment.StoragePath, "https://") {
		c.Redirect(http.StatusFound, attachment.StoragePath)
		return
	}
	c.Header("Cache-Control", "private, max-age=300")
	if strings.TrimSpace(attachment.ContentType) != "" {
		c.Header("Content-Type", attachment.ContentType)
	}
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, strings.ReplaceAll(attachment.Filename, `"`, "")))
	c.File(attachment.StoragePath)
}

type aiStudioRunRequest struct {
	Mode                string                     `json:"mode"`
	ConversationID      *int64                     `json:"conversation_id"`
	Prompt              string                     `json:"prompt"`
	Model               string                     `json:"model"`
	AttachmentIDs       []int64                    `json:"attachment_ids"`
	APIKeyID            *int64                     `json:"api_key_id"`
	StoreHistory        *bool                      `json:"store_history"`
	ResendFromMessageID *int64                     `json:"resend_from_message_id"`
	SystemPrompt        string                     `json:"system_prompt"`
	Temperature         *float64                   `json:"temperature"`
	MaxTokens           int                        `json:"max_tokens"`
	ThinkingEnabled     bool                       `json:"thinking_enabled"`
	ReasoningEffort     string                     `json:"reasoning_effort"`
	ImageSize           string                     `json:"image_size"`
	ImageQuality        string                     `json:"image_quality"`
	LocalContext        []aiStudioLocalMessageJSON `json:"local_context"`
}

type aiStudioLocalMessageJSON struct {
	Role    string `json:"role"`
	Kind    string `json:"kind"`
	Content string `json:"content"`
}

func (h *AIStudioHandler) Run(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "user not found")
		return
	}
	var req aiStudioRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	storeHistory := true
	if req.StoreHistory != nil {
		storeHistory = *req.StoreHistory
	}

	apiKey, err := h.studio.SelectOpenAIAPIKey(c.Request.Context(), userID, req.APIKeyID)
	if response.ErrorFrom(c, err) {
		return
	}
	if strings.TrimSpace(req.Model) == "" {
		req.Model = h.defaultRunModelForKey(c.Request.Context(), req.Mode, apiKey)
	}
	if !h.keyAllowsModel(c.Request.Context(), apiKey, req.Mode, req.Model) {
		response.BadRequest(c, "selected API key does not allow this model")
		return
	}
	if req.Mode == service.AIStudioModeImage && !h.keyAllowsImageRun(c.Request.Context(), apiKey) {
		response.BadRequest(c, "selected API key does not allow image generation")
		return
	}
	if !h.bindGatewayIdentity(c, apiKey) {
		return
	}
	h.studio.TouchSelectedAPIKey(c.Request.Context(), apiKey.ID)

	prepared, err := h.studio.PrepareRun(c.Request.Context(), service.AIStudioPrepareRunInput{
		UserID:              userID,
		Mode:                req.Mode,
		ConversationID:      req.ConversationID,
		Prompt:              req.Prompt,
		Model:               req.Model,
		AttachmentIDs:       req.AttachmentIDs,
		StoreHistory:        storeHistory,
		ResendFromMessageID: req.ResendFromMessageID,
		SystemPrompt:        req.SystemPrompt,
		MaxTokens:           req.MaxTokens,
		Temperature:         req.Temperature,
		LocalContext:        mapLocalMessages(req.LocalContext),
	})
	if response.ErrorFrom(c, err) {
		return
	}

	switch req.Mode {
	case service.AIStudioModeImage:
		h.runImage(c, userID, req, prepared)
	default:
		h.runChat(c, userID, req, prepared)
	}
}

func (h *AIStudioHandler) runChat(c *gin.Context, userID int64, req aiStudioRunRequest, prepared *service.AIStudioPreparedRun) {
	chatReq := apicompat.ChatCompletionsRequest{
		Model:    req.Model,
		Stream:   true,
		Messages: mapAIStudioChatMessages(prepared.Messages),
	}
	if req.MaxTokens > 0 {
		chatReq.MaxTokens = &req.MaxTokens
	}
	if req.Temperature != nil {
		chatReq.Temperature = req.Temperature
	}
	if req.ThinkingEnabled {
		if effort := normalizeAIStudioReasoningEffort(req.ReasoningEffort); effort != "" {
			chatReq.ReasoningEffort = effort
		}
	}
	responsesReq, err := apicompat.ChatCompletionsToResponses(&chatReq)
	if err != nil {
		response.InternalError(c, "failed to build request")
		return
	}
	responsesReq.Stream = true
	if req.ThinkingEnabled && responsesReq.Reasoning != nil {
		responsesReq.Reasoning.Summary = "auto"
	}
	bodyBytes, err := json.Marshal(responsesReq)
	if err != nil {
		response.InternalError(c, "failed to build request")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	c.Request.ContentLength = int64(len(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.URL.Path = "/v1/responses"

	writer := newAIStudioCaptureWriter(c.Writer)
	if prepared != nil && prepared.Conversation != nil {
		writer.SetPrefixEvent("conversation", mustJSON(mapConversation(*prepared.Conversation)))
	}
	c.Writer = writer
	h.openaiGateway.Responses(c)

	assistantText := parseOpenAIResponsesSSEContent(writer.Body())
	if prepared != nil && prepared.Conversation != nil && strings.TrimSpace(assistantText) != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = h.studio.CompleteAssistantMessage(ctx, userID, prepared, service.AIStudioModeChat, req.Model, assistantText, nil)
		h.studio.MaybeCompactSummary(ctx, userID, prepared.Conversation.ID)
	}
}

func (h *AIStudioHandler) runImage(c *gin.Context, userID int64, req aiStudioRunRequest, prepared *service.AIStudioPreparedRun) {
	rec := httptest.NewRecorder()
	gatewayCtx, _ := gin.CreateTestContext(rec)
	gatewayCtx.Request = cloneRequestForGateway(c.Request)
	copyGinKeys(c, gatewayCtx)

	endpoint := "/v1/images/generations"
	var bodyBytes []byte
	var contentType string
	var err error
	imageAttachments := filterImageAttachments(prepared.Attachments)
	if len(imageAttachments) > 0 {
		endpoint = "/v1/images/edits"
		bodyBytes, contentType, err = h.buildImageEditMultipart(req, imageAttachments)
	} else {
		bodyBytes, err = h.buildImageGenerationJSON(req)
		contentType = "application/json"
	}
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	gatewayCtx.Request.URL.Path = endpoint
	gatewayCtx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	gatewayCtx.Request.ContentLength = int64(len(bodyBytes))
	gatewayCtx.Request.Header.Set("Content-Type", contentType)

	h.openaiGateway.Images(gatewayCtx)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Accel-Buffering", "no")
	status := rec.Code
	if status == 0 {
		status = http.StatusOK
	}
	if status >= 400 {
		c.Status(status)
		writeSSE(c.Writer, "error", string(rec.Body.Bytes()))
		return
	}

	payload, generatedIDs := h.persistGeneratedImages(c.Request.Context(), userID, prepared, rec.Body.Bytes())
	if prepared != nil && prepared.Conversation != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		content := strings.TrimSpace(gjson.GetBytes(rec.Body.Bytes(), "data.0.revised_prompt").String())
		if content == "" {
			content = "图片生成完成"
		}
		msg, _ := h.studio.CompleteAssistantMessage(ctx, userID, prepared, service.AIStudioModeImage, req.Model, content, map[string]any{"image_attachment_ids": generatedIDs})
		if msg != nil && len(generatedIDs) > 0 {
			_ = h.studio.AttachGeneratedAttachmentsToMessage(ctx, userID, generatedIDs, prepared.Conversation.ID, msg.ID)
		}
	}
	c.Status(http.StatusOK)
	writeSSE(c.Writer, "image", mustJSON(payload))
	writeSSE(c.Writer, "done", "{}")
}

func (h *AIStudioHandler) buildImageGenerationJSON(req aiStudioRunRequest) ([]byte, error) {
	body := map[string]any{
		"model":           req.Model,
		"prompt":          req.Prompt,
		"response_format": "b64_json",
	}
	if strings.TrimSpace(req.ImageSize) != "" {
		body["size"] = strings.TrimSpace(req.ImageSize)
	}
	if strings.TrimSpace(req.ImageQuality) != "" {
		body["quality"] = strings.TrimSpace(req.ImageQuality)
	}
	return json.Marshal(body)
}

func (h *AIStudioHandler) buildImageEditMultipart(req aiStudioRunRequest, attachments []service.AIStudioAttachment) ([]byte, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("model", req.Model)
	_ = writer.WriteField("prompt", req.Prompt)
	_ = writer.WriteField("response_format", "b64_json")
	if strings.TrimSpace(req.ImageSize) != "" {
		_ = writer.WriteField("size", strings.TrimSpace(req.ImageSize))
	}
	if strings.TrimSpace(req.ImageQuality) != "" {
		_ = writer.WriteField("quality", strings.TrimSpace(req.ImageQuality))
	}
	for i, attachment := range attachments {
		data, err := h.studio.ReadAttachmentBytes(&attachment, 25<<20)
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		field := "image"
		if i > 0 {
			field = fmt.Sprintf("image[%d]", i)
		}
		part, err := writer.CreateFormFile(field, attachment.Filename)
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		if _, err := part.Write(data); err != nil {
			_ = writer.Close()
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return body.Bytes(), writer.FormDataContentType(), nil
}

func (h *AIStudioHandler) persistGeneratedImages(ctx context.Context, userID int64, prepared *service.AIStudioPreparedRun, body []byte) (gin.H, []int64) {
	items := make([]gin.H, 0)
	ids := make([]int64, 0)
	for i, result := range gjson.GetBytes(body, "data").Array() {
		if b64 := strings.TrimSpace(result.Get("b64_json").String()); b64 != "" {
			data, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				continue
			}
			attachment, err := h.studio.CompleteGeneratedImage(ctx, userID, prepared, fmt.Sprintf("ai-studio-%d.png", i+1), "image/png", data)
			if err != nil || attachment == nil {
				continue
			}
			ids = append(ids, attachment.ID)
			items = append(items, mapAttachment(*attachment))
			continue
		}
		if url := strings.TrimSpace(result.Get("url").String()); url != "" {
			items = append(items, gin.H{"url": url})
		}
	}
	payload := gin.H{
		"conversation": conversationPayload(prepared),
		"attachments":  items,
	}
	if len(items) == 0 && json.Valid(body) {
		payload["raw"] = json.RawMessage(body)
	}
	return payload, ids
}

func (h *AIStudioHandler) bindGatewayIdentity(c *gin.Context, apiKey *service.APIKey) bool {
	if apiKey == nil || apiKey.User == nil {
		response.Unauthorized(c, "api key user not found")
		return false
	}
	if h.cfg != nil && h.cfg.RunMode != config.RunModeSimple && apiKey.Group != nil && apiKey.Group.IsSubscriptionType() && h.subscriptionService != nil {
		sub, err := h.subscriptionService.GetActiveSubscription(c.Request.Context(), apiKey.User.ID, apiKey.Group.ID)
		if err != nil {
			response.Forbidden(c, "No active subscription found for this group")
			return false
		}
		needsMaintenance, err := h.subscriptionService.ValidateAndCheckLimits(sub, apiKey.Group)
		if err != nil {
			response.Forbidden(c, err.Error())
			return false
		}
		if needsMaintenance {
			cp := *sub
			h.subscriptionService.DoWindowMaintenance(&cp)
		}
		c.Set(string(middleware.ContextKeySubscription), sub)
	} else if h.cfg != nil && h.cfg.RunMode != config.RunModeSimple && apiKey.User.Balance <= 0 {
		response.Forbidden(c, "Insufficient account balance")
		return false
	}
	c.Set(string(middleware.ContextKeyAPIKey), apiKey)
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.User.ID, Concurrency: apiKey.User.Concurrency})
	c.Set(string(middleware.ContextKeyUserRole), apiKey.User.Role)
	return true
}

type aiStudioCaptureWriter struct {
	gin.ResponseWriter
	buffer        bytes.Buffer
	prefixEvent   string
	prefixData    string
	prefixWritten bool
}

func newAIStudioCaptureWriter(w gin.ResponseWriter) *aiStudioCaptureWriter {
	return &aiStudioCaptureWriter{ResponseWriter: w}
}

func (w *aiStudioCaptureWriter) SetPrefixEvent(event, data string) {
	w.prefixEvent = event
	w.prefixData = data
}

func (w *aiStudioCaptureWriter) Write(data []byte) (int, error) {
	if w.buffer.Len() < 2<<20 {
		_, _ = w.buffer.Write(data)
	}
	if err := w.writePrefixIfSSE(data); err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(data)
}

func (w *aiStudioCaptureWriter) WriteString(data string) (int, error) {
	if w.buffer.Len() < 2<<20 {
		_, _ = w.buffer.WriteString(data)
	}
	if err := w.writePrefixIfSSE([]byte(data)); err != nil {
		return 0, err
	}
	return w.ResponseWriter.WriteString(data)
}

func (w *aiStudioCaptureWriter) Flush() {
	w.ResponseWriter.Flush()
}

func (w *aiStudioCaptureWriter) Body() []byte {
	return w.buffer.Bytes()
}

func (w *aiStudioCaptureWriter) writePrefixIfSSE(data []byte) error {
	if w.prefixWritten || w.prefixEvent == "" {
		return nil
	}
	trimmed := bytes.TrimSpace(data)
	if !bytes.HasPrefix(trimmed, []byte("data:")) && !bytes.HasPrefix(trimmed, []byte("event:")) {
		return nil
	}
	w.prefixWritten = true
	writeSSE(w.ResponseWriter, w.prefixEvent, w.prefixData)
	return nil
}

func mapAIStudioChatMessages(items []map[string]any) []apicompat.ChatMessage {
	out := make([]apicompat.ChatMessage, 0, len(items))
	for _, item := range items {
		role, _ := item["role"].(string)
		if strings.TrimSpace(role) == "" {
			continue
		}
		raw, err := json.Marshal(item["content"])
		if err != nil {
			continue
		}
		out = append(out, apicompat.ChatMessage{
			Role:    role,
			Content: raw,
		})
	}
	return out
}

func parseOpenAIResponsesSSEContent(body []byte) string {
	var out strings.Builder
	var completedText string
	scanner := bufio.NewScanner(bytes.NewReader(body))
	scanner.Buffer(make([]byte, 0, 64*1024), 2<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" || !gjson.Valid(payload) {
			continue
		}
		if content := gjson.Get(payload, "delta").String(); gjson.Get(payload, "type").String() == "response.output_text.delta" && content != "" {
			out.WriteString(content)
			continue
		}
		if gjson.Get(payload, "type").String() == "response.completed" {
			completedText = extractResponsesCompletedText(payload)
		}
	}
	if out.Len() == 0 {
		return completedText
	}
	return out.String()
}

func extractResponsesCompletedText(payload string) string {
	var out strings.Builder
	for _, item := range gjson.Get(payload, "response.output").Array() {
		if item.Get("type").String() != "message" {
			continue
		}
		for _, part := range item.Get("content").Array() {
			if part.Get("type").String() == "output_text" {
				out.WriteString(part.Get("text").String())
			}
		}
	}
	return out.String()
}

func firstOpenAIChatTextDelta(payload string) string {
	for _, path := range []string{
		"choices.0.delta.content",
		"choices.0.message.content",
		"delta.content",
		"message.content",
		"content",
	} {
		if value := strings.TrimSpace(gjson.Get(payload, path).String()); value != "" {
			return value
		}
	}
	return ""
}

func mapConversations(items []service.AIStudioConversation) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, item := range items {
		out = append(out, mapConversation(item))
	}
	return out
}

func mapConversation(item service.AIStudioConversation) gin.H {
	return gin.H{
		"id":                    item.ID,
		"mode":                  item.Mode,
		"title":                 item.Title,
		"pinned":                item.Pinned,
		"last_active_at":        item.LastActiveAt,
		"created_at":            item.CreatedAt,
		"updated_at":            item.UpdatedAt,
		"summary":               item.Summary,
		"summary_message_id":    item.SummaryMessageID,
		"server_retention_days": service.AIStudioDefaultHistoryRetentionDays,
	}
}

func mapMessages(items []service.AIStudioMessage) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, item := range items {
		attachments := make([]gin.H, 0, len(item.Attachments))
		for _, attachment := range item.Attachments {
			attachments = append(attachments, mapAttachment(attachment))
		}
		out = append(out, gin.H{
			"id":              item.ID,
			"conversation_id": item.ConversationID,
			"role":            item.Role,
			"kind":            item.Kind,
			"content":         item.Content,
			"model":           item.Model,
			"status":          item.Status,
			"sequence":        item.Sequence,
			"metadata":        item.Metadata,
			"created_at":      item.CreatedAt,
			"updated_at":      item.UpdatedAt,
			"attachments":     attachments,
		})
	}
	return out
}

func mapAttachment(item service.AIStudioAttachment) gin.H {
	return gin.H{
		"id":              item.ID,
		"conversation_id": item.ConversationID,
		"message_id":      item.MessageID,
		"kind":            item.Kind,
		"filename":        item.Filename,
		"content_type":    item.ContentType,
		"byte_size":       item.ByteSize,
		"expires_at":      item.ExpiresAt,
		"created_at":      item.CreatedAt,
		"url":             fmt.Sprintf("/api/v1/ai-studio/attachments/%d/content", item.ID),
	}
}

func currentUserID(c *gin.Context) (int64, bool) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok || subject.UserID <= 0 {
		return 0, false
	}
	return subject.UserID, true
}

func parseIDParam(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(c.Param(name)), 10, 64)
}

func (h *AIStudioHandler) availableModels(ctx context.Context, groupID int64) []string {
	if h == nil || h.gatewayService == nil || groupID <= 0 {
		return openai.DefaultModelIDs()
	}
	return h.gatewayService.GetAvailableModels(ctx, &groupID, service.PlatformOpenAI)
}

func normalizeAIStudioModels(models []string) []string {
	if len(models) == 0 {
		return openai.DefaultModelIDs()
	}
	out := make([]string, 0, len(models))
	seen := make(map[string]struct{}, len(models))
	for _, model := range models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		if _, ok := seen[model]; ok {
			continue
		}
		seen[model] = struct{}{}
		out = append(out, model)
	}
	if len(out) == 0 {
		return openai.DefaultModelIDs()
	}
	return out
}

func chatCapableModels(models []string) []string {
	out := make([]string, 0, len(models))
	for _, model := range models {
		if !isImageModel(model) {
			out = append(out, model)
		}
	}
	return out
}

func imageCapableModels(models []string, allowImageGeneration bool) []string {
	if !allowImageGeneration {
		return []string{}
	}
	out := make([]string, 0)
	for _, model := range models {
		if isImageModel(model) {
			out = append(out, model)
		}
	}
	return out
}

func normalizeAIStudioReasoningEffort(effort string) string {
	value := strings.ToLower(strings.TrimSpace(effort))
	value = strings.NewReplacer("-", "", "_", "", " ", "").Replace(value)
	switch value {
	case "low", "medium", "high":
		return value
	case "xhigh", "extrahigh":
		return "xhigh"
	default:
		return "medium"
	}
}

func isImageModel(model string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	return strings.HasPrefix(model, "gpt-image-") || strings.Contains(model, "image")
}

func (h *AIStudioHandler) keyAllowsModel(ctx context.Context, apiKey *service.APIKey, mode, model string) bool {
	model = strings.TrimSpace(model)
	if apiKey == nil || apiKey.GroupID == nil || model == "" {
		return false
	}
	models := normalizeAIStudioModels(h.availableModels(ctx, *apiKey.GroupID))
	if mode == service.AIStudioModeImage {
		for _, allowed := range imageCapableModels(models, apiKey.Group != nil && apiKey.Group.AllowImageGeneration) {
			if allowed == model {
				return true
			}
		}
		return false
	}
	for _, allowed := range chatCapableModels(models) {
		if allowed == model {
			return true
		}
	}
	return false
}

func (h *AIStudioHandler) keyAllowsImageRun(ctx context.Context, apiKey *service.APIKey) bool {
	if apiKey == nil || apiKey.Group == nil || apiKey.GroupID == nil || !apiKey.Group.AllowImageGeneration {
		return false
	}
	models := normalizeAIStudioModels(h.availableModels(ctx, *apiKey.GroupID))
	return len(imageCapableModels(models, true)) > 0
}

func (h *AIStudioHandler) defaultRunModelForKey(ctx context.Context, mode string, apiKey *service.APIKey) string {
	groupID := int64(0)
	if apiKey != nil && apiKey.GroupID != nil {
		groupID = *apiKey.GroupID
	}
	models := normalizeAIStudioModels(h.availableModels(ctx, groupID))
	if mode == service.AIStudioModeImage {
		for _, model := range imageCapableModels(models, apiKey != nil && apiKey.Group != nil && apiKey.Group.AllowImageGeneration) {
			return model
		}
	}
	for _, model := range chatCapableModels(models) {
		return model
	}
	if len(models) > 0 {
		return models[0]
	}
	return ""
}

func mapLocalMessages(items []aiStudioLocalMessageJSON) []service.AIStudioMessage {
	out := make([]service.AIStudioMessage, 0, len(items))
	for _, item := range items {
		role := strings.TrimSpace(item.Role)
		if role != service.AIStudioRoleUser && role != service.AIStudioRoleAssistant && role != service.AIStudioRoleSystem {
			continue
		}
		out = append(out, service.AIStudioMessage{
			Role:    role,
			Kind:    defaultString(strings.TrimSpace(item.Kind), service.AIStudioModeChat),
			Content: item.Content,
		})
	}
	return out
}

func filterImageAttachments(items []service.AIStudioAttachment) []service.AIStudioAttachment {
	out := make([]service.AIStudioAttachment, 0, len(items))
	for _, item := range items {
		if strings.HasPrefix(strings.ToLower(item.ContentType), "image/") {
			out = append(out, item)
		}
	}
	return out
}

func cloneRequestForGateway(req *http.Request) *http.Request {
	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	return clone
}

func copyGinKeys(src, dst *gin.Context) {
	for _, key := range []middleware.ContextKey{
		middleware.ContextKeyAPIKey,
		middleware.ContextKeyUser,
		middleware.ContextKeyUserRole,
		middleware.ContextKeySubscription,
	} {
		if value, ok := src.Get(string(key)); ok {
			dst.Set(string(key), value)
		}
	}
}

func writeSSE(w io.Writer, event, data string) {
	if event != "" {
		_, _ = fmt.Fprintf(w, "event: %s\n", event)
	}
	for _, line := range strings.Split(data, "\n") {
		_, _ = fmt.Fprintf(w, "data: %s\n", line)
	}
	_, _ = io.WriteString(w, "\n")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func mustJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func conversationPayload(prepared *service.AIStudioPreparedRun) any {
	if prepared == nil || prepared.Conversation == nil {
		return nil
	}
	return mapConversation(*prepared.Conversation)
}

func defaultString(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
