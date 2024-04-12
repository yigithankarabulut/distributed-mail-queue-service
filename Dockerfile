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

ENV DB_NAME=DMQS
ENV DB_USER=ykarabul
ENV DB_PASS=yigitsh
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_MIGRATE=true
ENV PORT=8080
ENV JWT_SECRET=secret
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379

EXPOSE 8080

CMD ["./main"]