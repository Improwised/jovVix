package v1

import "testing"

func TestApplyRosterNamesPreservesLegacyCorrection(t *testing.T) {
	roster := []UserInfo{
		{UserId: "old", UserName: "Old first name", IsAlive: true},
		{UserId: "new", UserName: "New player", IsAlive: true},
	}
	applyRosterNames(roster, map[string]string{"old": "Quiz name"})
	if roster[0].UserName != "Quiz name" || !roster[0].IsAlive {
		t.Fatalf("cached player was not corrected: %+v", roster[0])
	}
	if roster[1].UserName != "New player" {
		t.Fatalf("new player was changed: %+v", roster[1])
	}
}
