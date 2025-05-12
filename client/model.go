package client

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

type Album struct {
	Cid        string  `json:"cid"`
	Name       string  `json:"name"`
	Intro      string  `json:"intro"`
	Belong     string  `json:"belong"`
	CoverUrl   string  `json:"coverUrl"`
	CoverDeUrl string  `json:"coverDeUrl"`
	Songs      []*Song `json:"songs"`
}

func (album *Album) IsExist() bool { return album != nil && album.Cid != "" }
