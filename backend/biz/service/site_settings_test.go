package service

import (
	"encoding/json"
	"testing"
)

func TestNormalizeSiteProfileSettings(t *testing.T) {
	t.Parallel()

	got := normalizeSiteProfileSettings(SiteProfileSettings{
		BrandText:          "  Miku Blog Starter  ",
		SiteTitle:          "",
		LogoAlt:            "",
		SiteURL:            "example.com/",
		DefaultDescription: "  hello world  ",
		DefaultSocialImage: "cdn.example.com/cover.png",
	})

	if got.BrandText != "Miku Blog Starter" {
		t.Fatalf("unexpected brand text: %q", got.BrandText)
	}

	if got.SiteTitle != "Miku Blog Starter" {
		t.Fatalf("unexpected site title: %q", got.SiteTitle)
	}

	if got.LogoAlt != "Miku Blog Starter logo" {
		t.Fatalf("unexpected logo alt: %q", got.LogoAlt)
	}

	if got.SiteURL != "https://example.com" {
		t.Fatalf("unexpected site url: %q", got.SiteURL)
	}

	if got.DefaultDescription != "hello world" {
		t.Fatalf("unexpected description: %q", got.DefaultDescription)
	}

	if got.DefaultSocialImage != "https://cdn.example.com/cover.png" {
		t.Fatalf("unexpected social image: %q", got.DefaultSocialImage)
	}
}

func TestNormalizeSiteAssetURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "relative asset", raw: "/picture/logo-64.webp", want: "/picture/logo-64.webp"},
		{name: "local asset without leading slash", raw: "picture/author.jpg", want: "/picture/author.jpg"},
		{name: "dot local asset path", raw: "./picture/author.jpg", want: "/picture/author.jpg"},
		{name: "absolute asset without scheme", raw: "cdn.example.com/cover.png", want: "https://cdn.example.com/cover.png"},
		{name: "legacy malformed picture host", raw: "https://picture/author.jpg", want: "/picture/author.jpg"},
		{name: "invalid empty", raw: "   ", want: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeSiteAssetURL(tc.raw); got != tc.want {
				t.Fatalf("normalizeSiteAssetURL(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

func TestNormalizeHomeHeroSettings(t *testing.T) {
	t.Parallel()

	got := normalizeHomeHeroSettings(HomeHeroSettings{
		HeroTitle:    "  创作入口  ",
		HeroSubtitle: "  把写作、作品与公开资料收拢到同一站点。  ",
	})

	if got.HeroTitle != "创作入口" {
		t.Fatalf("unexpected hero title: %q", got.HeroTitle)
	}

	if got.HeroSubtitle != "把写作、作品与公开资料收拢到同一站点。" {
		t.Fatalf("unexpected hero subtitle: %q", got.HeroSubtitle)
	}
}

func TestNormalizeHomeAssetsSettings(t *testing.T) {
	t.Parallel()

	got := normalizeHomeAssetsSettings(HomeAssetsSettings{
		HeroImages: []string{
			" /picture/fengmian/1.webp ",
			"cdn.example.com/home-hero.webp",
			"/picture/fengmian/1.webp",
			"   ",
		},
	})

	if len(got.HeroImages) != 2 {
		t.Fatalf("unexpected hero image count: %d", len(got.HeroImages))
	}

	if got.HeroImages[0] != "/picture/fengmian/1.webp" {
		t.Fatalf("unexpected first hero image: %q", got.HeroImages[0])
	}

	if got.HeroImages[1] != "https://cdn.example.com/home-hero.webp" {
		t.Fatalf("unexpected second hero image: %q", got.HeroImages[1])
	}
}

func TestNormalizeAuthorProfileSettings(t *testing.T) {
	t.Parallel()

	got := normalizeAuthorProfileSettings(AuthorProfileSettings{
		DisplayName:      "  Your Name  ",
		AvatarURL:        " /picture/diy/about/avatar-placeholder.svg ",
		Role:             "  Frontend / Full-stack Creator  ",
		Bio:              "  展示文章、作品与公开资料。  ",
		AboutDescription: "  长一点的介绍  ",
		Location:         "  Your City  ",
		Since:            "  Since 2026  ",
		Skills:           []string{" Astro ", "Vue", "Astro", ""},
		NowItems:         []string{" 整理默认值 ", "补齐模板", "整理默认值"},
		Quote:            "  让每次写作都能继续复用。  ",
		ContactEmail:     " mailto:test@example.com ",
		SocialLinks: []AuthorSocialLink{
			{Label: " GitHub ", Href: "github.com/yourname", IconKey: " github "},
			{Label: "", Href: "https://example.com"},
		},
		ContactLinks: []AuthorContactLink{
			{Label: " Portfolio ", Href: "https://example.com/portfolio "},
		},
	})

	if got.DisplayName != "Your Name" {
		t.Fatalf("unexpected display name: %q", got.DisplayName)
	}

	if got.AvatarURL != "/picture/diy/about/avatar-placeholder.svg" {
		t.Fatalf("unexpected avatar url: %q", got.AvatarURL)
	}

	if got.ContactEmail != "test@example.com" {
		t.Fatalf("unexpected contact email: %q", got.ContactEmail)
	}

	if len(got.Skills) != 2 {
		t.Fatalf("unexpected skills length: %d", len(got.Skills))
	}

	if len(got.SocialLinks) != 1 || got.SocialLinks[0].Href != "https://github.com/yourname" {
		t.Fatalf("unexpected social links: %+v", got.SocialLinks)
	}

	if len(got.ContactLinks) != 1 || got.ContactLinks[0].Href != "https://example.com/portfolio" {
		t.Fatalf("unexpected contact links: %+v", got.ContactLinks)
	}
}

func TestNormalizeAuthorProfileSettingsLegacyAvatar(t *testing.T) {
	t.Parallel()

	got := normalizeAuthorProfileSettings(AuthorProfileSettings{
		AvatarURL: " /picture/author.jpg ",
	})

	if got.AvatarURL != "/picture/author.jpg" {
		t.Fatalf("unexpected legacy avatar normalization: %q", got.AvatarURL)
	}
}

func TestNormalizeSiteIntegrationsSettings(t *testing.T) {
	t.Parallel()

	got := normalizeSiteIntegrationsSettings(SiteIntegrationsSettings{
		GitHubUsername:  "  yourname  ",
		WeatherLocation: "  Tokyo  ",
		ShowWeather:     true,
		ShowMusic:       false,
		ShowClock:       true,
	})

	if got.GitHubUsername != "yourname" {
		t.Fatalf("unexpected github username: %q", got.GitHubUsername)
	}

	if got.WeatherLocation != "Tokyo" {
		t.Fatalf("unexpected weather location: %q", got.WeatherLocation)
	}

	if !got.ShowWeather || got.ShowMusic || !got.ShowClock {
		t.Fatalf("unexpected integrations flags: %+v", got)
	}
}

func TestNormalizeFooterSettings(t *testing.T) {
	t.Parallel()

	got := normalizeFooterSettings(FooterSettings{
		ICPText: "  沪 ICP 123456  ",
		ICPLink: " https://beian.example.com ",
		CustomTexts: []string{
			" 第一行 ",
			"",
			"第二行",
			"第三行",
			"第四行",
			"第五行",
			"第六行",
			"第七行",
			"第八行",
			"第九行",
		},
	})

	if got.ICPText != "沪 ICP 123456" {
		t.Fatalf("unexpected icp text: %q", got.ICPText)
	}

	if got.ICPLink != "https://beian.example.com" {
		t.Fatalf("unexpected icp link: %q", got.ICPLink)
	}

	if len(got.CustomTexts) != 8 {
		t.Fatalf("unexpected custom text count: %d", len(got.CustomTexts))
	}

	if got.CustomTexts[0] != "第一行" || got.CustomTexts[7] != "第八行" {
		t.Fatalf("unexpected custom texts: %+v", got.CustomTexts)
	}
}

func TestDecodeEmptySettingsDocuments(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "footer",
			run: func(t *testing.T) {
				got, err := decodeFooterSettings(nil)
				if err != nil || got == nil {
					t.Fatalf("decodeFooterSettings() err=%v got=%+v", err, got)
				}
			},
		},
		{
			name: "site profile",
			run: func(t *testing.T) {
				got, err := decodeSiteProfileSettings(nil)
				if err != nil || got == nil {
					t.Fatalf("decodeSiteProfileSettings() err=%v got=%+v", err, got)
				}
			},
		},
		{
			name: "home hero",
			run: func(t *testing.T) {
				got, err := decodeHomeHeroSettings(nil)
				if err != nil || got == nil {
					t.Fatalf("decodeHomeHeroSettings() err=%v got=%+v", err, got)
				}
			},
		},
		{
			name: "home assets",
			run: func(t *testing.T) {
				got, err := decodeHomeAssetsSettings(nil)
				if err != nil || got == nil {
					t.Fatalf("decodeHomeAssetsSettings() err=%v got=%+v", err, got)
				}
			},
		},
		{
			name: "author profile",
			run: func(t *testing.T) {
				got, err := decodeAuthorProfileSettings(nil)
				if err != nil || got == nil {
					t.Fatalf("decodeAuthorProfileSettings() err=%v got=%+v", err, got)
				}
			},
		},
		{
			name: "site integrations",
			run: func(t *testing.T) {
				got, err := decodeSiteIntegrationsSettings(nil)
				if err != nil || got == nil {
					t.Fatalf("decodeSiteIntegrationsSettings() err=%v got=%+v", err, got)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}

func TestDecodeAuthorProfileSettingsNormalizesLegacyPayload(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{
		"display_name": "  Nana  ",
		"avatar_url": "cdn.example.com/avatar.png",
		"contact_email": "mailto:hello@example.com",
		"social_links": [{"label":" GitHub ","href":"github.com/yourname","icon_key":" github "}]
	}`)

	got, err := decodeAuthorProfileSettings(raw)
	if err != nil {
		t.Fatalf("decodeAuthorProfileSettings() error = %v", err)
	}

	if got.DisplayName != "Nana" {
		t.Fatalf("unexpected display name: %q", got.DisplayName)
	}
	if got.AvatarURL != "https://cdn.example.com/avatar.png" {
		t.Fatalf("unexpected avatar url: %q", got.AvatarURL)
	}
	if got.ContactEmail != "hello@example.com" {
		t.Fatalf("unexpected contact email: %q", got.ContactEmail)
	}
	if len(got.SocialLinks) != 1 || got.SocialLinks[0].Href != "https://github.com/yourname" {
		t.Fatalf("unexpected social links: %+v", got.SocialLinks)
	}
}

func TestDecodeAuthorProfileSettingsFixesMalformedLocalAvatar(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{
		"avatar_url": "https://picture/author.jpg"
	}`)

	got, err := decodeAuthorProfileSettings(raw)
	if err != nil {
		t.Fatalf("decodeAuthorProfileSettings() error = %v", err)
	}

	if got.AvatarURL != "/picture/author.jpg" {
		t.Fatalf("unexpected avatar url: %q", got.AvatarURL)
	}
}

func TestNormalizePublicLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "relative", raw: "/friends", want: "/friends"},
		{name: "mailto", raw: "mailto:test@example.com", want: "mailto:test@example.com"},
		{name: "bare domain", raw: "github.com/yourname", want: "https://github.com/yourname"},
		{name: "empty", raw: "   ", want: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizePublicLink(tc.raw); got != tc.want {
				t.Fatalf("normalizePublicLink(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

func TestNormalizeAboutPageSettings(t *testing.T) {
	t.Parallel()

	got := normalizeAboutPageSettings(AboutPageSettings{
		IntroCards: []AboutIntroCard{
			{Title: " 当前主线 ", Description: " 整理设置中心 "},
			{Title: " 当前主线 ", Description: " 整理设置中心 "},
		},
		Milestones: []AboutMilestone{
			{Year: " 2026 ", Title: " 收口 About ", Summary: " 接入后台 ", Result: " 链路打通 "},
		},
		CapabilityGroups: []AboutCapabilityGroup{
			{Title: " 前端体验 ", Desc: " 关注结构 ", Stack: []string{" Astro ", "Vue", "Astro"}},
		},
		FeaturedProjects: []AboutFeaturedProject{
			{Name: " Starter Site Kit ", Focus: " 模板 ", Role: " 配置治理 ", Metric: " 可安全开源 ", Href: "github.com/yourname/miku-blog-starter"},
		},
		MonthlyGoals: []string{" 补齐 About 页配置 ", "补齐 About 页配置"},
		ListeningNow: []string{" Lo-fi Focus Mix "},
		Signature: AboutSignatureSettings{
			Description: " 让内容维护更顺畅。 ",
			Footer:      " 持续迭代中。 ",
		},
	})

	if len(got.IntroCards) != 1 || got.IntroCards[0].Title != "当前主线" {
		t.Fatalf("unexpected intro cards: %+v", got.IntroCards)
	}
	if len(got.CapabilityGroups) != 1 || len(got.CapabilityGroups[0].Stack) != 2 {
		t.Fatalf("unexpected capability groups: %+v", got.CapabilityGroups)
	}
	if got.FeaturedProjects[0].Href != "https://github.com/yourname/miku-blog-starter" {
		t.Fatalf("unexpected project href: %+v", got.FeaturedProjects)
	}
	if len(got.MonthlyGoals) != 1 || got.MonthlyGoals[0] != "补齐 About 页配置" {
		t.Fatalf("unexpected monthly goals: %+v", got.MonthlyGoals)
	}
}

func TestDecodeAboutPageSettingsNormalizesLegacyPayload(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{
		"intro_cards":[{"title":" 当前主线 ","description":" 整理设置中心 "}],
		"featured_projects":[{"name":" Starter Site Kit ","focus":" 模板 ","role":" 配置治理 ","metric":" 可安全开源 ","href":"github.com/yourname/miku-blog-starter"}],
		"signature":{"description":" 让内容维护更顺畅。 ","footer":" 持续迭代中。 "}
	}`)

	got, err := decodeAboutPageSettings(raw)
	if err != nil {
		t.Fatalf("decodeAboutPageSettings() error = %v", err)
	}

	if len(got.IntroCards) != 1 || got.IntroCards[0].Title != "当前主线" {
		t.Fatalf("unexpected intro cards: %+v", got.IntroCards)
	}
	if len(got.FeaturedProjects) != 1 || got.FeaturedProjects[0].Href != "https://github.com/yourname/miku-blog-starter" {
		t.Fatalf("unexpected featured projects: %+v", got.FeaturedProjects)
	}
	if got.Signature.Description != "让内容维护更顺畅。" {
		t.Fatalf("unexpected signature: %+v", got.Signature)
	}
}
