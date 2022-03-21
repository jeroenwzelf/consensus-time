package server

import (
	"ConsensusTime/voting"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func newIndexHtmlTemplateExecutor() *template.Template {
	return template.Must(template.New("index.gohtml").Funcs(template.FuncMap{
		"GetConsensusDateDifferenceMillis": func() int64 {
			return int64(voting.GetConsensusTimeDifference().Milliseconds())
		},
	}).ParseFiles("html/index.gohtml"))
}

func AddWebPageRoutes(router *mux.Router) {
	// Root shows static html template `index.gohtml`
	templateExecutor := newIndexHtmlTemplateExecutor()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templateExecutor.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
