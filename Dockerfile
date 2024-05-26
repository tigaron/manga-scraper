FROM golang:1.22.3-bullseye

WORKDIR /workspace
 
# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download
 
# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch
 
COPY . ./
 
# generate .env file
RUN printenv > .env

# generate the Prisma Client Go client
RUN go run github.com/steebchen/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /app .
 
CMD ["/app"]