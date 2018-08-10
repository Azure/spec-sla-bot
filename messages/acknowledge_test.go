package messages

import (
	"testing"
	"time"
)

func Test_Add_To_DB(t *testing.T) {

}

func Test_slaDuration(t *testing.T) {
	const oneDay = 24 * time.Hour
	expected := map[time.Weekday]time.Duration{
		time.Monday:   oneDay,
		time.Thursday: oneDay,
		time.Friday:   3 * oneDay,
		time.Saturday: 2 * oneDay,
		time.Sunday:   oneDay,
	}

	for tc, want := range expected {
		t.Run(tc.String(), func(t *testing.T) {
			got := slaDuration(tc)
			if got != want {
				t.Logf("got: %v want: %v", got, want)
				t.Fail()
			}
		})
	}

}
