package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type Api struct {
	mdb   *MongoDB
	cache *Cache
}

func Run(host string, port string, mdb *MongoDB, cache *Cache) error {
	a := Api{
		mdb:   mdb,
		cache: cache,
	}
	r := chi.NewRouter()
	a.registerUrls(r)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r)
}

func (a *Api) registerUrls(r *chi.Mux) {
	r.Get("/", a.info)

	r.Route("/entity", func(r chi.Router) {
		r.Get("/", a.EntityGetList)
		r.Put("/{name}", a.EntityPut)
		r.Delete("/{name}", a.EntityDelete)
	})
}

func (a *Api) info(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"Status": "OK",
	})
}

func (a *Api) EntityGetList(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"entities": a.cache.GetAll(),
	})
}

func (a *Api) EntityPut(w http.ResponseWriter, r *http.Request) {
	oldName := chi.URLParam(r, "name")

	var newEntity Entity
	if err := json.NewDecoder(r.Body).Decode(&newEntity); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": err,
		})
		return
	}

	if newEntity.Name == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"error": "Name must not be empty",
		})
		return
	}

	entity, err := a.mdb.UpdateOrCreate(oldName, newEntity.Name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"error": err,
		})
		return
	}

	render.JSON(w, r, entity)
}

func (a *Api) EntityDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	err := a.mdb.Delete(name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"error": err,
		})
		return
	}
	render.Status(r, http.StatusNoContent)
}
