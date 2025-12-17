# Go Auth API

A secure User Authentication Microservice built with **Go**, implementing **Clean Architecture**, and **best security practices**.

## 🚀 Overview

This project serves as a case study for building a secure, scalable authentication system. It features:

- **Clean Architecture**: Separation of concerns (Domain, Usecase, Adapter, Infra).
- **Security First**:
    - **RSA**: Asymmetric encryption for password transmission.
    - **Bcrypt**: Secure password hashing for storage.
    - **JWT**: JSON Web Tokens for stateless session management.
- **Dockerized**: Fully containerized environment with MySQL and Redis.

## 🛠️ Tech Stack

- **Language**: [Go](https://golang.org/) (1.23+)
- **Database**: [MySQL](https://www.mysql.com/) (User Data)
- **Cache**: [Redis](https://redis.io/) (Session/Token Management)
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose

## 🏗️ Architecture

The project follows the standard Go project layout and Clean Architecture principles:

```
├── cmd/api/            # Main entry point
├── config/             # Configuration config
├── internal/
│   ├── domain/         # Enterprise business rules (Entities)
│   ├── usecase/        # Application business rules
│   ├── repository/     # Interface definitions
│   └── infra/          # External interfaces (DB, HTTP, Security)
└── pkg/                # Public shared libraries
```

## ⚙️ Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)

## 🏃‍♂️ Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/user-auth-api.git
cd user-auth-api
```

### 2. Run with Docker
The easiest way to run the application is using Docker Compose:

```bash
docker-compose up -d --build
```
This will start:
- **MySQL** on port `3306`
- **Redis** on port `6379`
- **Auth API** on port `8080`

### 3. Usage

#### Encryption Key
First, retrieve the Public Key to encrypt user passwords (simulating a frontend client).
```bash
GET /public-key
```

#### Register User
**POST** `/register`
```json
{
  "email": "user@example.com",
  "encrypted_password": "<BASE64_RSA_ENCRYPTED_PASSWORD>"
}
```

#### Login
**POST** `/login`
```json
{
  "email": "user@example.com",
  "encrypted_password": "<BASE64_RSA_ENCRYPTED_PASSWORD>"
}
```
Response:
```json
{
  "token": "ey...<JWT_TOKEN>"
}
```

## 🔐 Security Details

1. **Password Transmission**: The client requests the `public-key` and encrypts the password before sending it to the server. This protects the password even if TLS were to be terminated or intercepted upstream (Layer of Defense).
2. **Storage**: The server decrypts the payload using its `private-key`, then hashes the password using **Bcrypt** before storing it in MySQL.
3. **Authentication**: On login, the server verifies the hash and issues a **JWT** signed with a secret key.

## 📄 License

[MIT](LICENSE)
