# GinRESTfulAPIDemo
Golang Gin RESTful API Demo

### login

```
curl -X POST http://127.0.0.1:8082/auth/login \
     -H "Content-Type: application/json" \
     -d '{"nickname": "admin", "password": "your_password"}'
```

### get users

```
curl -X GET 'http://localhost:8082/api/users?page=1&limit=10&search=&status=' \
     -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```
