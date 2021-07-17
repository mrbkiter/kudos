package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"kudos-app.github.com/model"
	"kudos-app.github.com/repos"
	"kudos-app.github.com/slack/service"
	"kudos-app.github.com/slack/slackmodel"
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	myCtx := utils.NewMyContext(context.Background(), "")
	myCtx.GlobalConfig = globalConfig
	myCtx.Log.Infow("Payload", "request", request)
	req := service.ConvertAwsRequestToSlackCommandRequest(request)
	myCtx.Log.Infow("Converted Request", "request", req)
	resp := new(slackmodel.SlackResponse)
	var err error
	if req.Command == "/kudos" {
		resp, err = handleKudosCommand(myCtx, req)
	} else if req.Command == "/kudos-report" {
		resp, err = handleKudosReport(myCtx, req)
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Command not found"}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: utils.ConvertObjectToJSON(resp)}, nil
}

func handleKudosReport(ctx *model.MyContext, slackCommand *slackmodel.SlackCommandRequest) (*slackmodel.SlackResponse, error) {
	kudosFilter, reportType := service.ConvertToKudosReportFilter(slackCommand)
	reportRet := new(slackmodel.SlackResponse)
	if reportType == model.Report_detail {
		ret, err := repo.GetKudosReportDetails(ctx, kudosFilter)
		if err != nil {
			return nil, err
		}
		reportRet = service.BuildSlackResponseForReportDetail(ret)
	} else {
		ret, err := repo.GetKudosReport(ctx, kudosFilter)
		if err != nil {
			return nil, err
		}
		reportRet = service.BuildSlackResponseForReport(ret)
	}

	return reportRet, nil
}

func handleKudosCommand(ctx *model.MyContext, slackCommand *slackmodel.SlackCommandRequest) (*slackmodel.SlackResponse, error) {
	// return MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
	kudosData, err := service.ConvertToKudosData(slackCommand)
	if err != nil {
		return service.BuildQuickSlackResponse(err.Error()), nil
	}
	ctx.Log.Infof("Kudos Data", "kudosData", kudosData)
	err = repo.WriteKudosCommand(ctx, kudosData)
	if err != nil {
		return nil, err
	}

	return service.BuildSlackResponse(kudosData), nil
}

func main() {
	log.Println("Started Slack Kudos App")
	initEnvironment()
	repo = initRepo()
	lambda.Start(handler)
}
