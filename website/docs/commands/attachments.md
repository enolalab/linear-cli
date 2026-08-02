# Attachments

## `attachment upload`

Upload a local file and link it to an issue. Both the issue identifier and file path are required.

```bash
linear-cli attachment upload --issue ENG-123 --file ./screenshot.png --title "Login failure"
```

| Flag | Required | Description |
| --- | --- | --- |
| `--issue <identifier>` | Yes | Issue identifier, such as `ENG-123` |
| `--file <path>` | Yes | Local file to upload |
| `--title <text>` | No | Attachment title; defaults to the filename |

## Presigned upload flow

The command reads the file, detects its MIME type from the extension or content, resolves the issue, and then:

1. Calls Linear's `fileUpload` mutation to request a presigned upload URL and asset URL.
2. Sends the local file with an HTTP `PUT` to that presigned URL, using a 60-second client timeout.
3. Calls `attachmentCreate` with the resolved issue ID, title, and asset URL.

The final `data` contains the created attachment, asset URL, file name, size, content type, and original issue identifier. The file data is sent to the presigned storage URL, not included in the JSON response.
