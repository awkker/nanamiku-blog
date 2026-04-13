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

func TestNormalizeBlogIndexSettings(t *testing.T) {
	t.Parallel()

	got := normalizeBlogIndexSettings(BlogIndexSettings{
		HeroBadge:       " CREATOR SPACE ",
		HeroTitle:       " NanaMiku Blog ",
		HeroDescription: " 展示博客模板的首屏内容。 ",
		HeroActions: []BlogIndexHeroAction{
			{Label: " 看最新文章 ", Href: " #latest-posts "},
			{Label: " 看最新文章 ", Href: " #latest-posts "},
			{Label: "", Href: "/blog"},
		},
		QuickStats: []BlogIndexQuickStat{
			{Label: " 模板栈 ", Value: " Astro + Vue + Go "},
			{Label: " 模板栈 ", Value: " Astro + Vue + Go "},
		},
		FocusCard: BlogIndexFocusCard{
			Badge:       " 本月在做什么 ",
			Title:       " 打磨创作者模板 ",
			Description: " 继续收口默认配置。 ",
			Footnote:    " 文章来自后台管理面板发布 ",
		},
		ScrollCue: BlogIndexScrollCue{
			Label: " 向下阅读 ",
		},
	})

	if got.HeroBadge != "CREATOR SPACE" || got.HeroTitle != "NanaMiku Blog" {
		t.Fatalf("unexpected hero fields: %+v", got)
	}
	if len(got.HeroActions) != 1 || got.HeroActions[0].Href != "#latest-posts" {
		t.Fatalf("unexpected hero actions: %+v", got.HeroActions)
	}
	if len(got.QuickStats) != 1 || got.QuickStats[0].Value != "Astro + Vue + Go" {
		t.Fatalf("unexpected quick stats: %+v", got.QuickStats)
	}
	if got.ScrollCue.AriaLabel != "向下阅读" {
		t.Fatalf("unexpected scroll cue: %+v", got.ScrollCue)
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

func TestDecodeBlogIndexSettingsNormalizesPayload(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{
		"hero_badge":" CREATOR SPACE ",
		"hero_title":" NanaMiku Blog ",
		"hero_description":" 展示博客模板的首屏内容。 ",
		"hero_actions":[{"label":" 看最新文章 ","href":" /blog "}],
		"quick_stats":[{"label":" 默认主线 ","value":" 公开资料收口与可配置化 "}],
		"focus_card":{"badge":" 本月在做什么 ","title":" 打磨创作者模板 ","description":" 继续收口默认配置。 ","footnote":" 文章来自后台管理面板发布 "},
		"scroll_cue":{"label":" 向下阅读 "}
	}`)

	got, err := decodeBlogIndexSettings(raw)
	if err != nil {
		t.Fatalf("decodeBlogIndexSettings() error = %v", err)
	}

	if got.HeroTitle != "NanaMiku Blog" {
		t.Fatalf("unexpected hero title: %q", got.HeroTitle)
	}
	if len(got.HeroActions) != 1 || got.HeroActions[0].Href != "/blog" {
		t.Fatalf("unexpected hero actions: %+v", got.HeroActions)
	}
	if got.ScrollCue.AriaLabel != "向下阅读" {
		t.Fatalf("unexpected scroll cue: %+v", got.ScrollCue)
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
