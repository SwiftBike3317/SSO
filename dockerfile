FROM golang:1.21

WORKDIR /app

ADD . /app

COPY go.mod .
COPY main.go .


RUN go mod tidy 
RUN go build -o bin .

EXPOSE 3001

ENTRYPOINT ["/app/bin"]