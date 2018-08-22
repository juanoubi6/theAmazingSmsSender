package tasks

import (
	"encoding/json"
)

type PhoneCheckTask struct{
	Queue	 string
	PhoneCheckMessage PhoneCheckMessage
}

type PhoneCheckMessage struct{
	PhoneNumber string `json:"phone_number"`
}

type PhoneCheckTaskResponse struct {
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
	Error 		string `json:"error"`
}

func NewPhoneCheckTask (phoneNumber string) PhoneCheckTask{

	phoneCheckMessage := PhoneCheckMessage{
		PhoneNumber:phoneNumber,
	}

	return PhoneCheckTask{
		Queue:"phone_check_queue",
		PhoneCheckMessage:phoneCheckMessage,
	}
}


func (t PhoneCheckTask) GetMessageBytes () ([]byte,error){

	data,err := json.Marshal(t.PhoneCheckMessage)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t PhoneCheckTask) GetQueue () (queueName string){
	return t.Queue
}
