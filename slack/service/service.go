package service

import (
	b64 "encoding/base64"
	"fmt"
	"time"

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
func ConvertToKudosReportFilter(slackCommand *slackmodel.SlackCommandRequest) (*model.KudosReportFilter, model.ReportType) {
	filter := new(model.KudosReportFilter)
	filter.TeamId = *slackCommand.TeamId
	userIdMappings := utils.ExtractUserIdsFromText(*slackCommand.Text)
	userIds := make([]string, 0)
	for _, mapping := range userIdMappings {
		userIds = append(userIds, mapping.UserId)
	}
	filter.UserIds = userIds
	filter.ReportTime = utils.ExtractReportTime(*slackCommand.Text)
	reportType := utils.ExtractReportType(*slackCommand.Text)
	return filter, reportType
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
	ret.Username = *slackCommand.Username
	return ret
}

func BuildSlackResponseForReportDetail(reportDetail *model.KudosReportDetails) *slackmodel.SlackResponse {
	slackResp := new(slackmodel.SlackResponse)
	slackResp.ResponseType = slackmodel.In_channel
	if len(reportDetail.KudosTalk) == 0 {
		slackResp.Text = "No report found"
		return slackResp
	}
	slackResp.Blocks = make([]*slackmodel.BlockSection, 0)
	for _, tgId := range reportDetail.KudosTalk {
		targetUserTag := BuildUserTag(tgId.UserId)
		block := new(slackmodel.BlockSection)
		block.Type = utils.String("section")
		block.Text = new(slackmodel.TextSection)
		block.Text.Type = "mrkdwn"
		block.Text.Text = fmt.Sprintf("%v: %v (%v)", targetUserTag, tgId.Text, tgId.Timestamp.Format(time.RFC822))
		slackResp.Blocks = append(slackResp.Blocks, block)
	}

	return slackResp
}
func BuildSlackResponseForReport(kudosReportRet *model.KudosReportResult) *slackmodel.SlackResponse {
	slackResp := new(slackmodel.SlackResponse)
	slackResp.ResponseType = slackmodel.In_channel
	if len(kudosReportRet.UserReport) == 0 {
		slackResp.Text = "No report found"
		return slackResp
	}

	slackResp.Blocks = make([]*slackmodel.BlockSection, 0)
	for _, tgId := range kudosReportRet.UserReport {
		targetUserTag := BuildUserTag(tgId.UserId)
		block := new(slackmodel.BlockSection)
		block.Type = utils.String("section")
		block.Text = new(slackmodel.TextSection)
		block.Text.Type = "mrkdwn"
		block.Text.Text = fmt.Sprintf("%v: %v kudos", targetUserTag, tgId.Total)
		slackResp.Blocks = append(slackResp.Blocks, block)
	}

	return slackResp
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
		targetUserTags = targetUserTags + BuildUserTag(tgId.UserId)
	}
	slackResp.Text = fmt.Sprintf("+%v for %v", len(kudosData.TargetUserIds), targetUserTags)
	return slackResp
}

func BuildUserTag(userId string) string {
	return fmt.Sprintf("<@%s>", userId)
}
