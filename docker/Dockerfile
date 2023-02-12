FROM golang:1.19-alpine

WORKDIR /app

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go build -o /ethereum_subscriber ./cmd/api

EXPOSE 8080

CMD [ "/ethereum_subscriber" ]