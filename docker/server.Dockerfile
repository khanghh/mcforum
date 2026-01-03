# ---- Build server ----
FROM golang:1.25.1-alpine AS server-builder

WORKDIR /workdir

RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

COPY ./server/go.mod ./server/go.sum ./
RUN go mod download

COPY ./server ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -v -o bbs-go main.go

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

ENV TZ=Asia/Ho_Chi_Minh

RUN apk add --no-cache libc6-compat
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=server-builder /workdir/bbs-go /app/bbs-go
COPY --from=server-builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 3000

ENTRYPOINT ["/app/bbs-go"]
