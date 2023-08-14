FROM golang:1.20

ENV PATH $GOPATH/bin:$PATH
ENV CGO_ENABLED=1
ENV GO1111MODULE=on

WORKDIR /app

ADD . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/go-app

EXPOSE 8000 4000

CMD [ "/go/bin/dlv", "--listen=:4000", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/build/go-app" ]