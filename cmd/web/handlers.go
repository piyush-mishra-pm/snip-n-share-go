package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/piyush-mishra-pm/snip-n-share-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
	}
	snips, err := app.snips.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, snip := range snips {
		fmt.Fprintf(w, "%+v\n\n", snip)
	}

	// files := []string{
	// 	"./ui/html/base.tmpl.htm", // base must be first in list.
	// 	"./ui/html/partials/nav.tmpl.htm",
	// 	"./ui/html/pages/home.tmpl.htm",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// err = ts.ExecuteTemplate(w, "base", snips)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

func (app *application) snipView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

	fmt.Fprintf(w, "%+v", snip)
}

func (app *application) snipCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy Snip item:
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snips.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip/view?id=%d", id), http.StatusSeeOther)
}
