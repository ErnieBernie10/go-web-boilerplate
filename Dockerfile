FROM golang:1.23.2-bookworm

# Create a user for devcontainer
RUN  apt-get update && apt-get install -y make && \
    rm -rf /var/lib/apt/lists/*

# Set the work directory
WORKDIR /workspace

# Install necessary Go tools
RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/amacneil/dbmate@latest
