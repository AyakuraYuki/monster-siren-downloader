package monster_siren

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/panjf2000/ants/v2"
)

const baseURL = `https://monster-siren.hypergryph.com`

var version string

type MonsterSiren struct {
	client   *resty.Client
	pool     *ants.Pool
	progress progress.Writer
}

func New(versions ...string) *MonsterSiren {
	if len(versions) > 0 && versions[0] != "" {
		version = versions[0]
	} else {
		version = "build-" + strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "_")
	}

	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(20 * time.Minute).
		SetRetryCount(3).
		SetHeaders(map[string]string{
			"Accept":          "*/*",
			"Accept-Language": "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,en-GB;q=0.6,en-US;q=0.5",
			"Referer":         "https://monster-siren.hypergryph.com/",
			"User-Agent":      fmt.Sprintf("Go/%s monster-siren-downloader/%s", runtime.Version(), version),
		})

	pw := progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetMessageLength(120)
	pw.SetStyle(progress.StyleBlocks)
	pw.SetUpdateFrequency(100 * time.Millisecond)
	pw.ShowTime(false)
	pw.Style().Colors = progress.StyleColors{
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
	pw.Style().Options.DoneString = "下载完毕！"

	instance := &MonsterSiren{
		client:   client,
		progress: pw,
	}

	p, err := ants.NewPool(5, ants.WithPanicHandler(instance.antsPanicHandler))
	if err != nil {
		panic(err)
	}
	instance.pool = p

	return instance
}

func (m *MonsterSiren) newTracker(message string, total int64) (tracker *progress.Tracker) {
	tracker = &progress.Tracker{
		Message:            message,
		RemoveOnCompletion: true,
		Total:              total,
		Units:              progress.UnitsDefault,
	}
	m.progress.AppendTracker(tracker)
	return tracker
}

func (m *MonsterSiren) antsPanicHandler(_ interface{}) {
	buf := make([]byte, 4<<10) // 4K
	buf = buf[:runtime.Stack(buf, false)]
	m.progress.Log("panic: %s", string(buf))
}

// ----------------------------------------------------------------------------------------------------

func (m *MonsterSiren) Songs() (songs []*Song, autoplay string) {
	songs = make([]*Song, 0)
	autoplay = ""

	rsp := &SongsRsp{}

	reply, err := m.client.R().SetResult(rsp).Get(`/api/songs`)
	if err != nil {
		log.Printf("failed to get song list: %v", err)
		return
	}
	if reply.IsError() {
		log.Printf("failed to get song list: %v", reply.Error())
		return
	}

	if rsp.Data == nil {
		log.Printf("no song list found, got this: %s", reply.String())
		return
	}

	songs = rsp.Data.List
	autoplay = rsp.Data.Autoplay
	return
}

func (m *MonsterSiren) Song(songID string) *Song {
	rsp := &SongRsp{}
	path := fmt.Sprintf(`/api/song/%s`, songID)

	reply, err := m.client.R().SetResult(rsp).Get(path)
	if err != nil {
		log.Printf("failed to get song: %v", err)
		return nil
	}
	if reply.IsError() {
		log.Printf("failed to get song: %v", reply.Error())
		return nil
	}

	if !rsp.Data.IsExist() {
		log.Printf("song not found, got this: %s", reply.String())
		return nil
	}

	return rsp.Data
}

func (m *MonsterSiren) Albums() (albums []*Album) {
	albums = make([]*Album, 0)

	rsp := &AlbumsRsp{}

	reply, err := m.client.R().SetResult(rsp).Get(`/api/albums`)
	if err != nil {
		log.Printf("failed to get album list: %v", err)
		return
	}
	if reply.IsError() {
		log.Printf("failed to get album list: %v", reply.Error())
		return
	}

	if len(rsp.Data) > 0 {
		albums = rsp.Data
	}

	return albums
}

func (m *MonsterSiren) Album(albumID string) *Album {
	rsp := &AlbumRsp{}
	path := fmt.Sprintf(`/api/album/%s/detail`, albumID)

	reply, err := m.client.R().SetResult(rsp).Get(path)
	if err != nil {
		log.Printf("failed to get album: %v", err)
		return nil
	}
	if reply.IsError() {
		log.Printf("failed to get album: %v", reply.Error())
		return nil
	}

	if !rsp.Data.IsExist() {
		log.Printf("album not found, got this: %s", reply.String())
		return nil
	}

	return rsp.Data
}
