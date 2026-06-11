---
globs: ["*.go", "**/*.go"]
description: "Go database operations - GORM, migrations, transactions, N+1 prevention"
---

# GO_DATABASE.md - Go Database Operations

This file defines database operation rules for Go projects (primarily GORM, but principles apply to other ORMs).
These are high-level rules and do not cover deep details.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) Migration Rules (hard rules)
- 2) Context and Timeouts (general)
- 3) Transaction Guidelines
- 4) N+1 Query Prevention (mandatory)
- 5) Indexing Guidelines
- 6) GORM Usage Rules
- 7) Soft Deletes
- 8) Connection Pool Settings
- 9) When Uncertain (mandatory)

## 1) Migration Rules (hard rules)

### Irreversible operations require confirmation

If a migration includes any of the following, stop and ask:
- Dropping a column or table
- Changing a column type (data loss risk)
- Dropping an index (performance risk)
- Any operation that may cause data loss

### Migration execution principles

- You do not have to run migrations in Go; use existing tooling if the service already has a process
- If using go-database built-in migrate, confirm the workflow and risk controls first

### Migration file rules

- Use meaningful names (timestamp + description)
- Each migration must be independent (can run on its own)
- Provide Up and Down (rollback)

```go
// Correct: clear migration
// 20240111_add_user_status_column.go

func Up() error {
  return db.Exec(`
    ALTER TABLE users
    ADD COLUMN status VARCHAR(20) DEFAULT 'active'
  `).Error
}

func Down() error {
  return db.Exec(`
    ALTER TABLE users
    DROP COLUMN status
  `).Error
}
```

---

## 2) Context and Timeouts (general)

- All DB operations must include `context.Context` (except offline data maintenance)
- Always set timeouts to avoid long locks or blocking

---

## 3) Transaction Guidelines

### When to use transactions

You must use transactions for:
- Multiple related writes (atomicity)
- Read-then-write flows (consistency)
- Money/balance/order operations (always)

### Transaction implementation principles
- Do not enforce a single style; follow existing project patterns
- Provide guidance and risk notes without locking into one framework style

```go
// Correct: use transaction
func (r *OrderRepository) Create(ctx context.Context, order *Order) error {
  return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // Decrease inventory
    if err := tx.Model(&Product{}).
      Where("id = ? AND stock >= ?", order.ProductID, order.Quantity).
      Update("stock", gorm.Expr("stock - ?", order.Quantity)).Error; err != nil {
      return err
    }

    // Create order
    if err := tx.Create(order).Error; err != nil {
      return err
    }

    return nil
  })
}
```

### Transaction timeouts

- Keep transactions short
- Use `context.WithTimeout` to avoid long locks

---

## 4) N+1 Query Prevention (mandatory)

### Strategies to avoid N+1

- Use `JOIN` or `Preload` depending on readability and query needs
- No single enforced style; the goal is to prevent N+1 behavior

```go
// Incorrect: N+1 query
users, _ := db.Find(&users).Error
for _, user := range users {
  // Each user triggers another query for orders (N queries)
  db.Where("user_id = ?", user.ID).Find(&user.Orders)
}

// Correct: use Preload
db.Preload("Orders").Find(&users)
```

### Check query counts

- During development, watch SQL query counts in logs
- Use `db.Debug()` when you suspect N+1

---

## 5) Indexing Guidelines

### When indexes are needed

- Frequently used WHERE columns
- JOIN columns
- ORDER BY columns
- Unique constraint columns

### Composite indexes

- Order columns by query frequency and selectivity
- Put the most used column first

```go
// Correct: composite index
type User struct {
  gorm.Model
  Email  string `gorm:"index:idx_email_status"`
  Status string `gorm:"index:idx_email_status"`
}
```

---

## 6) GORM Usage Rules

### Context passing

- All DB operations must pass context

```go
// Correct
db.WithContext(ctx).Where("id = ?", id).First(&user)

// Incorrect: no context
db.Where("id = ?", id).First(&user)
```

### Error handling

```go
// Correct: check specific errors
result := db.WithContext(ctx).First(&user, id)
if result.Error != nil {
  if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    return nil, ErrUserNotFound
  }
  return nil, result.Error
}
```

---

## 7) Soft Deletes

- Use `gorm.DeletedAt` for soft deletes
- Be explicit about whether deleted records are included

```go
type User struct {
  gorm.Model // includes DeletedAt
}

// Default: exclude deleted
db.Find(&users)

// Include deleted
db.Unscoped().Find(&users)
```

---

## 8) Connection Pool Settings

```go
sqlDB, err := db.DB()

// Configure pool
sqlDB.SetMaxOpenConns(100)         // max open connections
sqlDB.SetMaxIdleConns(10)          // max idle connections
sqlDB.SetConnMaxLifetime(time.Hour) // max connection lifetime
```

---

## 9) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Whether a migration causes data loss
- Whether a transaction is needed
- How to optimize queries (N+1 or performance)
- Index strategy
- Rollback plan

---

Violating these rules is incorrect output.
