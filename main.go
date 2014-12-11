package main

import (
	"fmt"
	"github.com/blasphemy/urel"
	"github.com/spf13/viper"
	"github.com/blasphemy/np/getlasttrack"
	"os"
	"strings"
	"text/template"
)

var (
	debug bool
)

func RemoveSpaces(in string) string {
	return strings.Replace(in, " ", "", -1)

}

func main() {
	f := template.FuncMap{
		"RemoveSpaces": RemoveSpaces,
		"short":        urel.Short,
	}
	viper.AddConfigPath("$HOME/np")
	viper.AddConfigPath("$HOME/.np")
	viper.AddConfigPath("$HOME")
	wd, err := os.Getwd()
	if err != nil {
		viper.AddConfigPath(wd)
	}
	viper.SetConfigName("np")
	viper.SetDefault("template", "{{if .CurrentlyPlaying}}is now playing{{else}}last played{{end}} {{if .Artist}}[{{.Artist}}] - {{end}}{{if .Name}}[{{.Name}}] {{end}}{{if .Album}}On [{{.Album}}] {{end}}{{if .PlayCount}}[{{.PlayCount}} plays] {{end}}{{if .Tags}}[{{range $index, $element := .Tags}}{{if $index}} {{end}}#{{RemoveSpaces $element}}{{end}}]{{end}}{{if .SpotifyUrl}} [{{short .SpotifyUrl}}]{{end}}")
	viper.SetDefault("debug", false)
	viper.SetDefault("username", "zxo0oxz")
	viper.ReadInConfig()
	apikey := viper.GetString("apikey")
	username := viper.GetString("username")
	debug = viper.GetBool("debug")
	if apikey == "" {
		fmt.Println("no api key")
		os.Exit(1)
	}
	tmpl := template.Must(template.New("format").Funcs(f).Parse(viper.GetString("template")))
	track, err := getlasttrack.GetTrack(username, apikey)
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
