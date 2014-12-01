package htpasswd

import (
	"flag"
	"log"
	"net/http"

	"github.com/Tuxuri/pintu"
)

type (
	HtpasswdProvider struct {
		name          string
		start         string
		settings      *settings
		cookieFactory *pintu.CookieFactory
		ptype         string
		htpasswdfile  *HtpasswdFile
	}

	settings struct {
		path string
	}
)

const (
	optionPath = "htpasswd"
)

func NewHtpasswdProvider() *HtpasswdProvider {
	s := &settings{}
	flag.StringVar(&s.path, optionPath, "", "htpasswd file path")

	return &HtpasswdProvider{
		name:     "HTPasswd",
		start:    "/auth/htpasswd/start",
		settings: s,
		ptype:    "form",
	}
}

func (p *HtpasswdProvider) ParseSettings() {
	if p.settings.path == "" {
		pintu.EnvStringVar(&p.settings.path, optionPath, "")
		if p.settings.path == "" {
			log.Fatalf("missing param %s", optionPath)
		}
		var err error
		p.htpasswdfile, err = NewHtpasswdFromFile(p.settings.path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (p *HtpasswdProvider) RegisterCookie(f *pintu.CookieFactory) {
	p.cookieFactory = f
}

func (p *HtpasswdProvider) RegisterHandler(g *pintu.Guard) {
	g.HandleFunc(p.start, p.startHandler)
}

func (p *HtpasswdProvider) Partial(r *http.Request) string {
	hostURL := pintu.GetHostURL(r)
	action := hostURL + p.start
	partial := &pintu.LoginPartial{
		Action:   action,
		Redirect: r.URL.RequestURI(),
		Name:     p.name,
	}
	return partial.GetForm(r)
}

func (p *HtpasswdProvider) Type() string {
	return p.ptype
}

func (p *HtpasswdProvider) Name() string {
	return p.name
}

func (p *HtpasswdProvider) startHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	email := r.Form.Get("username")
	password := r.Form.Get("password")
	redirect := pintu.GetRedirect(r)

	if p.htpasswdfile.Validate(email, password) {
		p.cookieFactory.SetCookie(email, w, r)
		http.Redirect(w, r, redirect, 302)
		return
	}
	pintu.CustomError(w, r, pintu.ErrInvalidCredentials)
	return
}
