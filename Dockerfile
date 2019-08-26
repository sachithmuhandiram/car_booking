FROM golang:alpine as car_booking

RUN apk update
RUN apk add git

# ADD . /go/src/car_booking
RUN set -e \
    go get -u golang.org/x/lint/golint \
    && go get github.com/dgrijalva/jwt-go     

RUN  go get github.com/go-sql-driver/mysql

#RUN rc-service mariadb restart
RUN go get -u golang.org/x/crypto/...
COPY . /go/src/car_booking
WORKDIR /go/src/car_booking/

CMD ["go","run","main.go"]
# RUN pwd

# RUN go build -o main .

# RUN chmod 755 main

# CMD [ "./main" ]