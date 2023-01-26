package gclog

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/arafath-mk/gcstyle"
	"github.com/arafath-mk/gcstyle/wcolor"
)

type Writer struct {
	out io.Writer
	mut sync.Mutex
}

func (w *Writer) Write(data []byte) {
	w.mut.Lock()
	defer w.mut.Unlock()

	w.out.Write(data)
}

var loggerPool = &sync.Pool{
	New: func() interface{} {
		l := &Logger{
			w: &Writer{},
		}
		return l
	},
}

type Logger struct {
	w             *Writer // Writer is shared with children.
	canApplyStyle bool
	json          bool
	finished      bool
	context       *Line
}

func New(w io.Writer, json bool) *Logger {
	if w == nil {
		w = io.Discard
	}

	l := loggerPool.Get().(*Logger)
	l.w.out = w
	l.canApplyStyle = gcstyle.CanApplyStyle(w)
	l.json = json
	l.finished = false
	l.context = newLine(nil, l.canApplyStyle, l.json)
	return l
}

func (l *Logger) newChild() *Logger {
	newChild := loggerPool.Get().(*Logger)
	newChild.w = l.w // Writer is shared with children.
	newChild.canApplyStyle = l.canApplyStyle
	newChild.json = l.json
	// Allows to create new child of a finished logger. But, it should not output anything.
	newChild.finished = l.finished
	newChild.context = newLine(newChild, l.canApplyStyle, l.json)
	if !l.finished {
		newChild.context.buff = append(newChild.context.buff, l.context.buff...)
	}
	return newChild
}

func (l *Logger) With() *Line {
	// Create a new child to avoid any concurrency issues while updating the "prefix"
	newChild := l.newChild()
	return newChild.context
}

func (l *Logger) EndWith() {
	linePool.Put(l.context)
	l.finished = true
	loggerPool.Put(l)
}

func (l *Logger) ForceColor() {
	l.canApplyStyle = true
}

func (l *Logger) CanColorize() bool {
	return l.canApplyStyle
}

func (l *Logger) Println(a ...any) {
	l.Print(a...)
}

func (l *Logger) Print(a ...any) {
	l.print([]byte(l.msg(a...)))
}

func (l *Logger) print(msg []byte) {
	if !l.json {
		l.printText(msg)
		return
	}

	l.printJson(msg)
}

func (l *Logger) printText(msg []byte) {
	if l.finished {
		return
	}

	var prefix []byte = l.context.buff
	if len(l.context.buff) > 2 && l.context.buff[0] == ',' {
		prefix = l.context.buff[2:] // buff[2:] -> No need to print the ", " at the start.
	}

	line := newLine(nil, false, l.json)
	defer linePool.Put(line)

	// fmt.Sprintf("%s %s%s\n", now, l.context.buff, msg)
	if l.canApplyStyle {
		line.buff = append(line.buff, styleValStart...)
	}
	line.buff = time.Now().AppendFormat(line.buff, "2006/01/02 15:04:05")
	if l.canApplyStyle {
		line.buff = append(line.buff, styleValEnd...)
	}
	line.buff = append(line.buff, ' ')
	line.buff = append(line.buff, prefix...)
	if len(prefix) > 0 {
		line.buff = append(line.buff, ',')
		line.buff = append(line.buff, ' ')
	}
	line.buff = append(line.buff, msg...)
	line.buff = append(line.buff, '\n')
	l.w.Write(line.buff)
}

func (l *Logger) printJson(msg []byte) {
	if l.finished {
		return
	}

	var prefix []byte = l.context.buff
	if len(l.context.buff) > 2 && l.context.buff[0] == ',' {
		prefix = l.context.buff[2:] // buff[2:] -> No need to print the ", " at the start.
	}

	line := newLine(nil, false, l.json)
	defer linePool.Put(line)

	t := time.Now().UnixMicro()
	timeKey := l.ColorizeText("time", *styleKey.Color)
	// fmt.Sprintf("{\"%s\": \"%s\", %s%s}\n", timeKey, now, prefix, msg)
	line.buff = append(line.buff, '{', '"')
	line.buff = append(line.buff, timeKey...)
	line.buff = append(line.buff, '"', ':')
	if l.canApplyStyle {
		line.buff = append(line.buff, styleValStart...)
	}
	line.buff = strconv.AppendInt(line.buff, t, 10)
	if l.canApplyStyle {
		line.buff = append(line.buff, styleValEnd...)
	}
	line.buff = append(line.buff, ',', ' ')
	line.buff = append(line.buff, prefix...)
	if len(prefix) > 0 {
		line.buff = append(line.buff, ',', ' ')
	}
	line.buff = append(line.buff, msg...)
	line.buff = append(line.buff, '}', '\n')
	l.w.Write(line.buff)
}

func (l *Logger) Log(a ...any) {
	l.Println(a...)
}

func (l *Logger) Logf(format string, a ...any) {
	l.Print(fmt.Sprintf(format, a...))
}

func (l *Logger) Error(a ...any) {
	l.Println(l.ColorizeText(fmt.Sprint(a...), wcolor.Red))
}

func (l *Logger) LogHttpRequest(str string) {
	l.Println(str)
}

func (l *Logger) In(fn string) string {
	start := l.ColorizeText("[FN START] ", *styleVal.Color)

	fn = l.StyleText(fn, styleVal)
	l.Println(start, fn)

	end := l.ColorizeText("[FN ENDED]", *styleVal.Color)
	return end + " " + fn
}

func (l *Logger) Out(str string) {
	l.Println(str)
}

func (l *Logger) StartJson() *Line {
	return newLine(l, false, l.json)
}

func (l *Logger) msg(a ...any) string {
	if !l.json {
		return fmt.Sprint(a...)
	}
	return fmt.Sprintf("\"msg\":\"%s\"", fmt.Sprint(a...))
}
