import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import { createPinia } from "pinia";
import { createI18n } from "vue-i18n";
import SubscriptionPlanCard from "../SubscriptionPlanCard.vue";

const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackWarn: false,
  missingWarn: false,
  messages: {
    en: {
      payment: {
        days: () => "days",
        permanent: () => "Permanent",
        models: () => "Models",
        rechargeNow: () => "Top Up Now",
        planCard: {
          balanceRecharge: () => "Balance Top-Up",
          creditedBalance: () => "Credited Balance",
          quota: () => "Quota",
          rate: () => "Rate",
          unlimited: () => "Unlimited",
        },
        subscribeNow: () => "Subscribe now",
      },
    },
  },
});

const mountPlanCard = (groupPlatform: string) =>
  mount(SubscriptionPlanCard, {
    props: {
      plan: {
        id: 1,
        group_id: 10,
        group_platform: groupPlatform,
        name: "Pro",
        price: 10,
        features: [],
        rate_multiplier: 1,
        validity_days: 30,
        validity_unit: "day",
        supported_model_scopes: ["claude", "gemini_text", "gemini_image"],
        is_active: true,
      },
    },
    global: { plugins: [i18n, createPinia()] },
  });

describe("SubscriptionPlanCard", () => {
  it("does not show Antigravity model scopes for OpenAI plans", () => {
    const text = mountPlanCard("openai").text();

    expect(text).not.toContain("Claude");
    expect(text).not.toContain("Gemini");
    expect(text).not.toContain("Imagen");
  });

  it("shows model scopes for Antigravity plans", () => {
    const text = mountPlanCard("antigravity").text();

    expect(text).toContain("Claude");
    expect(text).toContain("Gemini");
    expect(text).toContain("Imagen");
  });

  it("renders balance top-up plans with balance labels", () => {
    const wrapper = mount(SubscriptionPlanCard, {
      props: {
        buttonMode: "balance",
        plan: {
          id: 1,
          group_id: 0,
          group_platform: "",
          name: "Balance 50",
          price: 50,
          features: [],
          validity_days: 30,
          validity_unit: "day",
          for_sale: true,
          sort_order: 1,
        },
      },
      global: { plugins: [i18n] },
    });

    const text = wrapper.text();

    expect(text).toContain("Balance Top-Up");
    expect(text).toContain("¥50");
    expect(text).toContain("Top Up Now");
    expect(text).not.toContain("payment.planCard.rate");
  });

  it("shows permanent validity when validity days is zero", () => {
    const wrapper = mount(SubscriptionPlanCard, {
      props: {
        plan: {
          id: 2,
          group_id: 10,
          group_platform: "openai",
          name: "Lifetime",
          price: 99,
          features: [],
          validity_days: 0,
          validity_unit: "",
          for_sale: true,
          sort_order: 1,
        },
      },
      global: { plugins: [i18n] },
    });

    expect(wrapper.text()).toContain("/ Permanent");
    expect(wrapper.text()).not.toContain("0days");
  });
});
