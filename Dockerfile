# Build the manager binary
FROM golang:1.15 as builder

# ENV TEMPLATE_PATH=/templates

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY manifests/ manifests/
COPY vendor/ vendor/

# Build
# https://www.reddit.com/r/golang/comments/9ai79z/correct_usage_of_go_modules_vendor_still_connects/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod vendor -a -o manager main.go \
    && rm -rf "${GOPATH}"

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
# COPY config/templates/deployment.yaml ${TEMPLATE_PATH}/deployment.yaml
# COPY config/templates/configmap.yaml ${TEMPLATE_PATH}/configmap.yaml
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
