package positions

import (
	"net/http"

	"github.com/tunakarakasoglu/finance-dashboard/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server.Positions(w, r)
}
