FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/hello-api .

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /out/hello-api /hello-api

USER nonroot:nonroot
EXPOSE 9090

ENTRYPOINT ["/hello-api"]
