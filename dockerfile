FROM golang:1.23.1 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -tags qa -o ./messageapi

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /build/messageapi ./messageapi
CMD [ "/messageapi" ]