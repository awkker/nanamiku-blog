package service

import "testing"

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
