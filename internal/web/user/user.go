package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ktsoator/connectify/internal/domain"
	"github.com/ktsoator/connectify/internal/service"
	"github.com/ktsoator/connectify/internal/web/resp"
)

type UserClaims struct {
	UserId    int64
	UserEmail string
	UserAgent string
	jwt.RegisteredClaims
}

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

	// rg.POST("/login", h.Login)
	rg.POST("/login_jwt", h.LoginJwt)

	// rg.GET("/profile", h.GetProfile)
	rg.GET("/profile_jwt", h.GetProfileJwt)

	rg.PUT("/profile", h.UpdateProfile)

	rg.POST("/logout", h.Logout)

}

func (h *UserHandler) Signup(c *gin.Context) {
	type SignUpRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	ok, err := ValidateEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "email format error",
			Data: nil,
		})
		return
	}

	ok, err = ValidatePassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "password format error",
			Data: nil,
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "passwords do not match",
			Data: nil,
		})
		return
	}

	err = h.svc.Signup(c.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			c.JSON(http.StatusOK, resp.Result{
				Code: resp.CodeUserExist,
				Msg:  "email already exists",
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "user registered successfully",
		Data: nil,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	user, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			c.JSON(http.StatusOK, resp.Result{
				Code: resp.CodeInvalidCreds,
				Msg:  "invalid email or password",
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("userId", user.ID)
	session.Set("email", user.Email)
	session.Save()

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "user logged in successfully",
		Data: nil,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	type ProfileResponse struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	session := sessions.Default(c)
	uidVal := session.Get("userId")
	uid, ok := uidVal.(int64)
	if !ok {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "login expired or invalid session",
			Data: nil,
		})
		return
	}

	user, err := h.svc.Profile(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusOK, resp.Result{
				Code: resp.CodeUserNotFound,
				Msg:  "user not found",
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "success",
		Data: ProfileResponse{
			Email:    user.Email,
			Nickname: user.Nickname,
			Intro:    user.Intro,
		},
	})
}

func (h *UserHandler) GetProfileJwt(c *gin.Context) {
	type ProfileResponse struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	claim := h.MustGetUserClaims(c)
	if c.IsAborted() {
		return
	}

	user, err := h.svc.Profile(c.Request.Context(), claim.UserId)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusOK, resp.Result{
				Code: resp.CodeUserNotFound,
				Msg:  "user not found",
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "success",
		Data: ProfileResponse{
			Email:    user.Email,
			Nickname: user.Nickname,
			Intro:    user.Intro,
		},
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear() // Clear all session data
	session.Save()

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "user logged out successfully",
		Data: nil,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	type UpdateProfileRequest struct {
		Nickname string `json:"nickname"`
		Intro    string `json:"intro"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	session := sessions.Default(c)
	uidVal := session.Get("userId")
	uid, ok := uidVal.(int64)
	if !ok {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "login expired or invalid session",
			Data: nil,
		})
		return
	}

	err := h.svc.Update(c.Request.Context(), domain.User{
		ID:       uid,
		Nickname: req.Nickname,
		Intro:    req.Intro,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusOK, resp.Result{
				Code: resp.CodeUserNotFound,
				Msg:  "user not found",
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "user profile updated successfully",
		Data: nil,
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

func (h *UserHandler) LoginJwt(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeInvalidParam,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	user, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			c.JSON(http.StatusOK, resp.Result{
				Code: resp.CodeInvalidCreds,
				Msg:  "invalid email or password",
				Data: nil,
			})
			return
		}

		c.JSON(http.StatusOK, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		return
	}

	h.SetJwtToken(c, domain.User{
		ID:    user.ID,
		Email: user.Email})

	c.JSON(http.StatusOK, resp.Result{
		Code: resp.CodeSuccess,
		Msg:  "user logged in successfully",
		Data: nil,
	})
}

func (u *UserHandler) MustGetUserClaims(c *gin.Context) UserClaims {
	// Get user information from claim stored in context by middleware
	claimAny, exists := c.Get("claim")
	if !exists {
		c.JSON(http.StatusInternalServerError, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		c.Abort()
		return UserClaims{}
	}

	claim, ok := claimAny.(UserClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		c.Abort()
		return UserClaims{}
	}

	return claim
}

func (u *UserHandler) SetJwtToken(c *gin.Context, user domain.User) {
	claims := UserClaims{
		UserId:    user.ID,
		UserEmail: user.Email,
		UserAgent: c.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("Ktsoator"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Result{
			Code: resp.CodeServerBusy,
			Msg:  "system error",
			Data: nil,
		})
		c.Abort()
		return
	}

	c.Header("Jwt-Token", tokenStr)
}
