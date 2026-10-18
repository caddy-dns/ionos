package ionos

import (
	"errors"
	"fmt"
	"strings"

	caddy "github.com/caddyserver/caddy/v2"
	caddyfile "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/ionos"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct {
	*ionos.Provider

	// APIToken replaces the deprecated auth_api_token.
	APIToken string `json:"api_token,omitempty"`
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.ionos",
		New: func() caddy.Module { return &Provider{Provider: new(ionos.Provider)} },
	}
}

// Before using the provider config, apply the api_token alias and resolve
// placeholders in the API token and allowed zones. A zone placeholder may
// expand to several whitespace-separated zones. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	if p.APIToken != "" {
		if p.AuthAPIToken != "" {
			return errors.New("only one of api_token and auth_api_token may be set")
		}
		p.AuthAPIToken, p.APIToken = p.APIToken, ""
	}
	repl := caddy.NewReplacer()
	p.AuthAPIToken = repl.ReplaceAll(p.AuthAPIToken, "")
	if p.AuthAPIToken == "" {
		return errors.New("missing API token")
	}
	var zones []string
	for _, zone := range p.AllowedZones {
		resolved := strings.Fields(repl.ReplaceAll(zone, ""))
		if len(resolved) == 0 {
			return fmt.Errorf("allowed_zones entry %q resolves to empty", zone)
		}
		zones = append(zones, resolved...)
	}
	p.AllowedZones = zones
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	ionos [<api_token>] {
//	    api_token <api_token>
//	    allowed_zones <zone> [<zone>...]
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.AuthAPIToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.AuthAPIToken != "" {
					return d.Err("API token already set")
				}
				if d.NextArg() {
					p.AuthAPIToken = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "allowed_zones":
				zones := d.RemainingArgs()
				if len(zones) == 0 {
					return d.ArgErr()
				}
				p.AllowedZones = append(p.AllowedZones, zones...)
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.AuthAPIToken == "" {
		return d.Err("missing API token")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
