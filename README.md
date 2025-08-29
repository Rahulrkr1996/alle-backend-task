# 📋 Task Management Service

A simple **Task Management System** built with **Go** following microservices principles.  
This service allows users to **create, read, update, and delete tasks**.  
It also supports **pagination** and **filtering by status**.

---

## 🚀 Features

- **CRUD APIs** for tasks
- **Pagination** for `GET /tasks`
- **Filtering** tasks by status (`Pending`, `InProgress`, `Completed`)
- **In-memory repository** (data is not persisted across restarts)
- **Separation of concerns**: repository, service, handler layers
- **Environment variable based config**
- **Optional seed data** for development

---

## 📂 Project Structure


├── main.go       # Entry point  
├── model.go      # Task model + repository interface  
├── repo.go       # In-memory repository implementation  
├── service.go    # Business logic layer  
├── handler.go    # HTTP handlers + routes  
└── README.md     # This file


---

## ⚙️ Setup & Run

### 1. Clone & Init
```bash
git clone https://github.com/Rahulrkr96/tasksvc.git
cd tasksvc
go mod tidy
2. Run Service
go run .

3. With Seed Data
SEED_DATA=true go run .

4. Change Port
PORT=9090 go run .
```

📡 API Endpoints
```bash
Create Task
POST /tasks
Content-Type: application/json

{
  "title": "Write docs",
  "description": "Complete the README",
  "status": "Pending"
}


✅ Response 201 Created:

{
  "id": "1",
  "title": "Write docs",
  "description": "Complete the README",
  "status": "Pending"
}

Get All Tasks (with Pagination & Filtering)
GET /tasks?page=1&limit=2&status=Pending


✅ Response:

{
  "tasks": [
    {
      "id": "1",
      "title": "Write docs",
      "description": "Complete the README",
      "status": "Pending"
    }
  ],
  "page": 1,
  "limit": 2,
  "total": 1
}

Get Task by ID
GET /tasks/1

Update Task
PUT /tasks/1
Content-Type: application/json

{
  "title": "Write docs",
  "description": "Update the README file",
  "status": "InProgress"
}

Delete Task
DELETE /tasks/1
```
## 🏗 Design Decisions

### 1. Microservices Style
- Clear separation into **Repository → Service → Handler**
- Easy to swap **in-memory DB** with a persistent store (Postgres, MongoDB, etc.)
- Each layer follows the **Single Responsibility Principle**

---

### 2. Scalability
- The service is **stateless** → can be scaled **horizontally**
- Multiple instances can run behind a **load balancer**
- A persistent DB (when added) will act as the **single source of truth**

---

### 3. Inter-Service Communication
- Current service exposes **REST APIs**
- Future services (e.g., **User Service**) can communicate via:
    - **REST** (simple & human-friendly)
    - **gRPC** (high performance, strongly typed contracts)
    - **Message queues** (Kafka, RabbitMQ) for **asynchronous workflows**

---

## 🧪 Testing APIs

```bash
Using curl:

# Create Task
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{"title":"Test","description":"First task","status":"Pending"}'

# List Tasks
curl "http://localhost:8080/tasks?page=1&limit=5"

# Filter Tasks
curl "http://localhost:8080/tasks?status=Completed"

# Get by ID
curl http://localhost:8080/tasks/1

# Update Task
curl -X PUT http://localhost:8080/tasks/1 \
-H "Content-Type: application/json" \
-d '{"title":"Updated","description":"Changed desc","status":"InProgress"}'

# Delete Task
curl -X DELETE http://localhost:8080/tasks/1
```