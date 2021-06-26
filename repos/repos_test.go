package repos_test

import (
	"fmt"
	"testing"
	"time"

	"kudos-app.github.com/model"
	"kudos-app.github.com/repos"
)

func Test_mapKudosDataToKudosCommandEtt(t *testing.T) {
	data := new(model.KudosData)
	tgUserIds := make([]string, 2)
	tgUserIds[0] = "1"
	tgUserIds[1] = "2"
	data.TargetUserIds = tgUserIds
	data.TeamId = "teamId"
	data.MessageId = "MessageId"
	ett := repos.MapKudosDataToKudosCommandEtt(data)
	monthFormat := time.Now().Format("2006-01")
	yearNumber, weekNumber := time.Now().ISOWeek()

	if len(ett) != 2 {
		t.Error("MapKudosDataToKudosCommandEtt should return 2 elements")
	}
	if ett[0].TeamId != "teamId" || ett[1].TeamId != "teamId" {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong teamId")
	}
	if ett[0].MsgId != "MessageId" || ett[1].MsgId != "MessageId" {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong MessageId")
	}
	if ett[0].Id1 != "teamId#1" || ett[1].Id1 != "teamId#2" {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong Id1")
	}
	if ett[0].TeamIdMonth != ("teamId#"+monthFormat) || ett[1].TeamIdMonth != ("teamId#"+monthFormat) {
		t.Error("MapKudosDataToKudosCommandEtt mapped wrong TeamIdMonth")
	}
	if ett[0].TeamIdWeek != fmt.Sprintf("teamId#%v#%v", yearNumber, weekNumber) || ett[0].TeamIdWeek != fmt.Sprintf("teamId#%v#%v", yearNumber, weekNumber) {
		t.Error("MapKudosDataToKudosCommandEtt mapped TeamIdWeek")
	}
}
