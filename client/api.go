package client

import (
	"fmt"
	"log"
)

type SongsRsp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data *SongsRspBody `json:"data"`
}

type SongsRspBody struct {
	List     []*Song `json:"list"`
	Autoplay string  `json:"autoplay"`
}

// Songs list released songs from Monster Siren Records
func (c *Client) Songs() (songs []*Song, autoplay string) {
	songs = make([]*Song, 0)

	rsp := &SongsRsp{}
	apiRsp, err := c.cli.R().SetResult(rsp).Get(`/api/songs`)
	if err != nil {
		log.Printf("failed to get song list: %v", err)
		return
	}
	if apiRsp.IsError() {
		log.Printf("failed to get song list: %v", apiRsp.Error())
		return
	}

	if rsp.Data == nil {
		log.Printf("no data reveived but got this: %s", apiRsp.String())
		return
	}

	songs = rsp.Data.List
	autoplay = rsp.Data.Autoplay
	return
}

type SongRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Song  `json:"data"`
}

// Song gets a single song detail includes song's source link and lyric's link
func (c *Client) Song(songID string) (song *Song) {
	rsp := &SongRsp{}
	uri := fmt.Sprintf(`/api/song/%s`, songID)
	apiRsp, err := c.cli.R().SetResult(rsp).Get(uri)
	if err != nil {
		log.Printf("failed to get song: %v", err)
		return nil
	}
	if apiRsp.IsError() {
		log.Printf("failed to get song: %v", apiRsp.Error())
		return nil
	}

	if rsp.Data == nil {
		log.Printf("no data reveived but got this: %s", apiRsp.String())
		return nil
	}

	song = rsp.Data
	return
}

type AlbumRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Album `json:"data"`
}

// Album gets an album detail includes a playlist
func (c *Client) Album(albumID string) (album *Album) {
	rsp := &AlbumRsp{}
	uri := fmt.Sprintf(`/api/album/%s/detail`, albumID)
	apiRsp, err := c.cli.R().SetResult(rsp).Get(uri)
	if err != nil {
		log.Printf("failed to get album: %v", err)
		return nil
	}
	if apiRsp.IsError() {
		log.Printf("failed to get album: %v", apiRsp.Error())
		return nil
	}

	if rsp.Data == nil {
		log.Printf("no data reveived but got this: %s", apiRsp.String())
		return nil
	}

	album = rsp.Data
	return
}
