FROM golang

RUN mkdir /app 
ADD . /app/
WORKDIR /app

RUN go get github.com/lib/pq
RUN go install github.com/lib/pq

RUN go build -o main .
CMD ["./main"]

