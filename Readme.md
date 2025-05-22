# 📁 FileNest

This project is a RESTful File Storage API built with **Golang**
It provides user authentication with JWT, secure file upload with storage quotas, and endpoints to retrieve uploaded files and remaining storage.

---

## 📌 Features

- ✅ User Registration & Login with JWT Authentication
- ✅ Secure file upload with storage size checks
- ✅ Per-user folder structure
- ✅ File metadata tracking (name, size, timestamp)
- ✅ Get remaining user storage
- ✅ Retrieve uploaded files
- 🔐 Password hashing using bcrypt

---

## 🛠 Tech Stack

- Go (Fiber web framework)
- JWT for user authentication
- bcrypt for password hashing
- Local file system for storage
- MongoDB for persistence

---

## 🧩 File Storage Structure
/storage/</br>
└── username/</br>
    &nbsp;├── id1</br>
    &nbsp;├── id2</br>

## Arcitecture
Source: [draw.io](https://pkg.go.dev/github.com/goccy/go-graphviz#section-readme)

## 🚀 How to Run the Server

1. Clone the repository:
```bash
git clone https://github.com/vedashruta/wobot-assignment
```
2. Change the working directory
```bash
cd wobot-assignment
```
3. Install dependencies
```bash
go mod tidy
```
4. Run the server
```bash
go run main.go
```

Building the binary executable (optional)
```bash
go build -o server && ./server
```