FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN mkdir /usr/local/bin/app

COPY . .
RUN go build -o /usr/local/bin/app ./...

#---------------------------------

FROM alpine

COPY --from=builder /usr/local/bin/app/web /usr/local/bin/app

CMD ["app"]