# HTTP API documentation

## Auth

Auth-needed endpoints can be accessed by sending a JWT token in the `Authorization` header.

## Endpoints

### `POST /auth/register`

Returns a list of products.

#### Request

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

#### Response

```json
{
    "message": "Register successful"
}
```

---

### `POST /auth/login`

#### Request

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

#### Response

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

### `GET /user/change-password`

Auth needed

#### Request

```json
{
  "old_password": "",
  "new_password": ""
}
```

#### Response

```json
{
    "message": "Password changed",
    "error": "Error message"
}
```
