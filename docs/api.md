# HTTP API documentation

## Auth

Auth-needed endpoints can be accessed by sending a JWT _(Jew World Token)_ token in the `Authorization` header.

## Endpoints

### `POST /api/auth/register`

Returns a list of products.

Phone must be >= 11 chars length before and == 11 after formatting

Phone number may be in any format, has endless amount of spaces (only spaces) and `(`, `)`, `-` characters, it can starts with either `+7` or `8`.

To properly store phone it will be formatted to `81234567890`.

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

### `POST /api/user/get-data`

Auth needed

Response

```json
{
  "email": "user@example.com",
  "first_name": "",
  "surname": "",
  "last_name": "",
  "phone_number": "123"
}
```

### `POST /api/user/update-company-data`

Auth needed

Request

```json
{
  "companyName": "Company name",
  "companyAddress": "Company address",
  "positionInCompany": "123",
  "companyINN": "123"
}
```

Response 200

```json
{
    "message": "Company data updated"
}
```

Response 400, 401, 500

```json
{
    "error": "Error message"
}
```

### `POST /api/user/get-company-data`

Auth needed

Response

```json
{
  "companyName": "Company name",
  "companyAddress": "Company address",
  "positionInCompany": "123",
  "companyINN": "123",
  "allRequiredFields": true
}
```

### `POST /api/orders/create`

Auth needed

Users must has all company data to create order

**Known issues:**

- Currently if there will be 2 equal products in request, there won't be errors, but in future this will be fixed.

Request

```json
{
  "products": [
    {
      "id": "e76aefbc-0fa1-4fe4-af90-b320438b03b4",
      "quantity": 2
    },
    {
      "id": "5aa5c1a0-61d8-418e-8346-6edd9decb864",
      "quantity": 3
    }
  ]
}
```

Response 200

```json
{
    "message": "Order created",
    "order_id": "e76aefbc-0fa1-4fe4-af90-b320438b03b4"
}
```

Response 400, 401, 404, 500

```json
{
    "error": "Error message"
}
```

### `GET /api/orders/get`

Auth needed

Currently not implemented

Response 200

```json
{
    "orders": [
        {
            "id": 12,
            "created_at": "2020-01-01 00:00:00",
            "products": [
                {
                    "id": 1,
                    "name": "Product 1",
                    "quantity": 2,
                    "price": 100
                },
                {
                    "id": 2,
                    "name": "Product 2",
                    "quantity": 3,
                    "price": 200
                }
            ],
        }
    ]
}
```

### `GET /api/products/:id`

Response 200

```json
{
    "code": 200,
    "product": {
        "id": "5aa5c1a0-61d8-418e-8346-6edd9decb864",
        "name": "Product 1",
        "description": "Description 1",
        "price": 100.2,
        "imagePath": "import_files/7c/7c63329c7b3711ee802be0b9a548d6d8_f11face0469f11f0ad0f000c292ac68f.jpg", 
        "articleNumber": "123",
        "count": 10,
        "manufacturer": "Manufacturer 1",
        "parameters": [
            {
                "name": "Цвет",
                "type":"list",
                "stringValue": "Красный",
                "numberValue": 0,
                "sliceValue": null
            },
            {
                "name": "Сечение однопроволочного проводника по (кв.мм.)",
                "type": "number",
                "stringValue": "",
                "numberValue": 10,
                "sliceValue": null
            }
        ]
    }
}
```

Response 400, 404, 500

```json
{
    "code": 400,
    "message": "Error message"
}
```

### `GET /api/products/all/:page`

`page` - page number, must be greater or equal to 0

**Known issues:**

- Pictures currently no implemented, so `imagePath` should be changed

Response 200

```json
{
    "code": 200,
    "products": [
        {
            "id": "5aa5c1a0-61d8-418e-8346-6edd9decb864",
            "name": "Product 1",
            "description": "Description 1",
            "price": 100.2,
            "imagePath": "import_files/7c/7c63329c7b3711ee802be0b9a548d6d8_f11face0469f11f0ad0f000c292ac68f.jpg", 
            "articleNumber": "123",
            "count": 10,
            "manufacturer": "Manufacturer 1",
            "parameters": [
                {
                    "name": "Цвет",
                    "type":"list",
                    "stringValue": "Красный",
                    "numberValue": 0,
                    "sliceValue": null
                },
                {
                    "name": "Сечение однопроволочного проводника по (кв.мм.)",
                    "type": "number",
                    "stringValue": "",
                    "numberValue": 4,
                    "sliceValue": null
                },
            ]
        }
    ]
}
```

### `GET /api/files/:file-path`

`file-path` - path to file, must be relative to `DATA_DIR`

Response 200 - Binary file

### `GET /api/filters`

Returns a list of parameters that should used as filters

Response 200

```json
{
    "code": 200,
    "parameters": [
        {
            "name": "name",
            "type": "list",
            "values": [
                "val1", "val2"
            ]
        },
        {
            "name": "name",
            "type": "number",
            "minValue": 1,
            "maxValue": 10
        }
    ]
}
```

Response 500

```json
{
    "code": 500
}
```
