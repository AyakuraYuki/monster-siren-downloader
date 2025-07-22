package monster_siren

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ayfile "github.com/AyakuraYuki/go-aybox/file"
	"github.com/jedib0t/go-pretty/v6/progress"
)

const saveTo = `./monster-siren`

func (m *MonsterSiren) downloadURL(mURL, path, name string) (err error) {
	absPath := filepath.Join(path, name)
	if ayfile.PathExist(absPath) {
		return nil // 跳过已下载的文件
	}

	tmpPath := absPath + ".tmp"
	_ = os.Remove(tmpPath)
	reply, err := m.client.R().SetOutput(tmpPath).Get(mURL)
	if err != nil {
		m.progress.Log("failed to download url %q, err: %v", mURL, err)
		return err
	}
	if reply.IsError() {
		m.progress.Log("failed to download url %q, err: %v", mURL, reply.Error())
		return fmt.Errorf("reply error: (code %d) %v", reply.StatusCode(), reply.Error())
	}
	_ = os.Rename(tmpPath, absPath)

	return nil
}

func (m *MonsterSiren) saveInfoFile(album *Album, infoPath string) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("专辑名称：%s\n", album.Name))
	builder.WriteString(fmt.Sprintf("专辑属于：%s\n", album.Belong))
	builder.WriteString(fmt.Sprintf("专辑作者：%s\n", strings.Join(album.Artistes, "、")))
	builder.WriteString(fmt.Sprintf("专辑介绍：\n%s\n\n", album.Intro))
	builder.WriteString("歌曲列表：\n")
	for index, song := range album.Songs {
		if !song.IsExist() {
			builder.WriteString(fmt.Sprintf("- %02d. %s\n", index+1, "<unknown: missing data>"))
			continue
		}
		builder.WriteString(fmt.Sprintf("- %02d. %s\n", index+1, song.Name))
		if len(song.Artists) > 0 {
			builder.WriteString(fmt.Sprintf("  作者：%s\n", strings.Join(song.Artists, "、")))
		} else {
			builder.WriteString(fmt.Sprintf("  作者：%s\n", strings.Join(song.Artistes, "、")))
		}
	}
	saveFile(infoPath, strings.TrimSpace(builder.String()))
}

func (m *MonsterSiren) DownloadTracks() (err error) {
	defer m.progress.Stop()
	defer m.pool.Release()
	go m.progress.Render()

	pwd, err := os.Getwd()
	if err != nil {
		m.progress.Log("cannot get workdir: %v", err)
		return err
	}

	firstPath := filepath.Join(pwd, saveTo)
	_ = os.MkdirAll(firstPath, os.ModePerm)

	albums := m.Albums()
	tracker := m.newTracker(fmt.Sprintf("下载塞壬唱片曲库，专辑数：%d", len(albums)), int64(len(albums)))
	tracker.Start()

	for albumIndex, album := range albums {
		backupArtistes := make([]string, len(album.Artistes))
		copy(backupArtistes, album.Artistes)

		album = m.AlbumWithSongs(album.Cid)
		if !album.IsExist() {
			m.progress.Log("cannot get detail of album: [%s] %s", album.Cid, album.Name)
			tracker.Increment(1)
			continue
		}

		for index, song := range album.Songs {
			album.Songs[index] = m.Song(song.Cid)
		}

		if len(album.Artistes) == 0 {
			album.Artistes = make([]string, len(backupArtistes))
			copy(album.Artistes, backupArtistes)
		}

		m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载中的专辑：《%s》", album.Name))
		songTracker := m.newTracker(fmt.Sprintf("下载专辑：《%s》（曲数：%d）", album.Name, len(album.Songs)), int64(len(album.Songs)))
		songTracker.Start()

		albumNo := len(albums) - albumIndex
		secondPath := filepath.Join(firstPath, fmt.Sprintf("%03d - %s", albumNo, album.FilenamifyName()))
		_ = os.MkdirAll(secondPath, os.ModePerm)

		infoPath := filepath.Join(secondPath, "info.txt")
		if !ayfile.PathExist(infoPath) {
			m.saveInfoFile(album, infoPath) // save album info
		}

		var wg sync.WaitGroup
		for index, song := range album.Songs {
			trackNo := index + 1
			song := song
			wg.Add(1)
			_ = m.pool.Submit(m.downloadSongsTaskWrapper(song, trackNo, secondPath, songTracker, &wg))
		}
		wg.Wait()
		songTracker.MarkAsDone()

		if album.CoverUrl != "" {
			ext := filepath.Ext(album.CoverUrl)
			m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载专辑封面：《%s》", album.Name))
			_ = m.downloadURL(album.CoverUrl, secondPath, fmt.Sprintf("专辑封面%s", ext))
		}
		if album.CoverDeUrl != "" {
			ext := filepath.Ext(album.CoverDeUrl)
			m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载封面：《%s》", album.Name))
			_ = m.downloadURL(album.CoverDeUrl, secondPath, fmt.Sprintf("封面%s", ext))
		}

		m.progress.Log("✅  《%s》", album.Name)
		tracker.Increment(1)
	}

	time.Sleep(500 * time.Millisecond)
	tracker.MarkAsDone()
	return nil
}

func (m *MonsterSiren) downloadSongsTaskWrapper(song *Song, trackNo int, path string, tracker *progress.Tracker, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()

		if !song.IsExist() {
			tracker.Increment(1)
			return
		}

		ext := filepath.Ext(song.SourceUrl)
		name := song.FilenamifyName()
		songName := fmt.Sprintf("%02d.%s%s", trackNo, name, ext)
		lyricName := fmt.Sprintf("%02d.%s.lrc", trackNo, name)
		if song.SourceUrl != "" {
			_ = m.downloadURL(song.SourceUrl, path, songName)
		}
		if song.LyricUrl != "" {
			_ = m.downloadURL(song.LyricUrl, path, lyricName)
		}

		tracker.Increment(1)
	}
}

func saveFile(path, text string) {
	fh, err := os.Create(path)
	if err != nil {
		return
	}
	defer func(fh *os.File) { _ = fh.Close() }(fh)
	buf := bufio.NewWriter(fh)
	_, _ = fmt.Fprintln(buf, text)
	_ = buf.Flush()
}

func saveFileAppend(path, text string) {
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer func(fh *os.File) { _ = fh.Close() }(fh)
	buf := bufio.NewWriter(fh)
	_, _ = fmt.Fprintln(buf, text)
	_ = buf.Flush()
}
