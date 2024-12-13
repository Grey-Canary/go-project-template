## Multistage build
FROM golang:1.23 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN  go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/project ./cmd/project

## Multistage deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /build/project .
COPY --from=build /build/.env.test .env

EXPOSE 5000

# Run the project api command
CMD ["./project", "api"]