package repos

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"kudos-app.github.com/model"
)

type DDBInterface interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
}

type Repo interface {
	WriteKudosCommand(ctx *model.MyContext, data *model.KudosData) error
}

type DDBRepo struct {
	Repo
	Ddb DDBInterface
}

func (me *DDBRepo) WriteKudosCommand(ctx *model.MyContext, data *model.KudosData) error {
	commands := mapKudosDataToKudosCommandEtt(data)

	commandWrites := make([]*dynamodb.WriteRequest, len(commands))

	for idx, cmd := range commands {
		commandWrites[idx] = ConvertToWriteRequest(cmd)
	}

	batchInputs := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			ctx.GlobalConfig.Ddb.TableName: commandWrites,
		},
	}
	_, err := me.Ddb.BatchWriteItem(batchInputs)
	if err != nil {
		ctx.Log.Errorf("Error writing result %v to ddb because [%v]", err)
		return err
	}
	return nil
}

func buildPartitionKey(p1 string, p2 string, p3 ...string) string {
	key := p1 + "#" + p2
	for _, s := range p3 {
		key = key + "#" + s
	}
	return key
}

func mapKudosDataToKudosCommandEtt(data *model.KudosData) []*KudosCommand {
	commands := make([]*KudosCommand, 0)
	now := time.Now().Unix()
	ttl := time.Now().Add(3 * 30 * 24 * time.Hour).Unix() //keep 90 days
	monthFormat := time.Now().Format("2006-01")
	yearNumber, weekNumber := time.Now().ISOWeek()

	for _, targetUserId := range data.TargetUserIds {
		ett := new(KudosCommand)
		ett.ChannelId = data.ChannelId
		ett.Id1 = buildPartitionKey(data.TeamId, targetUserId)
		ett.Id2 = data.MessageId
		ett.SourceUserId = data.SourceUserId
		ett.Text = data.Text
		ett.UserId = targetUserId
		ett.TeamId = data.TeamId
		ett.MsgId = data.MessageId
		ett.Timestamp = now
		ett.Ttl = ttl
		ett.TeamIdWeek = buildPartitionKey(data.TeamId, fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))
		ett.TeamIdMonth = buildPartitionKey(data.TeamId, monthFormat)
		commands = append(commands, ett)
	}
	return commands
}

type KudosCommand struct {
	Id1         string `json:"id1"`
	Id2         string `json:"id2"`
	TeamIdMonth string `json:"teamIdMonth"`
	TeamIdWeek  string `json:"teamIdWeek"`
	TeamId      string `json:"teamId"`
	ChannelId   string `json:"channelId"`
	// Command         string `json:"command"`
	Text         string `json:"text"`
	Timestamp    int64  `json:"timestamp"`
	UserId       string `json:"userId"`
	MsgId        string `json:"msgId"`
	SourceUserId string `json:"sourceId"`
	Ttl          int64  `json:"ttl"`
}

func ConvertToWriteRequest(method interface{}) *dynamodb.WriteRequest {
	attributes, _ := dynamodbattribute.MarshalMap(method)
	input := &dynamodb.PutRequest{Item: attributes}
	return &dynamodb.WriteRequest{PutRequest: input}

}
