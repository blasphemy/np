package getlasttrack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	api_base_url = "https://ws.audioscrobbler.com/2.0/"
)

type RecentTracks struct {
	Recenttracks struct {
		Track []struct {
			Artist struct {
				Text string `json:"#text"`
				Mbid string `json:"mbid"`
			} `json:"artist"`
			Name       string `json:"name"`
			Streamable string `json:"streamable"`
			Mbid       string `json:"mbid"`
			Album      struct {
				Text string `json:"#text"`
				Mbid string `json:"mbid"`
			} `json:"album"`
			URL   string `json:"url"`
			Image []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
			Attr struct {
				Nowplaying string `json:"nowplaying"`
			} `json:"@attr,omitempty"`
			Date struct {
				Uts  string `json:"uts"`
				Text string `json:"#text"`
			} `json:"date,omitempty"`
		} `json:"track"`
		Attr struct {
			User       string `json:"user"`
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			TotalPages string `json:"totalPages"`
			Total      string `json:"total"`
		} `json:"@attr"`
	} `json:"recenttracks"`
}

type TrackInfo struct {
	Track struct {
		Name       string `json:"name"`
		URL        string `json:"url"`
		Duration   string `json:"duration"`
		Streamable struct {
			Text      string `json:"#text"`
			Fulltrack string `json:"fulltrack"`
		} `json:"streamable"`
		Listeners string `json:"listeners"`
		Playcount string `json:"playcount"`
		Artist    struct {
			Name string `json:"name"`
			Mbid string `json:"mbid"`
			URL  string `json:"url"`
		} `json:"artist"`
		Album struct {
			Artist string `json:"artist"`
			Title  string `json:"title"`
			URL    string `json:"url"`
			Image  []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
		} `json:"album"`
		Userplaycount string `json:"userplaycount"`
		Userloved     string `json:"userloved"`
		Toptags       struct {
			Tag []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"tag"`
		} `json:"toptags"`
	} `json:"track"`
}

type Track struct {
	Album            string
	Artist           string
	Name             string
	PlayCount        int
	Tags             []string
	CurrentlyPlaying bool
}

func GetTrack(user string, ApiKey string) (Track, error) {
	resp, err := http.Get(fmt.Sprintf("%s?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=1", api_base_url, user, ApiKey))
	if err != nil {
		return Track{}, err
	}
	recentTrack := &RecentTracks{}
	err = json.NewDecoder(resp.Body).Decode(recentTrack)
	if err != nil {
		return Track{}, err
	}
	t := recentTrack.Recenttracks.Track[0]
	resp, err = http.Get(fmt.Sprintf("%s?method=track.getInfo&api_key=%s&artist=%s&track=%s&username=%s&format=json", api_base_url, ApiKey, t.Artist.Text, t.Name, user))
	if err != nil {
		return Track{}, err
	}
	ti := &TrackInfo{}
	err = json.NewDecoder(resp.Body).Decode(ti)
	if err != nil {
		return Track{}, err
	}
	var tags []string
	for _, x := range ti.Track.Toptags.Tag {
		tags = append(tags, x.Name)
	}
	playing := false
	if t.Attr.Nowplaying == "true" {
		playing = true
	}
	playcount, err := strconv.Atoi(ti.Track.Userplaycount)
	rt := Track{
		Album:            t.Album.Text,
		Artist:           t.Artist.Text,
		Name:             t.Name,
		Tags:             tags,
		CurrentlyPlaying: playing,
		PlayCount:        playcount,
	}
	return rt, nil
}
