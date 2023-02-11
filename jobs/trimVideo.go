package jobs

import (
	"errors"
	"time"
)

// TrimVideo accepts a video, timestamp and duration to trim the video around timestamp into a clip of length 'duration'
func TrimVideo(ts *time.Time) error {
	return errors.New("not implemented")
}
