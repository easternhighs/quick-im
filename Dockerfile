FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/quick-im-demo .

FROM alpine:3.21
WORKDIR /app
COPY --from=build /out/quick-im-demo ./quick-im-demo
EXPOSE 8080
CMD ["./quick-im-demo"]
