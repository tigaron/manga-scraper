FROM golang:1.22.3-alpine3.20 as builder

# Install tzdata and ca-certificates
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /workspace

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch

# Set arguments
ARG DATABASE_URL
ARG ENVIRONMENT
ARG HTTP_PORT
ARG VERSION
ARG ADMIN_SUB
ARG ROD_BROWSER_URL
ARG DATABASE_URL
ARG REDIS_URL
ARG SENTRY_DSN
ARG AUTH0_DOMAIN
ARG AUTH0_CLIENT_ID
ARG AUTH0_CLIENT_SECRET
ARG AUTH0_CALLBACK_URL

# Set environment variables
ENV DATABASE_URL {$DATABASE_URL}
ENV ENVIRONMENT {$ENVIRONMENT}
ENV HTTP_PORT {$HTTP_PORT}
ENV VERSION {$VERSION}
ENV ADMIN_SUB {$ADMIN_SUB}
ENV ROD_BROWSER_URL {$ROD_BROWSER_URL}
ENV DATABASE_URL {$DATABASE_URL}
ENV REDIS_URL {$REDIS_URL}
ENV SENTRY_DSN {$SENTRY_DSN}
ENV AUTH0_DOMAIN {$AUTH0_DOMAIN}
ENV AUTH0_CLIENT_ID {$AUTH0_CLIENT_ID}
ENV AUTH0_CLIENT_SECRET {$AUTH0_CLIENT_SECRET}
ENV AUTH0_CALLBACK_URL {$AUTH0_CALLBACK_URL}

# Generate .env file
RUN printenv > .env

COPY ./ ./
# generate the Prisma Client Go client
RUN go run github.com/steebchen/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...
 
# build a fully standalone binary with zero dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app ./cmd/manga-scraper/main.go

# use the scratch image for the smallest possible image size
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy .env file
COPY --from=builder /workspace/.env /.env

COPY --from=builder /workspace/app /app

ENTRYPOINT ["/app"]
