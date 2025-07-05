FROM golang:1.24-bookworm as builder

# Move to working directory /build
WORKDIR /home

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# copy to small linux
FROM alpine:3.22
RUN apk --no-cache add ca-certificates
WORKDIR /root/
RUN apk --no-cache add tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /home/app .

VOLUME /data

USER appuser

# Add provenance metadata label (minimal step)
LABEL org.opencontainers.image.source="https://github.com/erajuan/decolecta-auth"
LABEL org.opencontainers.image.created=$BUILD_DATE

CMD ["./app"]
