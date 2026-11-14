# Standalone artifact server

The standalone binary requires `ACT_ARTIFACT_PASSWORD` at startup. It refuses to start with an empty password.

```sh
ACT_ARTIFACT_PASSWORD='replace-with-a-long-random-password' artifact-server --dir ./artifacts --retention-days 7
```

Use HTTPS (for example, via a TLS reverse proxy) for remote access: HTTP Basic Authentication sends a reusable password with each request. Keep the server bound to localhost or a private network whenever possible.

The password protects the root artifact listing, `/artifacts`, local summary and download endpoints, legacy artifact browsing/download endpoints, and v4 artifact listing, signed-URL generation, and deletion. Browsers will prompt for the password. Automation calling these protected endpoints must send `Authorization: Basic ...` as well; otherwise it receives HTTP 401. The upload/finalization APIs and signed v4 blob transfer URLs remain accessible for workflow compatibility. This is **not** full API authentication: do not expose the artifact upload service directly to untrusted clients.

`--retention-days` now caps uploaded artifacts' lifetimes even if a client supplies a later `expires_at`. The cleanup routine runs at startup and every minute, including for previously uploaded artifacts whose saved metadata specifies a later expiration. Uploads that specify an earlier expiration keep their shorter lifetime. The cap is based on the metadata creation time. Older artifacts without metadata use the newest stored file modification time.

Signed blob URLs use a fresh cryptographic secret for each server process; existing signed URLs are invalidated by a restart.
