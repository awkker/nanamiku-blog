package service

import (
	"encoding/json"
	"testing"
)

func TestNormalizeSiteProfileSettings(t *testing.T) {
	t.Parallel()

	got := normalizeSiteProfileSettings(SiteProfileSettings{
		BrandText:          "  NanaMiku Blog  ",
		SiteTitle:          "",
		LogoAlt:            "",
		SiteURL:            "nanamiku.blog/",
		DefaultDescription: "  hello world  ",
		DefaultSocialImage: "cdn.nanamiku.blog/cover.png",
	})

	if got.BrandText != "NanaMiku Blog" {
		t.Fatalf("unexpected brand text: %q", got.BrandText)
	}

	if got.SiteTitle != "NanaMiku Blog" {
		t.Fatalf("unexpected site title: %q", got.SiteTitle)
	}

	if got.LogoAlt != "NanaMiku Blog logo" {
		t.Fatalf("unexpected logo alt: %q", got.LogoAlt)
	}

	if got.SiteURL != "https://nanamiku.blog" {
		t.Fatalf("unexpected site url: %q", got.SiteURL)
	}

	if got.DefaultDescription != "hello world" {
		t.Fatalf("unexpected description: %q", got.DefaultDescription)
	}

	if got.DefaultSocialImage != "https://cdn.nanamiku.blog/cover.png" {
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
		{name: "absolute asset without scheme", raw: "cdn.nanamiku.blog/cover.png", want: "https://cdn.nanamiku.blog/cover.png"},
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
		HeroTitle:    "  薰逸の猫窝  ",
		HeroSubtitle: "  「月が綺麗ですね」。  ",
	})

	if got.HeroTitle != "薰逸の猫窝" {
		t.Fatalf("unexpected hero title: %q", got.HeroTitle)
	}

	if got.HeroSubtitle != "「月が綺麗ですね」。" {
		t.Fatalf("unexpected hero subtitle: %q", got.HeroSubtitle)
	}
}

func TestNormalizeHomeAssetsSettings(t *testing.T) {
	t.Parallel()

	got := normalizeHomeAssetsSettings(HomeAssetsSettings{
		HeroImages: []string{
			" /picture/fengmian/1.webp ",
			"cdn.nanamiku.blog/home-hero.webp",
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

	if got.HeroImages[1] != "https://cdn.nanamiku.blog/home-hero.webp" {
		t.Fatalf("unexpected second hero image: %q", got.HeroImages[1])
	}
}

func TestNormalizeAuthorProfileSettings(t *testing.T) {
	t.Parallel()

	got := normalizeAuthorProfileSettings(AuthorProfileSettings{
		DisplayName:      "  Xunyi  ",
		AvatarURL:        " /picture/author.jpg ",
		Role:             "  Front-end Developer / Writer  ",
		Bio:              "  写前端、写系统、写日常。  ",
		AboutDescription: "  长一点的介绍  ",
		Location:         "  China  ",
		Since:            "  Since 2026  ",
		Skills:           []string{" Astro ", "Vue", "Astro", ""},
		NowItems:         []string{" 重构首页 ", "整理模板", "重构首页"},
		Quote:            "  让每篇文章都能复用。  ",
		ContactEmail:     " mailto:test@example.com ",
		SocialLinks: []AuthorSocialLink{
			{Label: " GitHub ", Href: "github.com/awkker", IconKey: " github "},
			{Label: "", Href: "https://example.com"},
		},
		ContactLinks: []AuthorContactLink{
			{Label: " QQ ", Href: "https://example.com/qq "},
		},
	})

	if got.DisplayName != "Xunyi" {
		t.Fatalf("unexpected display name: %q", got.DisplayName)
	}

	if got.AvatarURL != "/picture/author.jpg" {
		t.Fatalf("unexpected avatar url: %q", got.AvatarURL)
	}

	if got.ContactEmail != "test@example.com" {
		t.Fatalf("unexpected contact email: %q", got.ContactEmail)
	}

	if len(got.Skills) != 2 {
		t.Fatalf("unexpected skills length: %d", len(got.Skills))
	}

	if len(got.SocialLinks) != 1 || got.SocialLinks[0].Href != "https://github.com/awkker" {
		t.Fatalf("unexpected social links: %+v", got.SocialLinks)
	}

	if len(got.ContactLinks) != 1 || got.ContactLinks[0].Href != "https://example.com/qq" {
		t.Fatalf("unexpected contact links: %+v", got.ContactLinks)
	}
}

func TestNormalizeSiteIntegrationsSettings(t *testing.T) {
	t.Parallel()

	got := normalizeSiteIntegrationsSettings(SiteIntegrationsSettings{
		GitHubUsername:  "  awkker  ",
		WeatherLocation: "  Tokyo  ",
		ShowWeather:     true,
		ShowMusic:       false,
		ShowClock:       true,
	})

	if got.GitHubUsername != "awkker" {
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
		"avatar_url": "cdn.nanamiku.blog/avatar.png",
		"contact_email": "mailto:hello@example.com",
		"social_links": [{"label":" GitHub ","href":"github.com/awkker","icon_key":" github "}]
	}`)

	got, err := decodeAuthorProfileSettings(raw)
	if err != nil {
		t.Fatalf("decodeAuthorProfileSettings() error = %v", err)
	}

	if got.DisplayName != "Nana" {
		t.Fatalf("unexpected display name: %q", got.DisplayName)
	}
	if got.AvatarURL != "https://cdn.nanamiku.blog/avatar.png" {
		t.Fatalf("unexpected avatar url: %q", got.AvatarURL)
	}
	if got.ContactEmail != "hello@example.com" {
		t.Fatalf("unexpected contact email: %q", got.ContactEmail)
	}
	if len(got.SocialLinks) != 1 || got.SocialLinks[0].Href != "https://github.com/awkker" {
		t.Fatalf("unexpected social links: %+v", got.SocialLinks)
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
		{name: "bare domain", raw: "github.com/awkker", want: "https://github.com/awkker"},
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
