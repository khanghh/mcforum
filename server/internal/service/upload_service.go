package service

import (
	"bbs-go/internal/config"
	"bbs-go/internal/errs"
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"bbs-go/sqls"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang-jwt/jwt/v4"
)

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
	URL       string `json:"url"`
	FileName  string `json:"fileName"`
	ThumbName string `json:"thumbName"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
	CreatedAt int64  `json:"createdAt"`
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

func (s *uploadService) UploadStream(user *model.User, stream io.Reader, fileName string) (*model.Upload, error) {
	uploadURL := config.Instance().Uploader.SUpload.UploadURL
	if uploadURL == "" {
		slog.Error("Upload URL is not configured")
		return nil, errs.ErrInternalServer
	}
	uploadToken := s.getUploadToken(user)
	// build a streaming multipart request with field "file"
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	contentType := mw.FormDataContentType()

	// writer goroutine: write the multipart content and stream the file
	go func() {
		defer pw.Close()
		defer mw.Close()
		part, err := mw.CreateFormFile("file", fileName)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(part, stream); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	proxiedReq, err := http.NewRequest(
		http.MethodPost,
		uploadURL,
		pr,
	)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	proxiedReq.Header.Del("Host")
	proxiedReq.Header.Set("Authorization", "Bearer "+uploadToken)
	proxiedReq.Header.Set("Content-Type", contentType)

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
			FileName:  repsBody.Data.FileName,
			URL:       repsBody.Data.URL,
			Size:      repsBody.Data.Size,
			MimeType:  repsBody.Data.MimeType,
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

func (s *uploadService) UploadAvatar(user *model.User, data []byte) (*model.Upload, error) {
	reader := bytes.NewReader(data)
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, errs.ErrUnsupportedFileType
	}

	img = imaging.Thumbnail(img, 150, 150, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG); err != nil {
		return nil, errs.ErrInternalServer
	}

	filename := fmt.Sprintf("%d.jpg", time.Now().UnixMilli())
	uploadInfo, err := s.UploadStream(user, bytes.NewReader(buf.Bytes()), filename)
	if err != nil {
		slog.Error("Upload avatar failed", slog.Any("err", err))
		return nil, err
	}
	return uploadInfo, nil
}

func (s *uploadService) RecordUpload(upload *model.Upload) error {
	return repository.UploadRepository.Create(sqls.DB(), upload)
}
