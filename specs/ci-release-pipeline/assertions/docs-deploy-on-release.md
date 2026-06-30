---
id: docs-deploy-on-release
parent: ci-release-pipeline
created: 2026-06-30T00:00:00Z
priority: 2
status: in_progress
locked-by: builder-Paris-MacBook-Pro-2.local-67349-1782851453
depends-on: release-workflow-on-tag
---

# GitHub Pages Deploys Docs on Release

A GitHub Actions workflow builds the Zensical docs site and deploys to GitHub Pages
whenever a release is published. Also supports manual trigger for ad-hoc updates.

## Success Criteria

- `.github/workflows/docs.yml` exists
- Triggers on `release: types: [published]` and `workflow_dispatch`
- Permissions include `contents: read`, `pages: write`, `id-token: write`
- Uses `actions/configure-pages`, `actions/setup-python`, installs `zensical`,
  runs `zensical build --clean`
- Uploads `site/` via `actions/upload-pages-artifact` and deploys via
  `actions/deploy-pages`
- Docs source lives in `docs/` with Zensical markdown files
- A `zensical.toml` config exists at the project root
