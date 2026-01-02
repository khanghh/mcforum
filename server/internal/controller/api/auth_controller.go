package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"bbs-go/internal/config"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/service"
	"bbs-go/pkg/web"

	"github.com/kataras/iris/v12"
)

type AuthController struct {
	Ctx iris.Context
}

type userInfoResponse struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Picture  string `json:"picture,omitempty"`
}

type authenticationSuccess struct {
	User         *userInfoResponse `json:"user"`
	AccessToken  string            `json:"accessToken,omitempty"`
	RefreshToken string            `json:"refreshToken,omitempty"`
}

type authenticationFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AuthResponse struct {
	Data struct {
		Success *authenticationSuccess `json:"authenticationSuccess,omitempty"`
		Failure *authenticationFailure `json:"authenticationFailure,omitempty"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GetCallback posts service,ticket,client_id,client_secret as form values to validateURL
func (c *AuthController) PostLogin() *web.JsonResult {
	ticket := c.Ctx.FormValue("ticket")
	if ticket == "" {
		return web.JsonError(errs.ErrBadRequest)
	}

	validateURL := config.Instance().Auth.ValidateURL
	serviceURL := config.Instance().Auth.ServiceURL
	clientID := config.Instance().Auth.ClientID
	clientSecret := config.Instance().Auth.ClientSecret

	form := url.Values{}
	form.Set("service", serviceURL)
	form.Set("ticket", ticket)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	resp, err := http.PostForm(validateURL, form)
	if err != nil {
		slog.Error("Validate ticket request failed", "error", err)
		return web.JsonError(errs.ErrInternalServer)
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		slog.Error("Failed to decode auth response", slog.Any("err", err))
		return web.JsonError(errs.ErrInternalServer)
	}

	if resp.StatusCode != http.StatusOK {
		if authResp.Error.Message != "" {
			code := authResp.Error.Code
			message := authResp.Error.Message
			return web.JsonErrorCodeMsg(code, message)
		}
		slog.Error("Auth service returned non-200 status", "status", resp.StatusCode)
		return web.JsonError(errs.ErrInternalServer)
	}

	if authResp.Data.Failure != nil {
		message := authResp.Data.Failure.Message
		return web.JsonErrorCodeMsg(http.StatusUnauthorized, message)
	}

	if authResp.Data.Success != nil && authResp.Data.Success.User != nil {
		userInfo := authResp.Data.Success.User
		userID, err := strconv.ParseInt(userInfo.UserID, 10, 64)
		if err != nil {
			slog.Error("Invalid user ID format", slog.Any("err", err))
			return web.JsonError(errs.ErrInternalServer)
		}
		user, err := service.UserService.OAuthLogin("cas", userID, userInfo.Username, userInfo.FullName, userInfo.Email, userInfo.Picture)
		if err != nil {
			slog.Error("Login error", "error", err)
			return web.JsonError(errs.ErrInternalServer)
		}
		return payload.BuildLoginSuccess(c.Ctx, user, "/")
	}

	return web.JsonData(errs.ErrInternalServer)
}

func (c *AuthController) PostLogout() *web.JsonResult {
	err := service.UserTokenService.Signout(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
