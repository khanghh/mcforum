package api

import (
	"bbs-go/internal/config"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model"
	"bytes"
	"io"
	"net"
	"time"

	"bbs-go/pkg/web"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"

	"bbs-go/internal/service"
	"encoding/json"
	"net/http"
)

type UploadController struct {
	Ctx iris.Context
}

type UploadClaims struct {
	jwt.RegisteredClaims
	User string `json:"user,omitempty"`
	Role string `json:"role,omitempty"`
}

func (s *UploadController) getUploadToken(user *model.User) string {
	secret := config.Instance().Uploader.SUpload.Secret
	now := time.Now()
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	claims := &UploadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.Username.String,
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
		},
		User: user.Username.String,
		Role: roleName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return signedToken
}

func (c *UploadController) Post() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	uploadURL := config.Instance().Uploader.SUpload.UploadURL
	if uploadURL == "" {
		return web.JsonError(errs.ErrInternalServer)
	}
	uploadToken := c.getUploadToken(user)
	srcReq := c.Ctx.Request()
	ctx := srcReq.Context()
	proxiedReq, err := http.NewRequestWithContext(
		ctx,
		srcReq.Method,
		uploadURL,
		srcReq.Body, // streamed, NOT buffered
	)
	if err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}

	// Copy headers safely
	for k, vv := range srcReq.Header {
		for _, v := range vv {
			proxiedReq.Header.Add(k, v)
		}
	}
	// Optional: ensure Content-Type exists
	if proxiedReq.Header.Get("Content-Type") == "" {
		proxiedReq.Header.Set("Content-Type", "application/octet-stream")
	}

	proxiedReq.Header.Del("Host") // let client set it
	proxiedReq.Header.Set("Authorization", "Bearer "+uploadToken)

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 2 * time.Minute,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Do(proxiedReq)
	if err != nil {
		return web.JsonError(errs.ErrBadGateway)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	var repsBody struct {
		Data  payload.UploadResponse `json:"data"`
		Error interface{}            `json:"error"`
	}
	if err := json.NewDecoder(tee).Decode(&repsBody); err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}

	if resp.StatusCode == http.StatusCreated {
		service.UploadService.RecordUpload(
			&model.Upload{
				UserID:    user.ID,
				URL:       repsBody.Data.URL,
				Size:      repsBody.Data.Size,
				MimeType:  repsBody.Data.MimeType,
				FileName:  repsBody.Data.FileName,
				CreatedAt: time.Now().UnixMilli(),
			},
		)
	}

	// send response to client
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		c.Ctx.ResponseWriter().Header().Set("Content-Type", ct)
	}
	c.Ctx.ResponseWriter().WriteHeader(resp.StatusCode)
	io.Copy(c.Ctx.ResponseWriter(), &buf)
	return nil
}
