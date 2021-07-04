package repos_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"kudos-app.github.com/model"
	"kudos-app.github.com/repos"
	"kudos-app.github.com/utils"
)

func Test_mapKudosDataToKudosCommandEtt(t *testing.T) {
	data := new(model.KudosData)
	tgUserIds := make([]*model.UserNameIdMapping, 2)
	tgUserIds[0] = new(model.UserNameIdMapping)
	tgUserIds[1] = new(model.UserNameIdMapping)
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
	if ett[0].Id1 != "teamId#1" || ett[1].Id1 != "teamId#2" {
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
