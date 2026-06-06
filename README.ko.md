# ohchanwu.dev

_[Read this in English 🇬🇧](README.md)_

오찬우(Chanwu, Tyler Oh)의 개인 포트폴리오 사이트입니다.

Go로 작성되었으며 단일 바이너리(single binary)로 서비스됩니다. 소스 코드 자체가 곧 배포 단위입니다.

## 기술 스택(Stack)

- Go (HTTP 서버, html/template)
- 순수 HTML/CSS/JS (프레임워크 미사용)
- Docker / docker-compose
- AWS EC2에 배포, 앞단에 Cloudflare 배치

## 배포(Deploying)

오리진(origin) TLS는 바이너리 자체가 `golang.org/x/crypto/acme/autocert`를 통해 처리합니다(Let's Encrypt, HTTP-01 방식). 앞단의 Cloudflare는 Full (strict) 모드로 동작하며, Let's Encrypt(LE) 인증서가 오리진에서 Cloudflare(CF)의 요구 조건을 충족합니다.

컨테이너 내부에서 서버는 `:8080`(HTTP, ACME 챌린지 응답 및 HTTPS 리다이렉트도 함께 처리)과 `:8443`(HTTPS)에서 수신 대기합니다. 호스트 포트를 `80→8080`, `443→8443`으로 매핑하십시오.

### 인증서 캐시 (호스트 바인드 마운트 필수)

인증서 상태는 컨테이너 재빌드 사이에도 유지되어야 합니다. 캐시를 삭제하면 인증서를 새로 발급하게 되는데, 이는 동일한 이름 집합(exact name set)당 주 5회로 제한되는 Let's Encrypt의 발급 한도에 카운트됩니다.

호스트에서 최초 1회 실행:

```sh
sudo mkdir -p /var/cache/autocert
sudo chown 65532:65532 /var/cache/autocert
```

UID `65532`는 distroless 런타임 이미지 내부의 `nonroot` 사용자입니다. 이 소유권이 없으면 바이너리가 캐시에 쓸 수 없습니다.

실행 시 바인드 마운트:

```sh
docker run -d \
  -p 80:8080 -p 443:8443 \
  -v /var/cache/autocert:/var/cache/autocert \
  -e TLS_DOMAINS=ohchanwu.dev,www.ohchanwu.dev \
  -e ACME_EMAIL=you@example.com \
  ohchanwu-dev
```

첫 배포 시에는 `ACME_STAGING=true`를 설정하여, 점검(smoke-testing)하는 동안 Let's Encrypt의 스테이징(staging) 디렉터리를 사용하십시오 — 스테이징에는 의미 있는 발급 한도가 없습니다. 실서비스로 전환하기 전에 이 플래그를 제거하고 캐시 디렉터리에서 스테이징 인증서를 모두 삭제하십시오.

운영 중인 배포 환경에서 `/var/cache/autocert`를 삭제하지 마십시오. 이 디렉터리는 백업해도 안전하지만, 비워서는 안 됩니다.

### HSTS (보류)

`Strict-Transport-Security` 헤더는 의도적으로 설정하지 않았습니다. 브라우저는 이 정책을 `max-age` 기간 전체 — 보통 1년 — 동안 캐시하며, 그 기간 중 인증서에 문제가 생기면(만료, 잘못된 갱신 설정, ACME 파이프라인 장애) 운영자가 빠져나갈 방법이 없는, 클릭으로 우회할 수 없는 브라우저 차단 화면이 됩니다. 방문자에게 이미 캐시된 HSTS 정책을 취소할 방법은 없습니다. 기간이 끝날 때까지 기다리는 수밖에 없습니다.

autocert가 첫 자동 갱신을 문제없이 통과한 후(최초 발급으로부터 약 60일 후)에만 다시 검토하되, 그때에도 1년으로 확정하기 전에 작은 `max-age`부터 시작하십시오(예: 몇 주에 걸쳐 300초 → 86400 → 2592000 → 31536000).

## 라이선스(License)

MIT — [LICENSE](./LICENSE) 참조.
