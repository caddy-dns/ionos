# IONOS DNS module for Caddy

This package contains a DNS provider module for
[Caddy](https://github.com/caddyserver/caddy). It is used to manage DNS records
with the [IONOS DNS API](https://developer.hosting.ionos.com/docs/dns) using
[libdns-ionos](https://github.com/libdns/ionos)..

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

or with the Caddyfile:

```
your.domain.com {
  respond "Hello World"	# replace with whatever config you need...
  tls {
    dns ionos {env.YOUR_IONOS_AUTH_API_TOKEN}
  }
}
```

You can replace `{env.YOUR_IONOS_AUTH_API_TOKEN}` with the actual auth token if
you prefer to put it directly in your config instead of an environment
variable.

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/ionos#authenticating) 
for information about obtaining credentials.

## Author

(c) Copyright 2021 by Jan Delgado
License: MIT

