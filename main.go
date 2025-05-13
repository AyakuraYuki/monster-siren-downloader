package main

import monstersiren "github.com/AyakuraYuki/monster-siren-downloader/internal/monster-siren"

var version string

func main() {
	lib := monstersiren.New(version)
	_ = lib.DownloadTracks()
}
