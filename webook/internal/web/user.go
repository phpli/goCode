package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &UserHandler{
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
	ug.POST("/login", c.Login)
	ug.POST("/edit", c.Edit)
	ug.GET("/profile", c.Profile)
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
	ctx.String(http.StatusOK, "注册成功")
	//fmt.Printf()
}

// Login 用户登录接口
func (c *UserHandler) Login(ctx *gin.Context) {

}

// Edit 用户编译信息
func (c *UserHandler) Edit(ctx *gin.Context) {

}

// Profile 用户详情
func (c *UserHandler) Profile(ctx *gin.Context) {

}
