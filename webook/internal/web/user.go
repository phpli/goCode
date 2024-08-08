package web

import (
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const biz = "login"

// var _ handler = &UserHandler{}
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	ijwt.Handler
	cmd redis.Cmdable
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, cmd redis.Cmdable, jwtHdl ijwt.Handler) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None), //预编译
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		Handler:     jwtHdl,
		cmd:         cmd,
	}
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 直接注册
	//server.POST("/users/signup", c.SignUp)
	//server.POST("/users/login", c.Login)
	//server.POST("/users/edit", c.Edit)
	//server.GET("/users/profile", c.Profile)

	// 分组注册
	ug := server.Group("/users")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.LoginJWT)
	ug.POST("/logout", c.LogoutJWT)
	ug.POST("/edit", c.Edit)
	//ug.GET("/profile", c.Profile)
	ug.GET("/profile", c.ProfileJWT)
	ug.POST("/login_sms", c.loginSms)
	ug.POST("/login_sms/code/send", c.SendLoginSMSCode)
	ug.POST("/refresh_token", c.RefreshToken)
}

func (c *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := c.ClearToken(ctx)
	if err != nil {
		//要么redis有问题，要不已经退出登陆
		ctx.JSON(http.StatusOK, Result{
			Msg:  "退出登陆失败",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出成功",
	})
}

// RefreshToken 可以同时刷新长短 token，用 redis 来记录是否有效，即 refresh_token 是一次性的
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (c *UserHandler) RefreshToken(ctx *gin.Context) {

	refreshToken := c.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(*jwt.Token) (interface{}, error) {
		return ijwt.Rtkey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = c.CheckSession(ctx, rc.Ssid)
	////要是用redis session，没必要用这些
	//cnt, err := c.cmd.Exists(ctx, fmt.Sprintf("users:ssis:%s", rc.Ssid)).Result()
	if err != nil {
		//要么redis有问题，要不已经退出登陆
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = c.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}

func (c *UserHandler) loginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := c.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("codeSvc.Verify fail", zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	//使用jwt
	//使用jwt登陆态
	user, err := c.svc.FindOrCreate(ctx, req.Phone)
	if err = c.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "验证成功",
	})
}

func (c *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := c.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁,稍后重试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

// SignUp 用户注册接口
func (c *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `form:"email" json:"email"`
		Password        string `form:"password" json:"password"`
		ConfirmPassword string `form:"ConfirmPassword" json:"ConfirmPassword"`
	}
	var req SignUpReq
	//ShouldBind 方法尝试将请求体绑定到指定的结构体。如果绑定失败，
	//它不会立即返回错误，而是返回一个错误值，让你可以根据需要自行处理错误。
	//Bind可以直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := c.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusBadRequest, "邮箱错误")
		return
	}
	ok, err = c.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式有误")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入不一致")
		return
	}
	err = c.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrUserDuplicate) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	//fmt.Printf()
}

// Login 用户登录接口
func (c *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `form:"email" json:"email"`
		Password string `form:"password" json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: 30,
		//HttpOnly: true,
		//Secure: true,
	})
	session.Set("userId", u.Id)
	session.Save()
	ctx.String(http.StatusOK, "登陆成功")
	return
}

// jwt登陆
func (c *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `form:"email" json:"email"`
		Password string `form:"password" json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//使用jwt
	//使用jwt登陆态
	if err = c.SetLoginToken(ctx, u.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登陆成功")
	return
}

//func (c *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
//	claims := UserClaims{
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
//		},
//		Uid:       uid,
//		UserAgent: ctx.Request.UserAgent(),
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
//	tokenStr, err := token.SignedString([]byte("fb0e22c79ac75679e9881e6ba183b354"))
//	if err != nil {
//		return err
//	}
//
//	//fmt.Println(u)
//	ctx.Header("x-jwt-token", tokenStr)
//	return nil
//}

// Edit 用户编译信息
func (c *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Birthday    string `form:"birthday" json:"birthday"`
		Gender      int    `form:"gender" json:"gender"`
		Description string `form:"description" json:"description"`
		Nickname    string `form:"nickname" json:"nickname"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	session := sessions.Default(ctx)
	userID := session.Get("userId")
	id, err := strconv.ParseInt(fmt.Sprintf("%v", userID), 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "session is wrong")
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	err = c.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Birthday:    birthday,
		Gender:      req.Gender,
		Description: req.Description,
		Nickname:    req.Nickname,
		Id:          id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "修改失败")
		return
	}
	ctx.String(http.StatusOK, "登陆成功 %d", id)
}

// Profile 用户详情
func (c *UserHandler) Profile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("userId")
	id, err := strconv.ParseInt(fmt.Sprintf("%v", userID), 10, 64)
	if err != nil {
		ctx.String(http.StatusOK, "session is wrong")
	}
	user, err := c.svc.Profile(ctx, id)
	if errors.Is(err, service.ErrRecordNotFound) {
		ctx.String(http.StatusOK, "没有找到用户信息")
		return
	}
	ctx.String(http.StatusOK, "登陆成功 %d", user.Id)
}

func (c *UserHandler) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	//session.Delete("userId")
	session.Options(sessions.Options{
		MaxAge: -1,
	})
	session.Save()
	ctx.String(http.StatusOK, "退出成功")
}

// ProfileJWT 用户详情
func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	claim, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := claim.(*ijwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	id := claims.Uid
	_, err := c.svc.Profile(ctx, id)
	if errors.Is(err, service.ErrRecordNotFound) {
		ctx.String(http.StatusOK, "没有找到用户信息")
		return
	}
	ctx.String(http.StatusOK, "登陆成功 %d", id)
}

func (u *UserHandler) SessionLogout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录成功")
}
