# revoked — API

<div align="center">
  <h1>revoked</h1>
  <p>Self-host freely. Build on top. Don't sell it.</p>  

[![License: ELv2](https://img.shields.io/badge/License-Elastic_v2-blue.svg)](https://www.elastic.co/licensing/elastic-license)
[![API Docs](https://img.shields.io/badge/docs-revoked.link%2Fdocs-green)](https://revokedlink.github.io/docs)
[![Status: Experimental](https://img.shields.io/badge/status-experimental-yellow)](https://github.com/revokedlink/api)
</div>  

The backend API for [revoked](https://revoked.link), built with Go and PocketBase.

> [!WARNING]
> **Project Status:** `revoked` is currently in active development. It is **not ready for production use**, and
> backwards compatibility is not guaranteed until the v1.0.0 release.

## Development

```bash
# Clone
git clone https://github.com/revokedlink/api
cd api/cmd/revoked

# Run
go run main.go serve --http="0.0.0.0:7744"
```

## Branching

| Branch                                    | Purpose                                   |
|-------------------------------------------|-------------------------------------------|
| `main`                                    | Production — protected, triggers releases |
| `dev`                                     | Active development                        |
| `type/${issue-id}short-kebab-description` | Branch (e.g. feat/1-initial-setup)        |'

All work happens on `dev` or feature branches. Open a PR into `dev` when ready to release.

## Commit Convention

This repo uses [Conventional Commits](https://www.conventionalcommits.org):

| Prefix             | Version bump | Example                                |
|--------------------|--------------|----------------------------------------|
| `fix:`             | Patch        | `fix: handle nil certificate response` |
| `feat:`            | Minor        | `feat: add bulk revocation endpoint`   |
| `breaking:`        | Major        | `breaking: rename /check to /verify`   |
| `chore:` / `docs:` | None         | `chore: update dependencies`           |

You can reference issues by specifying the id.
```git
fix(ci): skip empty releases (#42)
```

The release workflow automatically determines the version bump from commits merged into `main`.

