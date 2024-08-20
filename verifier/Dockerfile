FROM golang:1.22.1-alpine3.18 AS builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# Copy the necessary files and directories to the /app directory in the container
COPY core /app/core
COPY api /app/api
COPY contracts /app/contracts
COPY go.mod /app
COPY go.sum /app
COPY keys /app/keys
COPY verifier /app/verifier

# # Set the working directory to /app/verifier where the Go project resides
WORKDIR /app/verifier

# # Build the Go project within the verifier directory
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/verifier

# # Use a smaller base image for the final stage
FROM alpine:3.18

# # Copy the built binary to the final image
COPY --from=builder /app/verifier/bin/verifier /usr/local/bin

# # Set the entry point to the built verifier binary
ENTRYPOINT ["verifier"]

CMD ["-config=config.yaml" ]
