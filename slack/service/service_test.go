package service_test

import (
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"kudos-app.github.com/model"
	"kudos-app.github.com/slack/service"
	"kudos-app.github.com/slack/slackmodel"
	"kudos-app.github.com/utils"
)

func Test_convertAwsRequestToSlackCommandRequest(t *testing.T) {
	payload, err := utils.ReadFile("../test-resource/payload1.json")
	if err != nil {
		t.Errorf(err.Error())
	}
	event := new(events.APIGatewayProxyRequest)
	utils.ConvertJSONReaderToObject(strings.NewReader(payload), event)

	command := service.ConvertAwsRequestToSlackCommandRequest(*event)
	if command.Text == "" {
		t.Errorf("convert slack message error")
	}

}

func Test_BuildSlackResponse(t *testing.T) {
	kudos := new(model.KudosData)
	targetUserIds := make([]*model.UserNameIdMapping, 0)
	targetUserIds = append(targetUserIds, &model.UserNameIdMapping{UserId: "tgUserId1"}, &model.UserNameIdMapping{UserId: "tgUserId2"})
	kudos.TargetUserIds = targetUserIds
	kudos.SourceUserId = "sourceUserId"
	kudos.Text = "This is new comment from abc"

	resp := service.BuildSlackResponse(kudos)
	payloadString := utils.ConvertObjectToJSON(resp)
	if !strings.Contains(payloadString, "in_channel") || !strings.Contains(payloadString, "+2") {
		t.Errorf("Converting to slack response incorrectly")
	}
}

func Test_BuildSlackResponse0(t *testing.T) {
	kudos := new(model.KudosData)
	targetUserIds := make([]*model.UserNameIdMapping, 0)
	kudos.TargetUserIds = targetUserIds
	kudos.SourceUserId = "sourceUserId"
	kudos.Text = "This is new comment from abc"

	resp := service.BuildSlackResponse(kudos)
	payloadString := utils.ConvertObjectToJSON(resp)
	if !strings.Contains(payloadString, "in_channel") || !strings.Contains(payloadString, "+0") {
		t.Errorf("Converting to slack response incorrectly")
	}
}
func Test_BuildUserTag(t *testing.T) {
	userTag := service.BuildUserTag("1234")

	if userTag != "<@1234>" {
		t.Error("BuildUserTag failed to build user tag")
	}

}

func Test_ConvertToKudosData(t *testing.T) {
	req := &slackmodel.SlackCommandRequest{
		UserId: "US1",
		Text:   "thanks <US2|mrbkiter> <US1|mrb> for helping us on tasks",
	}
	ret, err := service.ConvertToKudosData(req)
	if err != nil {
		t.Error("ConvertToKudosData has issue when converting to kudos data")
	}
	if len(ret.TargetUserIds) != 1 {
		t.Errorf("TargetUserIds not returned correctly")
	}
	if ret.TargetUserIds[0].UserId != "US2" {
		t.Error("TargetUserIds[0] not return correct id")
	}
}

func Test_ConvertToKudosDataError(t *testing.T) {
	req := &slackmodel.SlackCommandRequest{
		UserId: "US1",
		Text:   "<@U024D6VQX7Z|vu.yen.nguyen.88> <@U024U032H8A|vu.nguyen>",
	}
	_, err := service.ConvertToKudosData(req)
	if err == nil {
		t.Error("ConvertToKudosData should return error")
	}
}

func Test_ConvertToKudosReportFilter(t *testing.T) {
	req := &slackmodel.SlackCommandRequest{
		Text: "group-id LAST_MONTH",
	}
	f, reportType := service.ConvertToKudosReportFilter(req)
	if f.GroupId != "group-id" || reportType != model.Report_aggregate {
		t.Error("ConvertToKudosReportFilter extracted wrong group-id")
	}
	req.Text = "group-id"
	f, _ = service.ConvertToKudosReportFilter(req)
	if f.GroupId != "group-id" {
		t.Error("ConvertToKudosReportFilter extracted wrong group-id")
	}

	req.Text = "LAST_MONTH"
	f, reportType = service.ConvertToKudosReportFilter(req)
	if len(f.GroupId) > 0 || f.ReportTime != model.LAST_MONTH {
		t.Error("ConvertToKudosReportFilter extracted wrong group-id")
	}

}

func Test_ConvertToKudosSettingsInput(t *testing.T) {
	req := &slackmodel.SlackCommandRequest{
		Text: "add-member group-id <@UID1>",
	}
	input, err := service.ConvertToKudosSettingsInput(req)
	if err != nil {
		t.Error("ConvertToKudosSettingsInput should not return error")
	}
	if input.CommandType != slackmodel.Add_Member || input.GroupId != "group-id" || len(input.UserIds) != 1 {
		t.Error("ConvertToKudosSettingsInput wrong convert to add-member")
	}

	req.Text = "list-member group-id"
	input, err = service.ConvertToKudosSettingsInput(req)
	if err != nil {
		t.Error("ConvertToKudosSettingsInput should not return error")
	}
	if input.CommandType != slackmodel.List_Member {
		t.Error("ConvertToKudosSettingsInput wrong convert to list-member")
	}

	req.Text = "list-group"
	input, err = service.ConvertToKudosSettingsInput(req)
	if err != nil {
		t.Error("ConvertToKudosSettingsInput should not return error")
	}
	if input.CommandType != slackmodel.List_Group {
		t.Error("ConvertToKudosSettingsInput wrong convert to list-member")
	}

}

//(slackCommand *slackmodel.SlackCommandRequest) (*model.KudosData, error)
//{"blocks":[{"string":"section","text":{"type":"mrkdwn","response_type":"ephemeral","text":"+2"}}]}
/*
{
	"blocks":[
	{
		"type":"section",
		"text": {
			"type":"mrkdwn",
			"response_type":"ephemeral",
			"text":"+2"
		}
	}]
}


{
    "blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*It's 80 degrees right now.*"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Partly cloudy today and tomorrow"
			}
		}
	]
}
*/
