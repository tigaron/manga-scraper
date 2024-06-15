FROM golang:1.22.3-alpine3.20 as builder

# Install tzdata and ca-certificates
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /workspace

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch

# Set environment variables
ENV ENVIRONMENT {$ENVIRONMENT}
ENV HTTP_PORT {$HTTP_PORT}
ENV DATABASE_URL {$DATABASE_URL}
ENV ROD_BROWSER_URL {$ROD_BROWSER_URL}
ENV ADMIN_SUB {$ADMIN_SUB}
ENV SENTRY_DSN {$SENTRY_DSN}
ENV REDIS_URL {$REDIS_URL}
ENV VERSION {$VERSION}
ENV OPENSEARCH_URL {$OPENSEARCH_URL}
ENV CLERK_SECRET_KEY {$CLERK_SECRET_KEY}

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
