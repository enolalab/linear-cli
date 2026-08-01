# Configuration commands

## `config set <key> <value>`

Save a config value and return the key and value.

```bash
linear-cli config set default_team ENG
```

## `config get <key>`

Return a config value.

```bash
linear-cli config get default_team
```

## `config list`

Return all saved configuration. Saved `api_key` values are redacted.

```bash
linear-cli config list
```

The intended keys and `default_team` scope are documented in [configuration](../configuration.md).
