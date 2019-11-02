FROM golang:1.13-alpine
RUN apk add --no-cache tzdata
ENV TZ Europe/Minsk
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]