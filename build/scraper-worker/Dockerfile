FROM golang:1.22.3-bookworm AS builder

WORKDIR /build/

COPY . .
RUN go mod download

RUN go run github.com/steebchen/prisma-client-go prefetch

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
ENV KAFKA_URL {$KAFKA_URL}
ENV KAFKA_USERNAME {$KAFKA_USERNAME}
ENV KAFKA_PASSWORD {$KAFKA_PASSWORD}

RUN printenv > .env

COPY ./ ./

RUN go run github.com/steebchen/prisma-client-go generate
 
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-extldflags -static" \
  fourleaves.studio/manga-scraper/cmd/scraper-worker

FROM debian:12.5-slim
RUN set -x && \
  apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/.env .
COPY --from=builder /build/scraper-worker ./bin/scraper-worker

CMD ["scraper-worker"]
