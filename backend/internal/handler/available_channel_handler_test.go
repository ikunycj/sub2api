//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUserAvailableChannel_Unauthenticated401(t *testing.T) {
	// 没有 AuthSubject 注入时，handler 应返回 401 且不触达 service 依赖。
	gin.SetMode(gin.TestMode)
	h := &AvailableChannelHandler{} // nil services — 401 路径不会调用它们
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/channels/available", nil)

	h.List(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFilterUserVisibleGroups_IntersectionOnly(t *testing.T) {
	// 渠道挂在 {g1, g2, g3}，用户只允许 {g1, g3} —— 响应必须仅含 g1/g3。
	groups := []service.AvailableGroupRef{
		{ID: 1, Name: "g1", Platform: "anthropic"},
		{ID: 2, Name: "g2", Platform: "anthropic"},
		{ID: 3, Name: "g3", Platform: "openai"},
	}
	allowed := map[int64]struct{}{1: {}, 3: {}}

	visible := filterUserVisibleGroups(groups, allowed)
	require.Len(t, visible, 2)
	ids := []int64{visible[0].ID, visible[1].ID}
	require.ElementsMatch(t, []int64{1, 3}, ids)
}

func TestToUserSupportedModels_FiltersByAllowedPlatforms(t *testing.T) {
	// 用户可访问分组只覆盖 anthropic；anthropic 平台的模型保留，openai 模型被剔除。
	src := []service.SupportedModel{
		{Name: "claude-sonnet-4-6", Platform: "anthropic", Pricing: nil},
		{Name: "gpt-4o", Platform: "openai", Pricing: nil},
	}
	allowed := map[string]struct{}{"anthropic": {}}
	out := toUserSupportedModels(src, allowed)
	require.Len(t, out, 1)
	require.Equal(t, "claude-sonnet-4-6", out[0].Name)
}

func TestToUserSupportedModels_NilAllowedPlatformsKeepsAll(t *testing.T) {
	// 显式传 nil allowedPlatforms 表示不做过滤。
	src := []service.SupportedModel{
		{Name: "a", Platform: "anthropic"},
		{Name: "b", Platform: "openai"},
	}
	require.Len(t, toUserSupportedModels(src, nil), 2)
}

func TestUserAvailableChannel_FieldWhitelist(t *testing.T) {
	// 通过序列化 userAvailableChannel 结构体验证响应形状：
	// 只有 name / description / platforms；不含管理端字段。
	row := userAvailableChannel{
		Name:        "ch",
		Description: "d",
		Platforms: []userChannelPlatformSection{
			{
				Platform:        "anthropic",
				Groups:          []userAvailableGroup{{ID: 1, Name: "g1", Platform: "anthropic"}},
				SupportedModels: []userSupportedModel{},
			},
		},
	}
	raw, err := json.Marshal(row)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))

	for _, key := range []string{"id", "status", "billing_model_source", "restrict_models"} {
		_, exists := decoded[key]
		require.Falsef(t, exists, "user DTO must not expose %q", key)
	}
	for _, key := range []string{"name", "description", "platforms"} {
		_, exists := decoded[key]
		require.Truef(t, exists, "user DTO must expose %q", key)
	}

	// 验证 section 的字段（platform / groups / supported_models）。
	rawSection, err := json.Marshal(row.Platforms[0])
	require.NoError(t, err)
	var sectionDecoded map[string]any
	require.NoError(t, json.Unmarshal(rawSection, &sectionDecoded))
	for _, key := range []string{"platform", "groups", "supported_models"} {
		_, exists := sectionDecoded[key]
		require.Truef(t, exists, "platform section must expose %q", key)
	}

	// Group DTO 暴露区分专属/公开、订阅类型、默认倍率所需的字段，
	// 前端据此渲染 GroupBadge 并与 API 密钥页保持一致的视觉。
	rawGroup, err := json.Marshal(row.Platforms[0].Groups[0])
	require.NoError(t, err)
	var groupDecoded map[string]any
	require.NoError(t, json.Unmarshal(rawGroup, &groupDecoded))
	for _, key := range []string{"id", "name", "platform", "subscription_type", "rate_multiplier", "is_exclusive"} {
		_, exists := groupDecoded[key]
		require.Truef(t, exists, "group DTO must expose %q", key)
	}

	// pricing interval 白名单：不应暴露 id / sort_order。
	pricing := toUserPricing(&service.ChannelModelPricing{
		BillingMode: service.BillingModeToken,
		Intervals: []service.PricingInterval{
			{ID: 7, MinTokens: 0, MaxTokens: nil, SortOrder: 3},
		},
	})
	require.NotNil(t, pricing)
	require.Len(t, pricing.Intervals, 1)
	rawIv, err := json.Marshal(pricing.Intervals[0])
	require.NoError(t, err)
	var ivDecoded map[string]any
	require.NoError(t, json.Unmarshal(rawIv, &ivDecoded))
	for _, key := range []string{"id", "pricing_id", "sort_order"} {
		_, exists := ivDecoded[key]
		require.Falsef(t, exists, "user pricing interval must not expose %q", key)
	}
}

func TestBuildPlatformSections_GroupsByPlatform(t *testing.T) {
	// 一个渠道横跨 anthropic / openai / 空平台：应该生成 2 个 section，
	// 按 platform 字母序排序，各自 groups 和 supported_models 只含同平台条目。
	ch := service.AvailableChannel{
		Name: "ch",
		SupportedModels: []service.SupportedModel{
			{Name: "claude-sonnet-4-6", Platform: "anthropic"},
			{Name: "gpt-4o", Platform: "openai"},
		},
	}
	visible := []userAvailableGroup{
		{ID: 1, Name: "g-openai", Platform: "openai"},
		{ID: 2, Name: "g-ant", Platform: "anthropic"},
		{ID: 3, Name: "g-empty", Platform: ""},
	}
	sections := buildPlatformSections(ch, visible)
	require.Len(t, sections, 2)
	require.Equal(t, "anthropic", sections[0].Platform)
	require.Equal(t, "openai", sections[1].Platform)
	require.Len(t, sections[0].Groups, 1)
	require.Equal(t, int64(2), sections[0].Groups[0].ID)
	require.Len(t, sections[0].SupportedModels, 1)
	require.Equal(t, "claude-sonnet-4-6", sections[0].SupportedModels[0].Name)
}

func TestBuildPublicModelCatalog_StandardPublicGroupsOnly(t *testing.T) {
	groups := []service.Group{
		{
			ID:                    1,
			Name:                  "standard-public",
			Platform:              "openai",
			Status:                service.StatusActive,
			SubscriptionType:      service.SubscriptionTypeStandard,
			ModelsListConfig:      service.GroupModelsListConfig{Enabled: true, Models: []string{" gpt-4o ", "gpt-4o", ""}},
			RateMultiplier:        1.2,
			AllowMessagesDispatch: true,
		},
		{
			ID:               2,
			Name:             "subscription",
			Platform:         "openai",
			Status:           service.StatusActive,
			SubscriptionType: service.SubscriptionTypeSubscription,
			ModelsListConfig: service.GroupModelsListConfig{Enabled: true, Models: []string{"gpt-4.1"}},
		},
		{
			ID:               3,
			Name:             "exclusive",
			Platform:         "openai",
			Status:           service.StatusActive,
			SubscriptionType: service.SubscriptionTypeStandard,
			IsExclusive:      true,
			ModelsListConfig: service.GroupModelsListConfig{Enabled: true, Models: []string{"gpt-4.1"}},
		},
		{
			ID:               4,
			Name:             "anthropic",
			Platform:         "anthropic",
			Status:           service.StatusActive,
			SubscriptionType: service.SubscriptionTypeStandard,
			ModelsListConfig: service.GroupModelsListConfig{Enabled: true, Models: []string{"claude-sonnet-4-6"}},
		},
		{
			ID:               5,
			Name:             "disabled-group",
			Platform:         "openai",
			Status:           service.StatusDisabled,
			SubscriptionType: service.SubscriptionTypeStandard,
			ModelsListConfig: service.GroupModelsListConfig{Enabled: true, Models: []string{"disabled-model"}},
		},
		{
			ID:               6,
			Name:             "custom-list-disabled",
			Platform:         "openai",
			Status:           service.StatusActive,
			SubscriptionType: service.SubscriptionTypeStandard,
			ModelsListConfig: service.GroupModelsListConfig{Enabled: false, Models: []string{"configured-model"}},
		},
	}

	catalog := buildPublicModelCatalog(context.Background(), groups, nil)
	require.Len(t, catalog, 3)
	require.Equal(t, "standard-public", catalog[0].Name)
	require.Equal(t, "anthropic", catalog[1].Name)
	require.Equal(t, "custom-list-disabled", catalog[2].Name)
	require.Equal(t, float64(1.2), catalog[0].RateMultiplier)
	require.Equal(t, []publicModelCatalogItem{{Name: "gpt-4o", Platform: "openai"}}, catalog[0].Models)
	require.Equal(t, []publicModelCatalogItem{{Name: "claude-sonnet-4-6", Platform: "anthropic"}}, catalog[1].Models)
	require.Equal(t, []publicModelCatalogItem{{Name: "configured-model", Platform: "openai"}}, catalog[2].Models)
}

func TestToPublicResolvedPricing_TokenPricing(t *testing.T) {
	intervalInput := 1.5e-6
	pricing := toPublicResolvedPricing(&service.ResolvedPricing{
		Mode: service.BillingModeToken,
		BasePricing: &service.ModelPricing{
			InputPricePerToken:         2e-6,
			OutputPricePerToken:        8e-6,
			CacheCreationPricePerToken: 2.5e-6,
			CacheReadPricePerToken:     2e-7,
		},
		Intervals: []service.PricingInterval{
			{MinTokens: 1000, InputPrice: &intervalInput},
		},
	})

	require.NotNil(t, pricing)
	require.Equal(t, string(service.BillingModeToken), pricing.BillingMode)
	require.InDelta(t, 2e-6, *pricing.InputPrice, 1e-12)
	require.InDelta(t, 8e-6, *pricing.OutputPrice, 1e-12)
	require.InDelta(t, 2.5e-6, *pricing.CacheWritePrice, 1e-12)
	require.InDelta(t, 2e-7, *pricing.CacheReadPrice, 1e-12)
	require.Len(t, pricing.Intervals, 1)
	require.InDelta(t, intervalInput, *pricing.Intervals[0].InputPrice, 1e-12)
}

func TestToPublicResolvedPricing_RejectsEmptyPricing(t *testing.T) {
	require.Nil(t, toPublicResolvedPricing(&service.ResolvedPricing{
		Mode:        service.BillingModeToken,
		BasePricing: &service.ModelPricing{},
	}))
}

func TestBuildPublicModelCatalog_IncludesEffectiveAndOfficialPricing(t *testing.T) {
	billingService := service.NewBillingService(&config.Config{}, nil)
	pricingResolver := service.NewModelPricingResolver(nil, billingService)
	catalog := buildPublicModelCatalog(context.Background(), []service.Group{
		{
			ID:               7,
			Name:             "standard",
			Platform:         service.PlatformAnthropic,
			Status:           service.StatusActive,
			SubscriptionType: service.SubscriptionTypeStandard,
			RateMultiplier:   0.02,
			ModelsListConfig: service.GroupModelsListConfig{
				Enabled: true,
				Models:  []string{"claude-sonnet-4"},
			},
		},
	}, pricingResolver)

	require.Len(t, catalog, 1)
	require.Len(t, catalog[0].Models, 1)
	model := catalog[0].Models[0]
	require.NotNil(t, model.Pricing)
	require.NotNil(t, model.OfficialPricing)
	require.InDelta(t, 3e-6, *model.Pricing.InputPrice, 1e-12)
	require.InDelta(t, 3e-6, *model.OfficialPricing.InputPrice, 1e-12)
	require.Equal(t, float64(0.02), catalog[0].RateMultiplier)
}
