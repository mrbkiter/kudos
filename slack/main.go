package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	} else if req.Command == "/kudos-settings" {
		resp, err = handleKudosSettings(myCtx, req)
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

func handleKudosSettings(ctx *model.MyContext, slackCommand *slackmodel.SlackCommandRequest) (*slackmodel.SlackResponse, error) {
	kudosSettingsInput, err := service.ConvertToKudosSettingsInput(slackCommand)
	if err != nil {
		return service.BuildQuickSlackResponse(err.Error()), nil
	}
	ctx.Log.Infof("Kudos Settings", "settings", kudosSettingsInput)

	switch kudosSettingsInput.CommandType {
	case slackmodel.Add_Member:
		memberReq := convertToKudosGroupSettingsMembers(kudosSettingsInput)
		ctx.Log.Infow("Add_Members", "teamId", memberReq.TeamId, "groupId", memberReq.GroupId, "userId", slackCommand.UserId)
		if len(memberReq.GroupId) == 0 {
			return service.BuildQuickSlackResponse(GroupIdPatternMessage), nil
		}
		err = repo.AddMembersToTeamGroup(ctx, memberReq)
		return service.BuildQuickSlackResponse("Members have been added"), nil
	case slackmodel.Del_Member:
		memberReq := convertToKudosGroupSettingsMembers(kudosSettingsInput)
		ctx.Log.Infow("Delete_Members", "teamId", memberReq.TeamId, "groupId", memberReq.GroupId, "userId", slackCommand.UserId)
		if len(memberReq.GroupId) == 0 {
			return service.BuildQuickSlackResponse(GroupIdPatternMessage), nil
		}
		err = repo.DeleteMembersFromTeamGroup(ctx, memberReq)
		return service.BuildQuickSlackResponse("Members have been deleted"), nil
	case slackmodel.List_Member:
		req := convertToKudosTeamSettingsGroupAction(kudosSettingsInput)
		if len(req.GroupId) == 0 {
			return service.BuildQuickSlackResponse(GroupIdPatternMessage), nil
		}
		members, err := repo.ListMembersOfTeamGroup(ctx, req)
		if err != nil {
			return service.BuildQuickSlackResponse(err.Error()), nil
		}
		if len(members) == 0 {
			return service.BuildQuickSlackResponse("No member found"), nil
		}
		userTags := make([]string, 0)
		for _, userId := range members {
			userTags = append(userTags, service.BuildUserTag(userId))
		}
		return service.BuildQuickSlackResponse(strings.Join(userTags, ",")), nil
	case slackmodel.Add_Group:
		req := convertToKudosTeamSettingsGroupAction(kudosSettingsInput)
		ctx.Log.Infow("Add_Group", "teamId", req.TeamId, "groupId", req.GroupId, "userId", slackCommand.UserId)
		if len(req.GroupId) == 0 {
			return service.BuildQuickSlackResponse(GroupIdPatternMessage), nil
		}
		err = repo.AddTeamGroup(ctx, req)
		return service.BuildQuickSlackResponse("Group has been created"), nil
	case slackmodel.Del_Group:
		req := convertToKudosTeamSettingsGroupAction(kudosSettingsInput)
		ctx.Log.Infow("Delete_Group", "teamId", req.TeamId, "groupId", req.GroupId, "userId", slackCommand.UserId)
		err = repo.DeleteTeamGroup(ctx, req)
		return service.BuildQuickSlackResponse("Group has been deleted"), nil
	case slackmodel.List_Group:
		req := convertToKudosTeamSettingsListGroups(kudosSettingsInput)
		ctx.Log.Infow("List group", "req", req)
		groupIds, err := repo.ListAllTeamGroups(ctx, req)
		if err != nil {
			return service.BuildQuickSlackResponse(err.Error()), nil
		}
		if len(groupIds) == 0 {
			return service.BuildQuickSlackResponse("No group found"), nil
		}
		return service.BuildQuickSlackResponse(strings.Join(groupIds, ",")), nil
	}

	if err != nil {
		return service.BuildQuickSlackResponse(err.Error()), nil
	}
	return service.BuildQuickSlackResponse("Command not support. Use add-group, del-group, list-group, add-member, del-member, list-member"), nil
}

func convertToKudosTeamSettingsListGroups(kudosSettingsInput *slackmodel.KudosSettingsInput) *model.KudosTeamSettingsListGroups {
	req := new(model.KudosTeamSettingsListGroups)
	req.TeamId = kudosSettingsInput.TeamId
	return req
}
func convertToKudosTeamSettingsGroupAction(kudosSettingsInput *slackmodel.KudosSettingsInput) *model.KudosTeamSettingsGroupAction {
	req := new(model.KudosTeamSettingsGroupAction)
	req.GroupId = kudosSettingsInput.GroupId
	req.TeamId = kudosSettingsInput.TeamId
	return req
}

func convertToKudosGroupSettingsMembers(kudosSettingsInput *slackmodel.KudosSettingsInput) *model.KudosGroupSettingsMembers {
	memberReq := new(model.KudosGroupSettingsMembers)
	memberReq.TeamId = kudosSettingsInput.TeamId
	memberReq.GroupId = kudosSettingsInput.GroupId
	memberReq.TargetUserIds = kudosSettingsInput.UserIds
	return memberReq
}

func main() {
	log.Println("Started Slack Kudos App")
	initEnvironment()
	repo = initRepo()
	lambda.Start(handler)
}

const GroupIdPatternMessage = "Group Id has to follow pattern: [a-zA-Z0-9-_]+ and not reserved keywords"
