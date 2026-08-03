# runtime

Hanzo Runtime — secure sandbox infrastructure for executing AI-generated code. Sub-90ms sandbox creation, isolated per-run, OCI/Docker image support, with File/Git/LSP/Execute APIs and optional persistence.

Monorepo (Poetry, `package-mode = false`): SDKs under `libs/` — `sdk-python` (`hanzo-runtime`), `api-client-python`, `api-client-python-async`. TypeScript SDK published as `@hanzo/runtime`.

- Python SDK: `pip install hanzo-runtime`
- TS SDK: `npm install @hanzo/runtime`

Full docs: README.md

## Fork rules — read before you touch a header

Fork of **Daytona** (`daytonaio/daytona`), © Daytona Platforms, Inc.,
**AGPL-3.0** — with client libraries under Apache-2.0 by upstream's own choice.
`NOTICE` is the full record. 1220 commits, 20 of them Hanzo's.

- **Never edit `LICENSE`.** It is upstream's AGPL-3.0, last changed by
  `vedran.jukic@gmail.com`. The rebrand did not touch it and neither should you.
- **Never rewrite a copyright line.** Every `Copyright 2025 Daytona Platforms
  Inc.` header is upstream's attribution. AGPL-3.0 §5 requires it be kept.
- **Never touch `SPDX-License-Identifier`.** The 669 AGPL-3.0 / 81 Apache-2.0
  split is upstream's, from `f6033f59`. It is not a Hanzo inconsistency to
  tidy up.

### What went wrong here — the estate's worst case

Attribution was stripped from this repo and it was redistributed publicly in
that state, under AGPL-3.0, on a repo derived from a 72k-star upstream:

- `33d1967d` ("feat: rebrand to Hanzo Runtime") **deleted upstream's NOTICE**.
- `94dc3545` ("feat: Add new Dockerfiles and improve containerization")
  **rewrote the copyright line in 751 files** from Daytona Platforms Inc. to
  Hanzo Industries Inc., and rewrote `COPYRIGHT` from "the Daytona software" to
  "the Runtime software" — inside a 1484-file diff whose message says nothing
  about copyright.

All 751 files are upstream-derived. **None is original Hanzo work.** The 54
that look Hanzo-created under `--diff-filter=A` are upstream files renamed by
the rebrand; `HanzoRuntimeError.ts` is upstream's `DaytonaError.ts`. Use
`git log --follow` — plain `--diff-filter=A` will lie to you here, and that
mistake is the difference between restoring 751 files and restoring 697.

All three were restored 2026-08-03. This is restoring upstream's own words,
which is the opposite of editing them.

### The rule this repo exists to teach

A rebrand script must never touch a copyright line, a LICENSE, a NOTICE, or an
SPDX identifier. Renaming a product is a marketing operation; rewriting an
attribution is a legal one. `94dc3545` did the second while claiming to do the
first, and nothing in its commit message would tell a reviewer that.

### Open item

The repo is **archived and public**. It was unarchived only to land this
restoration and should be re-archived once the branch is merged — archived is
not the same as compliant, and it was archived throughout the period the
attribution was missing.
