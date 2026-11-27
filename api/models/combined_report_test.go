package models

import (
	"reflect"
	"testing"
)

func TestCombineLeaderboard(t *testing.T) {
	plays := []CombinedPlay{
		{PlayID: "1", UserID: "a", SessionID: "s1", Name: "Asha", Score: 50, ResponseTime: 3000},
		{PlayID: "2", UserID: "b", SessionID: "s1", Name: "Bilal", Score: 30, ResponseTime: 4000},
		{PlayID: "3", UserID: "a", SessionID: "s2", Name: "Asha", Score: 70, ResponseTime: 2000},
		{PlayID: "4", UserID: "c", SessionID: "s2", Name: "Chen", Score: 70, ResponseTime: 1000},
		{PlayID: "5", UserID: "d", SessionID: "s2", Name: "Dana", Score: 30, ResponseTime: 1500},
	}
	type row struct {
		Name    string
		Session string
		Score   int
		Rank    int
		Plays   int
	}
	var got []row
	for _, e := range CombineLeaderboard(plays) {
		got = append(got, row{e.Name, e.SessionID, e.Score, e.Rank, e.Plays})
	}
	// Asha keeps her better second play and appears once. Equal scores share a
	// rank and list the faster player first.
	want := []row{
		{"Chen", "s2", 70, 1, 1},
		{"Asha", "s2", 70, 1, 2},
		{"Dana", "s2", 30, 2, 1},
		{"Bilal", "s1", 30, 2, 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v\nwant %+v", got, want)
	}
}

func TestCombineLeaderboardBestPlayTieKeepsFasterPlay(t *testing.T) {
	plays := []CombinedPlay{
		{PlayID: "1", UserID: "a", SessionID: "s1", Score: 40, ResponseTime: 5000},
		{PlayID: "2", UserID: "a", SessionID: "s2", Score: 40, ResponseTime: 900},
	}
	got := CombineLeaderboard(plays)
	if len(got) != 1 || got[0].SessionID != "s2" || got[0].Plays != 2 {
		t.Fatalf("got %+v", got)
	}
	if len(CombineLeaderboard(nil)) != 0 {
		t.Fatal("no plays should give an empty leaderboard")
	}
}
