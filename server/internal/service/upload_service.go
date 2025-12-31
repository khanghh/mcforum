package service

import (
	"bbs-go/internal/config"
	"bbs-go/internal/errs"
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"bbs-go/sqls"
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const ()

var UploadService = newUploadService()

type uploadService struct {
}

func newUploadService() *uploadService {
	return &uploadService{}
}

type UploadClaims struct {
	jwt.RegisteredClaims
	User string `json:"user,omitempty"`
	Role string `json:"role,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UploadResponse struct {
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	FileName string `json:"file_name"`
}

func (s *uploadService) getUploadToken(user *model.User) string {
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

func (s *uploadService) Upload(user *model.User, stream io.Reader, mimeType string) (*model.Upload, error) {
	uploadURL := config.Instance().Uploader.SUpload.UploadURL
	if uploadURL == "" {
		return nil, errs.ErrInternalServer
	}
	uploadToken := s.getUploadToken(user)

	proxiedReq, err := http.NewRequest(
		http.MethodPost,
		uploadURL,
		stream,
	)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	if proxiedReq.Header.Get("Content-Type") == "" {
		proxiedReq.Header.Set("Content-Type", mimeType)
	}

	proxiedReq.Header.Del("Host")
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
		return nil, errs.ErrBadGateway
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	var repsBody struct {
		Data  UploadResponse `json:"data"`
		Error ErrorResponse  `json:"error"`
	}
	if err := json.NewDecoder(tee).Decode(&repsBody); err != nil {
		return nil, errs.ErrInternalServer
	}

	if resp.StatusCode == http.StatusCreated {
		record := &model.Upload{
			UserID:    user.ID,
			URL:       repsBody.Data.URL,
			Size:      repsBody.Data.Size,
			MimeType:  repsBody.Data.MimeType,
			FileName:  repsBody.Data.FileName,
			CreatedAt: time.Now().UnixMilli(),
		}
		s.RecordUpload(record)
		return record, nil
	}

	if repsBody.Error.Message != "" {
		return nil, errs.NewResponseError(repsBody.Error.Code, repsBody.Error.Message)
	}
	return nil, errs.ErrBadGateway
}

func (s *uploadService) RecordUpload(upload *model.Upload) error {
	return repository.UploadRepository.Create(sqls.DB(), upload)
}
