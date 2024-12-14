# chirpy

A web server for users to make and interact with others social media posts.

API Instructions:

## user resource

```json
{
  "id": "1a01c51d-50ec-405c-b6d8-7374699f725e",
  "created_at": "2024-12-13T10:39:59.477961Z",
  "updated_at": "2024-12-13T10:39:59.477961Z",
  "email": "mary@gmail.com",
  "is_chirpy_red": false
}
```

### POST /api/users

Creates a new user

Request body:

```json
{
  "password": "123456",
  "email": "mary@gmail.com"
}
```

Response body:

```json
{
  "id": "1a01c51d-50ec-405c-b6d8-7374699f725e",
  "created_at": "2024-12-13T10:39:59.477961Z",
  "updated_at": "2024-12-13T10:39:59.477961Z",
  "email": "mary@gmail.com",
  "is_chirpy_red": false
}
```

### POST /api/login

Login as a created user

Request body:

```json
{
  "password": "123456",
  "email": "mary@gmail.com"
}
```

Response body:

```json
{
  "id": "1a01c51d-50ec-405c-b6d8-7374699f725e",
  "created_at": "2024-12-13T10:39:59.477961Z",
  "updated_at": "2024-12-13T10:39:59.477961Z",
  "email": "mary@gmail.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOjI1YmM1YzA4NS0yNDI1LTQzMDAtYjk1MC02MDkyMjc3Nzg0MWIiLCJleHAiOjE3MzQxMDk0NTAsImlheCI6MTckNDEwNTg1MM0.YJAWtDCow1p2XOLBEJMTMnEsKMBGq34_i1TEvxSzSVg",
  "refresh_token": "79b67165a4848367f1479bd927d1e3f26abea0e4c80b10b1254a9756cb763930!"
}
```

### PUT /api/users

Update user email and password

Request body:

```json
{
  "password": "123456",
  "email": "mary@gmail.com"
}
```

Response body:

```json
{
  "id": "1a01c51d-50ec-405c-b6d8-7374699f725e",
  "created_at": "2024-12-13T10:39:59.477961Z",
  "updated_at": "2024-12-13T10:39:59.477961Z",
  "email": "mary@gmail.com",
  "is_chirpy_red": false
}
```

## chirp resource

```json
{
  "id": "5d22126d-f4ec-41cb-af6e-c77b42b4db85",
  "created_at": "2024-12-13T10:45:39.477961Z",
  "updated_at": "2024-12-13T10:45:39.477961Z",
  "body": "This is my chirp",
  "user_id": "1a01c51d-50ec-405c-b6d8-7374699f725e"
}
```

### POST /api/chirps

Create a chirp

Request body:

```json
{
  "body": "This is my chirp"
}
```

Response body:

```json
{
  "id": "5d22126d-f4ec-41cb-af6e-c77b42b4db85",
  "created_at": "2024-12-13T10:45:39.477961Z",
  "updated_at": "2024-12-13T10:45:39.477961Z",
  "body": "This is my chirp",
  "user_id": "1a01c51d-50ec-405c-b6d8-7374699f725e"
}
```

### GET /api/chirps/{chirpID}

Retrieve a chirp by the chirp ID specified in the endpoint. If no ID is specified, all
chirps will be returned.

Queries:

author_id: Returns all chirps from a specified user
sort: asc/desc sorts the resultsin the specified created_at order (default asc)

e.g. /api/chirps?sort=desc&author_id=1a01c51d-50ec-405c-b6d8-7374699f725e

### DELETE /api/chirps/{chirpID}

Delete chirp specified by ID at endpoint

### POST /api/refresh

Refresh user access token

Response body:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOjI1YmM1YzA4NS0yNDI1LTQzMDAtYjk1MC02MDkyMjc3Nzg0MWIiLCJleHAiOjE3MzQxMDk0NTAsImlheCI6MTckNDEwNTg1MM0.YJAWtDCow1p2XOLBEJMTMnEsKMBGq34_i1TEvxSzSVg"
}
```

### POST /api/revoke

Revokes the users refresh token

### POST /api/polka/webhooks

Webhook upgrade user to Chirpy Red

Request body:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "1a01c51d-50ec-405c-b6d8-7374699f725e"
  }
}
```
