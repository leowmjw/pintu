package google

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Tuxuri/pintu"
)

type (
	GoogleOauthProvider struct {
		name          string
		start         string
		redirect      string
		redemption    *url.URL
		login         *url.URL
		userInfo      *url.URL
		scopes        string
		settings      *settings
		cookieFactory *pintu.CookieFactory
		ptype         string
	}

	settings struct {
		clientId string
		secret   string
		domains  pintu.StringSlice
	}
)

const (
	optionGoogleClientId     = "google_client_id"
	optionGoogleClientSecret = "google_client_secret"
	optionGoogleDomain       = "google_domains"
)

var (
	errMissingCode    = errors.New("missing code")
	errDomainMismatch = errors.New("domain mismatch")
)

// NewGoogleOauthProvider bootstrap handler and authenticator
func NewGoogleOauthProvider() *GoogleOauthProvider {
	redemption, _ := url.Parse("https://accounts.google.com/o/oauth2/token")
	login, _ := url.Parse("https://accounts.google.com/o/oauth2/auth")
	userInfo, _ := url.Parse("https://www.googleapis.com/oauth2/v2/userinfo")
	scopes := "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email"

	s := &settings{}
	flag.StringVar(&s.clientId, optionGoogleClientId, "", "Google OAuth client ID: ie: \"123456.apps.googleusercontent.com\"")
	flag.StringVar(&s.secret, optionGoogleClientSecret, "", "Google OAuth client secret")
	flag.Var(&s.domains, optionGoogleDomain, "Google Apps domain")

	return &GoogleOauthProvider{
		name:       "Google",
		start:      "/oauth2/google/start",
		redirect:   "/oauth2/google/callback",
		redemption: redemption,
		login:      login,
		userInfo:   userInfo,
		scopes:     scopes,
		settings:   s,
		ptype:      "link",
	}
}

func (p *GoogleOauthProvider) ParseSettings() {
	if p.settings.clientId == "" {
		pintu.EnvStringVar(&p.settings.clientId, optionGoogleClientId, "")
		if p.settings.clientId == "" {
			log.Fatalf("missing param %s", optionGoogleClientId)
		}
	}

	if p.settings.secret == "" {
		pintu.EnvStringVar(&p.settings.secret, optionGoogleClientSecret, "")
		if p.settings.secret == "" {
			log.Fatalf("missing param %s", optionGoogleClientSecret)
		}
	}

	if len(p.settings.domains) == 0 {
		pintu.EnvStringSliceVar(&p.settings.domains, optionGoogleDomain)
	}
}

func (p *GoogleOauthProvider) RegisterCookie(f *pintu.CookieFactory) {
	p.cookieFactory = f
}

func (p *GoogleOauthProvider) RegisterHandler(g *pintu.Guard) {
	g.HandleFunc(p.start, p.startHandler)
	g.HandleFunc(p.redirect, p.redirectHandler)
}

// Partial provides AuthProxy rendered Partial of Login Form
func (p *GoogleOauthProvider) Partial(r *http.Request) string {
	hostURL := pintu.GetHostURL(r)
	action := hostURL + p.start
	partial := &pintu.LoginPartial{
		Action:   action,
		Redirect: r.URL.RequestURI(),
		Name:     p.name,
		Type:     "btn-google-plus",
		Btn:      "fa-google",
	}
	return partial.GetLink(r)
}

func (p *GoogleOauthProvider) Type() string {
	return p.ptype
}

func (p *GoogleOauthProvider) Name() string {
	return p.name
}

// GetLoginURL compile redirect url to provider's portal
func (p *GoogleOauthProvider) loginURL(r *http.Request) (string, error) {
	state := pintu.GetRedirect(r)
	callback := pintu.GetHostPath(r, p.redirect)

	params := url.Values{}
	params.Add("redirect_uri", callback)
	params.Add("prompt", "select_account")
	params.Add("scope", p.scopes)
	params.Add("client_id", p.settings.clientId)
	params.Add("response_type", "code")
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", p.login, params.Encode()), nil
}

// redeem consumes authorization code acquired
func (p *GoogleOauthProvider) redeem(r *http.Request) (string, error) {
	code := r.Form.Get("code")
	if code == "" {
		return "", errMissingCode
	}

	callback := pintu.GetHostPath(r, p.redirect)
	params := url.Values{}
	params.Add("redirect_uri", callback)
	params.Add("client_id", p.settings.clientId)
	params.Add("client_secret", p.settings.secret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", p.redemption.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Printf("failed building request %s", err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	json, err := pintu.APIRequest(req)
	if err != nil {
		log.Printf("failed making request %s", err.Error())
		return "", err
	}

	accessToken, err := json.Get("access_token").String()
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// GetUserInfo retrieves user information from provider
func (p *GoogleOauthProvider) getinfo(token string) (string, error) {
	params := url.Values{}
	params.Add("access_token", token)
	endpoint := fmt.Sprintf("%s?%s", p.userInfo.String(), params.Encode())

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("failed building request %s", err.Error())
		return "", err
	}

	json, err := pintu.APIRequest(req)
	if err != nil {
		log.Printf("failed making request %s", err.Error())
		return "", err
	}

	email, err := json.Get("email").String()
	if err != nil {
		log.Printf("failed getting email from response %s", err.Error())
		return "", err
	}
	return email, nil
}

func (p *GoogleOauthProvider) startHandler(w http.ResponseWriter, r *http.Request) {
	login, err := p.loginURL(r)
	if err != nil {
		pintu.CustomError(w, r, pintu.ErrAuthServerDown)
		return
	}
	http.Redirect(w, r, login, 302)
	return
}

func (p *GoogleOauthProvider) redirectHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		pintu.CustomError(w, r, err)
		return
	}

	errResp := r.Form.Get("error")
	if errResp != "" {
		pintu.Denied(w, r)
		return
	}

	token, err := p.redeem(r)
	if err != nil {
		log.Printf("error redeeming code %s", err.Error())
		pintu.CustomError(w, r, err)
		return
	}

	email, err := p.getinfo(token)
	if err != nil {
		log.Printf("error redeeming code %s", err.Error())
		pintu.CustomError(w, r, err)
		return
	}

	redirect := r.Form.Get("state")
	if redirect == "" {
		redirect = "/"
	}

	log.Printf("validating againsts domains %v", p.settings.domains)
	if !p.validate(email) {
		pintu.CustomError(w, r, errDomainMismatch)
		return
	}

	log.Printf("authenticating %s completed", email)
	p.cookieFactory.SetCookie(email, w, r)
	http.Redirect(w, r, redirect, 302)
	return
}

func (p *GoogleOauthProvider) validate(email string) bool {
	if len(p.settings.domains) > 0 {
		valid := false
		for _, domain := range p.settings.domains {
			suffix := fmt.Sprintf("@%s", domain)
			valid = valid || strings.HasSuffix(email, suffix)
		}
		return valid
	}
	return true
}
