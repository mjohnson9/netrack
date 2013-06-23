package app

import (
	"bytes"
	"html/template"
	"net/http"
)

// This is the data that should be available in all templates
type GlobalData struct {
	AppVersion string
	SiteName   string
	Flashes    []*Flash
}

func SafeExecuteTemplate(w http.ResponseWriter, t *template.Template, data interface{}) error {
	ourBuffer := new(bytes.Buffer)

	if err := t.Execute(ourBuffer, data); err != nil {
		return err
	}

	ourBuffer.WriteTo(w)
	return nil
}
