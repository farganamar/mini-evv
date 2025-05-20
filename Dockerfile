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

EXPOSE 3200
# Command to run the executable
CMD ["./main"]