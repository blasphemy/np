package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	username        = "zxo0oxz"
	apikey          string
	debug           = true
	defaultTemplate = "{{if .CurrentlyPlaying}}is now playing{{else}}last played{{end}} {{if .Artist}}[{{.Artist}}] - {{end}}{{if .Name}}[{{.Name}}] {{end}}{{if .Album}}On [{{.Album}}] {{end}}{{if .PlayCount}}[{{.PlayCount}} plays] {{end}}{{if .Tags}}[{{range $index, $element := .Tags}}{{if $index}} {{end}}#{{$element}}{{end}}]{{end}}"
)

func main() {
	tmpl, err := template.New("format").Parse(defaultTemplate)
	if err != nil {
		panic(err)
	}
	apitemp, apierror := ioutil.ReadFile("api.config")
	if apierror != nil {
		fmt.Println(apierror)
		os.Exit(1)
	}
	apikey = string(apitemp)
	apikey = strings.TrimSpace(apikey)

	track, err := GetTrack(username, apikey)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	err = tmpl.Execute(os.Stdout, track)
	if err != nil {
		panic(err)
	}
}
