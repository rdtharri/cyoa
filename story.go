package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	MapKeyError          = "Something went wrong..."
	ChapterNotFoundError = "Chapter not found."
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
    <title>Choose Your Own Adventure!</title>
  </head>
  <body>
		<section class="page">
	    <h1>{{.Title}}</h1>
	    {{range .Paragraphs}}
	    <p>{{.}}</p>
	    {{end}}
	    <ul>
	      {{range .Options}}
	        <li><a href="/{{.Chapter}}" >{{.Text}}</a></li>
	      {{end}}
	    </ul>
		</section>
		<style>
		      body {
		        font-family: helvetica, arial;
		      }
		      h1 {
		        text-align:center;
		        position:relative;
		      }
		      .page {
		        width: 80%;
		        max-width: 500px;
		        margin: auto;
		        margin-top: 40px;
		        margin-bottom: 40px;
		        padding: 80px;
		        background: #FCF6FC;
		        border: 1px solid #eee;
		        box-shadow: 0 10px 6px -6px #797;
		      }
		      ul {
		        border-top: 1px dotted #ccc;
		        padding: 10px 0 0 0;
		        -webkit-padding-start: 0;
		      }
		      li {
		        padding-top: 10px;
		      }
		      a,
		      a:visited {
		        text-decoration: underline;
		        color: #555;
		      }
		      a:active,
		      a:hover {
		        color: #222;
		      }
		      p {
		        text-indent: 1em;
		      }
		    </style>
  </body>
</html>
`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, MapKeyError, http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, ChapterNotFoundError, http.StatusNotFound)
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
