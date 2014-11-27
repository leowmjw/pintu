package pintu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CookieFactory struct {
	key    string
	secret string
	expiry time.Duration
}

func NewCookieFactory(key string, secret string, expires int64) *CookieFactory {
	// expires in hour
	expiresDuration := time.Duration(expires * int64(time.Hour))
	return &CookieFactory{
		key:    key,
		secret: secret,
		expiry: expiresDuration,
	}
}

// ValidateCookie checks the cookie validity
func (c *CookieFactory) ValidateCookie(cookie *http.Cookie) (string, bool) {
	// value, timestamp, signature
	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 3 {
		return "", false
	}
	sig := c.getCookieSignature(parts[0], parts[1])
	if parts[2] == sig {
		// it's a valid cookie. now get the contents
		ts, err := strconv.Atoi(parts[1])
		if err == nil && int64(ts) > time.Now().Add(c.expiry*-1).Unix() {
			rawValue, err := base64.URLEncoding.DecodeString(parts[0])
			if err == nil {
				return string(rawValue), true
			}
		}
	}
	return "", false
}

// getSignedCookieValue compiles user unique cookie string
func (c *CookieFactory) getSignedCookieValue(value string) string {
	encodedValue := base64.URLEncoding.EncodeToString([]byte(value))
	timeStr := fmt.Sprintf("%d", time.Now().Unix())
	sig := c.getCookieSignature(encodedValue, timeStr)
	return fmt.Sprintf("%s|%s|%s", encodedValue, timeStr, sig)
}

// getCookieSignature compiles base64 encoded cookie string
func (c *CookieFactory) getCookieSignature(args ...string) string {
	h := hmac.New(sha1.New, []byte(c.key))
	h.Write([]byte(c.secret))
	for _, arg := range args {
		h.Write([]byte(arg))
	}
	var b []byte
	b = h.Sum(b)
	return base64.URLEncoding.EncodeToString(b)
}

// SetCookie set cookie for the authenticated
func (c *CookieFactory) SetCookie(value string, w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:     c.key,
		Value:    c.getSignedCookieValue(value),
		Path:     "/",
		Domain:   GetDomain(req),
		Expires:  time.Now().Add(c.expiry),
		HttpOnly: true,
		Secure:   IsSecured(req),
	}
	http.SetCookie(w, cookie)
}

func (c *CookieFactory) ClearCookie(w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:     c.key,
		Value:    "",
		Path:     "/",
		Domain:   GetDomain(req),
		Expires:  time.Now().Add(time.Duration(1) * time.Hour * -1),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
