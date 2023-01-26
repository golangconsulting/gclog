package gclog

import (
	"math"
	"strconv"
	"time"
)

var styleKeyStart = styleKey.Start(true)
var styleKeyEnd = styleKey.End(true)

var styleValStart = styleVal.Start(true)
var styleValEnd = styleVal.End(true)

var styleValErrStart = styleValErr.Start(true)
var styleValErrEnd = styleValErr.End(true)

func (l *Line) appendKey(key string) {
	if !l.json {
		l.buff = append(l.buff, ',', ' ')
	} else {
		l.buff = append(l.buff, ',', ' ', '"')
	}

	if l.canColorize {
		l.buff = append(l.buff, styleKeyStart...)
	}

	if !l.json {
		l.buff = append(l.buff, key...)
	} else {
		l.appendSafeString(key)
	}

	if l.canColorize {
		l.buff = append(l.buff, styleKeyEnd...)
	}

	if !l.json {
		l.buff = append(l.buff, '=')
	} else {
		l.buff = append(l.buff, '"', ':')
	}
}

func (l *Line) appendInt(val int64) {
	if l.canColorize {
		l.buff = append(l.buff, styleValStart...)
	}
	l.buff = strconv.AppendInt(l.buff, val, 10)
	if l.canColorize {
		l.buff = append(l.buff, styleValEnd...)
	}
}

func (l *Line) appendUInt(val uint64) {
	if l.canColorize {
		l.buff = append(l.buff, styleValStart...)
	}
	l.buff = strconv.AppendUint(l.buff, val, 10)
	if l.canColorize {
		l.buff = append(l.buff, styleValEnd...)
	}
}

func (l *Line) appendStr(val string) {
	l.buff = append(l.buff, '"')
	if l.canColorize {
		l.buff = append(l.buff, styleValStart...)
	}

	if !l.json {
		l.buff = append(l.buff, val...)
	} else {
		l.appendSafeString(val)
	}

	if l.canColorize {
		l.buff = append(l.buff, styleValEnd...)
	}
	l.buff = append(l.buff, '"')
}

func (l *Line) appendBool(val bool) {
	if l.canColorize {
		l.buff = append(l.buff, styleValStart...)
	}
	if val {
		l.buff = append(l.buff, 't', 'r', 'u', 'e')
	} else {
		l.buff = append(l.buff, 'f', 'a', 'l', 's', 'e')
	}
	if l.canColorize {
		l.buff = append(l.buff, styleValEnd...)
	}
}

var floatErrStrings = []string{"null", "'NaN'", "'-∞'", "'∞'"}

func (l *Line) appendFloat(val float64, bitSize int) {
	// Error case.
	i := 0
	floatErr := false
	switch {
	case math.IsNaN(val):
		floatErr = true
		if !l.json {
			i = 1
		}
	case math.IsInf(val, -1):
		floatErr = true
		if !l.json {
			i = 2
		}
	case math.IsInf(val, 1):
		floatErr = true
		if !l.json {
			i = 3
		}
	}

	if floatErr {
		if l.canColorize {
			l.buff = append(l.buff, styleValErrStart...)
		}
		l.buff = append(l.buff, floatErrStrings[i]...)
		if l.canColorize {
			l.buff = append(l.buff, styleValErrEnd...)
		}
		return
	}

	// Normal case.
	if l.canColorize {
		l.buff = append(l.buff, styleValStart...)
	}
	l.buff = strconv.AppendFloat(l.buff, val, 'f', -1, bitSize)
	if l.canColorize {
		l.buff = append(l.buff, styleValEnd...)
	}
}

func (l *Line) appendTime(val time.Time) {
	if l.canColorize {
		l.buff = append(l.buff, styleValStart...)
	}
	l.buff = val.AppendFormat(l.buff, "2006/01/02 15:04:05")
	if l.canColorize {
		l.buff = append(l.buff, styleValEnd...)
	}
}
