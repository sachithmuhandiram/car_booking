FROM golang:alpine as notification

RUN apk update
ADD . /go/src/car_booking/notification_service/packages/notification

WORKDIR /go/src/car_booking/notification_service/packages/notification

# COPY . .

#RUN pwd

#RUN ls
RUN go run /go/src/car_booking/notification_service/packages/notification/notification.go
#build -o notification_sender notification.go
CMD [ "./notification_sender" ]