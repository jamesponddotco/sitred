sitredctl(1)

# NAME

sitredctl - control the HTTP server for the SitRed service

# SYNOPSIS

*sitredctl* [OPTIONS] COMMAND [OPTIONS]

# DESCRIPTION

*sitredctl* is a control application for the SitRed service. As such, it can be
used to control the state of a SitRed server.

# OPTIONS

*--server-pid*
	Path to the server PID file. This field is mandatory.

*-h*, *--help*
	Show help message and quit.

*-v*, *--version*
	Show version number and quit.

# COMMANDS

*start* [ARGUMENTS]
	Start a SitRed server.

	Options are:

	*--tls-certificate*
		Path to the TLS certificate. This field is mandatory.

	*--tls-key*
		Path to the TLS key. This field is mandatory.

	*--tls-version*
		Minimum TLS version to use. Defaults to 1.3.

	*--server-address*
		HTTP server address to listen on. Defaults to 127.0.0.1:1997.

	*--server-cache-ttl*
		How long to cache sitemaps for. Defaults to 30 minutes.

	*--server-access-log*
		Whether to log incoming HTTP requests. Defaults to false.

	*--service-name*
		Name of the service. Defaults to sitred.

	*--service-contact*
		Contact email address or URL. Defaults to
		https://sr.ht/~jamesponddotco/sitred.

	*--sitemap-url*
		URL of the sitemap to use when choosing a random URL to redirect
		the user to. This field is mandatory.

*stop* [ARGUMENTS]
	Stop a running SitRed server.

# ENVIRONMENT

SITRED_SERVER_PID
	Path to the server PID file.

SITRED_TLS_CERTIFICATE
	Path to the TLS certificate.

SITRED_TLS_KEY
	Path to the TLS key.

SITRED_TLS_VERSION
	Minimum TLS version to use.

SITRED_ADDRESS
	HTTP server address to listen on.

SITRED_CACHE_TTL
	How long to cache sitemaps for.

SITRED_ACCESS_LOG
	Whether to log incoming HTTP requests.

SITRED_SERVICE_NAME
	Name of the service.

SITRED_SERVICE_CONTACT
	Contact email address or URL.

SITRED_SITEMAP_URL
	URL of the sitemap to use when choosing a random URL to redirect
	the user to.

# AUTHORS

Maintained by James Pond <james@cipher.host>.
