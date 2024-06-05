# Stage 1: Update GeoIP Database
FROM maxmindinc/geoipupdate:latest AS geoipupdate
RUN geoipupdate

# Stage 2: ipwhitelist
FROM golang:1.22.3-alpine

# Install protoc
RUN apk add --no-cache protobuf protobuf-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Update PATH to include GOBIN (where go install puts binaries)
ENV PATH="$PATH:$(go env GOPATH)/bin" 

WORKDIR /app

COPY . .

# Generate gRPC code
RUN chmod +x ./generate_grpc.sh
RUN ./generate_grpc.sh

# Build go source code
RUN go mod tidy
RUN go build -o server ./cmd/server

EXPOSE 8080 50051

CMD ["./server"]