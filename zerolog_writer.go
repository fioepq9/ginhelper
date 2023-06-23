package ginhelper

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type zerologWriter struct {
	log zerolog.Logger
	lvl zerolog.Level
}

func NewZerologWriter(log zerolog.Logger, lvl zerolog.Level) *zerologWriter {
	return &zerologWriter{
		log: log,
		lvl: lvl,
	}
}

func (w *zerologWriter) Write(p []byte) (n int, err error) {
	s := string(p)

	s = strings.TrimPrefix(s, "[GIN-debug]")
	s = strings.TrimSpace(s)
	lvl := w.lvl

	if strings.HasPrefix(s, "[WARNING]") {
		s = strings.TrimPrefix(s, "[WARNING]")
		s = strings.TrimSpace(s)
		lvl = zerolog.WarnLevel
	}

	if strings.HasPrefix(s, "[ERROR]") {
		s = strings.TrimPrefix(s, "[ERROR]")
		s = strings.TrimSpace(s)
		lvl = zerolog.ErrorLevel
	}

	s = strings.ReplaceAll(s, `"`, `'`)
	s = strings.ReplaceAll(s, "\t", "  ")

	if strings.Contains(s, "\n") {
		ss := strings.Split(s, "\n")
		for i := range ss {
			ss[i] = strings.TrimSpace(ss[i])
		}
		w.log.WithLevel(lvl).Strs("message", ss).Send()
		return len(p), nil
	}

	w.log.WithLevel(lvl).Msg(s)
	return len(p), nil
}

func (w *zerologWriter) SetAll() *zerologWriter {
	return w.SetGinDefaultWriter().
		SetGinDefaultErrorWriter().
		SetGinDebugPrintRouteFunc()
}

func (w *zerologWriter) SetGinDefaultWriter() *zerologWriter {
	gin.DefaultWriter = w
	return w
}

func (w *zerologWriter) SetGinDefaultErrorWriter() *zerologWriter {
	gin.DefaultErrorWriter = NewZerologWriter(w.log, zerolog.ErrorLevel)
	return w
}

func (w *zerologWriter) SetGinDebugPrintRouteFunc() *zerologWriter {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		w.log.WithLevel(w.lvl).
			Str("method", httpMethod).
			Str("path", absolutePath).
			Str("handler", handlerName).
			Int("num_handlers", nuHandlers).
			Send()
	}
	return w
}
