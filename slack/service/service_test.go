package service_test

import (
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"kudos-app.github.com/model"
	"kudos-app.github.com/slack/service"
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
	if *command.Text == "" {
		t.Errorf("convert slack message error")
	}

}

func Test_BuildSlackResponse(t *testing.T) {
	kudos := new(model.KudosData)
	targetUserIds := make([]string, 0)
	targetUserIds = append(targetUserIds, "tgUserId1", "tgUserId2")
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
	targetUserIds := make([]string, 0)
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
