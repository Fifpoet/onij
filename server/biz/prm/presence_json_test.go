package prm

import (
	"encoding/json"
	"testing"
)

func TestPresenceTrackDTOAcceptsFloatTimeLong(t *testing.T) {
	var tr PresenceTrackDTO
	raw := []byte(`{"id":186016,"name":"晴天","artists":[],"time_long":269.848,"album_id":2,"album_name":"叶惠美"}`)
	if err := json.Unmarshal(raw, &tr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tr.Id != 186016 || tr.Name != "晴天" {
		t.Fatalf("track = %+v", tr)
	}
	if tr.TimeLong < 269 || tr.TimeLong > 270 {
		t.Fatalf("time_long = %v", tr.TimeLong)
	}
}
