# Keys

A key is 32 random bytes. Whoever has it can open every secret sealed with it; whoever loses it loses every secret sealed with it.

## Key file or environment

| | Key file (`OpenKeyFile`) | Environment (`NewFromBase64`) |
|---|---|---|
| Setup | none: created on the first start, mode 0600 | create once: `openssl rand -base64 32` |
| Good for | one server with a data directory (gpanel) | containers, platforms with a secret store |
| Backup | back up the file, apart from the database | keep it in your secret store |

`OpenKeyFile` never replaces an existing file. A file with the wrong size is an error, not a reason to make a new key, because that would lose every secret. When several processes start at the same time, all of them get the key of the one that created the file; the file never appears half written, also not after a crash.

## Backups

Back up the key **apart from the database**. A backup holding both together protects nothing; a backup of the database without any copy of the key cannot be restored.

## Changing the key

panzerschrank has no built-in rotation; with small numbers of secrets it is a short loop in the app:

1. Open the old box and create a new one (`GenerateKey` or a new file).
2. In one database transaction, for every sealed value: `Open` with the old box, `Seal` with the new one under the same name, write it back.
3. After the commit, switch the app to the new key and keep the old one until the next backup no longer needs it.

One key can seal about 4 billion values before it should be changed, which small secrets do not reach in practice.
