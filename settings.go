package pintu

import (
	"flag"
	"net/url"
	"os"
)

type (
	// Settings contains configuration required app's initiation
	Settings struct {
		HTTPAddress  string
		Upstream     string
		UpstreamURL  *url.URL
		CookieKey    string
		CookieSecret string
		CookieExpiry int64
	}
)

const (
	optionHTTPAddress  = "http"
	optionUpstream     = "upstream"
	optionCookieKey    = "cookie_key"
	optionCookieSecret = "cookie_secret"
	optionCookieExpiry = "cookie_expiry"

	defaultHTTPAddress            = "127.0.0.1:4180"
	defaultCookieKey              = "_helios"
	defaultCookieExpiryHour int64 = 168 // 7 days
)

// GetSettings returns settings from cli and env
// Settings priority CLI > ENV
func GetSettings() Settings {
	cli := cliSettings()
	env := envSettings()
	return merge(cli, env)
}

// cliSettings returns settings acquired from CLI
func cliSettings() Settings {
	settings := Settings{}
	flag.StringVar(&settings.HTTPAddress, optionHTTPAddress, "", "<addr>:<port> to listen on for HTTP clients")
	flag.StringVar(&settings.Upstream, optionUpstream, "", "the http url of the upstream endpoint")
	flag.StringVar(&settings.CookieKey, optionCookieKey, defaultCookieKey, "the name of the secure cookies")
	flag.StringVar(&settings.CookieSecret, optionCookieSecret, "", "the seed string for secure cookies")
	flag.Int64Var(&settings.CookieExpiry, optionCookieExpiry, defaultCookieExpiryHour, "cookie lifespan in hour")
	flag.Parse()
	return settings
}

// envSettings returns settings acquired from environment variables
func envSettings() Settings {
	settings := Settings{}
	EnvStringVar(&settings.HTTPAddress, optionHTTPAddress)
	EnvStringVar(&settings.Upstream, optionUpstream)
	return settings
}

// envStringVar is helper function for environment variables lookip
func EnvStringVar(option *string, field string) {
	if value := os.Getenv(field); value != "" {
		*option = value
	}
}

// merge filters settings values from CLI and ENV according to preset priority
func merge(cli Settings, env Settings) Settings {
	settings := Settings{}
	settings.HTTPAddress = TopVar(cli.HTTPAddress, env.HTTPAddress, defaultHTTPAddress)
	settings.Upstream = TopVar(cli.Upstream, env.Upstream, "")
	return settings
}

// TopVar applies the settings priority rule and filter the one with highest priority
// returns default value if available
func TopVar(cli, env, defaultval string) string {
	if cli != "" && cli != defaultHTTPAddress {
		return cli
	}
	if env != "" {
		return env
	}
	return defaultval
}
