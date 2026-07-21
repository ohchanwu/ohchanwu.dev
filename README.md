# ohchanwu.dev

_[한국어로 보기 🇰🇷](README.ko.md)_

Personal portfolio site for Chanwu (Tyler) Oh.

The site is plain HTML, CSS, and JavaScript published from `static/` by GitHub Pages.

Project documentation is indexed in [docs/README.md](docs/README.md).

## Stack

- HTML
- CSS
- JavaScript
- GitHub Actions
- GitHub Pages

## Local preview

```sh
python3 -m http.server 4173 --bind 127.0.0.1 --directory static
```

Open `http://127.0.0.1:4173/`.

## Deployment

Pushes to `main` deploy `static/` through `.github/workflows/pages.yml`.

The repository owner manages the GitHub Pages custom-domain setting and Cloudflare DNS records.

## License

[MIT](LICENSE)
