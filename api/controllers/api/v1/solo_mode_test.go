package v1

import "testing"

func TestHasOtherConnectedPlayer(t *testing.T) {
	for _, tc := range []struct {
		name   string
		roster []UserInfo
		want   bool
	}{
		{"empty", nil, false},
		{"host alone", []UserInfo{{UserId: "host", IsAlive: true}}, false},
		{"disconnected guest", []UserInfo{{UserId: "host", IsAlive: true}, {UserId: "guest", IsAlive: false}}, false},
		{"connected guest", []UserInfo{{UserId: "guest", IsAlive: true}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasOtherConnectedPlayer(tc.roster, "host"); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}
