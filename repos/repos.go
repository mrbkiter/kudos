package repos

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"kudos-app.github.com/ddb_entity"
	"kudos-app.github.com/model"
)

type DDBInterface interface {
	UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
	TransactWriteItems(input *dynamodb.TransactWriteItemsInput) (*dynamodb.TransactWriteItemsOutput, error)
}

type Repo interface {
	WriteKudosCommand(ctx *model.MyContext, data *model.KudosData) error
	IncreaseKudosCounter(ctx *model.MyContext, data *model.KudosCountUpdate) error
}

type DDBRepo struct {
	Repo
	Ddb DDBInterface
}

func (me *DDBRepo) IncreaseKudosCounter(ctx *model.MyContext, data *model.KudosCountUpdate) error {
	id1 := buildPartitionKey(data.TeamId, string(ddb_entity.ReportType))
	timestamp := time.Unix(0, data.Timestamp*int64(time.Millisecond))

	//build 2 counter
	monthFormat := timestamp.Format("2006-01")
	yearNumber, weekNumber := timestamp.ISOWeek()
	id2Month := buildPartitionKey(data.UserId, monthFormat)
	id2WeekNumber := buildPartitionKey(data.UserId, fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))

	transactInputs := &dynamodb.TransactWriteItemsInput{}

	monthInput := buildTransactionUpdateItem(ctx, id1, id2Month, data.Counter)
	weekInput := buildTransactionUpdateItem(ctx, id1, id2WeekNumber, data.Counter)
	transactInputs.TransactItems = make([]*dynamodb.TransactWriteItem, 2)
	transactInputs.TransactItems[0] = monthInput
	transactInputs.TransactItems[1] = weekInput

	_, err := me.Ddb.TransactWriteItems(transactInputs)
	if err != nil {
		//need to check if err is from record not exists
		putMonthInput := buildTransactionPutItem(ctx, id1, id2Month, data.Counter)
		putWeekInput := buildTransactionPutItem(ctx, id1, id2WeekNumber, data.Counter)
		transactInputs.TransactItems = make([]*dynamodb.TransactWriteItem, 2)
		transactInputs.TransactItems[0] = putMonthInput
		transactInputs.TransactItems[1] = putWeekInput

		_, err := me.Ddb.TransactWriteItems(transactInputs)
		return err
	}
	return err
}

func buildTransactionUpdateItem(ctx *model.MyContext, id1 string, id2 string, counter int) *dynamodb.TransactWriteItem {
	return &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":val": {
					N: aws.String(fmt.Sprintf("%v", counter)),
				},
			},
			TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
			Key: map[string]*dynamodb.AttributeValue{
				"id1": {
					N: aws.String(id1),
				},
				"id2": {
					S: aws.String(id2),
				},
			},
			UpdateExpression: aws.String("set count = count + :val"),
		},
	}
}
func buildTransactionPutItem(ctx *model.MyContext, id1 string, id2 string, counter int) *dynamodb.TransactWriteItem {
	return &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":val": {
					N: aws.String(string(counter)),
				},
			},
			TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
			Item: map[string]*dynamodb.AttributeValue{
				"id1": {
					S: aws.String(id1),
				},
				"id2": {
					S: aws.String(id2),
				},
				"count": {
					N: aws.String(fmt.Sprintf("%v", counter)),
				},
			},
		},
	}
}

func (me *DDBRepo) WriteKudosCommand(ctx *model.MyContext, data *model.KudosData) error {
	if len(data.TargetUserIds) == 0 {
		return nil
	}
	commands := MapKudosDataToKudosCommandEtt(data)

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

func MapKudosDataToKudosCommandEtt(data *model.KudosData) []*ddb_entity.KudosCommand {
	commands := make([]*ddb_entity.KudosCommand, 0)
	now := time.Now().Unix()
	ttl := time.Now().Add(3 * 30 * 24 * time.Hour).Unix() //keep 90 days
	monthFormat := time.Now().Format("2006-01")
	yearNumber, weekNumber := time.Now().ISOWeek()

	for _, targetUserId := range data.TargetUserIds {
		ett := new(ddb_entity.KudosCommand)
		ett.Type = ddb_entity.CommandType
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

func ConvertToWriteRequest(method interface{}) *dynamodb.WriteRequest {
	attributes, _ := dynamodbattribute.MarshalMap(method)
	input := &dynamodb.PutRequest{Item: attributes}
	return &dynamodb.WriteRequest{PutRequest: input}

}
