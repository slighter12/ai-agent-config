---
globs: ["*.go", "**/*.go"]
description: "Go architecture principles - layering, boundaries, folder responsibilities"
---

# GO_ARCHITECTURE.md - Go Architecture Principles

This file defines architecture and layering principles for Go projects.
These are principle-based rules and do not force a specific folder structure.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) Layering Principles (mandatory)
- 1) Boundary Rules (hard rules)
- 1) Transport Layer Responsibilities (mandatory)
- 1) Service/UseCase Layer Responsibilities (mandatory)
- 1) Repository/DAO Layer Responsibilities (mandatory)
- 1) Dependency Injection (mandatory)
- 1) Package Organization Principles
- 1) Interface Definition Location
- 1) Folder Structure Flexibility
- 1) When Uncertain (mandatory)

## 1) Layering Principles (mandatory)

### Keep layer boundaries clear

Regardless of architecture style (Clean Architecture, DDD, MVC), you must:

- Keep clear boundaries between layers
- Maintain a consistent dependency direction (inner layers do not depend on outer layers)
- Give each layer a clear responsibility

### Common layering patterns

Pattern A: Clean Architecture

```
delivery (HTTP/gRPC) -> usecase (business logic) -> repository (data access)
```

Pattern B: DDD

```
interface -> application -> domain -> infrastructure
```

Pattern C: Traditional 3-layer

```
handler/controller -> service -> dao/repository
```

Choose one pattern and keep it consistent.
When applying error-handling rules, map your project layers to the conceptual layers defined in `GO_ERROR_HANDLING.md`.

---

## 2) Boundary Rules (hard rules)

### Forbidden dependency directions

Do not:

- Put business logic in the Transport layer (HTTP handlers)
- Access the database directly from the Transport layer
- Have Service/UseCase depend on Transport
- Have Repository/DAO depend on Transport

Allowed:

- Transport -> Service -> Repository (one-way)
- Use interfaces to invert dependencies

### Cyclic dependencies

- Cyclic dependencies between packages are forbidden
- Use interfaces to decouple
- Refactor into independent packages when needed

---

## 3) Transport Layer Responsibilities (mandatory)

Transport layer (HTTP handlers, gRPC services) can only:

- Parse requests
- Validate input format (not business validation)
- Call Service/UseCase
- Map responses
- Handle HTTP-specific logic (status codes, headers)

Transport layer must not:

- Contain business rules
- Access the database directly
- Perform complex computations
- Contain long-running logic (delegate to Service)
- Return ad-hoc error responses (use the project's existing unified API error format, or optional `policy-api` detail when the external contract changes)

```go
// Correct
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
  var req CreateUserRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    handleError(w, NewBadRequestError("invalid request"))
    return
  }

  user, err := h.userService.Create(r.Context(), req.Email, req.Name)
  if err != nil {
    handleError(w, err)
    return
  }

  json.NewEncoder(w).Encode(user)
}
```

`handleError` and `NewBadRequestError` are illustrative placeholders. Use the project's existing error helpers for the actual response envelope and error mapping. Use optional `policy-api` detail when status codes, response schemas, or error-code contracts are being defined or changed.

```go
// Incorrect: handler contains business logic and DB access
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
  // ... parse request ...

  // Business logic should not be in handler
  if !isValidEmail(req.Email) {
    http.Error(w, "Invalid email", 400)
    return
  }

  // Do not access DB directly
  user := &User{Email: req.Email, Name: req.Name}
  if err := h.db.Create(user).Error; err != nil {
    http.Error(w, "Database error", 500)
    return
  }

  json.NewEncoder(w).Encode(user)
}
```

---

## 4) Service/UseCase Layer Responsibilities (mandatory)

Service layer should contain:

- Business logic
- Orchestration
- Domain rule validation
- Transaction boundaries
- Cross-repository operations

Service layer should not:

- Contain HTTP-specific logic (status codes, headers)
- Handle HTTP request/response directly
- Contain SQL queries (use Repository)

```go
// Correct
type UserService struct {
  userRepo UserRepository
  emailSvc EmailService
}

func (s *UserService) Create(ctx context.Context, email, name string) (*User, error) {
  // Business validation
  if !isValidEmail(email) {
    return nil, ErrInvalidEmail
  }

  // Check existing
  existing, err := s.userRepo.FindByEmail(ctx, email)
  if err != nil && !errors.Is(err, ErrNotFound) {
    return nil, err
  }
  if existing != nil {
    return nil, ErrEmailAlreadyExists
  }

  // Create user
  user := &User{
    Email: email,
    Name:  name,
  }

  if err := s.userRepo.Create(ctx, user); err != nil {
    return nil, err
  }

  // Send welcome email
  _ = s.emailSvc.SendWelcome(ctx, user.Email)

  return user, nil
}
```

---

## 5) Repository/DAO Layer Responsibilities (mandatory)

Repository layer can only:

- Persist data
- Query data
- Manage transactions (coordinated with Service)

Repository layer must not:

- Contain business logic
- Depend on Service or Handler layers
- Handle HTTP requests

---

## 6) Dependency Injection (mandatory)

- Prefer constructor injection
- Avoid global variables (except constants or logger)
- Depend on interfaces

```go
// Correct: interface definition
type UserRepository interface {
  FindByID(ctx context.Context, id string) (*User, error)
  Create(ctx context.Context, user *User) error
}

// Correct: constructor injection
type UserService struct {
  repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
  return &UserService{repo: repo}
}

// Incorrect: global variable
var globalDB *gorm.DB

func DoSomething() {
  globalDB.Create(&user) // Hard to test and trace dependencies
}
```

---

## 7) Package Organization Principles

### Feature-based grouping (recommended)

```
user/
|-- handler.go     # HTTP handlers
|-- service.go     # business logic
|-- repository.go  # data access
`-- model.go       # domain models
```

### Layer-based grouping (also acceptable)

```
handler/
`-- user.go
service/
`-- user.go
repository/
`-- user.go
```

Key point: choose one and keep it consistent.

---

## 8) Interface Definition Location

- Define interfaces on the consumer side
- Do not define interfaces in repository packages and import them into service
- Make dependencies point to abstractions

```go
// Correct: service defines the interface it needs
package service

type UserRepository interface {
  FindByID(ctx context.Context, id string) (*User, error)
}

type UserService struct {
  repo UserRepository // depend on abstraction
}

// repository implements the interface (but does not import service)
package repository

type UserRepo struct { /* ... */ }

func (r *UserRepo) FindByID(ctx context.Context, id string) (*User, error) {
  // implementation
}
```

---

## 9) Folder Structure Flexibility

This rule does not enforce a specific folder layout, but you must:

- Maintain clear layering
- Separate responsibilities
- Keep dependency direction consistent
- Document architectural decisions

---

## 10) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Which layer should hold a piece of logic
- Which layering pattern the project uses
- How to break a cyclic dependency
- Where to define interfaces

---

Violating these rules is incorrect output.
