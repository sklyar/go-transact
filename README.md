# Go Transaction Manager

[![pkg-img]][pkg-url]
[![version-img]][version-url]
[![license-img]][license-url]

Go Transaction Manager is a library for managing SQL transactions in Go applications, offering a simple interface for controlling both isolated and nested transactions. Through context-based transaction management, it ensures atomicity and ease of use.

## Features

- **Isolated Transactions**: Facilitates the management of isolated transactions across distinct parts of your application, ensuring data integrity.

- **Nested Transactions**: Enables nested transactions for more complex transactional operations. This allows transactions within transactions with individual commit or rollback controls.

- **Context-Aware Transactions**: Maintains and provides transactional context information, empowering your code with awareness of the current transaction state.

- **Manual and Automated Control**: Offers both manual transaction control via the `Begin` function, allowing commit or rollback on demand, and automated control through `BeginFunc`, which automates the commit or rollback based on the function's return value.

## Installation

To install the `go-transact` library, run the following command:

```
go get github.com/sklyar/go-transact
```
## Supported Adapters

`go-transact` is designed to be extensible and supports various database adapters. Here is a list of currently supported ones:

| Adapter                                                               | Description                                                                                                                                |
|-----------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| **Standard library adapter ([transactstd](./adapters/transactstd/))** | The standard SQL adapter provides an easy way to integrate `go-transact` with any database that conforms to Go's `database/sql` interface. |

## Usage

### Simple Initialization

Here's how to get started with `go-transact`:

```go
import (
    ...
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/sklyar/go-transact"
    "github.com/sklyar/go-transact/adapters/transactstd"
)

func main() {
    sqlDB, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
    checkErr(err)
    
    txManager, db, err := transact.NewManager(transactstd.Wrap(sqlDB))
    checkErr(err)
}
```

In the example above, a PostgreSQL database is initialized and wrapped using the provided standard SQL adapter. A new transaction manager is then created using this adapter.

### Usage Examples

Below are examples demonstrating the usage of `go-transact`.

#### Using BeginFunc for Automatic Transaction Control

```go
err = txManager.BeginFunc(ctx, func(ctx context.Context) error {
    // execute your transactional operations
    
    return nil // return an error to trigger rollback
})
checkErr(err)
```

In this example, `BeginFunc` manages the transaction life cycle. It automatically starts the transaction and commits it if no errors occur during the execution of the passed function. If an error is returned, `BeginFunc` triggers a rollback.

#### Using Begin for Manual Transaction Control

```go
ctx, tx, err := txManager.Begin(ctx)
checkErr(err)

defer tx.Rollback(ctx) // rollback if commit is not called

// execute your transactional operations
err = doSomething(ctx)
checkErr(err)

err = tx.Commit(ctx)
checkErr(err)
```

`Begin` provides more control over transactions. It starts a transaction and returns a transaction object (`tx`), which can be manually committed or rolled back.

#### Nesting Transactions

Both `Begin` and `BeginFunc` support nested transactions, allowing each transaction to be individually controlled:

```go
txManager.BeginFunc(ctx, func(ctx context.Context) error {
  // execute your transactional operations

  err = txManager.BeginFunc(ctx, func(ctx context.Context) error {
    // execute your nested transactional operations
    return nil // return an error to trigger rollback of this nested transaction
  })
  checkErr(err)

  return nil // return an error to trigger rollback of the entire transaction
})
```

In this example, nested transactions are created using `BeginFunc`. If an error is returned from the inner `BeginFunc`, it triggers a rollback of the nested transaction. If an error is returned from the outer `BeginFunc`, it triggers a rollback of the entire transaction.

## License
Go Transaction Manager is released under the MIT License. See the bundled LICENSE file for details.

[pkg-img]: https://pkg.go.dev/badge/sklyar/go-transact
[pkg-url]: https://pkg.go.dev/github.com/sklyar/go-transact
[version-img]: https://img.shields.io/github/v/release/sklyar/go-transact
[version-url]: https://github.com/sklyar/go-transact/releases
[license-img]: https://img.shields.io/github/license/sklyar/go-transact
[license-url]: https://raw.githubusercontent.com/sklyar/go-transact/master/LICENSE
