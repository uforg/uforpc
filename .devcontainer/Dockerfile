# To make sure we have the golang binaries
FROM golang:1.24.2-bookworm AS golang

# Set the base image
FROM debian:12.10

# Declare ARG to make it available in the RUN commands
ARG TARGETPLATFORM
RUN echo "Building for ${TARGETPLATFORM}"
RUN if [ "${TARGETPLATFORM}" != "linux/amd64" ] && [ "${TARGETPLATFORM}" != "linux/arm64" ]; then \
      echo "Unsupported architecture: ${TARGETPLATFORM}" && \
      exit 1; \
    fi

# Set the general environment variables, and move to temp dir
ENV DEBIAN_FRONTEND="noninteractive"
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/usr/local/dl/bin"
ENV GOBIN="/usr/local/go/bin"
RUN mkdir -p /app/temp /usr/local/dl/bin
WORKDIR /app/temp

# Install golang from docker image
COPY --from=golang /usr/local/go /usr/local/go

# Install system dependencies
RUN apt update && \
    apt install -y wget curl zip unzip p7zip-full tzdata git npm && \
    rm -rf /var/lib/apt/lists/*

# Install Node.js using n
RUN npm install -g n@latest && \
    n 22.15.0

# Install downloadable binaries
RUN set -e && \
  if [ "${TARGETPLATFORM}" = "linux/arm64" ]; then \
    echo "Downloading arm64 binaries" && \
    # Install task
    wget --no-verbose https://github.com/go-task/task/releases/download/v3.43.2/task_linux_arm64.tar.gz && \
    tar -xzf task_linux_arm64.tar.gz && \
    mv ./task /usr/local/dl/bin/task && \
    # Install golangci-lint
    wget --no-verbose https://github.com/golangci/golangci-lint/releases/download/v2.1.5/golangci-lint-2.1.5-linux-arm64.tar.gz && \
    tar -xzf golangci-lint-2.1.5-linux-arm64.tar.gz && \
    mv ./golangci-lint-2.1.5-linux-arm64/golangci-lint /usr/local/dl/bin/golangci-lint; \
  else \
    echo "Downloading amd64 binaries" && \
    # Install task
    wget --no-verbose https://github.com/go-task/task/releases/download/v3.43.2/task_linux_amd64.tar.gz && \
    tar -xzf task_linux_amd64.tar.gz && \
    mv ./task /usr/local/dl/bin/task && \
    # Install golangci-lint
    wget --no-verbose https://github.com/golangci/golangci-lint/releases/download/v2.1.5/golangci-lint-2.1.5-linux-amd64.tar.gz && \
    tar -xzf golangci-lint-2.1.5-linux-amd64.tar.gz && \
    mv ./golangci-lint-2.1.5-linux-amd64/golangci-lint /usr/local/dl/bin/golangci-lint; \
  fi && \
  # Make binaries executable
  chmod +x /usr/local/dl/bin/*

# Go install binaries
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Default git config
# https://github.com/golangci/golangci-lint/issues/4033
RUN git config --global --add safe.directory '*'

# Go to the app dir, delete the temporary dir and create backups dir
WORKDIR /app
RUN rm -rf /app/temp && \
  mkdir /backups && \
  chmod 777 /backups

##############
# START HERE #
##############

# Add the startup script on every bash session
RUN rm -rf /root/.bashrc
COPY .devcontainer/.bashrc /root/.bashrc

# Command just to keep the container running
CMD ["sleep", "infinity"]
