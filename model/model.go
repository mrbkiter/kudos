package model

import (
	"context"
	"time"

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
	ChannelName   string
	SourceUserId  string
	Text          string
	ApiAppId      string
	TargetUserIds []*UserNameIdMapping
	MessageId     string
	ResponseUrl   string
	AppId         AppId
	Username      string
}

//KudosGroupSettingsMembers add | delete members
type KudosGroupSettingsMembers struct {
	TeamId        string
	ChannelId     string
	SourceUserId  string
	GroupId       string
	TargetUserIds []*UserNameIdMapping
}

//KudosTeamSettingsListGroups list all groups of specific team
type KudosTeamSettingsListGroups struct {
	TeamId string
}

//KudosTeamSettingsGroupAction create | delete | get group
type KudosTeamSettingsGroupAction struct {
	TeamId  string
	GroupId string
}

type KudosCountUpdate struct {
	TeamId    string
	UserId    string
	Timestamp int64
	Counter   int
	Username  string
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

type UserNameIdMapping struct {
	Username string
	UserId   string
}

type KudosReportFilter struct {
	TeamId     string     `json:"teamId"`
	UserIds    []string   `json:"userIds"`
	ReportTime ReportTime `json:"reportTime"`
	GroupId    string     `json:"groupId"`
}

type KudosReportResult struct {
	FromTime   time.Time                `json:"fromTime"`
	ToTime     time.Time                `json:"toTime"`
	UserReport []*KudosUserReportResult `json:"userReport"`
}
type KudosUserReportResult struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Total    int    `json:"total"`
}

type KudosReportDetails struct {
	UserId    string                `json:"userId"`
	Username  string                `json:"username"`
	KudosTalk []*KudosSimpleCommand `json:"talks"`
}

type KudosSimpleCommand struct {
	Text      string
	UserId    string
	Timestamp time.Time
}
type ReportTime string

const (
	THIS_MONTH ReportTime = "THIS_MONTH"
	LAST_MONTH ReportTime = "LAST_MONTH"
	THIS_WEEK  ReportTime = "THIS_WEEK"
	LAST_WEEK  ReportTime = "LAST_WEEK"
)

type ReportType string

const (
	Report_detail    ReportType = "detail"
	Report_aggregate ReportType = "aggregate"
)
