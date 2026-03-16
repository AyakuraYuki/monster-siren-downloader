package monster_siren

import (
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

var _ io.Writer = (*trackerWrapper)(nil)

type trackerWrapper struct {
	dst     io.Writer
	tracker *progress.Tracker
}

func (wrap *trackerWrapper) Write(p []byte) (n int, err error) {
	n, err = wrap.dst.Write(p)
	if n > 0 {
		wrap.tracker.Increment(int64(n))
	}
	return n, err
}

func customizedProgress() progress.Writer {
	writer := progress.NewWriter()

	writer.SetAutoStop(false)
	writer.SetMessageLength(64)
	writer.SetStyle(progress.StyleBlocks)
	writer.SetUpdateFrequency(100 * time.Millisecond)
	writer.SetTrackerLength(30)
	writer.SetTrackerPosition(progress.PositionLeft)

	writer.Style().Colors = progress.StyleColors{
		Message: text.Colors{text.FgWhite},
		Error:   text.Colors{text.FgRed},
		Percent: text.Colors{text.FgHiGreen},
		Pinned:  text.Colors{text.BgHiBlack, text.FgHiWhite},
		Stats:   text.Colors{text.FgHiBlack},
		Time:    text.Colors{text.FgGreen},
		Tracker: text.Colors{text.FgCyan},
		Value:   text.Colors{text.FgCyan},
		Speed:   text.Colors{text.FgMagenta},
	}
	writer.Style().Options.DoneString = "下载完毕！"
	writer.Style().Options.PercentFormat = "%5.1f%%"
	writer.Style().Visibility.Time = false
	writer.Style().Visibility.Speed = false
	writer.Style().Visibility.ETA = false

	return writer
}

func (m *MonsterSiren) newTracker(message string, total int64, unit progress.Units) (tracker *progress.Tracker) {
	tracker = &progress.Tracker{
		Message:            message,
		RemoveOnCompletion: true,
		Total:              total,
		Units:              unit,
	}
	m.progress.AppendTracker(tracker)
	return tracker
}
