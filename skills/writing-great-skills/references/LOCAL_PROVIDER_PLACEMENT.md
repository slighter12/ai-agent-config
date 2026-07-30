# Local Provider Placement

Keep the canonical skill under this repository's `skills/<name>/`.

- Codex discovers shared skills through `~/.agents/skills`; user-only skills also require `agents/openai.yaml`.
- Claude uses the same linked skill and top-level `disable-model-invocation: true`.
- OpenCode uses the same linked skill and `metadata["opencode/autoinvoke"] = "false"`.
- Gemini receives the shared link on a best-effort basis; report runtime behavior as unverified unless it was probed.

The installer preserves user files and links from other repositories.
