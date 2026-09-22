package models

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

const maxDiagramRunes = 12000

var (
	diagramEventAttr = regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)`)
	diagramBanned    = regexp.MustCompile(`(?i)<\s*(script|foreignobject|iframe|object|embed|use|animate|set|handler|style)\b|javascript:|data:`)
)

// SanitizeDiagramSVG keeps a static SVG diagram or returns empty if none was sent.
func SanitizeDiagramSVG(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	if utf8.RuneCountInString(raw) > maxDiagramRunes {
		return "", errors.New("diagrama muito longo")
	}
	lower := strings.ToLower(raw)
	if !strings.HasPrefix(lower, "<svg") {
		return "", errors.New("diagrama deve ser SVG")
	}
	if !strings.Contains(lower, "</svg>") {
		return "", errors.New("diagrama SVG incompleto")
	}
	if diagramBanned.FindStringIndex(raw) != nil {
		return "", errors.New("diagrama SVG inválido")
	}
	raw = diagramEventAttr.ReplaceAllString(raw, "")
	return raw, nil
}
