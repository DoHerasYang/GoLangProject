package web

import (
	"GoLangProject/webhook/internal/domain"
	"GoLangProject/webhook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	emailRegxPattern   = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	passwordRegxPatter = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,}$"
)

type UserHandler struct {
	// Precompiled to lift efficiency
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	svc              *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(emailRegxPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegxPatter, regexp.None),
		svc:              svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", h.SignUp)
	userGroup.POST("/login", h.Login)
	userGroup.GET("/profile", h.Profile)
	userGroup.POST("/edit", h.Edit)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		panic("Error Sign Up Context........")
		return
	}
	// Check Email Validated
	isEmail, emailError := h.emailRegexExp.MatchString(req.Email)
	if emailError != nil {
		ctx.String(http.StatusOK, "邮箱校验存在错误")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
	}
	// Check Password
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
	}
	// Check Password Validated
	//isPassword, passwordError := h.passwordRegexExp.MatchString(req.Password)
	//if passwordError != nil {
	//	ctx.String(http.StatusOK, "检验密码出现错误")
	//}
	//if !isPassword {
	//	ctx.String(http.StatusOK, "密码强度不够")
	//}
	err := h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "注册成功")
		return
	case errors.Is(err, service.DuplicatedError):
		ctx.String(http.StatusOK, "邮箱重复,请重新选择邮箱")
		return
	default:
		ctx.String(http.StatusInternalServerError, "服务器内部错误")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	//
	// 1234@admin.com
	// 123
	type SignInReq struct {
		Email    string `json:"email"`
		Password string `jsom:"password"`
	}
	var req SignInReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或密码错误")
	case err == nil:
		sess := sessions.Default(ctx)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		sess.Set("userId", u.ID) // Type Int64
		err = sess.Save()
		if err != nil {
			fmt.Println(err)
			ctx.String(http.StatusNotImplemented, "服务器内部错误")
			return
		}
		ctx.String(http.StatusOK, "登陆成功")
	default:
		ctx.String(http.StatusNotImplemented, "服务器内部错误")
	}

}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type EditStruct struct {
		NickName string `json:"nickname"`
		Birthday string `json:"birthday" `
		AboutMe  string `json:"aboutMe"`
	}
	var editreq EditStruct
	if err := ctx.Bind(&editreq); err != nil {
		ctx.String(http.StatusOK, "服务器内部错误请稍后重试")
		return
	}
	// Check NickName Field Input Length
	if len(editreq.NickName) > 15 {
		ctx.String(http.StatusOK, "昵称字段超出限制")
		return
	}
	// Check Birthday Field Input Format
	birthdayConverted, err := time.Parse(time.DateOnly, editreq.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "请检查输入Birthday格式，日期格式错误")
		return
	}
	// Check AboutMe Filed Input Length
	if len(editreq.AboutMe) > 100 {
		ctx.String(http.StatusOK, "介绍字段超出限制")
		return
	}
	// Obtain Session ID
	sess := sessions.Default(ctx)
	userID, ok := sess.Get("userId").(int64) // Type Int64 Assertion
	if !ok {
		ctx.String(http.StatusOK, "服务器内部错误请稍后重试")
		return
	}
	err = h.svc.Edit(ctx, userID, editreq.NickName, birthdayConverted, editreq.AboutMe)
	if err != nil {
		ctx.String(http.StatusOK, "服务器内部错误请稍后重试")
		return
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	// Obtain User ID
	sess := sessions.Default(ctx)
	userID, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "服务器内部错误请稍后重试")
		return
	}
	// Declare New Return
	rawUserProfile, err := h.svc.Profile(ctx, userID)
	// Handle With Not Found UserID
	if errors.Is(err, service.ErrUserProfileNotFound) {
		ctx.String(http.StatusOK, "未找到该用户信息，请重新初始化")
	}
	if err != nil {
		ctx.String(http.StatusOK, "服务器内部错误请稍后重试")
		return
	}
	type returnUserProfile struct {
		NickName string
		Birthday time.Time
		AboutMe  string
	}
	ctx.JSON(http.StatusOK, returnUserProfile{
		NickName: rawUserProfile.Nickname,
		Birthday: rawUserProfile.Birthday,
		AboutMe:  rawUserProfile.AboutMe,
	})
}
