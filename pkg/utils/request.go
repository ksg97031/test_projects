package utils

import (
	"Yi/pkg/logging"
	"crypto/tls"
	"net/http"
	"net/url"

	"go.uber.org/ratelimit"
)

/**
  @author: yhy
  @since: 2023/1/12
  @desc: //TODO
**/

type Session struct {
	// Client is the current http client
	Client *http.Client
	// Rate limit instance
	RateLimiter ratelimit.Limiter // Restriction of request rate per second
}

func NewSession(proxy string) *Session {
	Transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		DisableKeepAlives:   true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Add proxy
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		if isSupportedProtocol(proxyURL.Scheme) {
			Transport.Proxy = http.ProxyURL(proxyURL)
		} else {
			logging.Logger.Warnln("Unsupported proxy protocol: %s", proxyURL.Scheme)
		}
	}

	client := &http.Client{
		Transport: Transport,
	}
	session := &Session{
		Client: client,
	}

	// Github API access plus the access rate of Token is 5,000 times per hour, and the average of more than one second is more than once per second.
	session.RateLimiter = ratelimit.New(1)

	return session

}
