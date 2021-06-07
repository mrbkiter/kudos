package model

import (
	"context"

	"go.uber.org/zap"
)

type GlobalConfig struct {
	Ddb DynamoConfig `json:"ddbConfig"`
}

type DynamoConfig struct {
	Region    string `json:"region"`
	TableName string `json:"tableName"`
}

type SlackCommandRequest struct {
	TeamId         *string `json:"team_id"`
	TeamDomain     *string `json:"team_domain"`
	EnterpriseId   *string `json:"enterprise_id"`
	EnterpriseName *string `json:"enterprise_name"`
	ChannelId      *string `json:"channel_id"`
	ChannelName    *string `json:"channel_name"`
	UserId         *string `json:"user_id"`
	Username       *string `json:"username"`
	Command        *string `json:"command"`
	Text           *string `json:"text"`
	ResponseUrl    *string `json:"response_url"`
	TriggerId      *string `json:"trigger_id"`
	ApiAppId       *string `json:"api_app_id"`
}

type SlackResponse struct {
}

type KudosData struct {
	TeamId        string
	ChannelId     string
	ChannelName   *string
	SourceUserId  string
	Text          string
	ApiAppId      *string
	TargetUserIds []string
	MessageId     string
	ResponseUrl   string
}

type MyContext struct {
	AwsRequestId           string
	Log                    *zap.SugaredLogger
	Ctx                    context.Context
	RequestTimeOutInMinute int64
	GlobalConfig           *GlobalConfig
	UserInternalId         string
	Username               string
	Testing                bool
}
