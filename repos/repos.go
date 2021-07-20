package repos

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/guregu/dynamo"
	"kudos-app.github.com/ddb_entity"
	"kudos-app.github.com/model"
	"kudos-app.github.com/utils"
)

type DDBInterface interface {
	UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
	TransactWriteItems(input *dynamodb.TransactWriteItemsInput) (*dynamodb.TransactWriteItemsOutput, error)
	DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

type Repo interface {
	WriteKudosCommand(ctx *model.MyContext, data *model.KudosData) error
	IncreaseKudosCounter(ctx *model.MyContext, data *model.KudosCountUpdate) error
	GetKudosReport(ctx *model.MyContext, filter *model.KudosReportFilter) (*model.KudosReportResult, error)
	GetKudosReportDetails(ctx *model.MyContext, filter *model.KudosReportFilter) (*model.KudosReportDetails, error)
	DeleteTeamGroup(ctx *model.MyContext, req *model.KudosTeamSettingsGroupAction) error
	AddTeamGroup(ctx *model.MyContext, req *model.KudosTeamSettingsGroupAction) error
	AddMembersToTeamGroup(ctx *model.MyContext, req *model.KudosGroupSettingsMembers) error
	DeleteMembersFromTeamGroup(ctx *model.MyContext, req *model.KudosGroupSettingsMembers) error
	ListMembersOfTeamGroup(ctx *model.MyContext, req *model.KudosTeamSettingsGroupAction) ([]string, error)
	ListAllTeamGroups(ctx *model.MyContext, req *model.KudosTeamSettingsListGroups) ([]string, error)
	doActionMembersToTeamGroup(ctx *model.MyContext, req *model.KudosGroupSettingsMembers, actionType string) error
}

type DDBRepo struct {
	Repo
	Ddb DDBInterface
}

func MapToGroupSettings(req *model.KudosTeamSettingsGroupAction) *ddb_entity.KudosGroupSettings {
	ettGroupSettings := new(ddb_entity.KudosGroupSettings)
	ettGroupSettings.Id1 = buildPartitionKey(req.TeamId, string(ddb_entity.GroupSettings))
	ettGroupSettings.Id2 = req.GroupId
	ettGroupSettings.GroupId = req.GroupId
	ettGroupSettings.TeamId = req.TeamId
	ettGroupSettings.Timestamp = time.Now().Unix()
	ettGroupSettings.Type = ddb_entity.GroupSettings

	return ettGroupSettings
}

func (me *DDBRepo) doActionMembersToTeamGroup(ctx *model.MyContext, req *model.KudosGroupSettingsMembers, actionType string) error {
	addedUserIds := make([]*string, 0)

	for _, val := range req.TargetUserIds {
		addedUserIds = append(addedUserIds, &val.UserId)
	}

	actionExpressionQuery := "ADD userIds :val"
	if actionType == "delete" {
		actionExpressionQuery = "DELETE userIds :val"
	}
	writeItem := &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":val": {
					SS: addedUserIds,
				},
			},
			TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
			Key: map[string]*dynamodb.AttributeValue{
				"id1": {
					S: aws.String(buildPartitionKey(req.TeamId, string(ddb_entity.GroupSettings))),
				},
				"id2": {
					S: aws.String(req.GroupId),
				},
			},
			UpdateExpression: aws.String(actionExpressionQuery),
		},
	}
	transactInputs := &dynamodb.TransactWriteItemsInput{}
	transactInputs.TransactItems = make([]*dynamodb.TransactWriteItem, 1)
	transactInputs.TransactItems[0] = writeItem
	_, err := me.Ddb.TransactWriteItems(transactInputs)
	return err
}

func (me *DDBRepo) AddMembersToTeamGroup(ctx *model.MyContext, req *model.KudosGroupSettingsMembers) error {
	return me.doActionMembersToTeamGroup(ctx, req, "add")
}

func (me *DDBRepo) DeleteMembersFromTeamGroup(ctx *model.MyContext, req *model.KudosGroupSettingsMembers) error {
	return me.doActionMembersToTeamGroup(ctx, req, "delete")
}

func (me *DDBRepo) ListMembersOfTeamGroup(ctx *model.MyContext, req *model.KudosTeamSettingsGroupAction) ([]string, error) {
	q := &dynamodb.QueryInput{
		TableName:              aws.String(ctx.GlobalConfig.Ddb.TableName),
		KeyConditionExpression: utils.String("id1 = :id1 AND id2 = :id2"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id1": {
				S: utils.String(buildPartitionKey(req.TeamId, string(ddb_entity.GroupSettings))),
			},
			":id2": {
				S: utils.String(req.GroupId),
			},
		},
		Limit: utils.Int64(1),
	}
	output, err := me.Ddb.Query(q)
	if err != nil {
		return nil, err
	}
	members := make([]string, 0)
	if len(output.Items) > 0 {
		groups := make([]*ddb_entity.KudosGroupSettings, 0)
		err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &groups)
		if err != nil {
			return nil, err
		}
		members = groups[0].UserIds
	}
	return members, nil
}
func (me *DDBRepo) ListAllTeamGroups(ctx *model.MyContext, req *model.KudosTeamSettingsListGroups) ([]string, error) {
	q := &dynamodb.QueryInput{
		TableName:              aws.String(ctx.GlobalConfig.Ddb.TableName),
		KeyConditionExpression: utils.String("#id1 = :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				S: utils.String(buildPartitionKey(req.TeamId, string(ddb_entity.GroupSettings))),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#id1": utils.String("id1"),
		},
		Limit:                utils.Int64(500),
		ProjectionExpression: utils.String("id1,groupId"),
	}
	output, err := me.Ddb.Query(q)
	if err != nil {
		return nil, err
	}
	groupIds := make([]string, 0)
	if len(output.Items) > 0 {
		groups := make([]*ddb_entity.KudosGroupSettings, 0)
		err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &groups)
		if err != nil {
			return nil, err
		}
		for _, id := range groups {
			groupIds = append(groupIds, id.GroupId)
		}
	}
	return groupIds, nil
}
func (me *DDBRepo) DeleteTeamGroup(ctx *model.MyContext, req *model.KudosTeamSettingsGroupAction) error {
	if len(req.GroupId) == 0 || len(req.TeamId) == 0 {
		return errors.New("GroupId And TeamId is required")
	}

	settings := MapToGroupSettings(req)

	delete := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id1": {
				S: utils.String(settings.Id1),
			},
			"id2": {
				S: utils.String(settings.Id2),
			},
		},
		TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
	}
	_, err := me.Ddb.DeleteItem(delete)
	if err != nil {
		return err
	}
	return nil
}

func (me *DDBRepo) AddTeamGroup(ctx *model.MyContext, req *model.KudosTeamSettingsGroupAction) error {
	if len(req.GroupId) == 0 || len(req.TeamId) == 0 {
		return errors.New("GroupId And TeamId is required")
	}

	settings := MapToGroupSettings(req)
	av, err := dynamo.MarshalItem(settings)
	if err != nil {
		return err
	}

	key := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id1": {
				S: utils.String(settings.Id1),
			},
			"id2": {
				S: utils.String(settings.Id2),
			},
		},
		TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
	}
	val, err := me.Ddb.GetItem(key)
	if val != nil && len(val.Item) > 0 {
		return errors.New("It does not allow to recreate TeamGroup")
	}

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
	}

	_, err = me.Ddb.PutItem(input)
	if err != nil {
		ctx.Log.Errorf("Error writing result %v to ddb because [%v]", req, err)
		return err
	}

	return nil
}

func (me *DDBRepo) GetKudosReportDetails(ctx *model.MyContext, filter *model.KudosReportFilter) (*model.KudosReportDetails, error) {
	if len(filter.UserIds) != 1 {
		return nil, errors.New("Only 1 user accepted at one time")
	}
	calculatedTime := time.Now()
	id1 := ""
	indexName := ""
	id1AttrName := ""
	switch filter.ReportTime {
	case model.THIS_MONTH:
		monthFormat := calculatedTime.Format("2006-01")
		id1 = buildPartitionKey(filter.TeamId, monthFormat)
		indexName = "teamIdMonth-userId-index"
		id1AttrName = "teamIdMonth"
	case model.LAST_MONTH:
		calculatedTime = calculatedTime.AddDate(0, -1, 0)
		monthFormat := calculatedTime.Format("2006-01")
		id1 = buildPartitionKey(filter.TeamId, monthFormat)
		indexName = "teamIdMonth-userId-index"
		id1AttrName = "teamIdMonth"
	case model.LAST_WEEK:
		calculatedTime = calculatedTime.Local().AddDate(0, 0, -7)
		yearNumber, weekNumber := calculatedTime.ISOWeek()
		id1 = buildPartitionKey(filter.TeamId, fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))
		indexName = "teamIdWeek-userId-index"
		id1AttrName = "teamIdWeek"
	case model.THIS_WEEK:
		yearNumber, weekNumber := calculatedTime.ISOWeek()
		id1 = buildPartitionKey(filter.TeamId, fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))
		indexName = "teamIdWeek-userId-index"
		id1AttrName = "teamIdWeek"
	}

	q := &dynamodb.QueryInput{
		TableName:              aws.String(ctx.GlobalConfig.Ddb.TableName),
		KeyConditionExpression: utils.String("#id1 = :id1 AND userId = :userId"),
		IndexName:              aws.String(indexName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id1": {
				S: utils.String(id1),
			},
			":userId": {
				S: utils.String(filter.UserIds[0]),
			},
			":type": {
				S: utils.String("command"),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#id1":  aws.String(id1AttrName),
			"#type": aws.String("type"),
		},
		FilterExpression: aws.String("#type = :type"),
		Limit:            utils.Int64(500),
	}
	output, err := me.Ddb.Query(q)
	if err != nil {
		return nil, err
	}

	ret := new(model.KudosReportDetails)

	if len(output.Items) > 0 {
		kudosCommands := make([]*ddb_entity.KudosCommand, 0)
		err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &kudosCommands)
		if err != nil {
			return nil, err
		}
		ret.UserId = filter.UserIds[0]

		ret.KudosTalk = make([]*model.KudosSimpleCommand, 0)
		for _, cmd := range kudosCommands {
			simpleCmd := new(model.KudosSimpleCommand)
			simpleCmd.Text = cmd.Text
			simpleCmd.UserId = cmd.SourceUserId
			simpleCmd.Timestamp = time.Unix(cmd.Timestamp, 0)
			ret.KudosTalk = append(ret.KudosTalk, simpleCmd)
		}
		sort.SliceStable(ret.KudosTalk, func(i, j int) bool {
			return ret.KudosTalk[i].Timestamp.After(ret.KudosTalk[j].Timestamp)
		})
	}
	return ret, nil

}
func (me *DDBRepo) GetKudosReport(ctx *model.MyContext, filter *model.KudosReportFilter) (*model.KudosReportResult, error) {
	id1 := buildPartitionKey(filter.TeamId, "report")
	q := &dynamodb.QueryInput{
		TableName:              aws.String(ctx.GlobalConfig.Ddb.TableName),
		KeyConditionExpression: utils.String("id1 = :id1"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id1": {
				S: utils.String(id1),
			},
		},
		Limit: utils.Int64(500),
	}
	timeFilter := "%v#"

	calculatedTime := time.Now()
	fromTime := time.Now()
	toTime := time.Now()

	switch filter.ReportTime {
	case model.THIS_MONTH:
		fromTime = time.Date(calculatedTime.Year(), calculatedTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		toTime = time.Now()
		monthFormat := calculatedTime.Format("2006-01")
		timeFilter = fmt.Sprintf(timeFilter, monthFormat)
	case model.THIS_WEEK:
		fromTime = utils.WeekStart(calculatedTime.UTC().ISOWeek())
		toTime = time.Now()
		yearNumber, weekNumber := calculatedTime.ISOWeek()
		id2WeekNumber := buildPartitionKey(fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))
		timeFilter = fmt.Sprintf(timeFilter, id2WeekNumber)
	case model.LAST_MONTH:
		calculatedTime = calculatedTime.AddDate(0, -1, 0)
		fromTime = time.Date(calculatedTime.Year(), calculatedTime.Month(), 1, 0, 0, 0, 0, time.UTC)
		toTime = calculatedTime.AddDate(0, 1, -1)
		monthFormat := calculatedTime.Format("2006-01")
		timeFilter = fmt.Sprintf(timeFilter, monthFormat)
	case model.LAST_WEEK:
		calculatedTime = calculatedTime.Local().AddDate(0, 0, -7)
		fromTime = utils.WeekStart(calculatedTime.UTC().ISOWeek())
		toTime = calculatedTime
		yearNumber, weekNumber := calculatedTime.ISOWeek()
		id2WeekNumber := buildPartitionKey(fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))
		timeFilter = fmt.Sprintf(timeFilter, id2WeekNumber)
	}
	if len(filter.GroupId) > 0 {
		req := new(model.KudosTeamSettingsGroupAction)
		req.GroupId = filter.GroupId
		req.TeamId = filter.TeamId
		members, err := me.ListMembersOfTeamGroup(ctx, req)
		if err != nil {
			return nil, err
		}
		filter.UserIds = members
	}
	if len(filter.UserIds) > 0 { //get specific users
		q.KeyConditionExpression = aws.String(fmt.Sprintf("%v AND begins_with(id2, :filterTime)", *q.KeyConditionExpression))
		q.ExpressionAttributeValues[":filterTime"] = &dynamodb.AttributeValue{
			S: utils.String(timeFilter),
		}

		userFilterExpr := `userId = %v`
		filterExpr := ""
		q.FilterExpression = aws.String("")
		for idx, id := range filter.UserIds {
			attrValue := fmt.Sprintf(":userId%v", idx)
			q.ExpressionAttributeValues[attrValue] = &dynamodb.AttributeValue{
				S: utils.String(id),
			}
			if idx == 0 {
				filterExpr = fmt.Sprintf(userFilterExpr, attrValue)
			} else {
				filterExpr = filterExpr + " OR " + fmt.Sprintf(userFilterExpr, attrValue)
			}
		}
		q.FilterExpression = aws.String(filterExpr)
	} else {
		q.KeyConditionExpression = aws.String(fmt.Sprintf("%v AND begins_with(id2, :filterTime)", *q.KeyConditionExpression))
		q.ExpressionAttributeValues[":filterTime"] = &dynamodb.AttributeValue{
			S: utils.String(timeFilter),
		}
	}

	output, err := me.Ddb.Query(q)
	if err != nil {
		return nil, err
	}

	ret := new(model.KudosReportResult)
	ret.FromTime = fromTime
	ret.ToTime = toTime

	if len(output.Items) > 0 {
		kudosUserReport := make([]*ddb_entity.KudosCounter, 0)
		err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &kudosUserReport)
		if err != nil {
			return nil, err
		}
		//convert
		result := make([]*model.KudosUserReportResult, 0)
		for _, r := range kudosUserReport {
			row := &model.KudosUserReportResult{
				UserId:   r.UserId,
				Total:    r.Count,
				Username: r.Username,
			}
			result = append(result, row)
		}
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Total >= result[j].Total
		})
		ret.UserReport = result
		return ret, nil
	}
	return ret, nil
}

func BuildIdMonthAndWeek(timestampInSec int64, userId string) (string, string) {
	timestamp := time.Unix(timestampInSec, 0)
	monthFormat := timestamp.Format("2006-01")
	yearNumber, weekNumber := timestamp.ISOWeek()
	id2Month := buildPartitionKey(monthFormat, userId)
	id2WeekNumber := buildPartitionKey(fmt.Sprint(yearNumber), fmt.Sprint(weekNumber), userId)
	return id2Month, id2WeekNumber
}

func (me *DDBRepo) IncreaseKudosCounter(ctx *model.MyContext, data *model.KudosCountUpdate) error {
	id1 := buildPartitionKey(data.TeamId, string(ddb_entity.ReportType))
	id2Month, id2WeekNumber := BuildIdMonthAndWeek(data.Timestamp, data.UserId)
	// timestamp := time.Unix(data.Timestamp, 0)

	//build month
	transactInputs := &dynamodb.TransactWriteItemsInput{}
	monthInput := buildTransactionUpdateItem(ctx, id1, id2Month, data)
	transactInputs.TransactItems = make([]*dynamodb.TransactWriteItem, 1)
	transactInputs.TransactItems[0] = monthInput
	_, err := me.Ddb.TransactWriteItems(transactInputs)
	if err != nil {
		ctx.Log.Warnw("Error when writing counter", "err", err)
		//need to check if err is from record not exists
		putMonthInput := buildTransactionPutItem(ctx, id1, id2Month, data)
		transactInputs.TransactItems = make([]*dynamodb.TransactWriteItem, 1)
		transactInputs.TransactItems[0] = putMonthInput

		_, err := me.Ddb.TransactWriteItems(transactInputs)
		if err != nil {
			return err
		}
	}

	//build week counter
	weekInput := buildTransactionUpdateItem(ctx, id1, id2WeekNumber, data)
	transactInputs.TransactItems[0] = weekInput
	_, err = me.Ddb.TransactWriteItems(transactInputs)
	if err != nil {
		ctx.Log.Warnw("Error when writing counter", "err", err)
		//need to check if err is from record not exists
		putWeekInput := buildTransactionPutItem(ctx, id1, id2WeekNumber, data)
		transactInputs.TransactItems = make([]*dynamodb.TransactWriteItem, 1)
		transactInputs.TransactItems[0] = putWeekInput
		_, err := me.Ddb.TransactWriteItems(transactInputs)
		return err
	}
	return nil
}

func buildTransactionUpdateItem(ctx *model.MyContext, id1 string, id2 string, counter *model.KudosCountUpdate) *dynamodb.TransactWriteItem {
	return &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":val": {
					N: aws.String(fmt.Sprintf("%v", counter.Counter)),
				},
				":userId": {
					S: aws.String(counter.UserId),
				},
				":username": {
					S: aws.String(counter.Username),
				},
				":teamId": {
					S: aws.String(counter.TeamId),
				},
			},
			TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
			Key: map[string]*dynamodb.AttributeValue{
				"id1": {
					S: aws.String(id1),
				},
				"id2": {
					S: aws.String(id2),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#count": aws.String("count"),
			},
			UpdateExpression: aws.String("set #count = #count + :val, userId=:userId, username = :username, teamId = :teamId"),
		},
	}
}
func buildTransactionPutItem(ctx *model.MyContext, id1 string, id2 string, counter *model.KudosCountUpdate) *dynamodb.TransactWriteItem {
	return &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName: aws.String(ctx.GlobalConfig.Ddb.TableName),
			Item: map[string]*dynamodb.AttributeValue{
				"id1": {
					S: aws.String(id1),
				},
				"id2": {
					S: aws.String(id2),
				},
				"count": {
					N: aws.String(fmt.Sprintf("%v", counter.Counter)),
				},
				"username": {
					S: aws.String(counter.Username),
				},
				"teamId": {
					S: aws.String(counter.TeamId),
				},
				"userId": {
					S: aws.String(counter.UserId),
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
		ett.Id1 = buildPartitionKey(data.TeamId, targetUserId.UserId)
		ett.Id2 = data.MessageId
		ett.SourceUserId = data.SourceUserId
		ett.Text = data.Text
		ett.UserId = targetUserId.UserId
		ett.Username = targetUserId.Username
		ett.TeamId = data.TeamId
		ett.MsgId = data.MessageId
		ett.Timestamp = now
		ett.Ttl = ttl
		ett.Username = targetUserId.Username
		ett.TeamIdWeek = buildPartitionKey(data.TeamId, fmt.Sprint(yearNumber), fmt.Sprint(weekNumber))
		ett.TeamIdMonth = buildPartitionKey(data.TeamId, monthFormat)
		ett.SourceUserName = data.Username
		commands = append(commands, ett)
	}
	return commands
}

func ConvertToWriteRequest(method interface{}) *dynamodb.WriteRequest {
	attributes, _ := dynamodbattribute.MarshalMap(method)
	input := &dynamodb.PutRequest{Item: attributes}
	return &dynamodb.WriteRequest{PutRequest: input}

}
