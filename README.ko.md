# ohchanwu.dev

_[Read this in English 🇬🇧](README.md)_

오찬우(Chanwu, Tyler Oh)의 개인 포트폴리오 사이트입니다.

`static/`의 순수 HTML, CSS, JavaScript를 GitHub Pages로 배포합니다.

프로젝트 문서는 [docs/README.md](docs/README.md)에 정리되어 있습니다.

## 기술 스택

- HTML
- CSS
- JavaScript
- GitHub Actions
- GitHub Pages

## 로컬 미리 보기

```sh
python3 -m http.server 4173 --bind 127.0.0.1 --directory static
```

`http://127.0.0.1:4173/`을 여십시오.

## 배포

`main`에 푸시된 변경 사항은 `.github/workflows/pages.yml`을 통해 `static/`을 배포합니다.

저장소 소유자가 GitHub Pages 사용자 지정 도메인 설정과 Cloudflare DNS 레코드를 관리합니다.

## 라이선스

[MIT](LICENSE)
