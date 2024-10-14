FROM golang:1.23.2 AS build

WORKDIR /app

COPY go.mod go.sum  ./

RUN go mod download

COPY *.go ./

RUN --mount=type=bind,target=. go build -o /repos .

FROM golang:1.23.2 AS base

WORKDIR /app

COPY  --from=build /repos /repos

LABEL authors="maliciousbucket"

EXPOSE 8080

ENTRYPOINT ["/repos"]