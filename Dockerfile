FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

FROM scratch

LABEL org.opencontainers.image.source=https://github.com/iammonika10/TP-DevOps_Rhina_LACMOU

COPY --from=builder /app/backend /backend

EXPOSE 8080

CMD ["/backend"]