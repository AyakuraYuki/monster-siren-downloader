package main

import monstersiren "github.com/AyakuraYuki/monster-siren-downloader/internal/monster-siren"

var version string

func main() {
	cli := monstersiren.New(version)
	_ = cli.Run()
}
