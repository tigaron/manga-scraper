FROM golang:1.22.3-alpine3.20 as builder

# Install tzdata and ca-certificates
RUN apk add --no-cache tzdata ca-certificates
 
WORKDIR /workspace
 
# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download
 
# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch
 
COPY ./ ./
# generate the Prisma Client Go client
RUN go run github.com/steebchen/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# Generate .env file
RUN printenv > .env
 
# build a fully standalone binary with zero dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app ./cmd/manga-scraper/main.go
 
# use the scratch image for the smallest possible image size
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy .env file
COPY --from=builder /workspace/.env /.env
 
COPY --from=builder /workspace/app /app
 
ENTRYPOINT ["/app"]
 