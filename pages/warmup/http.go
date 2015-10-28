package pages

import (
	"io"
	"net/http"

	"appengine"

	"app"
)

func init() {
	app.Get("/_ah/warmup", app.CreateHandler(warmup))
}

func warmup(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "Warmup successful")

	c.Debugf("Successfully performed warmup")
}
