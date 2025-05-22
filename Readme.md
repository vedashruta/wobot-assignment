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
