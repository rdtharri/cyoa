package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(DefaultHandlerTmpl))
}

var tpl *template.Template

var DefaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
      {{range .Options}}
        <li><a href="/{{.Chapter}}" >{{.Text}}</a></li>
      {{end}}
    </ul>
  </body>
</html>
`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			panic(err)
		}
	}

}

func JsonStory(file io.Reader) (Story, error) {
	decoder := json.NewDecoder(file)
	var story Story
	if err := decoder.Decode(&story); err != nil {
		panic(err)
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
