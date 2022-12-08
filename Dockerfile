FROM golang:1.19 as build

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o /go/bin/app cmd/notionboy/main.go

FROM gcr.io/distroless/static-debian11
WORKDIR /service/

COPY --from=build /go/bin/app .
COPY settings.yaml settings.yaml
CMD ["./app"]
