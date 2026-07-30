# Security Review Checklist

Use only for security-relevant diffs.

- Identify assets, actors, trust boundaries, and attacker-controlled inputs.
- Verify authentication and authorization at every state-changing boundary.
- Check secret, token, credential, PII, and session handling in storage, transport, errors, metrics, and logs.
- Check input validation, output encoding, injection paths, request forgery, replay, and confused-deputy behavior.
- Check cryptographic choices against current platform guidance; do not invent protocols.
- Check dependency and configuration changes for privilege or exposure expansion.
- State exploitability and evidence. Do not label a hypothetical concern as a vulnerability without a reachable path.
