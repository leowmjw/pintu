package ldap

import (
	"flag"
	"log"
	"net/http"

	"github.com/Tuxuri/pintu"
	"github.com/mqu/openldap"
)

type (
	LdapProvider struct {
		name          string
		start         string
		settings      *settings
		cookieFactory *pintu.CookieFactory
		ptype         string
	}

	settings struct {
		ldapServer string
		baseDN     string // Not required yet
	}
)

const (
	OptionLdapServer = "ldap_server"
	OptionBaseDN     = "ldap_base_dn"
)

func NewLdapProvider() *LdapProvider {
	s := &settings{}
	flag.StringVar(&s.ldapServer, OptionLdapServer, "", "LDAP Server URI ie ldap://127.0.0.1:389")
	return &LdapProvider{
		name:     "LDAP",
		start:    "/auth/ldap/start",
		settings: s,
		ptype:    "form",
	}
}

func (p *LdapProvider) ParseSettings() {
	if p.settings.ldapServer == "" {
		pintu.EnvStringVar(&p.settings.ldapServer, OptionLdapServer)
		if p.settings.ldapServer == "" {
			log.Fatalf("Please specify %s\n", OptionLdapServer)
		}
	}
}

func (p *LdapProvider) RegisterCookie(f *pintu.CookieFactory) {
	p.cookieFactory = f
}

func (p *LdapProvider) RegisterHandler(g *pintu.Guard) {
	g.HandleFunc(p.start, p.startHandler)
}

func (p *LdapProvider) Partial(r *http.Request) string {
	hostURL := pintu.GetHostURL(r)
	action := hostURL + p.start
	partial := &pintu.LoginPartial{
		Action:   action,
		Redirect: r.URL.RequestURI(),
		Name:     p.name,
	}
	return partial.GetForm(r)
}

func (p *LdapProvider) Type() string {
	return p.ptype
}

func (p *LdapProvider) Name() string {
	return p.name
}

func (p *LdapProvider) startHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	email := r.Form.Get("username")
	password := r.Form.Get("password")
	redirect := pintu.GetRedirect(r)

	ldap, err := openldap.Initialize(p.settings.ldapServer)
	if err != nil {
		log.Println(err)
		pintu.CustomError(w, r, pintu.ErrAuthServerDown)
		return
	}
	ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)
	err = ldap.Bind(email, password)
	if err != nil {
		pintu.CustomError(w, r, pintu.ErrInvalidCredentials)
		return
	} else {
		p.cookieFactory.SetCookie(email, w, r)
		http.Redirect(w, r, redirect, 302)
	}
	defer ldap.Close()
}
