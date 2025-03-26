package validate

import (
	"errors"
	"net/url"
	"regexp"
)

var (
	hostRegexp = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+([a-zA-Z]{2,}|xn--[a-zA-Z0-9]+)$`)
)

var (
	ErrValidateURLNil           = errors.New("validate: url is nil")
	ErrValidateURLNonAbs        = errors.New("validate: url is not absolute")
	ErrValidateURLHostInvalid   = errors.New("validate: url host invalid")
	ErrValidateURLSchemeInvalid = errors.New("validate: url scheme invalid")
)

const (
	schemeHttp  = "http"
	schemeHttps = "https"
	schemeFtp   = "ftp"
	schemeData  = "data"
	schemeBlob  = "blob"
	schemeFile  = "file"
)

func IsValidURL(rawURL string) bool {
	if len(rawURL) == 0 {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return ValidateURL(parsed) == nil
}

func IsValidHttpURL(rawURL string) bool {
	if len(rawURL) == 0 {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return ValidateHttpURL(parsed) == nil
}

func IsValidImageSrcUrlScheme(scheme string) (remote bool, local bool) {
	switch scheme {
	case schemeHttp, schemeHttps, schemeFtp:
		return true, false
	}
	switch scheme {
	case schemeBlob, schemeData, schemeFile:
		return false, true
	}
	return false, false
}

func ValidateURL(u *url.URL) error {
	if u == nil {
		return ErrValidateURLNil
	}
	if !u.IsAbs() {
		return ErrValidateURLNonAbs
	}
	if !hostRegexp.MatchString(u.Hostname()) {
		return ErrValidateURLHostInvalid
	}
	return nil
}

func ValidateHttpURL(u *url.URL) error {
	if err := ValidateURL(u); err != nil {
		return err
	}
	switch u.Scheme {
	case schemeHttps, schemeHttp:
	default:
		return ErrValidateURLSchemeInvalid
	}
	return nil
}
