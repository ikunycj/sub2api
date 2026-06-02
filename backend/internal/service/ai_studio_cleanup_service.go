package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

const (
	aiStudioCleanupLeaderLockKey = "ai_studio:cleanup:leader"
	aiStudioCleanupSchedule      = "0 3 * * *"
	aiStudioCleanupBatchSize     = 1000
	aiStudioCleanupRunTimeout    = 20 * time.Minute
	aiStudioCleanupLockTTL       = 30 * time.Minute
)

type AIStudioCleanupService struct {
	repo        AIStudioRepository
	db          *sql.DB
	redisClient *redis.Client
	cfg         *config.Config
	instanceID  string

	mu      sync.Mutex
	cron    *cron.Cron
	started bool
	stopped bool
}

func NewAIStudioCleanupService(repo AIStudioRepository, db *sql.DB, redisClient *redis.Client, cfg *config.Config) *AIStudioCleanupService {
	return &AIStudioCleanupService{
		repo:        repo,
		db:          db,
		redisClient: redisClient,
		cfg:         cfg,
		instanceID:  uuid.NewString(),
	}
}

func (s *AIStudioCleanupService) Start() {
	if s == nil || s.repo == nil || s.db == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started || s.stopped {
		return
	}
	loc := time.Local
	if s.cfg != nil && strings.TrimSpace(s.cfg.Timezone) != "" {
		if parsed, err := time.LoadLocation(strings.TrimSpace(s.cfg.Timezone)); err == nil && parsed != nil {
			loc = parsed
		}
	}
	c := cron.New(cron.WithParser(opsCleanupCronParser), cron.WithLocation(loc))
	if _, err := c.AddFunc(aiStudioCleanupSchedule, func() { s.runScheduled() }); err != nil {
		logger.LegacyPrintf("service.ai_studio_cleanup", "[AIStudioCleanup] not started: %v", err)
		return
	}
	c.Start()
	s.cron = c
	s.started = true
	logger.LegacyPrintf("service.ai_studio_cleanup", "[AIStudioCleanup] scheduled (schedule=%q tz=%s)", aiStudioCleanupSchedule, loc.String())
}

func (s *AIStudioCleanupService) Stop() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return
	}
	s.stopped = true
	if s.cron == nil {
		return
	}
	ctx := s.cron.Stop()
	select {
	case <-ctx.Done():
	case <-time.After(opsCleanupCronStopTimeout):
		logger.LegacyPrintf("service.ai_studio_cleanup", "[AIStudioCleanup] cron stop timed out")
	}
	s.cron = nil
}

func (s *AIStudioCleanupService) runScheduled() {
	if s == nil || s.repo == nil || s.db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), aiStudioCleanupRunTimeout)
	defer cancel()

	release, ok := s.tryAcquireLeaderLock(ctx)
	if !ok {
		return
	}
	if release != nil {
		defer release()
	}

	counts, err := s.RunOnce(ctx)
	if err != nil {
		logger.LegacyPrintf("service.ai_studio_cleanup", "[AIStudioCleanup] failed: %v", err)
		return
	}
	logger.LegacyPrintf("service.ai_studio_cleanup", "[AIStudioCleanup] complete: %s", counts)
}

type AIStudioCleanupCounts struct {
	Conversations int64
	Attachments   int64
	Files         int64
}

func (c AIStudioCleanupCounts) String() string {
	return fmt.Sprintf("conversations=%d attachments=%d files=%d", c.Conversations, c.Attachments, c.Files)
}

func (s *AIStudioCleanupService) RunOnce(ctx context.Context) (AIStudioCleanupCounts, error) {
	out := AIStudioCleanupCounts{}
	if s == nil || s.repo == nil {
		return out, nil
	}
	now := time.Now().UTC()
	conversationCutoff := now.AddDate(0, 0, -AIStudioDefaultHistoryRetentionDays)
	n, err := s.repo.DeleteExpiredConversations(ctx, conversationCutoff, aiStudioCleanupBatchSize)
	if err != nil {
		return out, err
	}
	out.Conversations = n

	for {
		expired, err := s.repo.ListExpiredAttachments(ctx, now, aiStudioCleanupBatchSize)
		if err != nil {
			return out, err
		}
		if len(expired) == 0 {
			break
		}
		ids := make([]int64, 0, len(expired))
		for _, item := range expired {
			ids = append(ids, item.ID)
			if strings.HasPrefix(item.StoragePath, "http://") || strings.HasPrefix(item.StoragePath, "https://") {
				continue
			}
			if err := os.Remove(item.StoragePath); err == nil {
				out.Files++
			}
		}
		deleted, err := s.repo.DeleteAttachmentsByIDs(ctx, ids)
		if err != nil {
			return out, err
		}
		out.Attachments += deleted
		if len(expired) < aiStudioCleanupBatchSize {
			break
		}
	}
	return out, nil
}

func (s *AIStudioCleanupService) tryAcquireLeaderLock(ctx context.Context) (func(), bool) {
	if s == nil {
		return nil, false
	}
	if s.cfg != nil && s.cfg.RunMode == config.RunModeSimple {
		return nil, true
	}
	if s.redisClient != nil {
		ok, err := s.redisClient.SetNX(ctx, aiStudioCleanupLeaderLockKey, s.instanceID, aiStudioCleanupLockTTL).Result()
		if err == nil {
			if !ok {
				return nil, false
			}
			return func() {
				_, _ = opsCleanupReleaseScript.Run(ctx, s.redisClient, []string{aiStudioCleanupLeaderLockKey}, s.instanceID).Result()
			}, true
		}
		logger.LegacyPrintf("service.ai_studio_cleanup", "[AIStudioCleanup] redis lock failed, fallback to DB advisory lock: %v", err)
	}
	return tryAcquireDBAdvisoryLock(ctx, s.db, hashAdvisoryLockID(aiStudioCleanupLeaderLockKey))
}

func ProvideAIStudioCleanupService(repo AIStudioRepository, db *sql.DB, redisClient *redis.Client, cfg *config.Config) *AIStudioCleanupService {
	svc := NewAIStudioCleanupService(repo, db, redisClient, cfg)
	svc.Start()
	return svc
}
