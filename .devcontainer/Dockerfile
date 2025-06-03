# Use the Go 1.24 devcontainer image
FROM mcr.microsoft.com/devcontainers/go:1.24

# Install PostgreSQL server & client libs, Neovim, git, curl, unzip
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  postgresql postgresql-contrib libpq-dev \
  neovim git curl unzip \
  && rm -rf /var/lib/apt/lists/* \
  \
  # Symlink Postgres binaries into /usr/local/bin so initdb, pg_ctl, etc. are on PATH
  && for pgbin in /usr/lib/postgresql/*/bin/*; do \
  ln -s "$pgbin" /usr/local/bin/; \
  done \
  \
  # Install a fixed migrate CLI version
  && go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

# Ensure Go and user binaries are on PATH
ENV GOPATH=/home/vscode/go
ENV PATH="/usr/local/go/bin:/home/vscode/go/bin:/usr/local/bin:${PATH}"

# Prepare Postgres data dir
USER root
RUN mkdir -p /var/lib/postgresql/data && chown -R postgres:postgres /var/lib/postgresql

# Switch to non-root user and set workdir
USER vscode
WORKDIR /workspace

# Copy init/start script (must be executable on host)
COPY init-postgres.sh /usr/local/bin/init-po
