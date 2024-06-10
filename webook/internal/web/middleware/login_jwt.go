package middleware

import (
	"encoding/gob"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(paths string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	//go的方式编码二进制
	gob.Register(time.Now())
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		tokenHeard := c.Request.Header.Get("Authorization")
		if tokenHeard == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//tokenStr := strings.SplitN(tokenHeard, " ", 2)
		segs := strings.Split(tokenHeard, " ")
		if len(segs) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("fb0e22c79ac75679e9881e6ba183b354"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//c.Set("userId", claims.Uid)
		c.Set("claims", claims)
	}
}
