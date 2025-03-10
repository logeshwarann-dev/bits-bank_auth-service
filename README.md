# 1️⃣ Auth Service (Authentication & User Management)

# Purpose: Handles user authentication, authorization, and session management.

# 🔹 Tech Stack:

# Backend: Golang
Database: PostgreSQL (Stores user details, hashed passwords, and session data)

Cache: Redis (For session store and rate limiting)

Auth Strategy: JWT-based authentication

# 🔹 Key Responsibilities:
✅ User Registration & Login (Email, Password, OTP-based login)

✅ Session Management (Token storage in Redis, auto-expiry)

✅ Role-Based Access Control (RBAC)

✅ Secure API Gateway Authentication (Protects microservices using JWT)

🔹 Endpoints Example:

POST `/signup` → Register a new user

POST `/login` → Authenticate user & generate JWT

GET `/get-user` → Fetch logged-in user details

POST `/logout` → Clear session

# 🔹 External Integrations:

API Gateway → To centralize authentication for all services

Stripe & Plaid Services → Verifies if the user is authenticated before processing requests
