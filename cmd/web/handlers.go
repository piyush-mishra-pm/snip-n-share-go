package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/piyush-mishra-pm/snip-n-share-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snips, err := app.snips.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snips = snips

	app.render(w, http.StatusOK, "home.tmpl.htm", data)
}

func (app *application) snipView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snip, err := app.snips.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snip = snip

	app.render(w, http.StatusOK, "view.tmpl.htm", data)
}

func (app *application) snipCreatePost(w http.ResponseWriter, r *http.Request) {
	// Dummy Snip item:
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snips.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)
}

func (app *application) snipCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display form to create new snip..."))
}
