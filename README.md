A simple proxy to by pass CORS restrictions when I use my web-based [RSS Reader](https://github.com/adhamsalama/simple-rss-reader) on my Kindle.

## Usage

```bash
go run main.go -port 8080
```

Then call it with `?url=` parameter:

```
http://localhost:8080/?url=https://example.com/feed.xml
```

## HTTPS (Mixed Content fix)

When calling the proxy from an HTTPS page (e.g. GitHub Pages), browsers block plain HTTP requests. Use a [cloudflared quick tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/do-more-with-tunnels/trycloudflare/) to expose the proxy over HTTPS without an account:

```bash
# Install on Termux
pkg install cloudflared

# Start tunnel
cloudflared tunnel --url http://localhost:8080
```

This prints a temporary `https://xxxx.trycloudflare.com` URL. Use that as your proxy base URL. The URL changes on every restart.
