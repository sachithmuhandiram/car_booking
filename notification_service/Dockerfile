FROM golang:alpine as notification

RUN apk update
#ADD . /go/src/car_booking/notification_service/packages/notification
COPY . /go/src/

#RUN go run src/packages/notification/notification.go
# RUN go build -o main src/main.go

EXPOSE 7070

# RUN chmod 755 main

# CMD [ "./main" ]