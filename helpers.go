package pintu

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
)

var ErrAPIError = errors.New("api request returned non 200 status code")

func IsSecured(req *http.Request) bool {
	if scheme := req.Header.Get("X-Forwarded-Proto"); scheme == "https" {
		return true
	}
	return false
}

func GetDomain(req *http.Request) string {
	return strings.Split(req.Host, ":")[0]
}

// GetHostURL retrieves host url from http request
func GetHostURL(req *http.Request) string {
	proto := "http"
	if req.Header.Get("X-Forwarded-Proto") != "" && len(req.Header["X-Forwarded-Proto"]) > 0 {
		proto = req.Header["X-Forwarded-Proto"][0]
	}
	hostStrings := []string{proto, "://", req.Host}
	host := strings.Join(hostStrings, "")
	return host
}

// GetHostPath merges host url with url path
func GetHostPath(req *http.Request, url string) string {
	return GetHostURL(req) + url
}

// GetRedirect checks request for referer and returns path for redirection
func GetRedirect(r *http.Request) string {
	redirect := r.FormValue("rd")
	if redirect == "" || strings.Contains(redirect, loginPromptPath) {
		redirect = "/"
	}
	return redirect
}

// GetRemoteIP retrieves visitor ip address
func GetRemoteIP(req *http.Request) string {
	remoteIP := req.Header.Get("X-Real-IP")
	if remoteIP == "" {
		remoteIP = req.RemoteAddr
	}
	return strings.Split(remoteIP, ":")[0]
}

// // APIRequest processes http requests and serializes response to json
func APIRequest(r *http.Request) (*simplejson.Json, error) {
	httpclient := &http.Client{}
	resp, err := httpclient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("got response code %d - %s", resp.StatusCode, body)
		return nil, ErrAPIError
	}

	data, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
