package pintu

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// Settings contains configuration required by Pintu
	Settings struct {
		HTTPAddress  string
		Upstream     string
		UpstreamURL  *url.URL
		CookieKey    string
		CookieSecret string
		CookieExpiry int64
	}

	StringSlice []string
)

const (
	optionHTTPAddress  = "http"
	optionUpstream     = "upstream"
	optionCookieKey    = "cookie_key"
	optionCookieSecret = "cookie_secret"
	optionCookieExpiry = "cookie_expiry"

	defaultHTTPAddress            = "127.0.0.1:4180"
	defaultUpstream               = ""
	defaultCookieKey              = "_pintu"
	defaultCookieSecret           = "randomly generated sha1 hash"
	defaultCookieExpiryHour int64 = 168 // 7 days
)

func (l *StringSlice) Set(s string) error {
	*l = append(*l, s)
	return nil
}

func (l *StringSlice) String() string {
	return fmt.Sprint(*l)
}

func EnvStringSliceVar(ss *StringSlice, field string) {
	values := os.Getenv(field)
	if values != "" {
		var l StringSlice
		for _, item := range strings.Split(values, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				l = append(l, item)
			}
		}
		if len(l) > 0 {
			*ss = l
		}
	}
}

// Warning, settings error is not tolerated, all the invalid input will result in os.Exit

// GetSettings returns settings from cli and env
// Settings priority CLI > ENV
func GetSettings() Settings {
	cli := cliSettings()
	env := envSettings()
	s := merge(cli, env)
	s.validateSettings()
	return s
}

// cliSettings returns settings acquired from CLI
func cliSettings() Settings {
	settings := Settings{}
	flag.StringVar(&settings.HTTPAddress, optionHTTPAddress, defaultHTTPAddress, "<addr>:<port> to listen on for HTTP clients")
	flag.StringVar(&settings.Upstream, optionUpstream, defaultUpstream, "the http url of the upstream endpoint")
	flag.StringVar(&settings.CookieKey, optionCookieKey, defaultCookieKey, "the name of the secure cookies")
	flag.StringVar(&settings.CookieSecret, optionCookieSecret, defaultCookieSecret, "the seed string for secure cookies")
	flag.Int64Var(&settings.CookieExpiry, optionCookieExpiry, defaultCookieExpiryHour, "cookie lifespan in hour")
	flag.Parse()
	return settings
}

// envSettings returns settings acquired from environment variables
func envSettings() Settings {
	settings := Settings{}
	EnvStringVar(&settings.HTTPAddress, optionHTTPAddress, defaultHTTPAddress)
	EnvStringVar(&settings.Upstream, optionUpstream, defaultUpstream)
	EnvStringVar(&settings.CookieKey, optionCookieKey, defaultCookieKey)
	EnvStringVar(&settings.CookieSecret, optionCookieSecret, defaultCookieSecret)
	EnvInt64Var(&settings.CookieExpiry, optionCookieExpiry, defaultCookieExpiryHour)
	return settings
}

// EnvStringVar is helper function for environment variables lookip
func EnvStringVar(option *string, field, defaultval string) {
	if value := os.Getenv(field); value != "" {
		*option = value
	} else {
		*option = defaultval
	}
}

func EnvInt64Var(option *int64, field string, defaultval int64) {
	stringval := os.Getenv(field)
	if stringval != "" {
		value, err := strconv.ParseInt(stringval, 10, 64)
		if err != nil {
			log.Fatalf("param %s requires integer character", field)
		}
		if value > 0 {
			*option = value
		} else {
			*option = defaultval
		}
	} else {
		*option = defaultval
	}
}

// merge filters settings values from CLI and ENV according to preset priority
func merge(cli Settings, env Settings) Settings {
	settings := Settings{}
	settings.HTTPAddress = TopString(cli.HTTPAddress, env.HTTPAddress, defaultHTTPAddress)
	settings.Upstream = TopString(cli.Upstream, env.Upstream, defaultUpstream)
	settings.CookieKey = TopString(cli.CookieKey, env.CookieKey, defaultCookieKey)
	settings.CookieSecret = TopString(cli.CookieSecret, env.CookieSecret, defaultCookieSecret)
	settings.CookieExpiry = TopInt64(cli.CookieExpiry, env.CookieExpiry, defaultCookieExpiryHour)
	return settings
}

// TopVar applies the settings priority rule and filter the one with highest priority
// returns default value if available
func TopString(cli, env, defaultval string) string {
	if cli != "" && cli != defaultval {
		return cli
	}
	if env != "" {
		return env
	}
	return defaultval
}

func TopInt64(cli, env, defaultval int64) int64 {
	if cli > 0 && cli != defaultval {
		return cli
	}
	if env > 0 {
		return env
	}
	return defaultval
}

func (s *Settings) validateSettings() {
	if s.HTTPAddress == "" {
		log.Fatalf("missing param %s", optionHTTPAddress)
	}

	if s.Upstream == "" {
		log.Fatalf("missing param %s", optionUpstream)
	}

	if s.CookieKey == "" {
		log.Fatalf("missing param %s", optionCookieKey)
	}

	if s.CookieSecret == "" || s.CookieSecret == defaultCookieSecret {
		h := sha1.New()
		io.WriteString(h, fmt.Sprintf("%d", time.Now().Unix()))
		s.CookieSecret = hex.EncodeToString(h.Sum(nil))
	}

	if s.CookieExpiry < 1 {
		log.Fatalf("missing param %s", optionCookieSecret)
	}
}
