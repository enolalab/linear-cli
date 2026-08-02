# Output and errors

All normal command responses are JSON on standard output. Use `--pretty` to indent the same envelope for people; do not enable it when a compact stream is preferable.

## Success envelope

```json
{
  "success": true,
  "data": {},
  "pagination": {
    "hasNextPage": true,
    "endCursor": "cursor-from-linear"
  }
}
```

`data` is command-specific. `pagination` appears only when the command passes pagination data to the output formatter.

## Error envelope

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "issue not found: ENG-999"
  }
}
```

| Exit code | JSON `error.code` | Meaning |
| ---: | --- | --- |
| 0 | N/A | Success |
| 1 | `GENERAL_ERROR` | Unspecified or internal operation error |
| 2 | `USAGE_ERROR` | Invalid arguments, required flags, or local input failures |
| 3 | `AUTH_ERROR` | Missing or invalid API key |
| 4 | `NOT_FOUND` | Requested resource was not found |
| 5 | `RATE_LIMITED` | Linear rate limit was reached |
| 6 | `NETWORK_ERROR` | Network connectivity failure |
| 7 | `PERMISSION_DENIED` | Insufficient permissions |

Exit status is the stable branch point for automation. Parse `error.code` and `error.message` only after checking a nonzero status.
