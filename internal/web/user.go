package web

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ok, err := ValidateEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "system error"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email format error"})
		return
	}

	ok, err = ValidatePassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "system error"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password format error"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	err = h.svc.Signup(c.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "system error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user logged in successfully"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	type ProfileResponse struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	c.JSON(http.StatusOK, ProfileResponse{
		Email:    "test@example.com",
		Nickname: "MockUser",
		Intro:    "This is a mock introduction.",
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	type UpdateProfileRequest struct {
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user profile updated successfully",
		"data":    req,
	})
}

func ValidatePassword(password string) (bool, error) {
	re := regexp2.MustCompile(passwordRegex, 0)
	return re.MatchString(password)
}

func ValidateEmail(email string) (bool, error) {
	re := regexp2.MustCompile(emailRegex, 0)
	return re.MatchString(email)
}
