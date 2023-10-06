package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.sr.ht/~kota/lists/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/", app.create)
	router.HandlerFunc(http.MethodGet, "/:name", app.view)
	router.HandlerFunc(http.MethodPost, "/:name", app.edit)
	return app.recoverPanic(app.logRequest(app.secureHeaders(router)))
}

type homePage struct {
	CSPNonce string
}

func (app *application) render(
	w http.ResponseWriter,
	status int,
	page string,
	data interface{},
) {
	ts, ok := app.templates[page]
	if !ok {
		app.serverError(w, fmt.Errorf(
			"the template %s is missing",
			page,
		))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

// home handles displaying the home page.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.tmpl", homePage{
		CSPNonce: app.cspNonce,
	})
}

// create handles creating a list.
func (app *application) create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	if name == "" {
		app.clientError(w, http.StatusBadRequest)
	}
	name = alphaNumeric(strings.ToLower(name))

	// Check if the list already exists.
	if _, err := app.lists.Get(name); !errors.Is(err, models.ErrNoRecord) {
		http.Redirect(w, r, fmt.Sprintf("/%s", name), http.StatusSeeOther)
		return
	}

	err = app.lists.Insert(name, "")
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", name), http.StatusSeeOther)
}

type viewPage struct {
	CSPNonce string

	Name string
	Body string
}

// view handles displaying a list page.
func (app *application) view(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	list, err := app.lists.Get(name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
		}
	}

	app.render(w, http.StatusOK, "list.tmpl", viewPage{
		CSPNonce: app.cspNonce,
		Name:     list.Name,
		Body:     list.Body,
	})
}

// edit handles editing a list.
func (app *application) edit(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	_, err := app.lists.Get(name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
		}
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	data := r.PostForm.Get("data")
	if data == "" {
		app.clientError(w, http.StatusBadRequest)
	}

	err = app.lists.Update(name, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", name), http.StatusSeeOther)
}
