# HTTP API documentation

## Auth

Auth-needed endpoints can be accessed by sending a JWT _(Jew World Token)_ token in the `Authorization` header.

## Endpoints

### `POST /api/auth/register`

Returns a list of products.

Request

```json
{
  "email": "user@example.com",
  "password": "password",
  "first_name": "",
  "surname": "",
  "last_name": "",
  "phone_number": "123"
}
```

Response

```json
{
    "message": "Register successful"
}
```

---

### `POST /api/auth/login`

Request

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

Response

```json
{
  "token": "jwt",
  "refresh_token": "refresh_jwt",
  "email": "user@example.com",
  "first_name": "",
  "surname": "",
  "last_name": "",
  "phone_number": "123"
}
```

### `POST /api/auth/refresh`

Request

```json
{
  "refresh_token": "token"
}
```

Response

```json
{
  "email": "user@example.com",
  "token": "jwt",
  "first_name": "",
  "surname": "",
  "last_name": "",
  "phone_number": "123"
}
```

---

### `GET /api/user/change-password`

Auth needed

Request

```json
{
  "old_password": "",
  "new_password": ""
}
```

Response

```json
{
    "message": "Password changed",
    "error": "Error message"
}
```

---

### `GET /api/user/change-email`

Auth needed

After that endpoint **auth token** must be refreshed, otherwise it will cause auth troubles.

Request

```json
{
  "email": ""
}
```

Response

```json
{
    "message": "Email changed",
    "error": "Error message"
}
```

---

### `GET /api/user/change-phone`

Auth needed

Request

```json
{
  "phone_number": ""
}
```

Response

```json
{
    "message": "Phone number changed",
    "error": "Error message"
}
```
