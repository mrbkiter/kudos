package repos_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	// "github.com/aws/aws-sdk-go/service/dynamodb"
	"kudos-app.github.com/model"
	"kudos-app.github.com/repos"
	"kudos-app.github.com/utils"
)

func Test_mapKudosDataToKudosCommandEtt(t *testing.T) {
	data := new(model.KudosData)
	tgUserIds := make([]*model.UserNameIdMapping, 2)
	tgUserIds[0] = new(model.UserNameIdMapping)
	tgUserIds[1] = new(model.UserNameIdMapping)
	tgUserIds[0].UserId = "UID1"
	tgUserIds[0].Username = "UN1"
	tgUserIds[1].UserId = "UID2"
	tgUserIds[1].Username = "UN2"

	data.TargetUserIds = tgUserIds
	data.TeamId = "teamId"
	data.MessageId = "MessageId"
	ett := repos.MapKudosDataToKudosCommandEtt(data)
	monthFormat := time.Now().Format("2006-01")
	yearNumber, weekNumber := time.Now().ISOWeek()

	if len(ett) != 2 {
		t.Error("MapKudosDataToKudosCommandEtt should return 2 elements")
	}
	if ett[0].TeamId != "teamId" || ett[1].TeamId != "teamId" {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong teamId")
	}
	if ett[0].MsgId != "MessageId" || ett[1].MsgId != "MessageId" {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong MessageId")
	}
	if ett[0].Id1 != "teamId#UID1" || ett[1].Id1 != "teamId#UID2" {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong Id1")
	}
	if ett[0].TeamIdMonth != ("teamId#"+monthFormat) || ett[1].TeamIdMonth != ("teamId#"+monthFormat) {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong TeamIdMonth")
	}
	if ett[0].TeamIdWeek != fmt.Sprintf("teamId#%v#%v", yearNumber, weekNumber) || ett[0].TeamIdWeek != fmt.Sprintf("teamId#%v#%v", yearNumber, weekNumber) {
		t.Error("MapKudosDataToKudosCommandEtt mapped TeamIdWeek")
	}
}

func Test_BuildIdMonthAndWeek(t *testing.T) {
	monthId2, weekNumber := repos.BuildIdMonthAndWeek(1625156312, "U1234")
	fmt.Println(monthId2)
	fmt.Println(weekNumber)
	if monthId2 != "2021-07#U1234" || weekNumber != "2021#26#U1234" {
		t.Errorf("Converting IdMonth2 & Id2WeekNUmber are incorrect")
	}
}

func Test_IncreaseKudosCounter(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
	awsConfig := &aws.Config{}
	ddb := dynamodb.New(sess, awsConfig)
	ddbRepo := new(repos.DDBRepo)
	ddbRepo.Ddb = ddb

	myCtx := utils.NewMyContext(context.Background(), "")
	globalConfig := new(model.GlobalConfig)
	globalConfig.Ddb = model.DynamoConfig{
		Region:    "ap-southeast-1",
		TableName: "dev-kudos",
	}
	myCtx.GlobalConfig = globalConfig
	data := new(model.KudosCountUpdate)
	data.UserId = "U123456"
	data.Counter = 1
	data.TeamId = "T98765"
	data.Timestamp = 1625156312
	data.Username = "UN123456"
	err := ddbRepo.IncreaseKudosCounter(myCtx, data)
	if err != nil {
		t.Errorf(fmt.Sprintf("Error %v when increase kudos counter", err.Error()))
	}
}

func Test_GetUserReport(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
	awsConfig := &aws.Config{}
	ddb := dynamodb.New(sess, awsConfig)
	ddbRepo := new(repos.DDBRepo)
	ddbRepo.Ddb = ddb

	myCtx := utils.NewMyContext(context.Background(), "")
	globalConfig := new(model.GlobalConfig)
	globalConfig.Ddb = model.DynamoConfig{
		Region:    "ap-southeast-1",
		TableName: "dev-kudos",
	}
	myCtx.GlobalConfig = globalConfig
	data := new(model.KudosReportFilter)
	data.TeamId = "T98765"
	data.ReportTime = model.THIS_MONTH
	userIds := make([]string, 2)
	userIds[0] = "U123456"
	userIds[1] = "U12345"
	data.UserIds = userIds
	ret, err := ddbRepo.GetKudosReport(myCtx, data)
	if err != nil {
		t.Errorf(fmt.Sprintf("Error %v when getting kudos report", err.Error()))
	} else if len(ret.UserReport) != 1 {
		t.Errorf("Report length should be 1")
	}
}

func Test_GetUserReportDetail(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
	awsConfig := &aws.Config{}
	ddb := dynamodb.New(sess, awsConfig)
	ddbRepo := new(repos.DDBRepo)
	ddbRepo.Ddb = ddb

	myCtx := utils.NewMyContext(context.Background(), "")
	globalConfig := new(model.GlobalConfig)
	globalConfig.Ddb = model.DynamoConfig{
		Region:    "ap-southeast-1",
		TableName: "dev-kudos",
	}
	myCtx.GlobalConfig = globalConfig
	data := new(model.KudosReportFilter)
	data.TeamId = "T022PA5N7KP"
	data.ReportTime = model.THIS_MONTH
	userIds := make([]string, 1)
	userIds[0] = "U024D6VQX7Z"
	data.UserIds = userIds
	_, err := ddbRepo.GetKudosReportDetails(myCtx, data)
	if err != nil {
		t.Errorf(fmt.Sprintf("Error %v when getting kudos report detail", err.Error()))
	}
}

func Test_AddMemberToTeamGroup(t *testing.T) {
	ctx, ddbRepo := initConfig()
	req := &model.KudosTeamSettingsGroupAction{
		TeamId:  "T01",
		GroupId: fmt.Sprintf("Group-%v", time.Now().Unix()),
	}
	err := ddbRepo.AddTeamGroup(ctx, req)
	if err != nil {
		t.Error("Add Team Group should return success status")
	}
	memReq := &model.KudosGroupSettingsMembers{
		TeamId:  req.TeamId,
		GroupId: req.GroupId,
		TargetUserIds: []*model.UserNameIdMapping{
			{Username: "UN1", UserId: "U1"},
			{Username: "UN2", UserId: "U2"},
		},
	}
	err = ddbRepo.AddMembersToTeamGroup(ctx, memReq)
	if err != nil {
		t.Error("Add members to group should return ok")
	}

	memReq2 := &model.KudosGroupSettingsMembers{
		TeamId:  req.TeamId,
		GroupId: req.GroupId,
		TargetUserIds: []*model.UserNameIdMapping{
			{Username: "UN4", UserId: "U4"},
			{Username: "UN3", UserId: "U3"},
			{Username: "UN2", UserId: "U2"},
		},
	}
	err = ddbRepo.AddMembersToTeamGroup(ctx, memReq2)
	if err != nil {
		t.Error("Add members second time to group should return ok")
	}
	err = ddbRepo.DeleteMembersFromTeamGroup(ctx, memReq)
	if err != nil {
		t.Error("Delete members from group should return ok")
	}

	//test list members in group
	members, err := ddbRepo.ListMembersOfTeamGroup(ctx, req)
	if err != nil {
		t.Error("List members should return items")
	}

	if len(members) != 2 {
		t.Error("List members not return correct number of items")
	}

	ddbRepo.DeleteTeamGroup(ctx, req)
}
func Test_AddDeleteTeamGroup(t *testing.T) {
	ctx, ddbRepo := initConfig()
	req := &model.KudosTeamSettingsGroupAction{
		TeamId:  "T01",
		GroupId: fmt.Sprintf("Group-%v", time.Now().Unix()),
	}
	err := ddbRepo.AddTeamGroup(ctx, req)
	if err != nil {
		t.Error("Add Team Group should return success status")
	}
	listGroupReq := new(model.KudosTeamSettingsListGroups)
	listGroupReq.TeamId = "T01"
	groups, err := ddbRepo.ListAllTeamGroups(ctx, listGroupReq)
	if err != nil {
		t.Error("List group should return value", err)
	}
	if len(groups) < 1 {
		t.Error("List group should return at list one newly added group")
	}
	err = ddbRepo.DeleteTeamGroup(ctx, req)
	if err != nil {
		t.Error("DeleteTeamGroup should return success")
	}
}

func Test_DeleteTeamGroupDup(t *testing.T) {
	ctx, ddbRepo := initConfig()
	req := &model.KudosTeamSettingsGroupAction{
		TeamId:  "T01",
		GroupId: "Group1",
	}
	err := ddbRepo.AddTeamGroup(ctx, req)
	if err == nil {
		t.Error("Add Team Group should return dup error")
	}
}

func Test_ListTemGroup(t *testing.T) {

}
func initConfig() (*model.MyContext, *repos.DDBRepo) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
	awsConfig := &aws.Config{}
	ddb := dynamodb.New(sess, awsConfig)
	ddbRepo := new(repos.DDBRepo)
	ddbRepo.Ddb = ddb

	myCtx := utils.NewMyContext(context.Background(), "")
	globalConfig := new(model.GlobalConfig)
	globalConfig.Ddb = model.DynamoConfig{
		Region:    "ap-southeast-1",
		TableName: "dev-kudos",
	}
	myCtx.GlobalConfig = globalConfig
	return myCtx, ddbRepo
}
