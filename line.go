package gclog

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var PrintCallStackForErr = false
var CallStackDepthToPrint = 2

const buffSize = 500
const msgKey = "msg"

var linePool = &sync.Pool{
	New: func() interface{} {
		return &Line{
			buff: make([]byte, 0, buffSize),
		}
	},
}

type Line struct {
	log         *Logger
	canColorize bool
	json        bool
	buff        []byte
}

func newLine(log *Logger, colorize bool, json bool) *Line {
	l := linePool.Get().(*Line)
	l.buff = l.buff[:0]
	l.log = log
	l.canColorize = colorize
	l.json = json
	return l
}

func (l *Line) Logger() *Logger {
	return l.log
}

func (l *Line) Finish() {
	if len(l.buff) > 2 {
		l.log.print(l.buff[2:]) // l.buff[2:] -> No need to print the ", " at the start.
	}
	linePool.Put(l)
}

func (l *Line) Send() {
	l.Finish()
}

func (l *Line) Int(key string, val int) *Line {
	l.appendKey(key)
	l.appendInt(int64(val))
	return l
}

func (l *Line) Ints(key string, val []int) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendInt(int64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendInt(int64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Int8(key string, val int8) *Line {
	l.appendKey(key)
	l.appendInt(int64(val))
	return l
}

func (l *Line) Ints8(key string, val []int8) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendInt(int64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendInt(int64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Int16(key string, val int16) *Line {
	l.appendKey(key)
	l.appendInt(int64(val))
	return l
}

func (l *Line) Ints16(key string, val []int16) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendInt(int64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendInt(int64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Int32(key string, val int32) *Line {
	l.appendKey(key)
	l.appendInt(int64(val))
	return l
}

func (l *Line) Ints32(key string, val []int32) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendInt(int64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendInt(int64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Int64(key string, val int64) *Line {
	l.appendKey(key)
	l.appendInt(int64(val))
	return l
}

func (l *Line) Ints64(key string, val []int64) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendInt(int64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendInt(int64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Uint(key string, val uint) *Line {
	l.appendKey(key)
	l.appendUInt(uint64(val))
	return l
}

func (l *Line) Uints(key string, val []uint) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendUInt(uint64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendUInt(uint64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Uint8(key string, val uint8) *Line {
	l.appendKey(key)
	l.appendUInt(uint64(val))
	return l
}

func (l *Line) Uints8(key string, val []uint8) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendUInt(uint64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendUInt(uint64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Uint16(key string, val uint16) *Line {
	l.appendKey(key)
	l.appendUInt(uint64(val))
	return l
}

func (l *Line) Uints16(key string, val []uint16) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendUInt(uint64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendUInt(uint64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}
func (l *Line) Uint32(key string, val uint32) *Line {
	l.appendKey(key)
	l.appendUInt(uint64(val))
	return l
}

func (l *Line) Uints32(key string, val []uint32) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendUInt(uint64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendUInt(uint64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Uint64(key string, val uint64) *Line {
	l.appendKey(key)
	l.appendUInt(uint64(val))
	return l
}

func (l *Line) Uints64(key string, val []uint64) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendUInt(uint64(val[0]))
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendUInt(uint64(val[i]))
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Bytes(key string, val []byte) *Line {
	l.appendKey(key)
	l.appendStr(string(val))
	return l
}

func (l *Line) Str(key string, val string) *Line {
	l.appendKey(key)
	l.appendStr(val)
	return l
}

func (l *Line) Strs(key string, val []string) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendStr(val[0])
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendStr(val[i])
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Bool(key string, val bool) *Line {
	l.appendKey(key)
	l.appendBool(val)
	return l
}

func (l *Line) Bools(key string, val []bool) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendBool(val[0])
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendBool(val[i])
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Msg(msg string) {
	l.Str(msgKey, msg)
	l.Finish()
}

func (l *Line) Msgf(f string, v ...any) {
	l.Msg(fmt.Sprintf(f, v...))
}

func (l *Line) Float32(key string, val float32) *Line {
	l.appendKey(key)
	l.appendFloat(float64(val), 32)
	return l
}

func (l *Line) Floats32(key string, val []float32) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendFloat(float64(val[0]), 32)
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendFloat(float64(val[i]), 32)
	}
	l.buff = append(l.buff, ']')
	return l
}

func (l *Line) Float64(key string, val float64) *Line {
	l.appendKey(key)
	l.appendFloat(val, 64)
	return l
}

func (l *Line) Floats64(key string, val []float64) *Line {
	l.appendKey(key)
	l.buff = append(l.buff, '[')
	l.appendFloat(val[0], 64)
	for i := 1; i < len(val); i++ {
		l.buff = append(l.buff, ',', ' ')
		l.appendFloat(val[i], 64)
	}
	l.buff = append(l.buff, ']')
	return l
}

// -->
func (l *Line) Err(err error) *Line {
	if err == nil {
		return l
	}

	l.appendKey("err")
	l.appendStr(err.Error())
	if PrintCallStackForErr {
		var file string
		var line int
		{
			var ok bool
			_, file, line, ok = runtime.Caller(CallStackDepthToPrint)
			if !ok {
				file = "???"
				line = 0
			}
		}
		l.appendKey("errLoggedFrom")
		l.buff = append(l.buff, '"')
		l.buff = append(l.buff, file...)
		l.buff = append(l.buff, ':')
		l.buff = strconv.AppendInt(l.buff, int64(line), 10)
		l.buff = append(l.buff, '"')
	}

	return l
}

func (l *Line) Time(key string, val time.Time) *Line {
	l.appendKey(key)
	l.appendTime(val)
	return l
}

func (l *Line) Dur(key string, val time.Duration) *Line {
	l.appendKey(key)
	if !l.json {
		l.appendStr(fmt.Sprintf("%v", val))
	} else {
		l.appendStr(fmt.Sprintf("%d", val))
	}
	return l
}

func (l *Line) Interface(key string, val any) *Line {
	l.appendKey(key)
	// l.appendStr(fmt.Sprintf("%v", val))
	b, _ := json.Marshal(val)
	l.appendStr(string(b))
	return l
}

func (l *Line) Type(key string, val any) *Line {
	l.appendKey(key)
	l.appendStr(reflect.TypeOf(val).String())
	return l
}

// <--
