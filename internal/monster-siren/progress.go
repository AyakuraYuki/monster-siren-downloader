package monster_siren

import (
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

var _ io.Writer = (*trackerWrapper)(nil)

type trackerWrapper struct {
	dst      io.Writer
	filename string
	tracker  *progress.Tracker
}

func (wrap *trackerWrapper) Write(p []byte) (n int, err error) {
	n, err = wrap.dst.Write(p)
	if n > 0 {
		wrap.tracker.Increment(int64(n))
	}
	return n, err
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

func customizedProgress() progress.Writer {
	writer := progress.NewWriter()

	writer.SetAutoStop(false)
	writer.SetMessageLength(60)
	writer.SetUpdateFrequency(100 * time.Millisecond)
	writer.SetTrackerLength(30)
	writer.SetTrackerPosition(progress.PositionLeft)

	writer.SetStyle(progress.StyleDefault)

	writer.Style().Chars = progress.StyleChars{
		BoxLeft:       "|",
		BoxRight:      "|",
		Finished:      "⣿",
		Finished25:    "⣀",
		Finished50:    "⣤",
		Finished75:    "⣶",
		Indeterminate: progress.IndeterminateIndicatorMovingBackAndForth("⣤⣿⣤", progress.DefaultUpdateFrequency/2),
		Unfinished:    "⠀",
	}

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
	writer.Style().Options.ErrorString = "错误！"
	writer.Style().Options.PercentFormat = "%5.1f%%"
	writer.Style().Options.Separator = " - "
	writer.Style().Options.SnipIndicator = "..."

	writer.Style().Visibility.ETA = false
	writer.Style().Visibility.Time = false
	writer.Style().Visibility.Speed = false

	return writer
}
