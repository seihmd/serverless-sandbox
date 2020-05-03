package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/guregu/dynamo"
	"os"
)

var (
	snsArn    string
	tableName string
	snsClient *sns.SNS
	dbClient  *dynamo.DB
)

func init() {
	snsArn = os.Getenv("ENV_SNS_ARN")
	tableName = os.Getenv("ENV_DYNAMO_TABLE")
	hostName := os.Getenv("LOCALSTACK_HOSTNAME")

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://" + hostName + ":4575"),
		//Region:   aws.String(endpoints.UsEast1RegionID),
	})

	if err != nil {
		panic(err)
	}

	snsClient = sns.New(sess)

	dbClient = dynamo.New(sess, &aws.Config{
		Endpoint: aws.String("http://" + hostName + ":4569"),
		Region:   aws.String(endpoints.UsEast1RegionID),
	})
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type RegisterRequest struct {
	UserId string `json:"user_id" dynamo:"user_id"`
	Token  string `json:"token" dynamo:"token"`
	OS     string `json:"os" dynamo:"os"`
	SNSArn string `dynamo:"sns_arn"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	var registerRequest RegisterRequest

	err := json.Unmarshal([]byte(request.Body), &registerRequest)
	if err != nil {
		return Response{StatusCode: 401}, err
	}

	if !isValidRequest(registerRequest) {
		return Response{StatusCode: 400}, errors.New("invalid request")
	}

	arn, err := createSNSEndPoint(registerRequest.Token)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	registerRequest.SNSArn = *arn
	err = registerToDynamoDB(registerRequest)

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	return Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

func isValidRequest(request RegisterRequest) bool {
	if request.Token == "" {
		return false
	}
	if request.OS != "android" && request.OS != "ios" {
		return false
	}
	if request.UserId == "" {
		return false
	}

	return true
}

func createSNSEndPoint(token string) (*string, error) {
	input := sns.CreatePlatformEndpointInput{
		Attributes:             map[string]*string{},
		CustomUserData:         aws.String(""),
		PlatformApplicationArn: aws.String(snsArn),
		Token:                  aws.String(token),
	}

	o, err := snsClient.CreatePlatformEndpoint(&input)
	if err != nil {
		return nil, err
	}

	return o.EndpointArn, nil
}

func registerToDynamoDB(request RegisterRequest) error {
	return dbClient.Table(tableName).Put(request).Run()
}
