FROM golang:1.24-bookworm as builder

# Move to working directory /build
WORKDIR /builder

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
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app/

COPY --from=builder /builder/app .
COPY ubigeos.json .
COPY scripts .

RUN chown -R appuser:appgroup /app && chmod 644 ubigeos.json

USER appuser

# Add provenance metadata label (minimal step)
LABEL org.opencontainers.image.source="https://github.com/erajuan/decolecta-ruc"
LABEL org.opencontainers.image.created=$BUILD_DATE

CMD ["./app"]
