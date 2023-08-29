FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . ./

RUN ./scripts/docs
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/main ./cmd/segment/main.go


FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /home/nonroot

COPY --from=build-stage /build/main  .

USER nonroot:nonroot

ENTRYPOINT ["./main"]
