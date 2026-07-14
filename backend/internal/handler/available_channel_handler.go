package handler

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AvailableChannelHandler 处理用户侧「可用渠道」查询。
//
// 用户侧接口委托 ChannelService.ListAvailable，并在返回前做三层过滤：
//  1. 行过滤：只保留状态为 Active 且与当前用户可访问分组有交集的渠道；
//  2. 分组过滤：渠道的 Groups 只保留用户可访问的那些；
//  3. 平台过滤：渠道的 SupportedModels 只保留平台在用户可见 Groups 中出现过的模型，
//     防止"渠道同时挂在 antigravity / anthropic 两个平台的分组上，用户只访问
//     antigravity，却看到 anthropic 模型"这类跨平台信息泄漏；
//  4. 字段白名单：仅返回用户需要的字段（省略 BillingModelSource / RestrictModels
//     / 内部 ID / Status 等管理字段）。
type AvailableChannelHandler struct {
	channelService  *service.ChannelService
	apiKeyService   *service.APIKeyService
	settingService  *service.SettingService
	groupService    *service.GroupService
	pricingResolver *service.ModelPricingResolver
}

// NewAvailableChannelHandler 创建用户侧可用渠道 handler。
func NewAvailableChannelHandler(
	channelService *service.ChannelService,
	apiKeyService *service.APIKeyService,
	settingService *service.SettingService,
	groupService *service.GroupService,
	pricingResolver *service.ModelPricingResolver,
) *AvailableChannelHandler {
	return &AvailableChannelHandler{
		channelService:  channelService,
		apiKeyService:   apiKeyService,
		settingService:  settingService,
		groupService:    groupService,
		pricingResolver: pricingResolver,
	}
}

// featureEnabled 返回 available-channels 开关是否启用。默认关闭（opt-in）。
func (h *AvailableChannelHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return false
	}
	return h.settingService.GetAvailableChannelsRuntime(c.Request.Context()).Enabled
}

// userAvailableGroup 用户可见的分组概要（白名单字段）。
//
// 前端据此区分专属 vs 公开分组（IsExclusive）、订阅 vs 标准分组（SubscriptionType，
// 订阅视觉加深），并用 RateMultiplier 作为默认倍率；用户专属倍率前端走
// /groups/rates，和 API 密钥页面保持一致。
type userAvailableGroup struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Platform         string  `json:"platform"`
	SubscriptionType string  `json:"subscription_type"`
	RateMultiplier   float64 `json:"rate_multiplier"`
	IsExclusive      bool    `json:"is_exclusive"`
}

// userSupportedModelPricing 用户可见的定价字段白名单。
type userSupportedModelPricing struct {
	BillingMode      string                   `json:"billing_mode"`
	InputPrice       *float64                 `json:"input_price"`
	OutputPrice      *float64                 `json:"output_price"`
	CacheWritePrice  *float64                 `json:"cache_write_price"`
	CacheReadPrice   *float64                 `json:"cache_read_price"`
	ImageOutputPrice *float64                 `json:"image_output_price"`
	PerRequestPrice  *float64                 `json:"per_request_price"`
	Intervals        []userPricingIntervalDTO `json:"intervals"`
}

// userPricingIntervalDTO 定价区间白名单（去掉内部 ID、SortOrder 等前端不渲染的字段）。
type userPricingIntervalDTO struct {
	MinTokens       int      `json:"min_tokens"`
	MaxTokens       *int     `json:"max_tokens"`
	TierLabel       string   `json:"tier_label,omitempty"`
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	CacheWritePrice *float64 `json:"cache_write_price"`
	CacheReadPrice  *float64 `json:"cache_read_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
}

// userSupportedModel 用户可见的支持模型条目。
type userSupportedModel struct {
	Name     string                     `json:"name"`
	Platform string                     `json:"platform"`
	Pricing  *userSupportedModelPricing `json:"pricing"`
}

// userChannelPlatformSection 单渠道内某个平台的子视图：用户可见的分组 + 该平台
// 支持的模型。按 platform 聚合后让前端可以把渠道名作为 row-group 一次渲染，
// 后面的平台行按 sections 顺序铺开。
type userChannelPlatformSection struct {
	Platform        string               `json:"platform"`
	Groups          []userAvailableGroup `json:"groups"`
	SupportedModels []userSupportedModel `json:"supported_models"`
}

// userAvailableChannel 用户可见的渠道条目（白名单字段）。
//
// 每个渠道聚合为一条记录，内嵌 platforms 子数组：每个 section 对应一个平台，
// 包含该平台的 groups 和 supported_models。
type userAvailableChannel struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Platforms   []userChannelPlatformSection `json:"platforms"`
}

// publicModelCatalogResponse 是公开模型广场使用的只读目录。
type publicModelCatalogResponse struct {
	Groups []publicModelCatalogGroup `json:"groups"`
}

// publicModelCatalogGroup 只暴露标准（余额）分组及该分组下可用模型。
type publicModelCatalogGroup struct {
	ID             int64                    `json:"id"`
	Name           string                   `json:"name"`
	Platform       string                   `json:"platform"`
	RateMultiplier float64                  `json:"rate_multiplier"`
	Models         []publicModelCatalogItem `json:"models"`
}

// publicModelCatalogItem 是公开模型广场的模型条目。
type publicModelCatalogItem struct {
	Name            string                     `json:"name"`
	Platform        string                     `json:"platform"`
	Pricing         *userSupportedModelPricing `json:"pricing"`
	OfficialPricing *userSupportedModelPricing `json:"official_pricing"`
}

// List 列出当前用户可见的「可用渠道」。
// GET /api/v1/channels/available
func (h *AvailableChannelHandler) List(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Feature 未启用时返回空数组（不暴露渠道信息）。检查放在认证之后，
	// 保持与未开关前的 401 行为一致：未登录先 401，登录后再按开关决定。
	if !h.featureEnabled(c) {
		response.Success(c, []userAvailableChannel{})
		return
	}

	userGroups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	allowedGroupIDs := make(map[int64]struct{}, len(userGroups))
	for i := range userGroups {
		allowedGroupIDs[userGroups[i].ID] = struct{}{}
	}

	channels, err := h.channelService.ListAvailable(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]userAvailableChannel, 0, len(channels))
	for _, ch := range channels {
		if ch.Status != service.StatusActive {
			continue
		}
		visibleGroups := filterUserVisibleGroups(ch.Groups, allowedGroupIDs)
		if len(visibleGroups) == 0 {
			continue
		}
		sections := buildPlatformSections(ch, visibleGroups)
		if len(sections) == 0 {
			continue
		}
		out = append(out, userAvailableChannel{
			Name:        ch.Name,
			Description: ch.Description,
			Platforms:   sections,
		})
	}

	response.Success(c, out)
}

// ListPublicModelCatalog 列出公开模型广场目录。
// GET /api/v1/models/catalog
//
// 公开页不具备用户上下文，因此仅展示非专属的标准（余额）分组；订阅分组和
// 专属标准分组都不暴露。模型来自分组的 models_list_config.models；其中
// enabled 仅控制 /v1/models 是否强制使用自定义列表，不作为模型广场可见性开关。
func (h *AvailableChannelHandler) ListPublicModelCatalog(c *gin.Context) {
	if h.groupService == nil {
		response.InternalError(c, "Group service is not configured")
		return
	}

	groups, err := h.groupService.ListActive(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, publicModelCatalogResponse{
		Groups: buildPublicModelCatalog(c.Request.Context(), groups, h.pricingResolver),
	})
}

func buildPublicModelCatalog(
	ctx context.Context,
	groups []service.Group,
	pricingResolver *service.ModelPricingResolver,
) []publicModelCatalogGroup {
	out := make([]publicModelCatalogGroup, 0, len(groups))
	for i := range groups {
		group := groups[i]
		if !isPublicStandardCatalogGroup(group) {
			continue
		}
		models := buildPublicModelCatalogItems(ctx, group, pricingResolver)
		if len(models) == 0 {
			continue
		}
		out = append(out, publicModelCatalogGroup{
			ID:             group.ID,
			Name:           group.Name,
			Platform:       group.Platform,
			RateMultiplier: group.RateMultiplier,
			Models:         models,
		})
	}
	return out
}

func isPublicStandardCatalogGroup(group service.Group) bool {
	status := strings.TrimSpace(group.Status)
	subscriptionType := strings.TrimSpace(group.SubscriptionType)
	return group.ID > 0 &&
		strings.TrimSpace(group.Name) != "" &&
		strings.TrimSpace(group.Platform) != "" &&
		!group.IsExclusive &&
		(status == "" || status == service.StatusActive) &&
		(subscriptionType == "" || subscriptionType == service.SubscriptionTypeStandard)
}

func buildPublicModelCatalogItems(
	ctx context.Context,
	group service.Group,
	pricingResolver *service.ModelPricingResolver,
) []publicModelCatalogItem {
	models := group.ModelsListConfig.Models
	if len(models) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(models))
	out := make([]publicModelCatalogItem, 0, len(models))
	for _, raw := range models {
		model := strings.TrimSpace(raw)
		if model == "" {
			continue
		}
		key := strings.ToLower(model)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		var pricing *userSupportedModelPricing
		var officialPricing *userSupportedModelPricing
		if pricingResolver != nil {
			groupID := group.ID
			pricing = toPublicResolvedPricing(pricingResolver.Resolve(ctx, service.PricingInput{
				Model:   model,
				GroupID: &groupID,
			}))
			officialPricing = toPublicResolvedPricing(pricingResolver.Resolve(ctx, service.PricingInput{
				Model: model,
			}))
		}

		out = append(out, publicModelCatalogItem{
			Name:            model,
			Platform:        group.Platform,
			Pricing:         pricing,
			OfficialPricing: officialPricing,
		})
	}
	return out
}

func toPublicResolvedPricing(resolved *service.ResolvedPricing) *userSupportedModelPricing {
	if resolved == nil {
		return nil
	}

	mode := resolved.Mode
	if mode == "" {
		mode = service.BillingModeToken
	}
	dto := &userSupportedModelPricing{
		BillingMode: string(mode),
	}

	if pricing := resolved.BasePricing; pricing != nil {
		dto.InputPrice = publicPricePtr(pricing.InputPricePerToken)
		dto.OutputPrice = publicPricePtr(pricing.OutputPricePerToken)
		dto.CacheWritePrice = publicPricePtr(pricing.CacheCreationPricePerToken)
		dto.CacheReadPrice = publicPricePtr(pricing.CacheReadPricePerToken)
		dto.ImageOutputPrice = publicPricePtr(pricing.ImageOutputPricePerToken)
	}

	switch mode {
	case service.BillingModePerRequest, service.BillingModeImage:
		dto.PerRequestPrice = publicPricePtr(resolved.DefaultPerRequestPrice)
		dto.Intervals = toUserPricingIntervals(resolved.RequestTiers)
	default:
		dto.Intervals = toUserPricingIntervals(resolved.Intervals)
	}

	if !hasUserPricingValues(dto) {
		return nil
	}
	return dto
}

func publicPricePtr(value float64) *float64 {
	if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return nil
	}
	return &value
}

func hasUserPricingValues(pricing *userSupportedModelPricing) bool {
	return pricing != nil && (pricing.InputPrice != nil ||
		pricing.OutputPrice != nil ||
		pricing.CacheWritePrice != nil ||
		pricing.CacheReadPrice != nil ||
		pricing.ImageOutputPrice != nil ||
		pricing.PerRequestPrice != nil ||
		len(pricing.Intervals) > 0)
}

// buildPlatformSections 把一个渠道按 visibleGroups 的平台集合拆成有序的 section 列表：
// 每个 section 对应一个平台，只包含该平台的 groups 和 supported_models。
// 输出按 platform 字母序稳定排序，便于前端等效比较与回归测试。
func buildPlatformSections(
	ch service.AvailableChannel,
	visibleGroups []userAvailableGroup,
) []userChannelPlatformSection {
	groupsByPlatform := make(map[string][]userAvailableGroup, 4)
	for _, g := range visibleGroups {
		if g.Platform == "" {
			continue
		}
		groupsByPlatform[g.Platform] = append(groupsByPlatform[g.Platform], g)
	}
	if len(groupsByPlatform) == 0 {
		return nil
	}

	platforms := make([]string, 0, len(groupsByPlatform))
	for p := range groupsByPlatform {
		platforms = append(platforms, p)
	}
	sort.Strings(platforms)

	sections := make([]userChannelPlatformSection, 0, len(platforms))
	for _, platform := range platforms {
		platformSet := map[string]struct{}{platform: {}}
		sections = append(sections, userChannelPlatformSection{
			Platform:        platform,
			Groups:          groupsByPlatform[platform],
			SupportedModels: toUserSupportedModels(ch.SupportedModels, platformSet),
		})
	}
	return sections
}

// filterUserVisibleGroups 仅保留用户可访问的分组。
func filterUserVisibleGroups(
	groups []service.AvailableGroupRef,
	allowed map[int64]struct{},
) []userAvailableGroup {
	visible := make([]userAvailableGroup, 0, len(groups))
	for _, g := range groups {
		if _, ok := allowed[g.ID]; !ok {
			continue
		}
		visible = append(visible, userAvailableGroup{
			ID:               g.ID,
			Name:             g.Name,
			Platform:         g.Platform,
			SubscriptionType: g.SubscriptionType,
			RateMultiplier:   g.RateMultiplier,
			IsExclusive:      g.IsExclusive,
		})
	}
	return visible
}

// toUserSupportedModels 将 service 层支持模型转换为用户 DTO（字段白名单）。
// 仅保留平台在 allowedPlatforms 中的条目，防止跨平台模型信息泄漏。
// allowedPlatforms 为 nil 时不做平台过滤（保留全部，供测试或明确无过滤场景使用）。
func toUserSupportedModels(
	src []service.SupportedModel,
	allowedPlatforms map[string]struct{},
) []userSupportedModel {
	out := make([]userSupportedModel, 0, len(src))
	for i := range src {
		m := src[i]
		if allowedPlatforms != nil {
			if _, ok := allowedPlatforms[m.Platform]; !ok {
				continue
			}
		}
		out = append(out, userSupportedModel{
			Name:     m.Name,
			Platform: m.Platform,
			Pricing:  toUserPricing(m.Pricing),
		})
	}
	return out
}

// toUserPricing 将 service 层定价转换为用户 DTO；入参为 nil 时返回 nil。
func toUserPricing(p *service.ChannelModelPricing) *userSupportedModelPricing {
	if p == nil {
		return nil
	}
	intervals := toUserPricingIntervals(p.Intervals)
	billingMode := string(p.BillingMode)
	if billingMode == "" {
		billingMode = string(service.BillingModeToken)
	}
	return &userSupportedModelPricing{
		BillingMode:      billingMode,
		InputPrice:       p.InputPrice,
		OutputPrice:      p.OutputPrice,
		CacheWritePrice:  p.CacheWritePrice,
		CacheReadPrice:   p.CacheReadPrice,
		ImageOutputPrice: p.ImageOutputPrice,
		PerRequestPrice:  p.PerRequestPrice,
		Intervals:        intervals,
	}
}

func toUserPricingIntervals(src []service.PricingInterval) []userPricingIntervalDTO {
	out := make([]userPricingIntervalDTO, 0, len(src))
	for _, iv := range src {
		out = append(out, userPricingIntervalDTO{
			MinTokens:       iv.MinTokens,
			MaxTokens:       iv.MaxTokens,
			TierLabel:       iv.TierLabel,
			InputPrice:      iv.InputPrice,
			OutputPrice:     iv.OutputPrice,
			CacheWritePrice: iv.CacheWritePrice,
			CacheReadPrice:  iv.CacheReadPrice,
			PerRequestPrice: iv.PerRequestPrice,
		})
	}
	return out
}
