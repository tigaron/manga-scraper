FROM golang:1.22.3-bullseye as builder
 
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
 
# build a fully standalone binary with zero dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .
 
# use the scratch image for the smallest possible image size
FROM scratch
 
COPY --from=builder /workspace/app /app
 
ENTRYPOINT ["/app"]
 