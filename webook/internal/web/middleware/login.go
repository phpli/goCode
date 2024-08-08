package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(paths string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	//go的方式编码二进制
	gob.Register(time.Now())
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		//if c.Request.URL.Path == "/users/login" ||
		//	c.Request.URL.Path == "/users/signup" {
		//	return
		//}
		session := sessions.Default(c)
		id := session.Get("userId")
		if id == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime := session.Get("updateTime")
		//说明没刷新过
		now := time.Now()
		if updateTime == nil {
			session.Set("updateTime", now)
			session.Save()
			return
		}
		//session.Set("update_time", now)
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		if now.Sub(updateTimeVal) > time.Minute {
			session.Set("updateTime", now)
			session.Save()
		}
	}
}
