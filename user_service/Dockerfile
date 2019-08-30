FROM golang:alpine as user_service

RUN apk update

RUN apk add git

# ADD . /go/src/car_booking
RUN set -e \
    go get -u golang.org/x/lint/golint \
    && go get github.com/dgrijalva/jwt-go     

RUN  go get github.com/go-sql-driver/mysql

#RUN rc-service mariadb restart
RUN go get -u golang.org/x/crypto/...
#RUN mysqladmin -u root password "${DB_ROOT_PASS}"
COPY . /go/src
WORKDIR /go/src/user_service

# RUN pwd

# CMD [ "go","run","main.go" ]



