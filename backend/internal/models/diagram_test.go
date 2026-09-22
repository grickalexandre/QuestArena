package models

import (
	"strings"
	"testing"
)

func TestSanitizeDiagramSVG(t *testing.T) {
	ok, err := SanitizeDiagramSVG(`<svg viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg"><rect x="1" y="1" width="8" height="8"/></svg>`)
	if err != nil || !strings.Contains(ok, "<rect") {
		t.Fatalf("valid svg: %q %v", ok, err)
	}
	cleaned, err := SanitizeDiagramSVG(`<svg onload="alert(1)"><text>ok</text></svg>`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(cleaned, "onload") {
		t.Fatalf("event handler leaked: %s", cleaned)
	}
	if _, err := SanitizeDiagramSVG(`<svg><script>x</script></svg>`); err == nil {
		t.Fatal("script must be rejected")
	}
	empty, err := SanitizeDiagramSVG("  ")
	if err != nil || empty != "" {
		t.Fatalf("empty: %q %v", empty, err)
	}
}
