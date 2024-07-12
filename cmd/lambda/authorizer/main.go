package main

import (
	"encoding/base64"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"strings"
)

/**************************************
MAIN HANDLER START
 ***************************************/
func StartHandler(req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	auth := req.AuthorizationToken

	if auth == "" {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(strings.Split(auth, " ")[1])
    if err != nil {
        return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
    }

	authString := strings.Split(string(rawDecodedText[:]), ":")
	username := authString[0]
	password := authString[1]

	if !(username == "admin" && password == "V05jAvQXt6w1%W^$y1") {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", req.MethodArn), nil
}

/**************************************
MAIN HANDLER END
 ***************************************/

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func main() {
	lambda.Start(StartHandler)
}
