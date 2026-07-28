package setting

import "testing"

func TestDefaultSiteConfigProvidesEditableHomepageContent(t *testing.T) {
	config := defaultSiteConfig()

	if config.SiteTheme != "syntro" {
		t.Fatalf("default theme = %q, want syntro", config.SiteTheme)
	}
	if config.Navigation.Items == nil || len(config.Navigation.Items) == 0 {
		t.Fatal("default navigation must contain editable items")
	}
	if config.Hero.Title == "" || config.Hero.PrimaryLabel == "" {
		t.Fatal("default hero must contain editable title and primary action")
	}
	if len(config.Features) < 3 || len(config.Stats.Items) < 3 {
		t.Fatal("default homepage must contain editable features and stats")
	}
	if config.Theme.Primary == "" || config.Theme.Accent == "" {
		t.Fatal("default theme must expose editable colors")
	}
	if config.Dashboard.Score == "" || config.SEO.Title == "" {
		t.Fatal("default dashboard and SEO config must be editable")
	}
}

func TestNormalizeThemeFallsBackToSyntro(t *testing.T) {
	if got := normalizeTheme("ember"); got != "syntro" {
		t.Fatalf("normalizeTheme(ember) = %q, want syntro", got)
	}
	if got := normalizeTheme("cruip"); got != "cruip" {
		t.Fatalf("normalizeTheme(cruip) = %q, want cruip", got)
	}
}

func TestValidateSiteConfigRejectsUnsafeThemeColors(t *testing.T) {
	config := defaultSiteConfig()
	config.Theme.Primary = "red; background:url(javascript:alert(1))"

	if err := ValidateSiteConfig(config); err == nil {
		t.Fatal("expected invalid theme color to be rejected")
	}
}
