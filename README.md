pintu
=====

Authentication reverse proxy server

## Architecture

```
    _______       _______       __________
    |Nginx| ----> |pintu| ----> |upstream|
    -------       -------       ----------

```

## Installation

Currently there's 3 flavours, to get the respective version

* Google OAuth `go get github.com/Tuxuri/pintu/cmd/pintud_google``
* LDAP `go get github.com/Tuxuri/pintu/cmd/pintud_ldap``
* HTPasswd `go get github.com/Tuxuri/pintu/cmd/pintud_htpasswd``
