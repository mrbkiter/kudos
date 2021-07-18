package slackmodel

import "kudos-app.github.com/model"

type SlackCommandRequest struct {
	TeamId         string `json:"team_id"`
	TeamDomain     string `json:"team_domain"`
	EnterpriseId   string `json:"enterprise_id"`
	EnterpriseName string `json:"enterprise_name"`
	ChannelId      string `json:"channel_id"`
	ChannelName    string `json:"channel_name"`
	UserId         string `json:"user_id"`
	Command        string `json:"command"`
	Text           string `json:"text"`
	ResponseUrl    string `json:"response_url"`
	TriggerId      string `json:"trigger_id"`
	ApiAppId       string `json:"api_app_id"`
	Username       string `json:"user_name"`
}

type SlackResponse struct {
	Blocks       []*BlockSection   `json:"blocks,omitempty"`
	ResponseType SlackResponseType `json:"response_type,omitempty"`
	Text         string            `json:"text,omitempty"`
}

type BlockSection struct {
	Type *string      `json:"type"`
	Text *TextSection `json:"text,omitempty"`
}

type TextSection struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type SlackResponseType string

const (
	In_channel SlackResponseType = "in_channel"
	Ephemeral  SlackResponseType = "ephemeral"
)

type SlackSettingsCommand string

const (
	Add_Member  SlackSettingsCommand = "add-member"
	Del_Member  SlackSettingsCommand = "del-member"
	List_Member SlackSettingsCommand = "list-member"
	Add_Group   SlackSettingsCommand = "add-group"
	Del_Group   SlackSettingsCommand = "del-group"
	List_Group  SlackSettingsCommand = "list-group"
)

type KudosSettingsInput struct {
	TeamId      string
	ChannelId   string
	UserIds     []*model.UserNameIdMapping
	CommandType SlackSettingsCommand
	GroupId     string
}
