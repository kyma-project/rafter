FROM golang:1.12.5 as builder

ENV BASE_APP_DIR=/go/src/github.com/kyma-project/rafter \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR ${BASE_APP_DIR}

# Copy the Go Modules manifests
COPY ./go.mod .
COPY ./go.sum .

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

#
# Copy files
#
COPY . ${BASE_APP_DIR}/

#
# Build app
#
RUN go build -a -o main cmd/extension/asyncapi/main.go \
    && mkdir /app \
    && mv ./main /app/main \
    && mv ${BASE_APP_DIR}/licenses /app/licenses

FROM gcr.io/distroless/static:latest

LABEL source = git@github.com:kyma-project/kyma.git

COPY --from=builder /app /app

CMD ["/app/main"]