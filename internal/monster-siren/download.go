package monster_siren

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ayfile "github.com/AyakuraYuki/go-aybox/files"
	msrModel "github.com/AyakuraYuki/monster-siren-api-go/model"
	"github.com/jedib0t/go-pretty/v6/progress"

	"github.com/AyakuraYuki/monster-siren-downloader/internal/filenames"
)

func (m *MonsterSiren) Run() (err error) {
	defer m.progress.Stop()
	defer m.pool.Release()
	go m.progress.Render()

	ctx := context.Background()

	pwd, err := os.Getwd()
	if err != nil {
		m.progress.Log("获取目录失败: %v", err)
		return err
	}

	firstPath := filepath.Join(pwd, saveTo)
	_ = os.MkdirAll(firstPath, os.ModePerm)

	albums, _ := m.client.Albums(ctx)
	tracker := m.newTracker(fmt.Sprintf("下载塞壬唱片曲库，专辑数：%d", len(albums)), int64(len(albums)), progress.UnitsDefault)
	tracker.Start()

	for index, album := range albums {
		artistes := make([]string, len(album.Artistes))
		copy(artistes, album.Artistes)

		album, _ = m.client.AlbumDetail(ctx, album.Cid)
		if !album.Exists() {
			m.progress.Log("查不到专辑详情: [%s] %s", album.Cid, album.Name)
			tracker.Increment(1)
			continue
		}

		for i, song := range album.Songs {
			// fill song's detail
			album.Songs[i], _ = m.client.Song(ctx, song.Cid)
		}

		if len(album.Artistes) == 0 {
			// write back artistes
			album.Artistes = make([]string, len(artistes))
			copy(album.Artistes, artistes)
		}

		m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载中的专辑：《%s》", album.Name))

		albumSerial := len(albums) - index
		secondPath := filepath.Join(firstPath, fmt.Sprintf("%03d - %s", albumSerial, filenames.AlbumName(album.Name)))
		_ = os.MkdirAll(secondPath, os.ModePerm)

		// save album info
		if infoPath := filepath.Join(secondPath, "info.txt"); !ayfile.PathExist(infoPath) {
			m.saveAlbumInfo(album, infoPath)
		}

		songTracker := m.newTracker(fmt.Sprintf("下载专辑：《%s》（曲数：%d）", album.Name, len(album.Songs)), int64(len(album.Songs)), progress.UnitsDefault)
		songTracker.Start()
		var wg sync.WaitGroup
		for i, song := range album.Songs {
			trackNo := i + 1
			wg.Add(1)
			_ = m.pool.Submit(m.newDownloadTask(song, trackNo, secondPath, songTracker, &wg))
		}
		wg.Wait()
		songTracker.MarkAsDone()

		if album.CoverUrl != "" {
			ext := filepath.Ext(album.CoverUrl)
			m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载专辑封面：《%s》", album.Name))
			_ = m.download(album.CoverUrl, secondPath, fmt.Sprintf("专辑封面%s", ext))
		}
		if album.CoverDeUrl != "" {
			ext := filepath.Ext(album.CoverDeUrl)
			m.progress.SetPinnedMessages(fmt.Sprintf(">>> 下载封面：《%s》", album.Name))
			_ = m.download(album.CoverDeUrl, secondPath, fmt.Sprintf("封面%s", ext))
		}

		m.progress.Log("✅  《%s》", album.Name)
		tracker.Increment(1)
	}

	time.Sleep(500 * time.Millisecond)
	tracker.MarkAsDone()
	return nil
}

func (m *MonsterSiren) newDownloadTask(song *msrModel.Song, trackNo int, path string, tracker *progress.Tracker, wg *sync.WaitGroup) func() {
	return func() {
		defer func() {
			tracker.Increment(1)
			wg.Done()
		}()

		if !song.Exists() {
			return
		}

		ext := filepath.Ext(song.SourceUrl)
		name := filenames.SongName(song.Name)
		songName := fmt.Sprintf("%02d.%s%s", trackNo, name, ext)
		lyricName := fmt.Sprintf("%02d.%s.lrc", trackNo, name)
		if song.SourceUrl != "" {
			_ = m.streamingDownload(song.SourceUrl, path, songName)
		}
		if song.LyricUrl != "" {
			_ = m.download(song.LyricUrl, path, lyricName)
		}
	}
}

func (m *MonsterSiren) download(link, dstDir, filename string) (err error) {
	dst := filepath.Join(dstDir, filename)
	if ayfile.PathExist(dst) {
		return nil // 跳过已下载的文件
	}

	tempDst := dst + ".tmp"
	_ = os.Remove(tempDst)

	response, err := m.downloader.R().SetOutput(tempDst).Get(link)
	if err != nil {
		m.progress.Log("下载失败 (%q)，错误: %v", link, err)
		return err
	}
	if response.IsError() {
		m.progress.Log("下载失败 (%q)，错误: %v", link, response.Error())
		return fmt.Errorf("download error: (code %d) %v", response.StatusCode(), response.Error())
	}

	_ = os.Rename(tempDst, dst)
	return nil
}

func (m *MonsterSiren) streamingDownload(link, dstDir, filename string) (err error) {
	dst := filepath.Join(dstDir, filename)
	if ayfile.PathExist(dst) {
		return nil // 跳过已下载的文件
	}

	// 1. HEAD 请求获取文件大小
	head, err := m.downloader.R().Head(link)
	if err != nil {
		return fmt.Errorf("fetch file size error: %w", err)
	}
	total := head.RawResponse.ContentLength

	// 2. 新建 tracker
	tracker := m.newTracker(filename, total, progress.UnitsBytes)

	// 3. 下载文件（获取响应体）
	response, err := m.downloader.R().SetDoNotParseResponse(true).Get(link)
	if err != nil {
		tracker.MarkAsErrored()
		m.progress.Log("下载失败 (%q)，错误: %v", link, err)
		return err
	}
	defer func(body io.ReadCloser) { _ = body.Close() }(response.RawBody())

	// 3-e. 处理错误响应
	if response.IsError() {
		tracker.MarkAsErrored()
		m.progress.Log("下载失败 (%q)，错误: %v", link, response.Status())
		return err
	}

	// 3.2. 修正文件尺寸
	if total <= 0 {
		total = response.RawResponse.ContentLength
		if total > 0 {
			tracker.UpdateTotal(total)
		}
	}

	// 4. 创建临时文件
	tempDst := dst + ".tmp"
	_ = os.Remove(tempDst)
	out, err := os.Create(tempDst)
	if err != nil {
		return err
	}
	defer func(f *os.File) { _ = f.Close() }(out)

	// 5. 创建下载进度跟踪器
	wrap := &trackerWrapper{
		dst:      out,
		filename: filename,
		tracker:  tracker,
	}

	// 6. 流式写入
	if _, err = io.Copy(wrap, response.RawBody()); err != nil {
		tracker.MarkAsErrored()
		m.progress.Log("下载失败 (%q)，错误: %v", link, err)
		return err
	}

	tracker.MarkAsDone()
	time.Sleep(300 * time.Millisecond)

	_ = os.Rename(tempDst, dst)
	return nil
}

func (m *MonsterSiren) saveAlbumInfo(album *msrModel.Album, infoPath string) {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, `专辑名称：%s
专辑属于：%s
专辑作者：%s
专辑介绍：
%s

歌曲列表：
`, album.Name, album.Belong, strings.Join(album.Artistes, "、"), album.Intro)

	for i, song := range album.Songs {
		if !song.Exists() {
			_, _ = fmt.Fprintf(&buf, "- %02d. <unknown song>\n", i+1)
			continue
		}

		_, _ = fmt.Fprintf(&buf, "- %02d. %s\n", i+1, song.Name)

		if len(song.Artists) > 0 {
			_, _ = fmt.Fprintf(&buf, "  作者：%s\n", strings.Join(song.Artists, "、"))
		} else {
			_, _ = fmt.Fprintf(&buf, "  作者：%s\n", strings.Join(song.Artistes, "、"))
		}
	}

	fh, err := os.Create(infoPath)
	if err != nil {
		return
	}
	defer func(fh *os.File) { _ = fh.Close() }(fh)
	fBuf := bufio.NewWriter(fh)
	_, _ = fmt.Fprintln(fBuf, strings.TrimSpace(buf.String()))
	_ = fBuf.Flush()
}
