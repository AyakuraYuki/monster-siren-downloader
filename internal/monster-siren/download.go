package monster_siren

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ayfile "github.com/AyakuraYuki/go-aybox/file"
)

const saveTo = `./monster-siren`

func (m *MonsterSiren) downloadURL(mURL, path, name string) (err error) {
	absPath := filepath.Join(path, name)
	if ayfile.PathExist(absPath) {
		return nil // 跳过已下载的文件
	}

	reply, err := m.client.R().SetOutput(absPath).Get(mURL)
	if err != nil {
		log.Printf("failed to download url %q, err: %v", mURL, err)
		return err
	}
	if reply.IsError() {
		log.Printf("failed to download url %q, err: %v", mURL, reply.Error())
		return fmt.Errorf("reply error: (code %d) %v", reply.StatusCode(), reply.Error())
	}

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
		builder.WriteString(fmt.Sprintf("  作者：%s\n", strings.Join(song.Artists, "、")))
	}
	saveFile(infoPath, strings.TrimSpace(builder.String()))
}

func (m *MonsterSiren) DownloadTracks() (err error) {
	go m.progress.Render()
	defer m.progress.Stop()

	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("cannot get workdir: %v", err)
		return err
	}

	firstPath := filepath.Join(pwd, saveTo)
	_ = os.MkdirAll(firstPath, os.ModePerm)

	albums := m.Albums()
	albumTracker := m.newTracker(fmt.Sprintf("下载塞壬唱片曲库，专辑数：%d", len(albums)), int64(len(albums)))

	for _, album := range albums {
		album = m.Album(album.Cid)
		if !album.IsExist() {
			log.Printf("cannot get detail of album: [%s] %s", album.Cid, album.Name)
			continue
		}

		for index, song := range album.Songs {
			album.Songs[index] = m.Song(song.Cid)
		}

		m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载中的专辑：《%s》", album.Name))
		songTracker := m.newTracker(fmt.Sprintf("下载专辑：《%s》（曲数：%d）", album.Name, len(album.Songs)), int64(len(album.Songs)))

		secondPath := filepath.Join(firstPath, album.FilenamifyName())
		_ = os.MkdirAll(secondPath, os.ModePerm)

		infoPath := filepath.Join(secondPath, "info.txt")
		if !ayfile.PathExist(infoPath) {
			m.saveInfoFile(album, infoPath) // save album info
		}

		for index, song := range album.Songs {
			if !song.IsExist() {
				songTracker.Increment(1)
				continue
			}

			ext := filepath.Ext(song.SourceUrl)
			songName := fmt.Sprintf("%02d.%s%s", index+1, song.FilenamifyName(), ext)
			lyricName := fmt.Sprintf("%02d.%s.lrc", index+1, song.FilenamifyName())
			if song.SourceUrl != "" {
				_ = m.downloadURL(song.SourceUrl, secondPath, songName)
			}
			if song.LyricUrl != "" {
				_ = m.downloadURL(song.LyricUrl, secondPath, lyricName)
			}

			songTracker.Increment(1)
		}

		if album.CoverUrl != "" {
			ext := filepath.Ext(album.CoverUrl)
			_ = m.downloadURL(album.CoverUrl, secondPath, fmt.Sprintf("专辑封面%s", ext))
		}
		if album.CoverDeUrl != "" {
			ext := filepath.Ext(album.CoverDeUrl)
			_ = m.downloadURL(album.CoverDeUrl, secondPath, fmt.Sprintf("封面%s", ext))
		}

		albumTracker.Increment(1)
		m.progress.Log("✅  《%s》", album.Name)
	}

	return nil
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
