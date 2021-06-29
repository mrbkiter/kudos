package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"kudos-app.github.com/ddb_entity"
	"kudos-app.github.com/model"
	repos "kudos-app.github.com/repos"
	"kudos-app.github.com/utils"
)

var globalConfig *model.GlobalConfig = new(model.GlobalConfig)
var repo repos.Repo

func initRepo() repos.Repo {
	region := globalConfig.Ddb.Region
	if len(region) == 0 {
		region = "ap-southeast-1"
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	log.Printf("Initialize dynamodb repo with region %v, table %v\n", globalConfig.Ddb.Region, globalConfig.Ddb.TableName)
	awsConfig := &aws.Config{}
	ddb := dynamodb.New(sess, awsConfig)
	ddbRepo := new(repos.DDBRepo)
	ddbRepo.Ddb = ddb
	return ddbRepo
}

func initEnvironment() {
	profile := os.Getenv("PROFILE")
	log.Println("PROFILE ", profile)
	//read config file
	globalConfigPath := fmt.Sprintf("./config/config-%v.json", profile)
	err := utils.ReadFromJSONFile(globalConfigPath, globalConfig)
	if err != nil {
		panic(err)
	}
	log.Printf("Global Config loaded %v\n", *globalConfig)
}

func main() {
	log.Println("Started Kudos Report")
	initEnvironment()
	repo = initRepo()
	lambda.Start(handler)
}

func handler(ctx context.Context, ddbEvent events.DynamoDBEvent) error {
	myCtx := utils.NewMyContext(context.Background(), "")
	myCtx.GlobalConfig = globalConfig
	for _, record := range ddbEvent.Records {
		// newctx.Log.Infof("[%s %s], old = %s, new = %s", record.EventName, record.Change.Keys, record.Change.OldImage, record.Change.NewImage)
		oldData := new(ddb_entity.KudosCommand)
		newData := new(ddb_entity.KudosCommand)

		utils.UnmarshalStreamImage(record.Change.OldImage, oldData)
		utils.UnmarshalStreamImage(record.Change.NewImage, newData)
		counter := new(model.KudosCountUpdate)
		if record.EventName == "INSERT" {

		} else if record.EventName == "MODIFY" {

		} else if record.EventName == "REMOVE" {

		}

		repo.IncreaseKudosCounter(myCtx, counter)

	}
	return nil
}
