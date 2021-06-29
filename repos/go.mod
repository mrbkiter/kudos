module kudos-app.github.com/repos

go 1.15

require (
	github.com/aws/aws-sdk-go v1.38.55
	kudos-app.github.com/ddb_entity v0.0.0-00010101000000-000000000000
	kudos-app.github.com/model v0.0.0-00010101000000-000000000000
)

replace kudos-app.github.com/model => ../model

replace kudos-app.github.com/ddb_entity => ../ddb_entity
