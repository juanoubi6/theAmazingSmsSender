package main

import (
	"theAmazingSmsSender/app"
	"theAmazingSmsSender/app/common"
)

func main() {
	common.ConnectToRabbitMQ()
	keep := make(chan bool)
	go app.ConsumeSmsQueue()
	go app.ConsumePhoneCheckQueue()
	<-keep
}
