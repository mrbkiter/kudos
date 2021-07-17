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
