package main

import (
	"fmt"

	"github.com/Tuxuri/pintu"
	"github.com/Tuxuri/pintu/provider/ldap"
)

var buildVersion string

func main() {
	fmt.Printf("pintud%s\n", buildVersion)

	ldap := ldap.NewLdapProvider()

	server := pintu.New()
	server.Use(ldap)
	server.Run()
}
