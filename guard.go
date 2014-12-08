package pintu

import (
	"html/template"
	"log"
	"net/http"
)

type (
	Guard struct {
		mux           *http.ServeMux
		cookieFactory *CookieFactory
		template      *template.Template
		providers     map[string]Provider
	}
)

func NewGuard() *Guard {
	mux := http.NewServeMux()
	t := GetTemplates()
	guard := &Guard{
		mux:       mux,
		template:  t,
		providers: make(map[string]Provider),
	}
	mux.HandleFunc(loginPromptPath, guard.LoginPrompt)
	return guard
}

func (g *Guard) Use(providers ...Provider) {
	for _, p := range providers {
		p.RegisterCookie(g.cookieFactory)
		p.RegisterHandler(g)
		p.ParseSettings()
		g.providers[p.Name()] = p
	}
}

func (g *Guard) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	g.mux.HandleFunc(pattern, handler)
}

func (g *Guard) CheckCookie(r *http.Request) (email string, ok bool) {
	cookie, err := r.Cookie(g.cookieFactory.key)
	if err == nil {
		email, ok = g.cookieFactory.ValidateCookie(cookie)
	}
	return email, ok
}

func (g *Guard) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.Header.Del("X-Forwarded-Email")
	email, ok := g.CheckCookie(r)

	if !ok {
		_, pattern := g.mux.Handler(r)
		if pattern == "" {
			log.Println("Please login")
			g.LoginPrompt(w, r)
		} else {
			log.Println("You're free to go")
			g.mux.ServeHTTP(w, r)
		}
		return
	}

	r.Header.Add("X-Forwarded-Email", email)
	next(w, r)
}

func (g *Guard) LoginPrompt(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("X-Forwarded-Email")
	g.cookieFactory.ClearCookie(w, r)

	partials := ""
	for _, p := range g.providers {
		partials += p.Partial(r)
	}
	g.template.ExecuteTemplate(w, "login.html", template.HTML(partials))
}
