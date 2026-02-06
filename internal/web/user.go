package web

import (
	"errors"
	"net/http"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/ktsoator/connectify/internal/domain"
	"github.com/ktsoator/connectify/internal/service"
)

const (
	emailRegex    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	passwordRegex = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&.])[A-Za-z\d@$!%*?&.]{8,}$`
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		svc: service,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	rg := r.Group("/user")

	rg.POST("/signup", h.Signup)
	rg.POST("/login", h.Login)
	rg.GET("/profile", h.GetProfile)
	rg.PUT("/profile", h.UpdateProfile)
}

func (h *UserHandler) Signup(c *gin.Context) {
	type SignUpRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: CodeInvalidParam, Msg: "invalid request", Data: nil})
		return
	}

	ok, err := ValidateEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: CodeServerBusy, Msg: "system error", Data: nil})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{Code: CodeInvalidParam, Msg: "email format error", Data: nil})
		return
	}

	ok, err = ValidatePassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: CodeServerBusy, Msg: "system error", Data: nil})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{Code: CodeInvalidParam, Msg: "password format error", Data: nil})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, Result{Code: CodeInvalidParam, Msg: "passwords do not match", Data: nil})
		return
	}

	err = h.svc.Signup(c.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			c.JSON(http.StatusOK, Result{Code: CodeUserExist, Msg: "email already exists", Data: nil})
			return
		}
		c.JSON(http.StatusOK, Result{Code: CodeServerBusy, Msg: "system error", Data: nil})
		return
	}

	c.JSON(http.StatusOK, Result{Code: CodeSuccess, Msg: "user registered successfully", Data: nil})
}

func (h *UserHandler) Login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: CodeInvalidParam, Msg: "invalid request", Data: nil})
		return
	}

	err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			c.JSON(http.StatusOK, Result{Code: CodeInvalidCreds, Msg: "invalid email or password", Data: nil})
			return
		}

		c.JSON(http.StatusOK, Result{Code: CodeServerBusy, Msg: "system error", Data: nil})
		return
	}

	c.JSON(http.StatusOK, Result{Code: CodeSuccess, Msg: "user logged in successfully", Data: nil})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	type ProfileResponse struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	c.JSON(http.StatusOK, Result{
		Code: CodeSuccess,
		Msg:  "success",
		Data: ProfileResponse{
			Email:    "test@example.com",
			Nickname: "MockUser",
			Intro:    "This is a mock introduction.",
		},
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	type UpdateProfileRequest struct {
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: CodeInvalidParam, Msg: "invalid request", Data: nil})
		return
	}

	c.JSON(http.StatusOK, Result{Code: CodeSuccess, Msg: "user profile updated successfully", Data: req})
}

func ValidatePassword(password string) (bool, error) {
	re := regexp2.MustCompile(passwordRegex, 0)
	return re.MatchString(password)
}

func ValidateEmail(email string) (bool, error) {
	re := regexp2.MustCompile(emailRegex, 0)
	return re.MatchString(email)
}
