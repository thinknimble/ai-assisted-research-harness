---
title: "Typeform Webhooks API Reference"
type: reference
status: active
created: 2026-06-30
updated: 2026-06-30
author: pari

source: web
source_url: "https://www.typeform.com/developers/webhooks/"

tags: [typeform, webhooks, api, surveys, integrations]
related: []
summary: "Complete Typeform Webhooks API reference — CRUD endpoints, payload schema, signature verification, and security. Covers all webhook lifecycle operations."
---

# Typeform Webhooks API Reference

Base URL: `https://api.typeform.com`

All requests require an access token passed via `Authorization: Bearer {your_access_token}`.

---

## Walkthroughs

Source: https://www.typeform.com/developers/webhooks/walkthroughs/

### Create a new webhook

You need the `form_id` for the typeform and your webhook URL.

```bash
curl --request PUT \
  --url https://api.typeform.com/forms/{form_id}/webhooks/{tag} \
  --header 'Authorization: Bearer {your_access_token}' \
  --header 'Content-Type: application/json' \
  -d '{"url":"https://{your_webhook_url}.com", "enabled":true}'
```

Set `enabled` to `true` to start pushing data immediately on submission.

### Change the URL for an existing webhook

Send the same `PUT` request with the new URL:

```bash
curl --request PUT \
  --url https://api.typeform.com/forms/{form_id}/webhooks/{tag} \
  --header 'Authorization: Bearer {your_access_token}' \
  --header 'Content-Type: application/json' \
  -d '{"url":"https://{your_new_webhook_url}.com", "enabled":true}'
```

> **NOTE:** You can also use the `POST` method to change a webhook URL.

### Delete a webhook

```bash
curl --request DELETE \
  --url https://api.typeform.com/forms/{form_id}/webhooks/{tag} \
  --header 'Authorization: Bearer {your_access_token}'
```

---

## Example Webhook Payload

Source: https://www.typeform.com/developers/webhooks/example-payload/

Webhooks deliver responses in JSON. The payload has this top-level structure:

```json
{
  "event_id": "LtWXD3crgy",
  "event_type": "form_response",
  "form_response": { ... }
}
```

### Event information

| Element | Type | Description |
|---|---|---|
| `event_id` | string | Unique ID for the webhook event. Automatically assigned by Typeform. |
| `event_type` | string | Reason the webhook is being sent. |
| `form_response` | object | Contains information about the typeform and the response. |

### form_response fields

| Element | Type | Description |
|---|---|---|
| `form_id` | string | Unique ID for the typeform (from the form URL). |
| `token` | string | Unique ID for the submission (identical to response_id in Responses API). |
| `response_url` | string | URL linking to the response in the Typeform results dashboard. |
| `submitted_at` | string | ISO 8601 UTC timestamp of submission. |
| `landed_at` | string | ISO 8601 UTC timestamp of form landing. |
| `calculated` | object | Contains `score` (integer) if the typeform includes score calculation. |
| `variables` | array | Array of typeform variables with `key` (string), `type` ("text"\|"number"), and the value field. |
| `hidden` | object | Hidden field values (if any). |
| `definition` | object | Lists questions in the typeform (see below). |
| `answers` | array | Respondent's answers (see below). |
| `ending` | object | Contains `id` and `ref` of the ending screen shown. |

### definition object

| Element | Type | Description |
|---|---|---|
| `id` | string | Unique ID for the typeform. |
| `title` | string | Title of the typeform. |
| `fields` | array | Questions in the typeform. Order matches the `answers` array. |
| `endings` | array | Ending screens defined in the typeform. |

Each field object may include:

| Element | Type | Description |
|---|---|---|
| `id` | string | Unique ID for the field. Use to match questions with answers. |
| `title` | string | Title of the question. |
| `type` | string | Question type (see answer types below). |
| `ref` | string | Reference name for the field. Can be set via Create API or auto-generated. Must be <255 chars, matching `^[a-zA-Z0-9_-]+$`. |
| `allow_multiple_selections` | boolean | `true` if respondents can select multiple choices. |
| `allow_other_choice` | boolean | `true` if the question includes an "Other" option. |
| `choices` | array | For `multiple_choice`, `picture_choice`, `dropdown` types. Each has `id`, `label`, `ref`. |

### answers array

Each answer object includes:

| Element | Type | Description |
|---|---|---|
| `type` | string | Answer type (see mapping below). |
| `{answer}` | varies | The actual answer value. Key name depends on type. |
| `answer_url` | string | URL to answer details in the Typeform dashboard. |
| `field` | object | Contains `id`, `type`, and optionally `ref` to match with the question. |

### Answer type mapping

| Answer type | Question types |
|---|---|
| `text` | `short_text`, `long_text` |
| `choice` | `dropdown`, `multiple_choice`, `picture_choice` (single selection) |
| `choices` | `dropdown`, `multiple_choice`, `picture_choice`, `ranking`, `checkbox` (multiple selection) |
| `date` | `date` |
| `boolean` | `legal`, `yes_no` |
| `number` | `rating`, `opinion_scale`, `number` |
| `file_url` | `file_upload` |
| `multi_format` | video/audio answers |
| `payment` | `payment` |
| `url` | `website`, `calendly`, `google_calendar` |
| `signature` | `signature` |

### Full payload example

```json
{
  "event_id": "LtWXD3crgy",
  "event_type": "form_response",
  "form_response": {
    "form_id": "lT4Z3j",
    "token": "a3a12ec67a1365927098a606107fac15",
    "response_url": "https://admin.typeform.com/form/lT4Z3j/results?responseId=a3a12ec67a1365927098a606107fac15#responses",
    "submitted_at": "2018-01-18T18:17:02Z",
    "landed_at": "2018-01-18T18:07:02Z",
    "calculated": { "score": 9 },
    "variables": [
      { "key": "score", "type": "number", "number": 4 },
      { "key": "name", "type": "text", "text": "typeform" }
    ],
    "hidden": { "user_id": "abc123456" },
    "definition": {
      "id": "lT4Z3j",
      "title": "Webhooks example",
      "fields": [
        {
          "id": "DlXFaesGBpoF",
          "title": "Thanks, {{answer_60906475}}! What's it like where you live?",
          "type": "long_text",
          "ref": "readable_ref_long_text",
          "allow_multiple_selections": false,
          "allow_other_choice": false
        }
      ],
      "endings": [
        {
          "id": "dN5FLyFpCMFo",
          "ref": "01GRC8GR2017M6WW347T86VV39",
          "title": "Bye!",
          "type": "thankyou_screen",
          "properties": {
            "button_text": "Create a typeform",
            "show_button": true,
            "share_icons": true,
            "button_mode": "default_redirect"
          }
        }
      ]
    },
    "answers": [
      {
        "type": "text",
        "text": "It's cold right now! I live in an older medium-sized city...",
        "answer_url": "https://admin.typeform.com/form/lT4Z3j/results?...",
        "field": { "id": "DlXFaesGBpoF", "type": "long_text" }
      },
      {
        "type": "email",
        "email": "laura@example.com",
        "field": { "id": "SMEUb7VJz92Q", "type": "email" }
      },
      {
        "type": "choice",
        "choice": { "id": "4WIlUvKOl0UB", "label": "London", "ref": "..." },
        "field": { "id": "k6TP9oLGgHjl", "type": "multiple_choice" }
      },
      {
        "type": "boolean",
        "boolean": true,
        "field": { "id": "gFFf3xAkJKsr", "type": "legal" }
      },
      {
        "type": "number",
        "number": 3,
        "field": { "id": "WOTdC00F8A3h", "type": "rating" }
      }
    ],
    "ending": {
      "id": "dN5FLyFpCMFo",
      "ref": "01GRC8GR2017M6WW347T86VV39"
    }
  }
}
```

---

## Secure Your Webhooks

Source: https://www.typeform.com/developers/webhooks/secure-your-webhooks/

Since webhook endpoints are publicly exposed URLs, you should verify that payloads actually come from Typeform. Typeform signs each payload with a secret using HMAC SHA-256, and includes the signature in the `Typeform-Signature` header.

### Set up your webhook secret

1. Generate a random string:
   ```bash
   ruby -rsecurerandom -e 'puts SecureRandom.hex(20)'
   ```
2. Update the webhook's `secret` field via a [PUT request to the Webhooks API](#create-or-update-webhook).

### Validate payload from Typeform

1. Using HMAC SHA-256, create a hash of the entire received payload (as binary) using the `secret` as the key.
2. Encode the binary hash in base64.
3. Prepend `sha256=` to the encoded hash.
4. Compare the result with the `Typeform-Signature` header value.

#### Node.js (Express)

```javascript
const crypto = require('crypto')

app.use(express.raw({ type: 'application/json' }))

app.post('/webhook', async (request, response) => {
  const signature = request.headers['typeform-signature']
  const isValid = verifySignature(signature, request.body.toString())
})

const verifySignature = function (receivedSignature, payload) {
  const hash = crypto
    .createHmac('sha256', process.env.SECRET_TOKEN)
    .update(payload)
    .digest('base64')
  return receivedSignature === `sha256=${hash}`
}
```

#### Node.js (Fastify)

```javascript
const crypto = require('crypto')
const fastify = require('fastify')()

await fastify.register(require('fastify-raw-body'))

fastify.post('/typeform/webhook', (request, reply) => {
  const signature = request.headers['typeform-signature']
  const isValid = verifySignature(signature, request.rawBody)
})

const verifySignature = function (receivedSignature, payload) {
  const hash = crypto
    .createHmac('sha256', process.env.SECRET_TOKEN)
    .update(payload)
    .digest('base64')
  return receivedSignature === `sha256=${hash}`
}
```

#### Ruby

```ruby
post '/webhook' do
  request.body.rewind
  payload_body = request.body.read
  verify_signature(request.env['HTTP_TYPEFORM_SIGNATURE'], payload_body)
  "Payload received: #{payload_body.inspect}"
end

def verify_signature(received_signature, payload_body)
  hash = OpenSSL::HMAC.digest(OpenSSL::Digest.new('sha256'), ENV['SECRET_TOKEN'], payload_body)
  actual_signature = 'sha256=' + Base64.strict_encode64(hash)
  return halt 500, "Signatures don't match!" unless Rack::Utils.secure_compare(actual_signature, received_signature)
end
```

#### Python (FastAPI)

```python
from fastapi import FastAPI, Request, HTTPException
import hashlib, hmac, base64, os

app = FastAPI()

@app.post("/hook")
async def recWebHook(req: Request):
    body = await req.json()
    raw = await req.body()
    receivedSignature = req.headers.get("typeform-signature")
    if receivedSignature is None:
        return HTTPException(403, detail="Permission denied.")
    sha_name, signature = receivedSignature.split('=', 1)
    if sha_name != 'sha256':
        return HTTPException(501, detail="Operation not supported.")
    is_valid = verifySignature(signature, raw)
    if not is_valid:
        return HTTPException(403, detail="Invalid signature.")

def verifySignature(receivedSignature: str, payload):
    WEBHOOK_SECRET = os.environ.get('TYPEFORM_SECRET_KEY')
    digest = hmac.new(WEBHOOK_SECRET.encode('utf-8'), payload, hashlib.sha256).digest()
    e = base64.b64encode(digest).decode()
    return e == receivedSignature
```

#### PHP

```php
<?php
$headers = getallheaders();
$header_signature = $headers["Typeform-Signature"];
$secret = getenv("TYPEFORM_WEBHOOK_SECRET");
$payload = @file_get_contents("php://input");
$hashed_payload = hash_hmac("sha256", $payload, $secret, true);
$base64encoded = "sha256=" . base64_encode($hashed_payload);

if ($header_signature === $base64encoded) {
    echo "success!\n";
}
```

#### Swift

```swift
import CryptoKit

func verifySig(receivedSig: String, payload: Request.Body) -> Bool {
    let secretString = "abc123" // replace with your own
    let payloadString = payload.string ?? ""
    let key = SymmetricKey(data: Data(secretString.utf8))
    let regenSig = HMAC<SHA256>.authenticationCode(for: Data(payloadString.utf8), using: key)
    let sigData = Data(regenSig)
    let sigBase64 = sigData.base64EncodedString()
    let final = "sha256=\(sigBase64)"
    return final == receivedSig
}
```

### HTTPS requirement

All new webhook URLs **must** use `https`. Attempting to create/update with `http` returns `400 Bad Request` with error code `webhook_url_https_required`. SSL/TLS certificates must be valid (self-signed will not work).

**Existing HTTP webhooks:** Legacy `http` webhooks continue to work. You can update other fields (`enabled`, `secret`, `event_types`) without changing the URL, but you cannot change the URL to another `http` address -- only to `https`.

**The `verify_ssl` field** is now read-only and auto-derived from the URL scheme:
- `https` URLs -> `verify_ssl: true`
- `http` URLs (legacy) -> `verify_ssl: false`

If included in a request, it will be ignored.

> **NOTE:** Typeform is hosted on AWS with dynamic IP addresses. Static IP or IP ranges cannot be guaranteed for webhook requests.

---

## API Reference

### Retrieve webhooks

Source: https://www.typeform.com/developers/webhooks/reference/retrieve-webhooks/

```
GET /forms/{form_id}/webhooks
```

Retrieve all webhooks for the specified typeform.

**Path Parameters:**

| Parameter | Type | Required | Description |
|---|---|---|---|
| `form_id` | string | yes | Unique ID for the form. Found in form URL (e.g. `u6nXL7` from `https://mysite.typeform.com/to/u6nXL7`). |

**Response:** `200 OK`

```json
{
  "items": [
    {
      "created_at": "2016-11-21T12:23:28.000Z",
      "enabled": true,
      "event_types": {
        "form_response": true,
        "form_response_partial": true
      },
      "form_id": "abc123",
      "id": "yRtagDm8AT",
      "tag": "phoenix",
      "updated_at": "2016-11-21T12:23:28.000Z",
      "url": "https://test.com",
      "verify_ssl": true
    }
  ]
}
```

### Retrieve single webhook

Source: https://www.typeform.com/developers/webhooks/reference/retrieve-single-webhook/

```
GET /forms/{form_id}/webhooks/{tag}
```

Retrieve a single webhook.

**Path Parameters:**

| Parameter | Type | Required | Description |
|---|---|---|---|
| `form_id` | string | yes | Unique ID for the form. |
| `tag` | string | yes | Unique name of the webhook. |

**Response:** `200 OK`

```json
{
  "created_at": "2016-11-21T12:23:28.000Z",
  "enabled": true,
  "event_types": {
    "form_response": true,
    "form_response_partial": true
  },
  "form_id": "abc123",
  "id": "yRtagDm8AT",
  "tag": "phoenix",
  "updated_at": "2016-11-21T12:23:28.000Z",
  "url": "https://test.com",
  "verify_ssl": true
}
```

### Create or update webhook

Source: https://www.typeform.com/developers/webhooks/reference/create-or-update-webhook/

```
PUT /forms/{form_id}/webhooks/{tag}
```

Create or update a webhook.

**Path Parameters:**

| Parameter | Type | Required | Description |
|---|---|---|---|
| `form_id` | string | yes | Unique ID for the form. |
| `tag` | string | yes | Unique name for the webhook. |

**Request Body:**

| Field | Type | Description |
|---|---|---|
| `enabled` | boolean | `true` to send responses immediately. |
| `event_types` | object | Event types the webhook subscribes to (e.g. `{"form_response": true, "form_response_partial": true}`). |
| `secret` | string | Used to sign payload with HMAC SHA-256 for verification. |
| `url` | string | Webhook URL (must be `https`). |
| `verify_ssl` | boolean | Read-only. Auto-derived from URL scheme. |

**Request Example:**

```json
{
  "enabled": true,
  "event_types": {
    "form_response_partial": true
  },
  "url": "https://test.com"
}
```

**Response:** `200 OK` -- returns the full webhook object (same schema as retrieve).

### Delete webhook

Source: https://www.typeform.com/developers/webhooks/reference/delete-webhook/

```
DELETE /forms/{form_id}/webhooks/{tag}
```

Delete a webhook.

**Path Parameters:**

| Parameter | Type | Required | Description |
|---|---|---|---|
| `form_id` | string | yes | Unique ID for the form. |
| `tag` | string | yes | Unique name of the webhook. |

**Responses:**

- `204 No Content` -- webhook deleted successfully.
- `404 Not Found` -- webhook not found. Error response schema:

| Field | Type | Description |
|---|---|---|
| `code` | string | Snake_case error key. |
| `description` | string | Developer-readable error description. |
| `details` | array | Optional array with field-level error info (`code`, `description`, `field`, `help`, `in`). |
| `help` | string | URL linking to help content. |

---

## Webhook Object Schema

All webhook responses share this schema:

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique ID for the webhook. |
| `form_id` | string | Unique ID for the typeform. |
| `tag` | string | Unique name for the webhook. |
| `url` | string | Webhook URL. |
| `enabled` | boolean | Whether responses are sent to the webhook immediately. |
| `secret` | string | HMAC SHA-256 signing secret (if set). |
| `event_types` | object | Subscribed event types (`form_response`, `form_response_partial`). |
| `verify_ssl` | boolean | Whether SSL certificates are verified (read-only, derived from URL scheme). |
| `created_at` | string | ISO 8601 creation timestamp. |
| `updated_at` | string | ISO 8601 last-update timestamp. |
