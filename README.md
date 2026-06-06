# ohchanwu.dev

_[한국어로 보기 🇰🇷](README.ko.md)_

Personal portfolio site for Chanwu (Tyler) Oh.

Built with Go and served from a single binary. Source code is the deployment.

## Stack

- Go (HTTP server, html/template)
- Vanilla HTML/CSS/JS (no framework)
- Docker / docker-compose
- Deployed on AWS EC2, fronted by Cloudflare

## Deploying

Origin TLS is handled by the binary itself via `golang.org/x/crypto/acme/autocert` (Let's Encrypt, HTTP-01). Cloudflare sits in front with Full (strict) mode; the LE cert satisfies CF at the origin.

Inside the container the server listens on `:8080` (HTTP, also serves ACME challenges and redirects to HTTPS) and `:8443` (HTTPS). Map host `80→8080` and `443→8443`.

### Cert cache (host bind-mount required)

Cert state must persist across container rebuilds. A wipe forces fresh issuance, which counts against Let's Encrypt's 5-certs-per-week limit per exact name set.

On the host, once:

```sh
sudo mkdir -p /var/cache/autocert
sudo chown 65532:65532 /var/cache/autocert
```

UID `65532` is the `nonroot` user inside the distroless runtime image; the binary cannot write to the cache without that ownership.

Bind-mount at run time:

```sh
docker run -d \
  -p 80:8080 -p 443:8443 \
  -v /var/cache/autocert:/var/cache/autocert \
  -e TLS_DOMAINS=ohchanwu.dev,www.ohchanwu.dev \
  -e ACME_EMAIL=you@example.com \
  ohchanwu-dev
```

Set `ACME_STAGING=true` on the first deploy to use Let's Encrypt's staging directory while smoke-testing — staging has no meaningful rate limit. Drop the flag and delete any staging certs from the cache dir before going live.

Do not delete `/var/cache/autocert` on a running deploy. The directory is safe to back up, not safe to wipe.

### HSTS (deferred)

The `Strict-Transport-Security` header is intentionally not set. Browsers cache the policy for the full `max-age` window — commonly a year — and any cert breakage during that window (expiry, misconfigured renewal, ACME pipeline failure) becomes an un-clickable browser wall with no operator escape hatch. There is no way to revoke an HSTS policy that's already been cached by a visitor; you can only wait it out.

Revisit only after autocert has survived its first auto-renewal cleanly (roughly 60 days after initial issuance), and even then start with a small `max-age` (e.g. 300 seconds → 86400 → 2592000 → 31536000 over a few weeks) before committing to a year.

## License

MIT — see [LICENSE](./LICENSE).
