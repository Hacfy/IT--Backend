# 🖥️ IT Inventory Management – Backend

A serverless-ready backend system for managing IT inventory in organizations. Built with **Go**, **Echo**, and **PostgreSQL**, this system allows organizations to manage branches, departments, components, and issues across different access levels.

---

## 📚 Purpose

This backend solves the problem of managing IT infrastructure within organizations:

* 📦 Track components and their units
* 🛠️ Monitor maintenance costs
* 📍 Know where each unit is placed (workspaces)
* ⚠️ Raise issues when problems occur
* 🔒 Role-based access control across the organization hierarchy

---

## 🛠️ Tech Stack

* **Go** (Golang)
* **Echo** (Web Framework)
* **PostgreSQL**
* **JWT** (Authentication)
* **Planned Serverless Deployment** via AWS Lambda / Azure Functions

---

## 🧩 Project Structure

```plaintext
├── cmd/                # Application entry point
├── internals/
│   ├── handlers/       # API route handlers
│   ├── middleware/     # JWT, logging, etc.
│   └── models/         # Request/response structs
├── pkg/
│   ├── database/       # DB connection logic & queries
│   ├── templates/      # Response templates (if any)
│   └── utils/          # Helper functions (JWT, validation)
└── repository/         # Login & authentication logic
```

---

## 👤 Roles and Permissions

### 🛠️ Main Admin

* Create and delete organizations
* Create other Main Admins
* View all data of organizations they created

### 🏢 Organization

* Represents a single client organization
* Can create and delete Super Admins
* View and manage all data within their organization

### 👑 Super Admin

* Can create and delete branches
* Assign Branch Heads to specific branches
* View and manage data only for branches they created

### 🏬 Branch Head

* Manages a specific branch
* Can create departments and warehouses
* View data of departments and warehouses they created

### 🧑‍🏫 Department Head

* Responsible for managing a department
* Can create workspaces within their department
* Can raise issues for units in their department's workspaces

### 🧰 Warehouse

* Represents the technical staff for a branch
* Can create components and their units
* Can resolve issues raised by department heads
* Maintain and delete components/units they created

---

## 🚀 Key Features

* 🔐 Secure login/authentication with JWT
* 🏢 Role-based organization and access
* 🧰 Unit assignment and tracking
* 📄 Maintenance tracking per component/unit
* 📍 Workspace management
* ⚠️ Issue reporting and resolution
* 📊 Clean RESTful API structure

---

## 📦 Setup Instructions

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/it-inventory-backend.git
cd it-inventory-backend
```

### 2. Create `.env`

```env
DB_URL=postgres://user:password@localhost:5432/yourdb?sslmode=disable
PORT=8080
JWT_SECRET=your_jwt_secret
```

### 3. Install dependencies and run

```bash
go mod tidy
go run cmd/main.go
```

---

## 📬 API Example (Raise Issue)

```http
POST /issues
Authorization: Bearer <token>
Content-Type: application/json

{
  "department_id": 2,
  "workspace_id": 8,
  "unit_id": 17,
  "unit_prefix": "PC",
  "issue": "Monitor not turning on"
}
```

---

## 📌 Notes

* Backend-only; frontend will be developed separately.
* This project will be containerized or deployed serverlessly later.
* Use Swagger or Postman collections for testing.

---

## 📃 License

MIT License

---

## ✍️ Author

**Ashith Kumar Gowda**
[LinkedIn](https://www.linkedin.com/in/ashith-kumar-gowda-446685297) | [GitHub](https://github.com/isAshithKumarGowda)
