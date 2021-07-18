module kudos-app.github.com/slack

go 1.15

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.70
	github.com/joncalhoun/qson v0.0.0-20200422171543-84433dcd3da0
	kudos-app.github.com/model v0.0.0-00010101000000-000000000000
	kudos-app.github.com/repos v0.0.0-00010101000000-000000000000
	kudos-app.github.com/utils v0.0.0-00010101000000-000000000000
)

replace kudos-app.github.com/model => ../model

replace kudos-app.github.com/repos => ../repos

replace kudos-app.github.com/utils => ../utils

replace kudos-app.github.com/ddb_entity => ../ddb_entity
