# Hosting the service

Start by [building and installing
**SitRed**](https://git.sr.ht/~jamesponddotco/sitred#installation), and then [grab a TLS certificate for your instance](https://certbot.eff.org/).

Now, to start **SitRed**, run this command:

```bash
sitredctl --server-pid '/path/to/your/pid-file.pid' start \
  --tls-certificate '/path/to/your/tls-certificate.pem' \
  --tls-key '/path/to/your/tls-key.pem' \
  --sitemap-url 'https://example.com/sitemap.xml'
```

You can run `siteredctl start --help` for more options.

For production you'll probably want to have a `systemd` service to run
that command for you. Here's a simple example of one.

```bash
[Unit]
Description=Redirect incoming requests to a random URL
Documentation=https://sr.ht/~jamesponddotco/sitred/
ConditionFileIsExecutable=/usr/bin/sitredctl
After=network.target nss-lookup.target

[Service]
Type=simple
UMask=117
Environment="SITRED_SERVER_PID=/path/to/your/pid-file.pid"
Environment="SITRED_TLS_CERTIFICATE=/path/to/your/tls-certificate.pem"
Environment="SITRED_TLS_KEY=/path/to/your/tls-key.pem"
ExecStart=/usr/bin/sitredctl start
ExecStop=/usr/bin/sitredctl stop
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

You'll want to improve your `systemd` service with sandbox and security
features, but that's beyond the scope of this documentation.

You'll also need to have a server such as NGINX in front of the service,
as it was written to sit behind one. Here's an example `location` for
NGINX.

```nginx
location / {
  proxy_pass https://random.example.com:1997;
  proxy_set_header Host $host;
  proxy_http_version 1.1;

  proxy_ssl_server_name on;
  proxy_ssl_protocols TLSv1.2 TLSv1.3;

  proxy_set_header X-Real-IP         $remote_addr;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header X-Forwarded-Host  $host;
  proxy_set_header X-Forwarded-Port  $server_port;

  proxy_connect_timeout 60s;
  proxy_send_timeout 60s;
  proxy_read_timeout 60s;
}
```

Again, for production you'll want to improve this `location` and have a
proper NGINX configuration file in place with rate limiting and other
security features, since the service itself doesn't implement any.

With everything up and running, you can now access the service at
`https://${ADDRESS}/` and it should redirect you to a random URL from
the sitemap you provided.
