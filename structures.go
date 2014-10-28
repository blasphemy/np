package main

type Loader2 struct {
	Recent RT2 `json:"recenttracks"`
}
type RT2 struct {
	TRay []Track `json:"track"`
}

type Loader struct {
	Recent  RT    `json:"recenttracks"`
	Track_v Track `json:"track"`
}

type RT struct {
	Track_v Track `json:"track"`
}

type Track struct {
	Name_v   string `json:"name"`
	Artist_v Artist `json:"artist"`
	Album_v  Album  `json:"album"`
}

type Album struct {
	Name_v string `json:"title"`
}
type Artist struct {
	Name_v string `json:"#text"`
	Name_2 string `json:"name"`
}

type RealTrack struct {
	Name      string
	Artist    string
	Album     string
	PlayCount int
	tags      []string
}
