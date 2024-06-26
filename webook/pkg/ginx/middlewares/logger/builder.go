package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

// 注意点
// 1.小心日志内容过多。url请求过长，请求体响应体过大
// 2.考虑问题1，以及用户切换不同框架，足够灵活
// 3.考虑动态开关
type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc:   fn,
		allowReqBody: atomic.NewBool(false),
	}
}
func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			IP:     ctx.ClientIP(),
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			//io里的数据 是一次性的，要重复用的话 需要 放回去
			body, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(body) > 1024*10 {
				body = body[:1024*10]
			}
			al.ReqBody = string(body)
		}
		if b.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()
		ctx.Next()
	}
}

func (b *MiddlewareBuilder) AllowReqBody(allowReqBody bool) *MiddlewareBuilder {
	b.allowReqBody.Store(true)
	return b
}

func (b *MiddlewareBuilder) AllowRespBody(allowRespBody bool) *MiddlewareBuilder {
	b.allowRespBody = allowRespBody
	return b
}

type AccessLog struct {
	IP       string
	Method   string
	Url      string
	ReqBody  string
	RespBody string
	Duration string
	Status   int
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (rw responseWriter) Write(b []byte) (int, error) {
	rw.al.RespBody = string(b)
	return rw.ResponseWriter.Write(b)
}

func (rw responseWriter) WriteHeader(statusCode int) {
	rw.al.Status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw responseWriter) WriterString(data string) (int, error) {
	rw.al.RespBody = data
	return rw.ResponseWriter.WriteString(data)
}
