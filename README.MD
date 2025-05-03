var readmeContent = `
# Ticketing System API

This repository contains a RESTful API for event ticketing built with Go (Golang), Gin, GORM, MySQL, and JWT authentication, following a clean architecture structure.

## Base URL

All endpoints are prefixed with:

\`\`\`
http://localhost:8080/api
\`\`\`

## Setup

1. **Clone the repository**

   \`\`\`bash
   git clone https://github.com/yourusername/ticketing-api.git
   cd ticketing-api
   \`\`\`

2. **Environment Variables** Create a \`.env\` file in the root directory with the following:

   \`\`\`ini
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=password
   DB_NAME=ticketing
   SERVER_ADDRESS=:8080
   JWT_SECRET=your-secret-key
   JWT_EXPIRATION_HOURS=24
   \`\`\`

3. **Database** Make sure MySQL is running and a database named \`ticketing\` exists. The application will auto-migrate tables on startup.

4. **Run the server**

   \`\`\`bash
   go run main.go
   \`\`\`

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

## Postman Collection

Import \`postman_collection.json\` and use the following environment values:

* \`base_url\`: \`http://localhost:8080/api\`
* \`token\`: \`<JWT token from login>\`

### Example: Purchase Ticket

1. **Request**

   \`\`\`
   POST {{base_url}}/tickets
   Authorization: Bearer {{token}}
   Content-Type: application/json

   {
     "event_id": 1
   }
   \`\`\`

2. **Response**

   \`\`\`json
   HTTP/1.1 201 Created
   {
     "message": "Ticket purchased successfully",
     "ticket": {
       "id": 10,
       "event_id": 1,
       "user_id": 2,
       "status": "active",
       "price": 50.0,
       "purchased_at": "2025-05-02T10:00:00Z",
       "created_at": "2025-05-02T10:00:00Z",
       "updated_at": "2025-05-02T10:00:00Z"
     }
   }
   \`\`\`

## Testing

Use the provided Postman collection to test all endpoints. Ensure the \`token\` variable is set after login.

---

*Generated on May 2, 2025*
`
