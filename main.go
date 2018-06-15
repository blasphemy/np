package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/blasphemy/np/getlasttrack"
	"github.com/spf13/viper"
)

func RemoveSpaces(in string) string {
	return strings.Replace(in, " ", "", -1)

}

func main() {
	f := template.FuncMap{
		"RemoveSpaces": RemoveSpaces,
	}
	viper.AddConfigPath(".")
	viper.SetConfigName("np")
	viper.SetDefault("template", "{{if .CurrentlyPlaying}}is now playing{{else}}last played{{end}} {{if .Artist}}[{{.Artist}}] - {{end}}{{if .Name}}[{{.Name}}] {{end}}{{if .Album}}On [{{.Album}}] {{end}}{{if .PlayCount}}[{{.PlayCount}} plays] {{end}}{{if .Tags}}[{{range $index, $element := .Tags}}{{if $index}} {{end}}#{{RemoveSpaces $element}}{{end}}]{{end}}")
	viper.SetDefault("debug", false)
	viper.ReadInConfig()
	apikey := viper.GetString("apikey")
	username := viper.GetString("username")
	if apikey == "" {
		fmt.Println("no api key")
		os.Exit(1)
	}
	tmpl := template.Must(template.New("format").Funcs(f).Parse(viper.GetString("template")))
	track, err := getlasttrack.GetTrack(username, apikey)
	if err != nil {
		os.Exit(1)
	}
	err = tmpl.Execute(os.Stdout, track)
	if err != nil {
		panic(err)
	}
}
