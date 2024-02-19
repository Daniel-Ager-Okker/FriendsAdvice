FROM golang:latest

WORKDIR /Application

# Set ENVs for getting right connection to the DB
ENV DATABASE_HOST_NAME=172.25.32.1
ENV DATABASE_USER=postgres
ENV DATABASE_PASSWORD=postgres
ENV DATABASE_NAME=postgres
ENV DATABASE_PORT=5432

# Build application
COPY cmd ./cmd
COPY internal ./internal
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download
RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go build ./cmd/app/main.go

# CMD ["./main"]