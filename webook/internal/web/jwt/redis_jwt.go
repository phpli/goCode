package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	AtKey = []byte("fb0e22c79ac75679e9881e6ba183b354")
	Rtkey = []byte("fb0e22c79ac7567929881e6ba183b354")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (r *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
		Ssid:      ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}

	//fmt.Println(u)
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	//claim, ok := ctx.Get("claims")
	claims := ctx.MustGet("claims").(*UserClaims)
	//if !ok {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	//claims, ok := claim.(*UserClaims)
	//if !ok {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	ssid := claims.Ssid
	return r.cmd.Set(ctx, fmt.Sprintf("users:ssis:%s", ssid), "", time.Hour*24*7).Err()
}

func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := r.cmd.Exists(ctx, fmt.Sprintf("users:ssis:%s", ssid)).Result()
	return err
}

func (r *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeard := ctx.GetHeader("Authorization")
	if tokenHeard == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.Split(tokenHeard, " ")
	//tokenStr := strings.SplitN(tokenHeard, " ", 2)
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (r *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, rc)
	tokenStr, err := token.SignedString(Rtkey)
	if err != nil {
		return err
	}

	//fmt.Println(u)
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}
