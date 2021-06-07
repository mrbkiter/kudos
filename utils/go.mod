module kudos-app.github.com/utils

go 1.15

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.17.0
	kudos-app.github.com/model v0.0.0-00010101000000-000000000000
)

replace kudos-app.github.com/model => ../model
