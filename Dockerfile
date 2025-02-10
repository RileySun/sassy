FROM golang:1.24rc2-bookworm AS base

# Builder
# =============================================================================

FROM base AS builder

WORKDIR /build

COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o sassy


# Production stage
# =============================================================================

FROM scratch AS production

WORKDIR /prod

COPY ./admin/html/static ./admin/html/static

COPY --from=builder /build/sassy ./

EXPOSE 7070
EXPOSE 8080
EXPOSE 9090

# Start the application
CMD ["/prod/sassy"]