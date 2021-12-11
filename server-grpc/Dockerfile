FROM golang:latest 

WORKDIR / 

COPY . . 

EXPOSE 5051 

RUN go build server.go 

CMD ["./server"]