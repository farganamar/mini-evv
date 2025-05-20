ARG GO_VERSION=1.23
#BUILDER
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && \ 
  apk --update add git make build-base

WORKDIR /app
COPY go.mod go.sum /
RUN go mod download

COPY . .

RUN GOFLAGS="-buildvcs=false" go generate ./...
RUN GOFLAGS="-buildvcs=false" go build -o main .


# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && apk --no-cache add ca-certificates && \
  apk --update --no-cache add tzdata

WORKDIR /app 

COPY --from=builder /app/main .

COPY --from=builder /app/makefile .


# Create the db directory and copy the SQLite database file
RUN mkdir -p /app/db
COPY --from=builder /app/db/mini-evv.db /app/db/

# Create a default .env file with minimal required settings
RUN echo 'APP_CORS_ALLOWCREDENTIALS=true\n\
  APP_CORS_ALLOWEDHEADERS=Accept,Authorization,Content-Type\n\
  APP_CORS_ALLOWEDMETHODS=GET,PUT,POST,PATCH,DELETE,OPTIONS\n\
  APP_CORS_ALLOWEDORIGINS=http://localhost:3000,http://localhost:8080,http://127.0.01:8080,http://192.168.1.10:8080,http://localhost:5173\n\
  APP_CORS_ENABLE=true\n\
  APP_CORS_MAXAGESECONDS=300\n\
  APP_TZ=UTC\n\
  APP_URL=http://localhost:3200\n\
  \n\
  SERVER_ENV=production\n\
  SERVER_LOGLEVEL=info\n\
  SERVER_PORT=3200\n\
  SERVER_SHUTDOWN_CLEANUP_PERIOD_SECONDS=15\n\
  SERVER_SHUTDOWN_GRACE_PERIOD_SECONDS=15\n\
  \n\
  ACCESSTOKEN_SECRET=b63a00ecd0c76b088e14a5a53a211c40417bf72a8b513bb75285b13ed77d5f03\n\
  ACCESSTOKEN_EXPIRYINHOUR=720\n\
  REFRESHTOKEN_SECRET=b63a00ecd0c76b088e14a5a53a211c40417bf72a8b513bb75285b13ed77d5f03\n\
  REFRESHTOKEN_EXPIRYINHOUR=4320\n\
  \n\
  # SQLite Configuration\n\
  DB_SQLITE_PATH=db/mini-evv.db\n\
  \n\
  # Redis Configuration\n\
  CACHE_REDIS_PRIMARY_HOST=redis-acty-jbde-462157.leapcell.cloud\n\
  CACHE_REDIS_PRIMARY_PORT=6379\n\
  CACHE_REDIS_PRIMARY_PREFIX=USER:\n\
  CACHE_REDIS_PRIMARY_RETRY_COUNT=5\n\
  CACHE_REDIS_PRIMARY_PASSWORD=Ae00000QYJcj2WqM8U89ptXxrh2w/H6Hbq9CQWmHazo7Mhb0CMImWKWnKcCUoph33rkaahs' > .env

EXPOSE 3200
# Command to run the executable
CMD ["./main"]