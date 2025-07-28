# Go CRUD

A Golang project for CRUD demonstration using Gorm and Gin.

## Architecture

This project is built using a Clean Architecture approach, which emphasizes a separation of concerns. This is achieved by dividing
the application into several layers:

1.  **Domain Layer**: This is the core of the application. It contains the business logic and entities. It has no dependencies on any
other layer.
2.  **Application Layer**: This layer orchestrates the use cases of the application. It depends on the domain layer but not on the
infrastructure layer. It defines interfaces that are implemented by the infrastructure layer.
3.  **Infrastructure Layer**: This layer is responsible for external concerns such as databases, web frameworks, and other external
services. It depends on the application layer to implement its interfaces.

## Getting Started

### Prerequisites

- **Go**: Make sure you have Go installed on your system. You can download it from the [official Go website](https://golang.org/dl/).
- **Docker**: Make sure you have Docker installed on your system. You can download it from the [official Docker website](
https://www.docker.com/get-started).

### Running with Docker

Start all containers by running the following command

```bash
docker compose up -d
```

then access the Swagger document via http://localhost:8080/swagger/index.html

## API Endpoints

The following are the main API endpoints available:

### Accounts

-   **GET /accounts**: Get a list of all accounts.
-   **GET /accounts/{id}**: Get a single account by its ID.
-   **POST /accounts**: Create a new account.
-   **PUT /accounts/{id}**: Update an existing account.
-   **DELETE /accounts/{id}**: Delete an account.

### Transactions

-   **GET /transactions**: Get a list of all transactions.
-   **GET /transactions/{id}**: Get a single transaction by its ID.
-   **POST /transactions**: Create a new transaction.
-   **POST /transactions/{id}/process**: Process a transaction.
-   **DELETE /transactions/{id}**: Cancel a transaction.