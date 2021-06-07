package utils_test

import (
	"testing"

	"kudos-app.github.com/utils"
)

func Test_ExtractUserIdsFromText(t *testing.T) {
	test1 := "ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL|mrbkiter> in <#C012ABCDE>"
	userIds := utils.ExtractUserIdsFromText(test1)
	if len(userIds) != 2 {
		t.Error("Errror when extracting userIds")
	}
	if userIds[0] != "U012ABCDEF" || userIds[1] != "U345GHIJKL" {
		t.Error("Incorrectly extract user id")
	}
}
