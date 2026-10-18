// some unit/explanatory tests for the IONOS DNS plugin
// (c) copyright 2021 by Jan Delgado
package ionos

import (
	"encoding/json"
	"fmt"
	"slices"
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
			p := Provider{Provider: &ionos.Provider{}}
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

func TestUnmarshalCaddyFileExtractsAllowedZones(t *testing.T) {
	tests := []struct {
		test     string
		expected []string
	}{
		{"ionos token", nil},
		{"ionos token {\n allowed_zones a.com\n }", []string{"a.com"}},
		{"ionos token {\n allowed_zones a.com b.com\n }", []string{"a.com", "b.com"}},
		{`ionos {
			api_token token
			allowed_zones a.com
			allowed_zones b.com
		  }`, []string{"a.com", "b.com"}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			dispenser := caddyfile.NewTestDispenser(tc.test)
			p := Provider{Provider: &ionos.Provider{}}
			// when
			err := p.UnmarshalCaddyfile(dispenser)
			// then
			if err != nil {
				t.Errorf("UnmarshalCaddyfile failed with %v", err)
				return
			}
			if !slices.Equal(tc.expected, p.AllowedZones) {
				t.Errorf("expected AllowedZones to be %v but got %v", tc.expected, p.AllowedZones)
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
		{"ionos token {\n allowed_zones\n }", "wrong argument count"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			dispenser := caddyfile.NewTestDispenser(tc.test)
			p := Provider{Provider: &ionos.Provider{}}
			// when
			err := p.UnmarshalCaddyfile(dispenser)
			// then
			if err == nil {
				t.Fatalf("expected error with '%s' but got nil", tc.expected)
			}
			if !strings.Contains(strings.ToLower(err.Error()), tc.expected) {
				t.Errorf("expected error with '%s' but got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestProvisionTransformsAPIToken(t *testing.T) {
	// given
	expected := "{value}"
	p := Provider{Provider: &ionos.Provider{}}
	p.AuthAPIToken = "\\{value\\}"
	// when
	_ = p.Provision(caddy.Context{})
	// then
	actual := p.AuthAPIToken
	if expected != actual {
		t.Errorf("expected AuthAPIToken to be %s but got %s", expected, actual)
	}
}

func TestProvisionTransformsAllowedZones(t *testing.T) {
	// given
	expected := []string{"{value}"}
	p := Provider{Provider: &ionos.Provider{AuthAPIToken: "token"}}
	p.AllowedZones = []string{"\\{value\\}"}
	// when
	_ = p.Provision(caddy.Context{})
	// then
	if !slices.Equal(expected, p.AllowedZones) {
		t.Errorf("expected AllowedZones to be %v but got %v", expected, p.AllowedZones)
	}
}

func TestJSONConfig(t *testing.T) {
	tests := []struct {
		test     string
		expected []string
	}{
		{`{"auth_api_token":"token"}`, nil},
		{`{"auth_api_token":"token","allowed_zones":["a.com","b.com"]}`, []string{"a.com", "b.com"}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			p := Provider{Provider: &ionos.Provider{}}
			// when
			err := json.Unmarshal([]byte(tc.test), &p)
			// then
			if err != nil {
				t.Errorf("json.Unmarshal failed with %v", err)
				return
			}
			if p.AuthAPIToken != "token" {
				t.Errorf("expected AuthAPIToken to be 'token' but got '%s'", p.AuthAPIToken)
			}
			if !slices.Equal(tc.expected, p.AllowedZones) {
				t.Errorf("expected AllowedZones to be %v but got %v", tc.expected, p.AllowedZones)
			}
		})
	}
}

func TestProvisionAppliesAPITokenAlias(t *testing.T) {
	// given
	p := Provider{Provider: &ionos.Provider{}}
	dec := json.NewDecoder(strings.NewReader(`{"api_token":"token"}`))
	dec.DisallowUnknownFields() // as caddy does
	if err := dec.Decode(&p); err != nil {
		t.Fatalf("decode failed with %v", err)
	}
	// when
	err := p.Provision(caddy.Context{})
	// then
	if err != nil {
		t.Fatalf("Provision failed with %v", err)
	}
	if p.AuthAPIToken != "token" {
		t.Errorf("expected AuthAPIToken to be 'token' but got '%s'", p.AuthAPIToken)
	}
}

func TestProvisionRejectsBothTokens(t *testing.T) {
	// given
	p := Provider{Provider: &ionos.Provider{AuthAPIToken: "a"}, APIToken: "b"}
	// when
	err := p.Provision(caddy.Context{})
	// then
	if err == nil {
		t.Error("expected error when both api_token and auth_api_token are set")
	}
}

func TestProvisionSplitsAllowedZonesPlaceholder(t *testing.T) {
	// given
	t.Setenv("TEST_IONOS_ZONES", "a.com  b.com")
	p := Provider{Provider: &ionos.Provider{
		AuthAPIToken: "token",
		AllowedZones: []string{"{env.TEST_IONOS_ZONES}", "c.com"},
	}}
	// when
	err := p.Provision(caddy.Context{})
	// then
	if err != nil {
		t.Fatalf("Provision failed with %v", err)
	}
	expected := []string{"a.com", "b.com", "c.com"}
	if !slices.Equal(expected, p.AllowedZones) {
		t.Errorf("expected AllowedZones to be %v but got %v", expected, p.AllowedZones)
	}
}

func TestProvisionRejectsEmptyAllowedZone(t *testing.T) {
	// given
	p := Provider{Provider: &ionos.Provider{
		AuthAPIToken: "token",
		AllowedZones: []string{"a.com", "{env.TEST_IONOS_UNSET}"},
	}}
	// when
	err := p.Provision(caddy.Context{})
	// then
	if err == nil {
		t.Error("expected error for allowed_zones entry resolving to empty")
	}
}

func TestProvisionRejectsMissingAPIToken(t *testing.T) {
	tests := []string{"", "{env.TEST_IONOS_UNSET}"}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			p := Provider{Provider: &ionos.Provider{AuthAPIToken: tc}}
			// when
			err := p.Provision(caddy.Context{})
			// then
			if err == nil {
				t.Error("expected error for missing API token")
			}
		})
	}
}

func TestProvisionIsIdempotent(t *testing.T) {
	// given
	p := Provider{Provider: &ionos.Provider{}, APIToken: "token"}
	// when
	err1 := p.Provision(caddy.Context{})
	err2 := p.Provision(caddy.Context{})
	// then
	if err1 != nil || err2 != nil {
		t.Fatalf("Provision failed with %v, %v", err1, err2)
	}
	if p.AuthAPIToken != "token" {
		t.Errorf("expected AuthAPIToken to be 'token' but got '%s'", p.AuthAPIToken)
	}
}
