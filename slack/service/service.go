package service

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joncalhoun/qson"
	"kudos-app.github.com/model"
	"kudos-app.github.com/slack/slackmodel"
	"kudos-app.github.com/utils"
)

func ConvertAwsRequestToSlackCommandRequest(request events.APIGatewayProxyRequest) *slackmodel.SlackCommandRequest {
	bodyForm := request.Body
	req := new(slackmodel.SlackCommandRequest)
	if request.IsBase64Encoded {
		content, _ := b64.StdEncoding.DecodeString(request.Body)
		// vals, _ := url.ParseQuery(string(bodyForm))
		bodyForm = string(content)
	} else {
		bodyForm = string(request.Body)
	}
	qson.Unmarshal(req, string(bodyForm))
	return req
}

func ConvertToKudosData(slackCommand *slackmodel.SlackCommandRequest) *model.KudosData {
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
	ret.AppId = model.Slack
	return ret
}

func BuildSlackResponse(kudosData *model.KudosData) *slackmodel.SlackResponse {

	slackResp := new(slackmodel.SlackResponse)
	slackResp.ResponseType = slackmodel.In_channel
	if len(kudosData.TargetUserIds) == 0 {
		slackResp.Text = "+0 kudos"
		return slackResp
	}

	targetUserTags := ""
	for _, tgId := range kudosData.TargetUserIds {
		targetUserTags = targetUserTags + BuildUserTag(tgId)
	}
	slackResp.Text = fmt.Sprintf("+%v for %v", len(kudosData.TargetUserIds), targetUserTags)
	return slackResp
}

func BuildUserTag(userId string) string {
	return fmt.Sprintf("<@%s>", userId)
}
