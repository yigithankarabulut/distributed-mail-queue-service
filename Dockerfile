FROM golang:alpine as builder

LABEL maintainer="Yigithan Karabulut <yigithannkarabulutt@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .


ENV DB_NAME=your-secret
ENV DB_USER=your-secret
ENV DB_PASS=your-secret
ENV DB_HOST=your-secret
ENV DB_PORT=your-secret
ENV DB_MIGRATE=true
ENV PORT=your-secret
ENV JWT_SECRET=your-secret
ENV REDIS_HOST=your-secret
ENV REDIS_PORT=your-secret

EXPOSE 8080

CMD ["./main"]