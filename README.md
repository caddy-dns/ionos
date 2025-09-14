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
        "auth_api_token": "YOUR_IONOS_AUTH_API_TOKEN"
      }
    }
  }
}
```

or with the Caddyfile:

```
your.domain.com {
  respond "Hello World"	# replace with whatever config you need...
  tls {
    dns ionos {env.YOUR_IONOS_AUTH_API_TOKEN}
  }
}
```

The IONOS webinterface will give you a so called `public prefix` and a `secret`
those two values concatenated by a dot (`.`) represent your API key.

For example
`cfc9247a69084db483bbfd4548350805.nFE1uW9G78kLl8siNFM-4y0jVut7Thx85jcWJt_qh48W2eYA8d079kcbx1K2HT9OhPVZGnavYlMsJIrCdut6Dg`

You can replace `{env.YOUR_IONOS_AUTH_API_TOKEN}` with the actual auth token if
you prefer to put it directly in your config instead of an environment
variable.

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/ionos#authenticating) 
for information about obtaining credentials.

## Author

(c) Copyright 2021 by [Jan Delgado](https://github.com/jandelgado)

License: MIT
