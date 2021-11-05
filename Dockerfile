FROM golang:1.17.2

WORKDIR /home
COPY ./pkg /home

RUN cd /home && go build -o library

CMD ["/home/library"]