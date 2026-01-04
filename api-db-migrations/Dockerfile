FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/api /app/api
COPY migrations /app/migrations
ENV HTTP_ADDR=":8080"
EXPOSE 8080
CMD ["/app/api"]
