package model

type GlobalConfig struct {
	Ddb DynamoConfig `json:"ddbConfig"`
}

type DynamoConfig struct {
	Region    string `json:"region"`
	TableName string `json:"tableName"`
}
