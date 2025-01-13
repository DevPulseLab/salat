FROM golang:1.23.4-alpine3.20

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

RUN apk add --no-cache --update \
    nodejs \
    npm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/server/main.go

WORKDIR /app/frontend
RUN npm install
RUN npm run build
WORKDIR /app

EXPOSE 8088

CMD ["./main"]
