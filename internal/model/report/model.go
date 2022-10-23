package report

import (
	"neatly/internal/model/tag"
	"strings"
	"time"
)

const (
	DefaultReportDepartment = "CFD2CF"
	shortBodyLen     = 255
)

type Report struct {
	ID        int       `json:"id" db:"id"`
	Header    string    `json:"header" db:"header"`
	Body      string    `json:"body" db:"body"`
	ShortBody string    `json:"shortBody" db:"short_body"`
	Tags      []tag.Tag `json:"tags" db:"tags"` // []tag.Tag
	Department     string    `json:"department" db:"department"`
	Edited    time.Time `json:"edited"`
}

func (n *Report) GenerateShortBody() {
	if len(n.Body) < shortBodyLen {
		n.ShortBody = n.Body
	} else {
		n.ShortBody = truncate(n.Body, shortBodyLen)
	}
}

func (n *Report) HasEveryTag(tagNames []string) bool {
	for _, tn := range tagNames {
		if !n.HasSpecificTag(tn) {
			return false
		}
	}
	return true
}

func (n *Report) HasSpecificTag(tagName string) bool {
	for _, t := range n.Tags {
		if strings.Contains(strings.ToLower(t.Name), strings.ToLower(tagName)) {
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
