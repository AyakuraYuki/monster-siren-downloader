package monster_siren

import (
	"strings"

	"github.com/flytam/filenamify"

	"github.com/AyakuraYuki/monster-siren-downloader/internal/str"
)

type Song struct {
	Cid        string   `json:"cid"`
	Name       string   `json:"name"`
	AlbumCid   string   `json:"albumCid"`
	SourceUrl  string   `json:"sourceUrl,omitempty"`
	LyricUrl   string   `json:"lyricUrl,omitempty"`
	MvUrl      string   `json:"mvUrl,omitempty"`
	MvCoverUrl string   `json:"mvCoverUrl,omitempty"`
	Artists    []string `json:"artists"`
}

func (song *Song) IsExist() bool { return song != nil && song.Cid != "" }

func (song *Song) FilenamifyName() string {
	if !song.IsExist() {
		return ""
	}
	name, err := filenamify.Filenamify(song.Name, filenamify.Options{Replacement: "_"})
	if err != nil {
		return strings.TrimSpace(song.Name)
	}
	return strings.TrimSpace(name)
}

type Album struct {
	Cid        string   `json:"cid"`
	Name       string   `json:"name"`
	Intro      string   `json:"intro,omitempty"`
	Belong     string   `json:"belong,omitempty"`
	CoverUrl   string   `json:"coverUrl"`
	CoverDeUrl string   `json:"coverDeUrl,omitempty"`
	Artistes   []string `json:"artistes,omitempty"`
	Songs      []*Song  `json:"songs,omitempty"`
}

func (album *Album) IsExist() bool { return album != nil && album.Cid != "" }

func (album *Album) FilenamifyName() string {
	if !album.IsExist() {
		return ""
	}
	name, err := filenamify.Filenamify(album.Name, filenamify.Options{Replacement: "_"})
	if err != nil {
		return strings.TrimSpace(album.Name)
	}
	name = strings.TrimSpace(name)
	name = str.ReplaceDotSuffixRune(name)
	return strings.TrimSpace(name)
}

type SongsRsp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data *SongsRspBody `json:"data"`
}

type SongsRspBody struct {
	List     []*Song `json:"list"`
	Autoplay string  `json:"autoplay"`
}

type SongRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Song  `json:"data"`
}

type AlbumsRsp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []*Album `json:"data"`
}

type AlbumRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Album `json:"data"`
}
