package main

import (
	"context"
	"log"

	b64 "encoding/base64"

	"com.github.kudos-writer/model"
	"com.github.kudos-writer/repos"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	qson "github.com/joncalhoun/qson"
)

var globalConfig *model.GlobalConfig = new(model.GlobalConfig)

func initRepo() repos.Repo {
	region := globalConfig.Ddb.Region
	if len(region) == 0 {
		region = "ap-southeast-1"
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	log.Printf("Initialize dynamodb repo with region %v, table %v\n", globalConfig.Ddb.Region, globalConfig.Ddb.TableName)
	awsConfig := &aws.Config{}
	ddb := dynamodb.New(sess, awsConfig)
	ddbRepo := new(repos.DDBRepo)
	ddbRepo.Ddb = ddb
	return ddbRepo
}

/*
&team_id=T0001
&team_domain=example
&enterprise_id=E0001
&enterprise_name=Globular%20Construct%20Inc
&channel_id=C2147483705
&channel_name=test
&user_id=U2147483697
&user_name=Steve
&command=/weather
&text=94070
&response_url=https://hooks.slack.com/commands/1234/5678
&trigger_id=13345224609.738474920.8088930838d88f008e0
&api_app_id=A123456
*/
type SlackCommandRequest struct {
	TeamId         *string `json:"team_id"`
	TeamDomain     *string `json:"team_domain"`
	EnterpriseId   *string `json:"enterprise_id"`
	EnterpriseName *string `json:"enterprise_name"`
	ChannelId      *string `json:"channel_id"`
	ChannelName    *string `json:"channel_name"`
	UserId         *string `json:"user_id"`
	Username       *string `json:"username"`
	Command        *string `json:"command"`
	Text           *string `json:"text"`
	ResponseUrl    *string `json:"response_url"`
	TriggerId      *string `json:"trigger_id"`
	ApiAppId       *string `json:"api_app_id"`
}

type SlackResponse struct {
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req := new(SlackCommandRequest)
	bodyForm := request.Body
	if request.IsBase64Encoded {
		content, _ := b64.StdEncoding.DecodeString(request.Body)
		// vals, _ := url.ParseQuery(string(bodyForm))
		bodyForm = string(content)
	}
	qson.Unmarshal(req, string(bodyForm))
	_, err := handleKudosCommand(req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: *req.Command + " " + *req.Text}, nil
}
func handleKudosCommand(slackCommand *SlackCommandRequest) (*SlackResponse, error) {
	// return MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil

}

func transformLambdaFunc() {

}

func main() {
	lambda.Start(handler)
}
