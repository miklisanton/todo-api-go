# TODO API

#### Description
Simple API that allows you to create, read, update and delete tasks.
It also tracks if the task is overdue.

#### Endpoints
- GET /tasks
- POST /tasks
- PUT /tasks/{id}
- DELETE /tasks/{id}
- PATCH /tasks/{id}/complete

#### Usage
```bash
docker build -t todo-api .
docker run -v data:/build/data -p 8080:8080 todo-api

```

#### Testing
```bash
go test ./internal/db/repository -v
```
