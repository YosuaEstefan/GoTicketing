# Go Ticketing API

GoTicketing is a robust, scalable RESTful API built with Go (Golang), Gin, GORM, MySQL, and JWT authentication. Designed for mid-to-large scale ticketing across industries — from entertainment and transportation to corporate services — it offers secure user management, role-based access control, event orchestration, ticketing workflows, and comprehensive reporting.This repository contains a RESTful API for event ticketing built with Go (Golang), Gin, GORM, MySQL, and JWT authentication, following a clean architecture structure.

---

## 🚀 Technology Stack

- **Backend:** Gin (Go)
- **ORM:** GORM
- **Database:** MySQL
- **Authentication:** JWT
- **Containerization:** Docker & Docker Compose
- **Go Version:** 1.20+

---

## ✨ Features

- **User Management:**
  - Registration, login, and role-based access control (Admin, Customer).
- **Event Management:**
  - CRUD operations on Event.
- **Ticket Management:**
  - CRUD operations on Ticket.
- **Reports and Analytics:**
  - Summary Revenue report by event or time period.
---

## 📁 Project Structure

```
go-ticketing/
├── config/
│   ├── config.go
│   └── database.go
├── controllers/
│   ├── event_controller.go
│   ├── report_controller.go
│   ├── ticket_controller.go
│   └── user_controller.go
├── middleware/
│   ├── cors.go
│   ├── jwt_auth.go
│   └── rbac.go
├── models/
│   ├── event.go
│   ├── ticket.go
│   └── user.go
├── repository/
│   ├── event_repository.go
│   ├── ticket_repository.go
│   └── user_repository.go
├── routers/
│   └── router.go
├── service/
│   ├── event_service.go
│   ├── jwt_service.go
│   ├── report_service.go
│   ├── ticket_service.go
│   └── user_service.go
├── utils/
│   └── jwt_util.go
├── .env
├── .gitignore
├── README.md
├── go.mod
└── go.sum
```

---

## 🗄️ Database Schema

- **One-to-Many:**
  - `Users → Ticket` 
  - `Event → Ticket`
- **One-to-One:**
  - `Ticket ↔ User`

---

## 🔄 Business Flow

Brief Business Logic Flow
Event Creation: Admin creates event → recorded in events table.
Ticket Purchase: User selects event → creates a new row in tickets with user_id, event_id, price, and purchased_at.
Event Status:
Before an update, the BeforeUpdate hook checks the current date.
If it has passed end_date → status becomes completed.
If it has passed start_date but not end_date → status is ongoing.
Ticket Status:
Default active.
Can be changed to used or cancelled.

---

## 📬 API Endpoints

> **Full documentation** in  [📑 Postman LMS API Docs](https://documenter.getpostman.com/view/39611346/2sB2j4grLy)

## Authentication

* **Register**: \`POST /api/register\`
* **Login**: \`POST /api/login\`

On login, you receive a JWT token. Include it in \`Authorization\` header for protected routes:

\`\`\`
Authorization: Bearer <token>
\`\`\`

## Endpoints

### Public

| Method | Endpoint        | Description                 |
| ------ | --------------- | --------------------------- |
| POST   | \`/api/register\` | Register a new user         |
| POST   | \`/api/login\`    | Login and receive JWT       |
| GET    | \`/api/events\`   | List all events (paginated) |

### Protected (User & Admin)

| Method | Endpoint           | Description                 |
| ------ | ------------------ | --------------------------- |
| GET    | \`/api/tickets\`     | Get tickets of current user |
| POST   | \`/api/tickets\`     | Purchase a ticket           |
| GET    | \`/api/tickets/:id\` | Get ticket by ID            |
| PATCH  | \`/api/tickets/:id\` | Cancel ticket by ID         |

### Admin Only

| Method | Endpoint                 | Description                    |
| ------ | ------------------------ | ------------------------------ |
| POST   | \`/api/events\`            | Create a new event             |
| PUT    | \`/api/events/:id\`        | Update an existing event       |
| DELETE | \`/api/events/:id\`        | Delete an event                |
| GET    | \`/api/reports/summary\`   | Get overall sales summary      |
| GET    | \`/api/reports/event/:id\` | Get sales report for one event |

   ```
4. **Access the API** at `http://localhost:8080`

## ✅ Run the server**
   \`\`\`bash
   go run main.go
   \`\`\`


## 🔒 Security

- Passwords are hashed with **bcrypt**.
- **JWT** tokens are used for authentication.
- Role-based access control enforced on all endpoints.
.

