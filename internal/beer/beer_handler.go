package beer

import (
	"net/http"
)

type BeerHandler struct {
	BelgianHandler *BelgianHandler
}

func (h *BeerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {
	case "latvian":
		http.Error(res, "Latvian not supported", http.StatusNotImplemented)
	case "belgian":
		h.BelgianHandler.ServeHTTP(res, req)
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}