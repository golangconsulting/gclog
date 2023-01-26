package gclog

import "unicode/utf8"

const hex = "0123456789abcdef"

func (l *Line) appendSafeString(s string) {
	escapeHTML := true
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				l.buff = append(l.buff, s[start:i]...)
			}

			l.buff = append(l.buff, '\\')
			switch b {
			case '\\', '"':
				l.buff = append(l.buff, b)
			case '\n':
				l.buff = append(l.buff, '\\', 'n')
			case '\r':
				l.buff = append(l.buff, '\\', 'r')
			case '\t':
				l.buff = append(l.buff, '\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				l.buff = append(l.buff, 'u', '0', '0', hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				l.buff = append(l.buff, s[start:i]...)
			}
			l.buff = append(l.buff, `\ufffd`...)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				l.buff = append(l.buff, s[start:i]...)
			}
			l.buff = append(l.buff, `\u202`...)
			l.buff = append(l.buff, hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		l.buff = append(l.buff, s[start:]...)
	}
}
