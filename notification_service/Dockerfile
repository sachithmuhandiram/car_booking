FROM golang:alpine as notification

RUN apk update
#ADD . /go/src/car_booking/notification_service/packages/notification
COPY . /go/src/

#RUN go run src/packages/notification/notification.go
# RUN go build -o notification_sender src/packages/notification/notification.go
# CMD [ "./notification_sender" ]