FROM golang:1.17.8-alpine3.15

RUN mkdir seatReservation

WORKDIR /seatReservation

COPY . .

RUN export GO111MODULE=on
RUN cd /seatReservation
RUN go build -o main.exe

EXPOSE 8080

CMD [ "/seatReservation/main.exe" ]

