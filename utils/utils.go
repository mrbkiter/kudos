package utils

import (
	"context"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"kudos-app.github.com/model"
)

func NewMyContext(ctx context.Context, userInternalId string) *model.MyContext {
	lc, _ := lambdacontext.FromContext(ctx)
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	myctx := new(model.MyContext)
	myctx.Testing = false
	sugarLogger := logger.Sugar().With(zap.String("userInternalId", userInternalId))
	if lc != nil {
		sugarLogger = sugarLogger.With(zap.String("rqId", lc.AwsRequestID))
		myctx.AwsRequestId = lc.AwsRequestID
	}
	myctx.UserInternalId = userInternalId
	myctx.Log = sugarLogger
	myctx.Ctx = ctx
	return myctx
}

func GenerateUuid() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func ParseDateFormatISO(dateValue string) (time.Time, error) {
	return time.Parse("2006-01-02", dateValue)
}

func ValidateDateFormatISO(dateValue string) bool {
	_, err := ParseDateFormatISO(dateValue)
	if err != nil {
		return false
	}
	return true
}

func ParseDateTimeFormatISO(dateValue string) (*time.Time, error) {
	if len(dateValue) == 0 {
		return nil, nil
	}
	v, err := time.Parse("2006-01-02T15:04:05.000Z", dateValue)
	return &v, err
}

func ValidateDateTimeFormatISO(dateValue string) bool {
	if len(dateValue) == 0 {
		return true
	}
	_, err := ParseDateTimeFormatISO(dateValue)
	if err != nil {
		return false
	}
	return true
}

var textUserIdCoveredRegex = regexp.MustCompile(`<[A-Za-z0-9@|]*>|@[^ ]*`)
var userIdRegex = regexp.MustCompile(`[^@<>|]+`)

func ExtractUserIdsFromText(text string) []string {
	// userIds := textUserIdRegex.FindAllString(text, -1)
	userIdSet := make(map[string]bool)
	userIds := textUserIdCoveredRegex.FindAllString(text, -1)
	for _, userId := range userIds {
		userId1 := userIdRegex.FindAllString(userId, 1)[0]
		userIdSet[userId1] = true

	}
	v := make([]string, 0, len(userIdSet))
	for userId, _ := range userIdSet {
		v = append(v, userId)
	}
	return v
}

func ReadFile(fileName string) (string, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
