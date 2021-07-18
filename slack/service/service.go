package service

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"strings"
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
	filter.TeamId = slackCommand.TeamId
	userIdMappings := utils.ExtractUserIdsFromText(slackCommand.Text, "")
	userIds := make([]string, 0)
	for _, mapping := range userIdMappings {
		userIds = append(userIds, mapping.UserId)
	}
	filter.UserIds = userIds
	filter.ReportTime = utils.ExtractReportTime(slackCommand.Text)
	reportType := utils.ExtractReportType(slackCommand.Text)
	if reportType == model.Report_aggregate {
		if len(userIds) == 0 {
			temps := strings.Split(slackCommand.Text, " ")
			if !strings.Contains(temps[0], "_") { //this should be report type
				filter.GroupId = utils.ExtractGroupId(temps[0])
			}
		}
	}
	return filter, reportType
}

//ConvertToKudosSettingsFilter convert slack to kudos settings input
func ConvertToKudosSettingsInput(slackCommand *slackmodel.SlackCommandRequest) (*slackmodel.KudosSettingsInput, error) {
	ret := new(slackmodel.KudosSettingsInput)
	ret.TeamId = slackCommand.TeamId
	ret.ChannelId = slackCommand.ChannelId
	ret.UserIds = utils.ExtractUserIdsFromText(slackCommand.Text, "")

	remainingText := slackCommand.Text
	if strings.HasPrefix(slackCommand.Text, string(slackmodel.Add_Member)) {
		ret.CommandType = slackmodel.Add_Member
		remainingText = strings.Replace(remainingText, string(slackmodel.Add_Member), "", 1)
	} else if strings.HasPrefix(slackCommand.Text, string(slackmodel.Del_Member)) {
		ret.CommandType = slackmodel.Del_Member
		remainingText = strings.Replace(remainingText, string(slackmodel.Del_Member), "", 1)
	} else if strings.HasPrefix(slackCommand.Text, string(slackmodel.List_Member)) {
		ret.CommandType = slackmodel.List_Member
		remainingText = strings.Replace(remainingText, string(slackmodel.List_Member), "", 1)
	} else if strings.HasPrefix(slackCommand.Text, string(slackmodel.Add_Group)) {
		ret.CommandType = slackmodel.Add_Group
		remainingText = strings.Replace(remainingText, string(slackmodel.Add_Group), "", 1)
	} else if strings.HasPrefix(slackCommand.Text, string(slackmodel.Del_Group)) {
		ret.CommandType = slackmodel.Del_Group
		remainingText = strings.Replace(remainingText, string(slackmodel.Del_Group), "", 1)
	} else if strings.HasPrefix(slackCommand.Text, string(slackmodel.List_Group)) {
		ret.CommandType = slackmodel.List_Group
		remainingText = strings.Replace(remainingText, string(slackmodel.List_Group), "", 1)
	}

	ret.GroupId = utils.ExtractGroupId(strings.Trim(remainingText, " "))
	return ret, nil
}

func ConvertToKudosData(slackCommand *slackmodel.SlackCommandRequest) (*model.KudosData, error) {
	ret := new(model.KudosData)
	ret.ApiAppId = slackCommand.ApiAppId
	ret.ChannelId = slackCommand.ChannelId
	ret.ChannelName = slackCommand.ChannelName
	ret.MessageId = slackCommand.TriggerId
	ret.TeamId = slackCommand.TeamId
	ret.SourceUserId = slackCommand.UserId
	ret.ResponseUrl = slackCommand.ResponseUrl
	ret.Text = slackCommand.Command + " " + slackCommand.Text
	//analyze text
	kudosText := utils.AnalyzeKudosText(slackCommand.Text)
	if kudosText == "" {
		return nil, errors.New("Unable to detect kudos text. Please check allowed kudos syntax again")
	}
	ret.TargetUserIds = utils.ExtractUserIdsFromText(kudosText, slackCommand.UserId)
	ret.AppId = model.Slack
	ret.Username = slackCommand.Username
	return ret, nil
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

func BuildQuickSlackResponse(text string) *slackmodel.SlackResponse {
	slackResp := new(slackmodel.SlackResponse)
	slackResp.ResponseType = slackmodel.In_channel
	slackResp.Text = text
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
	slackResp.Text = fmt.Sprintf("+1 for %v", targetUserTags)
	return slackResp
}

func BuildUserTag(userId string) string {
	return fmt.Sprintf("<@%s>", userId)
}
