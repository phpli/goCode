package web

import (
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	svc         *service.UserService
	codeSvc     *service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
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
	ug.POST("/edit", c.Edit)
	//ug.GET("/profile", c.Profile)
	ug.GET("/profile", c.ProfileJWT)
	ug.GET("/login_sms", c.loginSms)
	ug.POST("/login_sms/code/send", c.SendLoginSMSCode)
}

func (c *UserHandler) loginSms(ctx *gin.Context) {

}

func (c *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	const biz = "login"
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := c.codeSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "发送成功",
	})
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

	if errors.Is(err, service.ErrUserDuplicateEmail) {
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
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       u.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("fb0e22c79ac75679e9881e6ba183b354"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	//fmt.Println(u)
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登陆成功")
	return
}

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

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

// ProfileJWT 用户详情
func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	claim, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := claim.(*UserClaims)
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
