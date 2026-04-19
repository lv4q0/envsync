# envsync

> CLI tool to diff and sync `.env` files across environments with secret masking

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git
cd envsync && go build -o envsync .
```

---

## Usage

**Diff two `.env` files:**

```bash
envsync diff .env.local .env.production
```

**Sync missing keys from one file to another:**

```bash
envsync sync .env.local .env.production
```

**Mask secret values during diff:**

```bash
envsync diff .env.local .env.production --mask-secrets
```

Example output:

```
~ API_URL        development.api.com → production.api.com
+ NEW_FEATURE_FLAG  (missing in production)
- LEGACY_KEY        (not present in local)
```

Secrets matching patterns like `*_KEY`, `*_SECRET`, or `*_TOKEN` are automatically masked as `****` when `--mask-secrets` is enabled.

---

## License

[MIT](LICENSE)