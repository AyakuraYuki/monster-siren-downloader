package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Songs(t *testing.T) {
	client := New()
	songs, autoplay := client.Songs()
	assert.NotEmpty(t, songs)
	assert.NotEmpty(t, autoplay)
}

func TestClient_Song(t *testing.T) {
	tests := []string{
		"048794", // Warm and Small Light
		"779487", // Warm and Small Light (Instrument)
	}
	client := New()
	for _, tt := range tests {
		song := client.Song(tt)
		assert.True(t, song.IsExist())
	}
}
