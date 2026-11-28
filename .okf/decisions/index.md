# Decisions

* [One shared package](one-shared-package.md) - Replace the copies in gpanel and static-hosting with one tested package that opens their data unchanged.
* [The format is frozen](format-is-frozen.md) - Stored data must keep opening; a new format needs a new function and a major version.
* [No secret material in errors](no-secrets-in-errors.md) - Every open failure is ErrDecrypt; key errors never quote the key.
* [Standard library only](stdlib-only.md) - No dependencies and no own cryptography.
* [Random nonces only](random-nonces.md) - A fresh nonce from crypto/rand for every seal; callers never supply one.
* [The key file is never overwritten](key-file-never-overwritten.md) - Wrong size is an error; creation is atomic via temp file and hard link.
* [Release by hand from main](release-from-main.md) - Releases are tagged by a manual workflow run on main; the release branch was dropped.
