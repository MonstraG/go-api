package ytDlp

import (
	"testing"
)

func TestSanitizeUrl(t *testing.T) {
	songUrl := "https://music.youtube.com/watch?v=UW6a0KRC7_8&list=RDAMVMUW6a0KRC7_8"
	got := SanitizeUrl(songUrl)
	want := "https://music.youtube.com/watch?v=UW6a0KRC7_8"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
