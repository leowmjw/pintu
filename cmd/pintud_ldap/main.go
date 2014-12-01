package main

import (
	"fmt"

	"github.com/Tuxuri/pintu"
	"github.com/Tuxuri/pintu/provider/google"
	"github.com/Tuxuri/pintu/provider/htpasswd"
	"github.com/Tuxuri/pintu/provider/ldap"
)

var buildVersion string

func main() {
	fmt.Printf("pintud%s\n", buildVersion)

	ldap := ldap.NewLdapProvider()
	google := google.NewGoogleOauthProvider()
	htpasswd := htpasswd.NewHtpasswdProvider()

	server := pintu.NewPintu()
	server.Use(htpasswd)
	server.Use(ldap)
	server.Use(google)
	server.Run()
}
