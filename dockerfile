FROM golang:1.24.1-alpine

RUN apk add --no-cache git gcc musl-dev

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

WORKDIR /app

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /bank-api ./cmd/server/

EXPOSE 8080
CMD ["/bank-api"]
