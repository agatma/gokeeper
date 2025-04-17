# GOPHKEEPER 🔒

Secure client-server password manager with end-to-end encryption

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?logo=postgresql)
![AES-256](https://img.shields.io/badge/AES-256_GCM-5C3EE6?logo=openssl)

## Features ✨

**Server**:
- AES-256-GCM encryption with random nonces
- JWT authentication with refresh tokens
- CRUD operations for secrets
- Paginated data retrieval (limit/offset)
- PostgreSQL storage with Docker support
- Automatic data versioning

**Client**:
- Interactive CLI interface
- Support for multiple data types:
    - 🔑 Login/Password pairs
    - 💳 Credit card information
    - 📝 Text notes
    - 📁 Binary files
- Local storage fallback mode
- Bulk upload capability
- Cross-platform compatibility

## Installation 📦

### Prerequisites
- Go 1.20+
- Docker 20.10+
- PostgreSQL 15+
