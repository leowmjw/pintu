package pintu

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/codegangsta/negroni"
)

const loginPromptPath = "/auth"

type (
	Pintu struct {
		proxy     *negroni.Negroni
		providers []Provider
	}

	Provider interface {
		RegisterCookie(*CookieFactory)
		RegisterHandler(*Guard)
		Partial(*http.Request) string
		ParseSettings()
		Name() string
		Type() string
	}
)

func New() *Pintu {
	logger := negroni.NewLogger()
	recovery := negroni.NewRecovery()
	proxy := negroni.New(logger, recovery)
	return &Pintu{
		proxy: proxy,
	}
}

func (p *Pintu) Use(providers ...Provider) {
	p.providers = append(p.providers, providers...)
}

func (p *Pintu) Run() {
	settings := GetSettings()

	guard := NewGuard()
	cookieFactory := NewCookieFactory("_pintu", "123", int64(14*24)*time.Hour.Nanoseconds())
	guard.cookieFactory = cookieFactory

	guard.Use(p.providers...)

	// Put this to the settings validator
	upstreamurl, err := url.Parse(settings.Upstream)
	if err != nil {
		log.Fatal(err)
	}

	// Warning, the route declaration follows the order strictly
	mux := http.NewServeMux()
	mux.Handle("/", httputil.NewSingleHostReverseProxy(upstreamurl))
	mux.HandleFunc(loginPromptPath, guard.LoginPrompt)

	p.proxy.Use(guard)
	p.proxy.UseHandler(mux)
	p.proxy.Run(settings.HTTPAddress)
}
