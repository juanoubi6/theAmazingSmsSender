package twilio

import (
	"github.com/gin-gonic/gin/json"
	"github.com/subosito/twilio"
	"net/http"
	"theAmazingSmsSender/app/config"
)

var (
	AccountSid   = config.GetConfig().TWILIO_SID
	AuthToken    = config.GetConfig().TWILIO_AUTH_TOKEN
	AccountPhone = config.GetConfig().TWILIO_ACC_PHONE
)

type CheckPhoneResult struct {
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
	Error       string `json:"error"`
}

func SendVerificationSMS(verificationCode string, to string) error {

	// Initialize twilio client
	c := twilio.NewClient(AccountSid, AuthToken, nil)

	// Send Message
	params := twilio.MessageParams{
		Body: "Your verification code is: " + verificationCode,
	}
	_, _, err := c.Messages.Send(AccountPhone, to, params)
	if err != nil {
		return err
	}

	return nil

}

func ValidatePhoneNumber(number string) CheckPhoneResult {

	//Create client
	var result CheckPhoneResult

	client := &http.Client{}

	//Create request
	request, err := http.NewRequest(http.MethodGet, "https://"+AccountSid+":"+AuthToken+"@lookups.twilio.com/v1/PhoneNumbers/"+number, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	//Fetch Request
	response, err := client.Do(request)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	//Check response
	if response.StatusCode != http.StatusOK {
		result.Error = "Invalid phone"
		return result
	} else {
		return result
	}

}
