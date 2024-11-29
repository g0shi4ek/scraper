FROM golang

WORKDIR /scraper
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o scraper .


EXPOSE 2000

CMD ["./scraper"]