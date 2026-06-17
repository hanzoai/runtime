# runtime

Hanzo Runtime — secure sandbox infrastructure for executing AI-generated code. Sub-90ms sandbox creation, isolated per-run, OCI/Docker image support, with File/Git/LSP/Execute APIs and optional persistence.

Monorepo (Poetry, `package-mode = false`): SDKs under `libs/` — `sdk-python` (`hanzo-runtime`), `api-client-python`, `api-client-python-async`. TypeScript SDK published as `@hanzo/runtime`.

- Python SDK: `pip install hanzo-runtime`
- TS SDK: `npm install @hanzo/runtime`

Full docs: README.md
