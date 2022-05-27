package https

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type KTP struct {
	*Server
}

func (k *KTP) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Name string `json:"name"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			k.render.JSON(w, http.StatusBadRequest, Map{
				"error": err.Error(),
			})
			return
		}

		log.Error().Str("name", req.Name).Msg("")

		if err := k.NameChecker.Check(req.Name); err != nil {
			k.render.JSON(w, http.StatusBadRequest, Map{
				"error": err.Error(),
			})
			return
		}
		k.render.JSON(w, http.StatusOK, Map{
			"name": req.Name,
		})
	}
}
