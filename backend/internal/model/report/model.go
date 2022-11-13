package report

import (
	"reports_system/internal/model/label"
	"strings"
	"time"
)

const (
	shortBodyLen = 255
)

type Report struct {
	ID        int           `json:"id" db:"id"`
	Header    string        `json:"header" db:"header"`
	Body      string        `json:"body" db:"body"`
	ShortBody string        `json:"shortBody" db:"short_body"`
	Labels    []label.Label `json:"labels" db:"labels"` // []label.Label
	Edited    time.Time     `json:"edited"`
}

func (n *Report) GenerateShortBody() {
	if len(n.Body) < shortBodyLen {
		n.ShortBody = n.Body
	} else {
		n.ShortBody = truncate(n.Body, shortBodyLen)
	}
}

func (n *Report) HasEveryLabel(labelNames []string) bool {
	for _, tn := range labelNames {
		if !n.HasSpecificLabel(tn) {
			return false
		}
	}
	return true
}

func (n *Report) HasSpecificLabel(labelName string) bool {
	for _, t := range n.Labels {
		if strings.Contains(strings.ToLower(t.Name), strings.ToLower(labelName)) {
			return true
		}
	}
	return false
}

func truncate(text string, width int) string {
	r := []rune(text)
	trunc := r[:width]
	return string(trunc)
}
