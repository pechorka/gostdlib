package httpx

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pechorka/gostdlib/pkg/errs"
)

func ReadJSON[T any](body io.Reader) (T, error) {
	var t T
	if body == nil {
		return t, errs.New("body is nil")
	}

	err := json.NewDecoder(body).Decode(&t)
	if err != nil {
		return t, errs.Wrap(err, "failed to parse JSON")
	}

	return t, nil
}

func WriteJSON[T any](w http.ResponseWriter, v T) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return errs.Wrap(err, "failed to write JSON")
	}
	return nil
}
