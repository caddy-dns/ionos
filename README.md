# IONOS DNS module for Caddy

[![run tests](https://github.com/caddy-dns/ionos/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/caddy-dns/ionos/actions/workflows/test.yml)

This package contains a DNS provider module for
[Caddy](https://github.com/caddyserver/caddy). It is used to manage DNS records
with the [IONOS DNS API](https://developer.hosting.ionos.com/docs/dns) using
[libdns-ionos](https://github.com/libdns/ionos).

## Caddy module name

```
dns.providers.ionos
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "ionos",
        "api_token": "YOUR_IONOS_AUTH_API_TOKEN"
      }
    }
  }
}
```

or with the Caddyfile. Syntax:

```
ionos [<api_token>] {
    api_token <api_token>
    allowed_zones <zone> [<zone>...]
}
```

This allows two forms. The short form takes the token as argument:

```
your.domain.com {
  respond "Hello World"	# replace with whatever config you need...
  tls {
    dns ionos {env.YOUR_IONOS_AUTH_API_TOKEN}
  }
}
```

The block form sets the token with `api_token` and allows further options
like `allowed_zones` (see below):

```
tls {
  dns ionos {
    api_token {env.YOUR_IONOS_AUTH_API_TOKEN}
    allowed_zones "mydomain.net"   # optional, see below
  }
}
```

Set the token in only one place, either as argument or as `api_token`.

`auth_api_token` is accepted as an alias for `api_token`. It is deprecated and
will be removed in version 2.

The IONOS webinterface will give you a so called `public prefix` and a `secret`
those two values concatenated by a dot (`.`) represent your API key.

For example
`cfc9247a69084db483bbfd4548350805.nFE1uW9G78kLl8siNFM-4y0jVut7Thx85jcWJt_qh48W2eYA8d079kcbx1K2HT9OhPVZGnavYlMsJIrCdut6Dg`

You can replace `{env.YOUR_IONOS_AUTH_API_TOKEN}` with the actual auth token if
you prefer to put it directly in your config instead of an environment
variable.

## Restricting writes to zones

IONOS API tokens grant write access to all zones of an account. The optional
`allowed_zones` setting limits write operations to the listed zones. Without
it, all zones are writable.

```
tls {
  dns ionos {
    api_token {env.YOUR_IONOS_AUTH_API_TOKEN}
    allowed_zones example.com example.org
  }
}
```

JSON: `"allowed_zones": ["example.com", "example.org"]`.

A placeholder such as `{env.ZONES}` may expand to several space-separated
zones. An entry that expands to nothing is an error.

This check runs client-side. It guards against misconfiguration, not against
a leaked token.

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/ionos#authenticating)
for information about obtaining credentials.

## Author

(c) Copyright 2021 by [Jan Delgado](https://github.com/jandelgado)

License: MIT
