package monster_siren

import (
	"testing"

	mjson "github.com/AyakuraYuki/monster-siren-downloader/internal/json"
)

func TestMonsterSiren_Songs(t *testing.T) {
	m := New("test_cases")
	songs, autoplay := m.Songs()
	if len(songs) == 0 {
		t.Fatal("API error")
	}
	t.Logf("autoplay: %s", autoplay)
	t.Logf("songs: %s\n...", mjson.Prettify(songs)[:1024])
}

func TestMonsterSiren_Song(t *testing.T) {
	m := New("test_cases")
	cid := "697699" // Grow on My Time
	song := m.Song(cid)
	if !song.IsExist() {
		t.Fatal("API error")
	}
	t.Log(mjson.Prettify(song))
}

func TestMonsterSiren_Albums(t *testing.T) {
	m := New("test_cases")
	t.Logf("albums: %s\n...", mjson.Prettify(m.Albums())[:1024])
}

func TestMonsterSiren_Album(t *testing.T) {
	m := New("test_cases")
	cid := "1010" // Grow on My Time
	album := m.Album(cid)
	if !album.IsExist() {
		t.Fatal("API error")
	}
	t.Log(mjson.Prettify(album))
}

func TestMonsterSiren_AlbumWithSongs(t *testing.T) {
	m := New("test_cases")
	cid := "6660" // Warm and Small Light
	album := m.AlbumWithSongs(cid)
	if !album.IsExist() {
		t.Fatal("API error")
	}
	t.Log(mjson.Prettify(album))
}
