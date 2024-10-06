FROM golang:1.23.2-bookworm

# Create a user for devcontainer
RUN apt-get update && apt-get install -y make stow postgresql-common
RUN yes | /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh
RUN apt update && apt install -y postgresql-client libpq-dev

# Set the work directory
WORKDIR /workspace

# Install necessary Go tools
RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/amacneil/dbmate@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
