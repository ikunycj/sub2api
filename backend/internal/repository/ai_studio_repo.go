package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type aiStudioRepository struct {
	db *sql.DB
}

func NewAIStudioRepository(db *sql.DB) service.AIStudioRepository {
	return &aiStudioRepository{db: db}
}

func (r *aiStudioRepository) ListConversations(ctx context.Context, userID int64, limit, offset int) ([]service.AIStudioConversation, int64, error) {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT id, user_id, mode, title, pinned, summary, summary_message_id, last_active_at, created_at, updated_at,
       COUNT(*) OVER() AS total_count
FROM ai_conversations
WHERE user_id = $1
ORDER BY pinned DESC, last_active_at DESC, id DESC
LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		if isMissingRelationErrorRepo(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	defer rows.Close()

	var out []service.AIStudioConversation
	var total int64
	for rows.Next() {
		var item service.AIStudioConversation
		if err := scanConversationRow(rows, &item, &total); err != nil {
			return nil, 0, err
		}
		out = append(out, item)
	}
	return out, total, rows.Err()
}

func (r *aiStudioRepository) GetConversation(ctx context.Context, userID, conversationID int64) (*service.AIStudioConversation, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, user_id, mode, title, pinned, summary, summary_message_id, last_active_at, created_at, updated_at
FROM ai_conversations
WHERE id = $1 AND user_id = $2`, conversationID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, service.ErrAIStudioConversationNotFound
	}
	var item service.AIStudioConversation
	if err := scanConversationOnly(rows, &item); err != nil {
		return nil, err
	}
	return &item, rows.Err()
}

func (r *aiStudioRepository) CreateConversation(ctx context.Context, conversation *service.AIStudioConversation) error {
	if conversation == nil {
		return nil
	}
	if conversation.LastActiveAt.IsZero() {
		conversation.LastActiveAt = time.Now().UTC()
	}
	if strings.TrimSpace(conversation.Title) == "" {
		conversation.Title = "新对话"
	}
	err := scanSingleRow(ctx, r.db, `
INSERT INTO ai_conversations (user_id, mode, title, pinned, summary, summary_message_id, last_active_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at`,
		[]any{conversation.UserID, conversation.Mode, conversation.Title, conversation.Pinned, conversation.Summary, conversation.SummaryMessageID, conversation.LastActiveAt},
		&conversation.ID, &conversation.CreatedAt, &conversation.UpdatedAt)
	return err
}

func (r *aiStudioRepository) UpdateConversation(ctx context.Context, conversation *service.AIStudioConversation) error {
	if conversation == nil {
		return nil
	}
	if conversation.LastActiveAt.IsZero() {
		conversation.LastActiveAt = time.Now().UTC()
	}
	res, err := r.db.ExecContext(ctx, `
UPDATE ai_conversations
SET mode = $1, title = $2, pinned = $3, summary = $4, summary_message_id = $5, last_active_at = $6, updated_at = NOW()
WHERE id = $7 AND user_id = $8`,
		conversation.Mode, conversation.Title, conversation.Pinned, conversation.Summary, conversation.SummaryMessageID, conversation.LastActiveAt, conversation.ID, conversation.UserID)
	if err != nil {
		return err
	}
	return requireAffected(res, service.ErrAIStudioConversationNotFound)
}

func (r *aiStudioRepository) DeleteConversation(ctx context.Context, userID, conversationID int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM ai_conversations WHERE id = $1 AND user_id = $2`, conversationID, userID)
	if err != nil {
		return err
	}
	return requireAffected(res, service.ErrAIStudioConversationNotFound)
}

func (r *aiStudioRepository) TouchConversation(ctx context.Context, userID, conversationID int64, at time.Time) error {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	res, err := r.db.ExecContext(ctx, `
UPDATE ai_conversations
SET last_active_at = $1, updated_at = NOW()
WHERE id = $2 AND user_id = $3`, at, conversationID, userID)
	if err != nil {
		return err
	}
	return requireAffected(res, service.ErrAIStudioConversationNotFound)
}

func (r *aiStudioRepository) ListMessages(ctx context.Context, userID, conversationID int64) ([]service.AIStudioMessage, error) {
	return r.listMessages(ctx, userID, conversationID, 0, false)
}

func (r *aiStudioRepository) ListRecentMessages(ctx context.Context, userID, conversationID int64, limit int) ([]service.AIStudioMessage, error) {
	return r.listMessages(ctx, userID, conversationID, limit, true)
}

func (r *aiStudioRepository) listMessages(ctx context.Context, userID, conversationID int64, limit int, recent bool) ([]service.AIStudioMessage, error) {
	query := `
SELECT m.id, m.conversation_id, m.user_id, m.role, m.kind, m.content, m.model, m.status, m.sequence, m.metadata::text, m.created_at, m.updated_at
FROM ai_messages m
JOIN ai_conversations c ON c.id = m.conversation_id
WHERE m.conversation_id = $1 AND m.user_id = $2 AND c.user_id = $2`
	args := []any{conversationID, userID}
	if recent && limit > 0 {
		query += ` ORDER BY m.sequence DESC LIMIT $3`
		args = append(args, limit)
	} else {
		query += ` ORDER BY m.sequence ASC`
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []service.AIStudioMessage
	for rows.Next() {
		var item service.AIStudioMessage
		if err := scanMessage(rows, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if recent {
		for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
			out[i], out[j] = out[j], out[i]
		}
	}
	return out, rows.Err()
}

func (r *aiStudioRepository) CountMessages(ctx context.Context, userID, conversationID int64) (int64, error) {
	var count int64
	err := scanSingleRow(ctx, r.db, `
SELECT COUNT(*)
FROM ai_messages m
JOIN ai_conversations c ON c.id = m.conversation_id
WHERE m.conversation_id = $1 AND m.user_id = $2 AND c.user_id = $2`,
		[]any{conversationID, userID}, &count)
	return count, err
}

func (r *aiStudioRepository) CreateMessage(ctx context.Context, message *service.AIStudioMessage) error {
	if message == nil {
		return nil
	}
	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		if message.Sequence <= 0 {
			if err := scanSingleRow(ctx, tx, `
SELECT COALESCE(MAX(sequence), 0) + 1
FROM ai_messages
WHERE conversation_id = $1 AND user_id = $2`, []any{message.ConversationID, message.UserID}, &message.Sequence); err != nil {
				return err
			}
		}
		err := scanSingleRow(ctx, tx, `
INSERT INTO ai_messages (conversation_id, user_id, role, kind, content, model, status, sequence, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb)
RETURNING id, created_at, updated_at`,
			[]any{
				message.ConversationID,
				message.UserID,
				message.Role,
				message.Kind,
				message.Content,
				message.Model,
				defaultString(message.Status, "completed"),
				message.Sequence,
				service.EncodeAIStudioMetadata(message.Metadata),
			},
			&message.ID, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
UPDATE ai_conversations
SET last_active_at = NOW(), updated_at = NOW()
WHERE id = $1 AND user_id = $2`, message.ConversationID, message.UserID)
		return err
	})
}

func (r *aiStudioRepository) UpdateUserMessageAndDeleteAfter(ctx context.Context, userID, messageID int64, content string, metadata map[string]any) (*service.AIStudioMessage, error) {
	var out service.AIStudioMessage
	err := withTx(ctx, r.db, func(tx *sql.Tx) error {
		err := scanSingleRow(ctx, tx, `
SELECT id, conversation_id, user_id, role, kind, content, model, status, sequence, metadata::text, created_at, updated_at
FROM ai_messages
WHERE id = $1 AND user_id = $2 AND role = 'user'
FOR UPDATE`, []any{messageID, userID}, scanMessageDest(&out)...)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return service.ErrAIStudioMessageNotFound
			}
			return err
		}
		if _, err := tx.ExecContext(ctx, `
DELETE FROM ai_messages
WHERE conversation_id = $1 AND user_id = $2 AND sequence > $3`,
			out.ConversationID, userID, out.Sequence); err != nil {
			return err
		}
		err = scanSingleRow(ctx, tx, `
UPDATE ai_messages
SET content = $1, metadata = $2::jsonb, updated_at = NOW()
WHERE id = $3 AND user_id = $4
RETURNING id, conversation_id, user_id, role, kind, content, model, status, sequence, metadata::text, created_at, updated_at`,
			[]any{content, service.EncodeAIStudioMetadata(metadata), messageID, userID}, scanMessageDest(&out)...)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
UPDATE ai_conversations
SET last_active_at = NOW(), updated_at = NOW()
WHERE id = $1 AND user_id = $2`, out.ConversationID, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *aiStudioRepository) AttachToMessage(ctx context.Context, userID int64, attachmentIDs []int64, conversationID, messageID int64) error {
	if len(attachmentIDs) == 0 {
		return nil
	}
	res, err := r.db.ExecContext(ctx, `
UPDATE ai_attachments
SET conversation_id = $1, message_id = $2
WHERE user_id = $3 AND id = ANY($4) AND expires_at > NOW()`,
		conversationID, messageID, userID, pq.Array(attachmentIDs))
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != int64(len(uniqueInt64sRepo(attachmentIDs))) {
		return service.ErrAIStudioAttachmentNotFound
	}
	return nil
}

func (r *aiStudioRepository) UpdateSummary(ctx context.Context, userID, conversationID int64, summary string, summaryMessageID *int64) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE ai_conversations
SET summary = $1, summary_message_id = $2, updated_at = NOW()
WHERE id = $3 AND user_id = $4`,
		summary, summaryMessageID, conversationID, userID)
	if err != nil {
		return err
	}
	return requireAffected(res, service.ErrAIStudioConversationNotFound)
}

func (r *aiStudioRepository) CreateAttachment(ctx context.Context, attachment *service.AIStudioAttachment) error {
	if attachment == nil {
		return nil
	}
	err := scanSingleRow(ctx, r.db, `
INSERT INTO ai_attachments (user_id, conversation_id, message_id, kind, filename, storage_path, content_type, byte_size, sha256, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, created_at`,
		[]any{
			attachment.UserID,
			attachment.ConversationID,
			attachment.MessageID,
			attachment.Kind,
			attachment.Filename,
			attachment.StoragePath,
			attachment.ContentType,
			attachment.ByteSize,
			attachment.SHA256,
			attachment.ExpiresAt,
		},
		&attachment.ID, &attachment.CreatedAt)
	return err
}

func (r *aiStudioRepository) GetAttachment(ctx context.Context, userID, attachmentID int64) (*service.AIStudioAttachment, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, user_id, conversation_id, message_id, kind, filename, storage_path, content_type, byte_size, sha256, expires_at, created_at
FROM ai_attachments
WHERE id = $1 AND user_id = $2`, attachmentID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, service.ErrAIStudioAttachmentNotFound
	}
	var out service.AIStudioAttachment
	if err := scanAttachment(rows, &out); err != nil {
		return nil, err
	}
	return &out, rows.Err()
}

func (r *aiStudioRepository) ListAttachmentsByIDs(ctx context.Context, userID int64, ids []int64) ([]service.AIStudioAttachment, error) {
	ids = uniqueInt64sRepo(ids)
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT id, user_id, conversation_id, message_id, kind, filename, storage_path, content_type, byte_size, sha256, expires_at, created_at
FROM ai_attachments
WHERE user_id = $1 AND id = ANY($2)
ORDER BY id`, userID, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]service.AIStudioAttachment, 0)
	for rows.Next() {
		var item service.AIStudioAttachment
		if err := scanAttachment(rows, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *aiStudioRepository) ListAttachmentsByMessageIDs(ctx context.Context, userID int64, messageIDs []int64) (map[int64][]service.AIStudioAttachment, error) {
	out := map[int64][]service.AIStudioAttachment{}
	messageIDs = uniqueInt64sRepo(messageIDs)
	if len(messageIDs) == 0 {
		return out, nil
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT id, user_id, conversation_id, message_id, kind, filename, storage_path, content_type, byte_size, sha256, expires_at, created_at
FROM ai_attachments
WHERE user_id = $1 AND message_id = ANY($2)
ORDER BY id`, userID, pq.Array(messageIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item service.AIStudioAttachment
		if err := scanAttachment(rows, &item); err != nil {
			return nil, err
		}
		if item.MessageID != nil {
			out[*item.MessageID] = append(out[*item.MessageID], item)
		}
	}
	return out, rows.Err()
}

func (r *aiStudioRepository) DeleteExpiredConversations(ctx context.Context, cutoff time.Time, batchSize int) (int64, error) {
	if batchSize <= 0 {
		batchSize = 1000
	}
	var total int64
	for {
		res, err := r.db.ExecContext(ctx, `
WITH batch AS (
  SELECT id FROM ai_conversations
  WHERE last_active_at < $1
  ORDER BY id
  LIMIT $2
)
DELETE FROM ai_conversations
WHERE id IN (SELECT id FROM batch)`, cutoff, batchSize)
		if err != nil {
			if isMissingRelationErrorRepo(err) {
				return total, nil
			}
			return total, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return total, err
		}
		total += n
		if n == 0 {
			return total, nil
		}
	}
}

func (r *aiStudioRepository) ListExpiredAttachments(ctx context.Context, now time.Time, limit int) ([]service.AIStudioExpiredAttachment, error) {
	if limit <= 0 {
		limit = 1000
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT id, storage_path
FROM ai_attachments
WHERE expires_at < $1
ORDER BY id
LIMIT $2`, now, limit)
	if err != nil {
		if isMissingRelationErrorRepo(err) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	out := make([]service.AIStudioExpiredAttachment, 0)
	for rows.Next() {
		var item service.AIStudioExpiredAttachment
		if err := rows.Scan(&item.ID, &item.StoragePath); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *aiStudioRepository) DeleteAttachmentsByIDs(ctx context.Context, ids []int64) (int64, error) {
	ids = uniqueInt64sRepo(ids)
	if len(ids) == 0 {
		return 0, nil
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM ai_attachments WHERE id = ANY($1)`, pq.Array(ids))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type conversationScanner interface {
	Scan(dest ...any) error
}

func scanConversationOnly(scanner conversationScanner, item *service.AIStudioConversation) error {
	return scanner.Scan(
		&item.ID,
		&item.UserID,
		&item.Mode,
		&item.Title,
		&item.Pinned,
		&item.Summary,
		&item.SummaryMessageID,
		&item.LastActiveAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
}

func scanConversationRow(scanner conversationScanner, item *service.AIStudioConversation, total *int64) error {
	return scanner.Scan(
		&item.ID,
		&item.UserID,
		&item.Mode,
		&item.Title,
		&item.Pinned,
		&item.Summary,
		&item.SummaryMessageID,
		&item.LastActiveAt,
		&item.CreatedAt,
		&item.UpdatedAt,
		total,
	)
}

func scanMessage(rows *sql.Rows, item *service.AIStudioMessage) error {
	return rows.Scan(scanMessageDest(item)...)
}

func scanMessageDest(item *service.AIStudioMessage) []any {
	var metadataRaw string
	return []any{
		&item.ID,
		&item.ConversationID,
		&item.UserID,
		&item.Role,
		&item.Kind,
		&item.Content,
		&item.Model,
		&item.Status,
		&item.Sequence,
		metadataScanner{target: &item.Metadata, raw: &metadataRaw},
		&item.CreatedAt,
		&item.UpdatedAt,
	}
}

type metadataScanner struct {
	target *map[string]any
	raw    *string
}

func (m metadataScanner) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*m.target = map[string]any{}
	case string:
		*m.target = service.DecodeAIStudioMetadata(v)
	case []byte:
		*m.target = service.DecodeAIStudioMetadata(string(v))
	default:
		*m.target = map[string]any{}
	}
	return nil
}

func scanAttachment(rows *sql.Rows, item *service.AIStudioAttachment) error {
	return rows.Scan(
		&item.ID,
		&item.UserID,
		&item.ConversationID,
		&item.MessageID,
		&item.Kind,
		&item.Filename,
		&item.StoragePath,
		&item.ContentType,
		&item.ByteSize,
		&item.SHA256,
		&item.ExpiresAt,
		&item.CreatedAt,
	)
}

func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func requireAffected(res sql.Result, missing error) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return missing
	}
	return nil
}

func uniqueInt64sRepo(ids []int64) []int64 {
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

func defaultString(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func isMissingRelationErrorRepo(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "does not exist") && strings.Contains(s, "relation")
}
