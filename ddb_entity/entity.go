package ddb_entity

type KudosCommand struct {
	Id1         string `json:"id1"`
	Id2         string `json:"id2"`
	TeamIdMonth string `json:"teamIdMonth"`
	TeamIdWeek  string `json:"teamIdWeek"`
	TeamId      string `json:"teamId"`
	ChannelId   string `json:"channelId"`
	// Command         string `json:"command"`
	Text         string    `json:"text"`
	Timestamp    int64     `json:"timestamp"`
	UserId       string    `json:"userId"`
	MsgId        string    `json:"msgId"`
	SourceUserId string    `json:"sourceUserId"`
	Ttl          int64     `json:"ttl"`
	Type         KudosType `json:"type"`
}

type KudosCounter struct {
	Id1       string    `json:"id1"` //teamId#report
	Id2       string    `json:"id2"` //<weekNumber | monthNumber>#userId
	Timestamp int64     `json:"timestamp"`
	Count     int64     `json:"count"`
	UserId    string    `json:"userId"`
	Type      KudosType `json:"type"`
}

type KudosType string

const (
	CommandType KudosType = "command"
	ReportType  KudosType = "report"
)
