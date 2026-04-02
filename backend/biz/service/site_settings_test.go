package service

import "testing"

func TestNormalizeFooterSettings(t *testing.T) {
	settings := normalizeFooterSettings(FooterSettings{
		ICPText: "  沪ICP备12345678号-1  ",
		ICPLink: " https://beian.miit.gov.cn/ ",
		CustomTexts: []string{
			"  第一行  ",
			"",
			"   ",
			"第二行",
		},
	})

	if settings.ICPText != "沪ICP备12345678号-1" {
		t.Fatalf("unexpected icp text: %q", settings.ICPText)
	}
	if settings.ICPLink != "https://beian.miit.gov.cn/" {
		t.Fatalf("unexpected icp link: %q", settings.ICPLink)
	}
	if len(settings.CustomTexts) != 2 {
		t.Fatalf("expected 2 custom texts, got %d", len(settings.CustomTexts))
	}
	if settings.CustomTexts[0] != "第一行" || settings.CustomTexts[1] != "第二行" {
		t.Fatalf("unexpected custom texts: %#v", settings.CustomTexts)
	}
}

func TestSanitizeFooterCustomTextsLimit(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	result := sanitizeFooterCustomTexts(lines)

	if len(result) != maxFooterCustomTexts {
		t.Fatalf("expected %d custom texts, got %d", maxFooterCustomTexts, len(result))
	}
	if result[maxFooterCustomTexts-1] != "8" {
		t.Fatalf("unexpected last custom text: %q", result[maxFooterCustomTexts-1])
	}
}
