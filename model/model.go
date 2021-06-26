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
	AppId         AppId
	Username      string
}

type KudosCountUpdate struct {
	TeamId    string
	UserId    string
	Timestamp int64
	Counter   int
}

type AppId string

const (
	Slack AppId = "SLACK"
)

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
