package ddb_entity

type KudosTeamSettings struct {
	Id1            string `json:"id1"`
	Id2            string `json:"id2"`
	TeamId         string `json:"teamId"`
	AdminChannelId string `json:"adminChannelId"`
	Timestamp      int64  `json:"timestamp"`
}

type KudosGroupSettings struct {
	Id1       string    `dynamo:"id1"`
	Id2       string    `dynamo:"id2"`
	TeamId    string    `dynamo:"teamId"`
	GroupId   string    `dynamo:"groupId"`
	Timestamp int64     `dynamo:"timestamp"`
	Type      KudosType `dynamo:"type"`
	UserIds   []string  `dynamo:"userIds,set"`
}

type KudosCommand struct {
	Id1         string `json:"id1"`
	Id2         string `json:"id2"`
	TeamIdMonth string `json:"teamIdMonth"`
	TeamIdWeek  string `json:"teamIdWeek"`
	TeamId      string `json:"teamId"`
	ChannelId   string `json:"channelId"`
	// Command         string `json:"command"`
	Text           string    `json:"text"`
	Timestamp      int64     `json:"timestamp"`
	UserId         string    `json:"userId"`
	Username       string    `json:"username"`
	MsgId          string    `json:"msgId"`
	SourceUserId   string    `json:"sourceUserId"`
	SourceUserName string    `json:"sourceUsername"`
	Ttl            int64     `json:"ttl"`
	Type           KudosType `json:"type"`
}

type KudosCounter struct {
	Id1       string    `json:"id1"` //teamId#report
	Id2       string    `json:"id2"` //<weekNumber | monthNumber>#userId
	Timestamp int64     `json:"timestamp"`
	Count     int       `json:"count"`
	UserId    string    `json:"userId"`
	Type      KudosType `json:"type"`
	Username  string    `json:"username"`
}

type KudosType string

const (
	CommandType   KudosType = "command"
	ReportType    KudosType = "report"
	GroupSettings KudosType = "group_settings"
)
