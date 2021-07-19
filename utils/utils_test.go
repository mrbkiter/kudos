package utils_test

import (
	"testing"

	"kudos-app.github.com/utils"
)

func Test_ExtractUserIdsFromText(t *testing.T) {
	test1 := "ask <@U012ABCDEF> to bake a birthday cake for @mrbkiter and <@U345GHIJKL|mrbkiter> in <#C012ABCDE>"
	userIds := utils.ExtractUserIdsFromText(test1, "U345GHIJKL")
	if len(userIds) != 2 {
		t.Error("Errror when extracting userIds")
	}
	if userIds[0].UserId != "U012ABCDEF" || userIds[1].UserId != "mrbkiter" {
		t.Error("Incorrectly extract user id")
	}
}

func Test_ExtractKudosText(t *testing.T) {
	test1 := "good job <@U012ABCDEF> <@U345GHIJKL|mrbkiter> for helping"
	kudosText := utils.AnalyzeKudosText(test1)
	if kudosText != "good job <@U012ABCDEF> <@U345GHIJKL|mrbkiter> " {
		t.Error("Extract kudos text error")
	}

	test2 := "thanks <@U012ABCDEF> <@U345GHIJKL|mrbkiter> for helping <@12233> on something"
	kudosText = utils.AnalyzeKudosText(test2)
	if kudosText != "thanks <@U012ABCDEF> <@U345GHIJKL|mrbkiter> " {
		t.Error("Extract kudos text error")
	}

	test3 := "great work <@U012ABCDEF> <@U345GHIJKL|mrbkiter>. You did it well"
	kudosText = utils.AnalyzeKudosText(test3)
	if kudosText != "great work <@U012ABCDEF> <@U345GHIJKL|mrbkiter>" {
		t.Error("Extract kudos text error")
	}

	test4 := " <@U012ABCDEF> <@U345GHIJKL|mrbkiter>. You did it well"
	kudosText = utils.AnalyzeKudosText(test4)
	if kudosText != " <@U012ABCDEF> <@U345GHIJKL|mrbkiter>" {
		t.Error("Extract kudos text error")
	}

	test5 := "<@U012ABCDEF> <@U345GHIJKL|mrbkiter>. You did it well"
	kudosText = utils.AnalyzeKudosText(test5)
	if kudosText != "<@U012ABCDEF> <@U345GHIJKL|mrbkiter>" {
		t.Error("Extract kudos text error")
	}

	test6 := "Thanks <@U012ABCDEF> <@U345GHIJKL|mrbkiter>. You did it well"
	kudosText = utils.AnalyzeKudosText(test6)
	if kudosText != "Thanks <@U012ABCDEF> <@U345GHIJKL|mrbkiter>" {
		t.Error("Extract kudos text error")
	}

	test7 := "<@U024D6VQX7Z|vu.yen.nguyen.88> <@U024U032H8A|vu.nguyen>"
	kudosText = utils.AnalyzeKudosText(test7)
	if kudosText != "<@U024D6VQX7Z|vu.yen.nguyen.88> <@U024U032H8A|vu.nguyen>" {
		t.Error("Extract kudos text error")
	}
}

func Test_ExtractGroupId(t *testing.T) {
	groupId := utils.ExtractGroupId("group-Id <@U012ABCDEF|mrbkiter> ")
	if groupId != "group-Id" {
		t.Error("Error when extracting groupId")
	}
	groupId = utils.ExtractGroupId(" <@U012ABCDEF|mrbkiter> ")
	if groupId != "" {
		t.Error("Error when extracting groupId")
	}
	groupId = utils.ExtractGroupId("group_id")
	if groupId != "group_id" {
		t.Error("Error when extracting groupId")
	}
	groupId = utils.ExtractGroupId("THIS_MONTH")
	if groupId != "" {
		t.Error("Error when extracting groupId")
	}

}
