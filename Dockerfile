# Got this idea from https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# I wanted to just build from a scratch image, but I figured being more secure never harmed anyone...

############################
# STEP 1 build executable binary
############################
FROM golang:1.21.0-alpine3.18 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/connectrn/eval/
COPY . .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-ldflags='-w -s -extldflags "-static"' -a \
	-o /go/bin/eval cmd/connectrn/evaluation.go

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/eval /go/bin/eval

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 8080

ENV DB_HOST=host.docker.internal

# Run the eval binary.
ENTRYPOINT ["/go/bin/eval", "server"]
