FROM golang:1.8

RUN apt-get update && \
    apt-get install -y python3.9 python3-pip
COPY ./requirements.txt /app/requirements.txt
RUN pip3 install -r /app/requirements.txt

COPY . /go/src/app
WORKDIR /go/src/app

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]