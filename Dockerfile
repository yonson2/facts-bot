FROM golang:bullseye
ENV GOPROXY http://proxy.golang.org
ENV GO_ENV=production
WORKDIR /var/bot/
COPY . ./
RUN go mod download
RUN go build -o /bot ./main.go

CMD ["/bot"]
