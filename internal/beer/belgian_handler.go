package beer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type BelgianHandler struct {
}

var beers = []string {"leffe", "chimay", "hoegaarden"}

func (b *BelgianHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/" {
		bytes, err := json.Marshal(beers)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Can't encode response"))
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		res.Write(bytes)
		return
	}

	head, _ := ShiftPath(req.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		b.handleGet(res, id)
	default:
		http.Error(res, "Only GET is allowed", http.StatusMethodNotAllowed)
	}
}

func (b *BelgianHandler) handleGet(res http.ResponseWriter, id int) {
	bytes, err := json.Marshal(beers[id])
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Can't encode response"))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}