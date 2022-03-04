package render

import (
	"YP-metrics-and-alerting/internal/config"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var pathToTemplates = "../../internal/templates"

var app *config.Application

func NewRenderer(appConfig *config.Application) {
	app = appConfig
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	log.Println(fmt.Sprintf("%s/*.gohtml", pathToTemplates))
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	log.Println(pages)
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}
	return myCache, nil
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td interface{}) error {
	var tc map[string]*template.Template = app.TemplateCache

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}