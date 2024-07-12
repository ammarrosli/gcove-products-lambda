package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func InitEnv() {
	_ = os.Setenv("DDBTABLE_PRODUCT", "dev-gc-products-ver2")
	_ = os.Setenv("MG_DOMAIN", "mailgun.digitalsymphony.it")
	_ = os.Setenv("MG_ADMIN_MAIL", "misyani@digitalsymphony.it")
	_ = os.Setenv("MG_ADMIN_SUBJECT", "New Registration Alert - Gamuda Cove Product")
	_ = os.Setenv("SENTINO_ENDPOINT", "https://www.sentinocrm.com/index.php/service2/register?ws=1")
	_ = os.Setenv("SENTINO_PROJECT_ID", "Jtd9tL1820180423115313")
	_ = os.Setenv("SENTINO_SOURCE_ID", "14")
	_ = os.Setenv("TELEGRAM_CHAT_ID", "-1001267327365")
	_ = os.Setenv("FUNCTION_NAME", "test-gc-product-landingpage")
}

func TestEDM(t *testing.T) {
	InitEnv()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	request := events.APIGatewayProxyRequest{
		Headers: headers,
		Body: `{
			"name":"TEST",
			"email":"misyani@digitalsymphony.it",
			"phone":"0123456789",
			"project":"No Preference",
			"source":"",
			"utm_sources":""
		}`,
	}

	response, _ := StartHandler(request)
	fmt.Println(response.Body)
	assert.Equal(t, 200, response.StatusCode)
}
