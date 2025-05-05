// some unit/explanatory tests for the IONOS DNS plugin
// (c) copyright 2021 by Jan Delgado
package ionos

import (
	"fmt"
	"strings"
	"testing"

	caddy "github.com/caddyserver/caddy/v2"
	caddyfile "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/ionos"
)

func TestUnmarshalCaddyFileExtractsApiToken(t *testing.T) {
	tests := []string{
		"ionos token { }",
		`ionos {
			     api_token token
		       }`}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			dispenser := caddyfile.NewTestDispenser(tc)
			p := Provider{&ionos.Provider{}}
			// when
			err := p.UnmarshalCaddyfile(dispenser)
			// then
			if err != nil {
				t.Errorf("UnmarshalCaddyfile failed with %v", err)
				return
			}

			expected := "token"
			actual := p.AuthAPIToken
			if expected != actual {
				t.Errorf("Expected AuthAPIToken to be '%s' but got '%s'", expected, actual)
			}
		})
	}
}

func TestUnmarshalCaddyFileReportsErrorConditions(t *testing.T) {

	tests := []struct{ test, expected string }{
		{"ionos token invalid", "wrong argument count"},
		{"ionos { }", "missing api token"},
		{`ionos token { api_token token }`, "api token already set"},
		{`ionos { api_token token invalid }`, "wrong argument count"},
		{`ionos token { invalid token }`, "unrecognized subdirective 'invalid'"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			dispenser := caddyfile.NewTestDispenser(tc.test)
			p := Provider{&ionos.Provider{}}
			// when
			err := p.UnmarshalCaddyfile(dispenser)
			// then
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), tc.expected) {
				t.Errorf("expected error with '%s' but got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestProvisionTransformsAPIToken(t *testing.T) {
	// given
	expected := "{value}"
	p := Provider{&ionos.Provider{}}
	p.AuthAPIToken = "\\{value\\}"
	// when
	_ = p.Provision(caddy.Context{})
	// then
	actual := p.AuthAPIToken
	if expected != actual {
		t.Errorf("expected AuthAPIToken to be %s but got %s", expected, actual)
	}
}