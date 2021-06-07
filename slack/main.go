package slack

import (
	"context"
	"fmt"
	"log"
	"os"

	b64 "encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	qson "github.com/joncalhoun/qson"
	"kudos-app.github.com/model"
	"kudos-app.github.com/repos"
	"kudos-app.github.com/utils"
)

var globalConfig *model.GlobalConfig = new(model.GlobalConfig)
var repo repos.Repo

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

func initEnvironment() {
	profile := os.Getenv("PROFILE")
	log.Println("PROFILE ", profile)
	//read config file
	globalConfigPath := fmt.Sprintf("./config/config-%v.json", profile)
	err := utils.ReadFromJSONFile(globalConfigPath, globalConfig)
	if err != nil {
		panic(err)
	}
	log.Printf("Global Config loaded %v\n", *globalConfig)
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	myCtx := utils.NewMyContext(context.Background(), "")
	myCtx.GlobalConfig = globalConfig
	req := new(model.SlackCommandRequest)
	bodyForm := request.Body
	if request.IsBase64Encoded {
		content, _ := b64.StdEncoding.DecodeString(request.Body)
		// vals, _ := url.ParseQuery(string(bodyForm))
		bodyForm = string(content)
	}
	qson.Unmarshal(req, string(bodyForm))
	_, err := handleKudosCommand(myCtx, req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: *req.Command + " " + *req.Text}, nil
}
func handleKudosCommand(ctx *model.MyContext, slackCommand *model.SlackCommandRequest) (*model.SlackResponse, error) {
	// return MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
	kudosData := convertToKudosData(slackCommand)
	err := repo.WriteKudosCommand(ctx, kudosData)
	if err != nil {
		return nil, err
	}
	return &model.SlackResponse{}, nil
}

func convertToKudosData(slackCommand *model.SlackCommandRequest) *model.KudosData {
	ret := new(model.KudosData)
	ret.ApiAppId = slackCommand.ApiAppId
	ret.ChannelId = *slackCommand.ChannelId
	ret.ChannelName = slackCommand.ChannelName
	ret.MessageId = *slackCommand.TriggerId
	ret.TeamId = *slackCommand.TeamId
	ret.SourceUserId = *slackCommand.UserId
	ret.ResponseUrl = *slackCommand.ResponseUrl
	ret.Text = *slackCommand.Command + " " + *slackCommand.Text
	ret.TargetUserIds = utils.ExtractUserIdsFromText(*slackCommand.Text)
	return ret
}

func main() {
	log.Println("Started Slack Kudos App")
	initEnvironment()
	repo = initRepo()
	//init engine
	// initBotEngine(retChan)
	// ctx, cancelFunc := context.WithCancel(context.Background())
	lambda.Start(handler)
}
