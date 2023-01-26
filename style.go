package gclog

import (
	"github.com/arafath-mk/gcstyle"
	"github.com/arafath-mk/gcstyle/wcolor"
)

func (l *Logger) ColorizeText(text string, c wcolor.Color) string {
	if l.canApplyStyle {
		return gcstyle.ApplyTo(text, l.canApplyStyle).Color(c).String()
	}
	return text
}

func (l *Logger) StyleText(text string, style gcstyle.Style) string {
	if l.canApplyStyle {
		return style.ApplyTo(text, l.canApplyStyle).String()
	}
	return text
}
