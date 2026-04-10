# ======================== BUILD STAGE ========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ewerton-go .

# ======================== RUN STAGE ========================
FROM alpine

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/ewerton-go .

USER appuser

EXPOSE 3000

CMD ["./ewerton-go"]
