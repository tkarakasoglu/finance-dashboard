package healthz

import (
	"net/http"

	"github.com/tunakarakasoglu/finance-dashboard/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server.Healthz(w, r)
}

