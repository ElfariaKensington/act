# Standalone artifact server

The standalone binary requires `ACT_ARTIFACT_PASSWORD` at startup. It refuses to start with an empty password.

```sh
ACT_ARTIFACT_PASSWORD='replace-with-a-long-random-password' artifact-server --dir ./artifacts --retention-days 7
```

Use HTTPS (for example, via a TLS reverse proxy) for remote access: HTTP Basic Authentication sends a reusable password with each request. Keep the server bound to localhost or a private network whenever possible.

The password protects only the human-facing artifact index (`/` and `/artifacts`). Browsers will prompt when opening either page. Other routes — including GitHub Actions' artifact creation, upload, finalization, listing, deletion, signed-URL generation, summary images, and download endpoints — do **not** require HTTP Basic Authentication. This is intentional for GitHub Actions compatibility, including artifact overwrite operations. **Important:** the password does not protect individual artifact files or API access; anyone able to reach the server and obtain or guess a download URL may access artifacts. Restrict access to the server through a trusted network or reverse proxy if artifacts are private.

`--retention-days` now caps uploaded artifacts' lifetimes even if a client supplies a later `expires_at`. The cleanup routine runs at startup and every minute, including for previously uploaded artifacts whose saved metadata specifies a later expiration. Uploads that specify an earlier expiration keep their shorter lifetime. The cap is based on the metadata creation time. Older artifacts without metadata use the newest stored file modification time.

Signed blob URLs use a fresh cryptographic secret for each server process; existing signed URLs are invalidated by a restart.
