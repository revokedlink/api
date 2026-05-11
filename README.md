# revoked — API

The backend API for [revoked](https://revoked.link), built with Go and PocketBase.

## Development

```bash
# Clone
git clone https://github.com/revokedlink/api
cd api

# Run
go run main.go serve --http="0.0.0.0:7744"
```

## Branching

| Branch | Purpose |
|--------|---------|
| `main` | Production — protected, triggers releases |
| `dev`  | Active development |

All work happens on `dev` or feature branches. Open a PR into `dev` when ready to release.

## Commit Convention

This repo uses [Conventional Commits](https://www.conventionalcommits.org):

| Prefix | Version bump | Example |
|--------|-------------|---------|
| `fix:` | Patch | `fix: handle nil certificate response` |
| `feat:` | Minor | `feat: add bulk revocation endpoint` |
| `breaking:` | Major | `breaking: rename /check to /verify` |
| `chore:` / `docs:` | None | `chore: update dependencies` |

The release workflow automatically determines the version bump from commits merged into `main`.

## License

[Elastic License 2.0](https://www.elastic.co/licensing/elastic-license) — free to use and self-host.
