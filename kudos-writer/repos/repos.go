package repos

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DDBInterface interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

type Repo interface {
}

type DDBRepo struct {
	Repo
	Ddb DDBInterface
}
