# build stage
FROM golang:1.23 as builder
WORKDIR /app

# copying dependencies description and downloading them
COPY go.mod go.sum ./
RUN go mod download

# compiling service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o image-previewer cmd/image-previewer/main.go

# compiling service
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/image-previewer .
COPY config/local.yml ./config/local.yml
RUN mkdir -p /app/tmp /app/storage
CMD ["./image-previewer"]