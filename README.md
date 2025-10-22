# 🐱 About API

A simple Go/Chi-based API that analyzes a string and returns information about the string.
This project demonstrates a clean and minimal structure for building small REST APIs in Go.

---

## 🚀 Features

- Built using [Chi](https://go-chi.io/) — a fast, lightweight Go web framework.
- Automatically loads environment variables using [godotenv](https://github.com/joho/godotenv).

---

## 🧱 Project Structure

```bash
# 🐱 About API

A simple Go/Chi-based API that analyzes a string and returns information about the string.
This project demonstrates a clean and minimal structure for building small REST APIs in Go.

---

## 🚀 Features

- Built using [Gin](https://github.com/gin-gonic/gin) — a fast, lightweight Go web framework.
- Automatically loads environment variables using [godotenv](https://github.com/joho/godotenv).

---

## 🧱 Project Structure

```bash
├── cmd
│   └── stringAnalyzer
│       └── main.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── application
│   │   └── application.go
│   ├── config
│   │   └── config.go
│   ├── database
│   │   ├── database.go
│   │   ├── migrations
│   │   │   ├── 001_setup.sql
│   │   │   └── 002_alter_table.sql
│   │   └── migrator.go
│   ├── dto
│   │   └── dto.go
│   ├── errs
│   │   └── error.go
│   ├── handler
│   │   └── string.go
│   ├── logger
│   │   └── logger.go
│   ├── middleware
│   ├── model
│   │   └── string.go
│   ├── repository
│   │   └── repository.go
│   ├── routes
│   │   └── string.go
│   ├── server
│   │   └── server.go
│   └── util
│       └── util.go
├── README.md
├── Taskfile.yml
└── tmp
    ├── build-errors.log
    └── main
```


## ⚙️ Setup Instructions

### 1. Clone the Repository

```bash
  git clone https://github.com/justinndidit/stringAnalyzer.git
  cd stringAnalyzer
```

### 2. Install Dependencies

Make sure you have Go 1.21+ installed, then run:

```bash
  go mod tidy
```

### 3. 🌍 Environment Variables

```bash
  cp .env.sample .env
```

### 4. ▶️ Run Locally

Using go run:

```bash
  task run:dev
```
### Ensure you have setup the .env properly and have a running postgres server

