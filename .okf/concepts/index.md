# Concepts

* [Stored format](sealed-format.md) - Sealed data is nonce(12) || ciphertext || tag(16), the name is the AAD, base64 is StdEncoding; frozen.
* [Name binding](name-binding.md) - The secret's name is authenticated but not stored, so a value only opens under the name it was sealed with.
* [Keys](keys.md) - Key file or environment, backups apart from the database, and rotation by hand.
* [Threat model](threat-model.md) - What panzerschrank protects against and what is up to the app.
