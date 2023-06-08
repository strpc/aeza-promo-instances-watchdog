[Aeza provider](https://aeza.net/?ref=369322) promo-instances watchdog. *Can be killed if they consume a lot of memory.*
We come and startup. Lightweight, written on Go with zero-dependencies. Github action configuration ready.

## Usage

0. [Generate API KEY](https://my.aeza.net/settings/apikeys).

### github action

1. [fork](https://github.com/strpc/aeza-promo-instances-watchdog/fork) this repository.
2. Add secrets in repository: Settings -> Secrets and variables -> Actions:
   `AEZA_API_KEY`(from step 0), `IP_ADDRESSES`(ip addresses join with ",". example: `1.1.1.1,8.8.8.8`)

### run inside docker

```shell
docker run --restart unless-stopped \
  -e AEZA_API_KEY="superapitoken" \
  -e IP_ADDRESSES="1.1.1.1,8.8.8.8" \
  strpc/aeza-promo-instances-watchdog:latest
```

### run on host

```shell
go install github.com/strpc/aeza-promo-instances-watchdog
aeza-promo-instances-watchdog

# or

git clone https://github.com/strpc/aeza-promo-instances-watchdog.git 
cd aeza-promo-instances-watchdog
go build -o aeza-promo-instances-watchdog ./cmd/aeza-promo-instances-watchdog
./aeza-promo-instances-watchdog
```

### configuration

| Key                  | Description                                                                | Default  value        |
|----------------------|----------------------------------------------------------------------------|-----------------------|
| `AEZA_API_KEY`       | API key for access. [Generate](https://my.aeza.net/settings/apikeys)       | ❌ requiered           |
| `IP_ADDRESSES`       | IP addresses for monitoring. Join with `,`. <br>Example: `1.1.1.1,8.8.8.8` | ❌ requiered           |
| `WATCH_DELAY`        | Delay for retry watch.                                                     | `5m`                  |
| `AEZA_API_BASE_URL`  | Aeza base API url.                                                         | https://core.aeza.net |
| `AEZA_VM_BASE_URL`   | Aeza base VM-panel url.                                                    | https://vm.aeza.net   |
| `NTFY_CHANNEL`       | [ntfy](https://ntfy.sh/)-channel for notify about start/errors.            | ""                    |
| `GITHUB_ACTION_MODE` | Run in github action cron.                                                 | ""                    |
