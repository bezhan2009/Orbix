package handlers

import (
	"goCmd/src/Orbix"
	"net/http"
)

func HandlerWebOrbix(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	cmd := r.Form.Get("cmd")
	if cmd == "" {
		cmd = r.URL.Query().Get("cmd")
	}

	Orbix.ExecLtCommand(cmd)
	http.Redirect(w, r, "/", 301)
}
