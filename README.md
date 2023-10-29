# Auth Service
Protect user identities with multi-factor authentication options, empowering you with the flexibility to choose the most suitable method. Paired with Casbin authorization, our service guarantees robust access control, granting only the necessary permissions. Elevate your security standards with our comprehensive authentication solution.

### Benchmark
[Benchmark Result](/cmd/benchmark/benchmark_result.md)

### Generate Repository

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run test:
    1. Install gomock
    ```bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

    2. Generate mock
    ```bash
    make mock
    ```

    3. Run test
    ```bash
    make test
    ```

## Endpoints

### **POST /otp**

One-Time-Password. Will deliver a sms otp to the user depending on whether the request body contains an "phone" key.

- Env key: `OTP_EXPIRED_DURATION`
- Expire duration: 5m

```json
{
    "phone": "+84909000999" // has country code
}

```

Response:
```json
{}
```

### **POST /token**

This is an OAuth2 endpoint that currently implements
the `otp` grant types

Auto create new user if not exist

Query params:

```
?grant_type=otp
```

Body:

```json
{
  "phone": "+84909000999",
  "otp": "1234"
}
```

Response:

```json
{
    "meta": null,
    "data": {
        "id": "a3450777-9e18-495f-945c-cd3663203f95",
        "createdAt": "2023-10-28T17:51:51.08104+07:00",
        "updatedAt": "2023-10-28T17:51:51.08104+07:00",
        "fullName": "",
        "phone": "+84918919314",
        "country": "",
        "email": null,
        "birthday": null,
        "avatar": null,
        "address": null,
        "gender": null,
        "status": 1,
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2MjIzYzExLWRiMjYtNGVkOC1iOTBlLTI0MmIwOGY4ZjVmYyIsImlhdCI6IjIwMjMtMTAtMjlUMTU6MTk6MTIuNDM4MzU2NTk4KzA3OjAwIiwiZXhwIjoiMjAyMy0xMC0zMFQxNToxOToxMi40MzgzNTY2NTQrMDc6MDAiLCJkYXRhIjp7InBob25lIjoiKzg0OTE4OTE5MzE0IiwidXNlcklkIjoiYTM0NTA3NzctOWUxOC00OTVmLTk0NWMtY2QzNjYzMjAzZjk1In19.I8cCtOQBhGvjyohLPoMjcoq_1TieTUGU8jrb3hvuXLM",
        "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2MjIzYzExLWRiMjYtNGVkOC1iOTBlLTI0MmIwOGY4ZjVmYyIsImlhdCI6IjIwMjMtMTAtMjlUMTU6MTk6MTIuNDM4NDQ3OTg2KzA3OjAwIiwiZXhwIjoiMjAyMy0xMC0zMVQxNToxOToxMi40Mzg0NDgxMDMrMDc6MDAiLCJkYXRhIjp7InBob25lIjoiKzg0OTE4OTE5MzE0IiwidXNlcklkIjoiYTM0NTA3NzctOWUxOC00OTVmLTk0NWMtY2QzNjYzMjAzZjk1In19.3PZOKdOdSY7H4QnI8gdbSzmJcv-g0LHUS0IwiNJecu4"
    }
}
```

### **POST /token/refresh**

Body:

```json
{
  "refreshToken": "refresh-token"
}
```

Once you have an access token, you can access the methods requiring authentication
by settings the `Authorization: Bearer YOUR_ACCESS_TOKEN_HERE` header.

Response:

```json
{
    "meta": null,
    "data": {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImUxMjA3ZTJjLWQ0YzgtNDkwMi1hOTFiLWZiYmRhNzJhNDMxOSIsImlhdCI6IjIwMjMtMTAtMjlUMTc6Mjk6NTUuNDA1NTI4NDIrMDc6MDAiLCJleHAiOiIyMDIzLTEwLTMwVDE3OjI5OjU1LjQwNTUyODQ2KzA3OjAwIiwiZGF0YSI6eyJwaG9uZSI6Iis4NDkxODkxOTMxNCIsInVzZXJJZCI6ImEzNDUwNzc3LTllMTgtNDk1Zi05NDVjLWNkMzY2MzIwM2Y5NSJ9fQ.2f-KJj15lOxMcxseK6nHfAXbTZx3Q_UGzj25lH-DtGs",
        "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImUxMjA3ZTJjLWQ0YzgtNDkwMi1hOTFiLWZiYmRhNzJhNDMxOSIsImlhdCI6IjIwMjMtMTAtMjlUMTc6Mjk6NTUuNDA1NTY4MjE0KzA3OjAwIiwiZXhwIjoiMjAyMy0xMC0zMVQxNzoyOTo1NS40MDU1NjgyNTkrMDc6MDAiLCJkYXRhIjp7InBob25lIjoiKzg0OTE4OTE5MzE0IiwidXNlcklkIjoiYTM0NTA3NzctOWUxOC00OTVmLTk0NWMtY2QzNjYzMjAzZjk1In19.ufCtXhm2_di_Kny0DAxJ88ocr4oOg0L4_91VyFBeu60"
    }
}
```

### **GET /token/validate**

Validate access token - API Gateway can call to check token.
Logic:
1. Verify token with secret key
2. Check token exists in blacklist in Redis

Query Params:
- `token` the token that want to validate


```json
/token/validate?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2MjIzYzExLWRiMjYtNGVkOC1iOTBlLTI0MmIwOGY4ZjVmYyIsImlhdCI6IjIwMjMtMTAtMjlUMTU6MTk6MTIuNDM4MzU2NTk4KzA3OjAwIiwiZXhwIjoiMjAyMy0xMC0zMFQxNToxOToxMi40MzgzNTY2NTQrMDc6MDAiLCJkYXRhIjp7InBob25lIjoiKzg0OTE4OTE5MzE0IiwidXNlcklkIjoiYTM0NTA3NzctOWUxOC00OTVmLTk0NWMtY2QzNjYzMjAzZjk1In19.I8cCtOQBhGvjyohLPoMjcoq_1TieTUGU8jrb3hvuXLM
```

Response:

```json
{
    "meta": null,
    "data": true
}
```

### **GET /user**

***(Requires authentication)***

Get the JSON object for the logged in user information

### **PUT /user**

***(Requires authentication)***

Update a user information

