package events

import (
	"testing"
	"time"
)

func TestNormalizeOptionalEndAndUTC(t *testing.T) {
	input := sampleInput()
	input.StartsAt = time.Date(2026, 11, 21, 17, 0, 0, 123, time.FixedZone("Warsaw", 3600))
	input.EndsAt = time.Time{}
	got, err := normalize(input)
	if err != nil {
		t.Fatal(err)
	}
	if got.StartsAt.Format(time.RFC3339Nano) != "2026-11-21T16:00:00Z" || !got.EndsAt.Equal(got.StartsAt) {
		t.Fatalf("optional end / UTC normalization: %+v", got)
	}
}
