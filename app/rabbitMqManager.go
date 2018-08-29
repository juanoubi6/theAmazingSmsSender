package app

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
	"theAmazingSmsSender/app/common"
	"theAmazingSmsSender/app/config"
	"theAmazingSmsSender/app/helpers/twilio"
	"theAmazingSmsSender/app/communications/rabbitMQ/tasks"
)

var workerAmount, _ = strconv.Atoi(config.GetConfig().WORKER_AMOUNT)

func ConsumeSmsQueue() {

	for i := 0; i < workerAmount; i++ {
		go consumeSmsQueueWorker(i)
	}

}

func ConsumePhoneCheckQueue() {

	ch := common.GetRabbitMQChannel()
	defer ch.Close()

	//Queue declared but not needed if created previously
	queue, err := ch.QueueDeclare(
		"phone_check_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "declaring exchange",
		}).Info(err.Error())
	}

	//Allow rabbitMQ to send me as many messages as workers I have
	err = ch.Qos(
		workerAmount,
		0,
		false,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "prefetch count",
		}).Info(err.Error())
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "consuming message",
		}).Info(err.Error())
	}

	forever := make(chan bool)

	for i := 0; i < workerAmount; i++ {
		go checkPhone(msgs, ch)
	}

	log.Printf("Waiting for phone check tasks")
	<-forever

}

func smsSend(messageChannel <-chan amqp.Delivery) {
	for d := range messageChannel {
		println("Received SMS task")
		var smsMessageData tasks.SmsMessage
		err := json.Unmarshal(d.Body, &smsMessageData)
		if err != nil {
			log.WithFields(log.Fields{
				"place": "decoding message body",
			}).Info(err.Error())
		}
		err = twilio.SendVerificationSMS(smsMessageData.MessageInfo, smsMessageData.PhoneNumber)
		if err != nil {
			log.WithFields(log.Fields{
				"place": "sending sms",
			}).Info(err.Error())
		}
		d.Ack(false)
	}
}

func checkPhone(messageChannel <-chan amqp.Delivery, ch *amqp.Channel) {
	for d := range messageChannel {
		println("Received phone check task")
		var phoneCheckMessageData tasks.PhoneCheckMessage
		err := json.Unmarshal(d.Body, &phoneCheckMessageData)
		if err != nil {
			log.WithFields(log.Fields{
				"place": "decoding message body",
			}).Info(err.Error())
		}

		checkPhoneResult := twilio.ValidatePhoneNumber(phoneCheckMessageData.PhoneNumber)
		jsonResponse, err := json.Marshal(checkPhoneResult)
		if err != nil {
			log.WithFields(log.Fields{
				"place": "marshaling json response",
			}).Info(err.Error())
		}

		err = ch.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          jsonResponse,
			})
		if err != nil {
			log.WithFields(log.Fields{
				"place": "sending rpc response",
			}).Info(err.Error())
		}

		d.Ack(false)
	}
}

func consumeSmsQueueWorker(number int){

	ch := common.GetRabbitMQChannel()

	//Queue declared but not needed if created previously
	queue, err := ch.QueueDeclare(
		"sms_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "declaring exchange",
		}).Info(err.Error())
	}

	//Allow rabbitMQ to send me as many messages as workers I have
	err = ch.Qos(
		workerAmount,
		0,
		false,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "prefetch count",
		}).Info(err.Error())
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "consuming message",
		}).Info(err.Error())
	}

	forever := make(chan bool)

	for d := range msgs {
		var smsMessageData tasks.SmsMessage
		err := json.Unmarshal(d.Body, &smsMessageData)
		if err != nil {
			log.WithFields(log.Fields{
				"place": "decoding message body",
			}).Info(err.Error())
		}
		err = twilio.SendVerificationSMS(smsMessageData.MessageInfo, smsMessageData.PhoneNumber)
		if err != nil {
			log.WithFields(log.Fields{
				"place": "sending sms",
			}).Info(err.Error())
		}
		d.Ack(false)
	}

	log.Printf("Waiting for SMS tasks on queue " + strconv.Itoa(number))
	<-forever

}