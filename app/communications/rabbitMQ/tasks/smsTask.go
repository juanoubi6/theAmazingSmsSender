package tasks

import (
	"encoding/json"
)

type SmsTask struct{
	Queue	 string
	SmsMessage SmsMessage
}

type SmsMessage struct{
	PhoneNumber string `json:"phone_number"`
	MessageInfo	string `json:"message_info"`
}

func NewSmsTask (phoneNumber string, message string) SmsTask{

	smsMessage := SmsMessage{
		PhoneNumber:phoneNumber,
		MessageInfo:message,
	}

	return SmsTask{
		Queue:"sms_queue",
		SmsMessage:smsMessage,
	}
}


func (t SmsTask) GetMessageBytes () ([]byte,error){

	data,err := json.Marshal(t.SmsMessage)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t SmsTask) GetQueue () (queueName string){
	return t.Queue
}
