# Complete Guide to Context in Golang

A comprehensive guide to understanding and effectively using Context in Go programming with visual diagrams and practical examples.

## 📋 Table of Contents

1. [What is Context?](#1-what-is-context)
2. [Creating a Context](#2-creating-a-context)
3. [Propagating Context](#3-propagating-context)
4. [Retrieving Values from Context](#4-retrieving-values-from-context)
5. [Cancelling Context](#5-cancelling-context)
6. [Timeouts and Deadlines](#6-timeouts-and-deadlines)
7. [Context in HTTP Requests](#7-context-in-http-requests)
8. [Context in Database Operations](#8-context-in-database-operations)
9. [Best Practices](#9-best-practices)
10. [Common Pitfalls to Avoid](#10-common-pitfalls-to-avoid)
11. [Context and Goroutine Leaks](#11-context-and-goroutine-leaks)
12. [Using Context with Third-Party Libraries](#12-using-context-with-third-party-libraries)
13. [New Features in Go 1.21](#13-new-features-in-go-121)
14. [FAQs](#14-faqs)

---

## 1. What is Context?

Context is a built-in package in the Go standard library that provides a powerful toolset for managing concurrent operations. It enables the propagation of cancellation signals, deadlines, and values across goroutines, ensuring that related operations can gracefully terminate when necessary.

### 🎯 Key Features

```mermaid
graph TD
    A[Context Package] --> B[Manage Concurrent Operations]
    A --> C[Propagate Cancellation Signals]
    A --> D[Handle Deadlines & Timeouts]
    A --> E[Pass Values Across Goroutines]
    
    B --> F[Coordinate Multiple Goroutines]
    C --> G[Graceful Termination]
    D --> H[Prevent Indefinite Waits]
    E --> I[Share Request-Scoped Data]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style D fill:#F44336,color:#fff
    style E fill:#9C27B0,color:#fff
```
![alt text](image.png)
### 📝 Example: Managing Concurrent API Requests

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    urls := []string{
        "https://api.example.com/users",
        "https://api.example.com/products",
        "https://api.example.com/orders",
    }

    results := make(chan string)

    for _, url := range urls {
        go fetchAPI(ctx, url, results)
    }

    for range urls {
        fmt.Println(<-results)
    }
}

func fetchAPI(ctx context.Context, url string, results chan<- string) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        results <- fmt.Sprintf("Error creating request for %s: %s", url, err.Error())
        return
    }

    client := http.DefaultClient
    resp, err := client.Do(req)
    if err != nil {
        results <- fmt.Sprintf("Error making request to %s: %s", url, err.Error())
        return
    }
    defer resp.Body.Close()

    results <- fmt.Sprintf("Response from %s: %d", url, resp.StatusCode)
}
```

**Output:**
```
Response from https://api.example.com/users: 200
Response from https://api.example.com/products: 200
Response from https://api.example.com/orders: 200
```

### 🔄 Concurrent API Request Flow
![alt text](image-1.png)
```mermaid
sequenceDiagram
    participant Main
    participant Context
    participant API1 as API Request 1
    participant API2 as API Request 2
    participant API3 as API Request 3
    
    Main->>Context: Create WithTimeout(5s)
    Main->>API1: Launch goroutine
    Main->>API2: Launch goroutine
    Main->>API3: Launch goroutine
    
    API1->>API1: Make HTTP Request
    API2->>API2: Make HTTP Request
    API3->>API3: Make HTTP Request
    
    alt Request exceeds timeout
        Context-->>API1: Cancel signal
        Context-->>API2: Cancel signal
        Context-->>API3: Cancel signal
    else Request completes
        API1-->>Main: Response
        API2-->>Main: Response
        API3-->>Main: Response
    end
```

---

## 2. Creating a Context

To create a context, you can use various functions from the context package. Each serves a different purpose.

### 🏗️ Context Creation Methods

![alt text](image-2.png)
```mermaid
graph TB
    A[Context Creation Methods] --> B[context.Background]
    A --> C[context.TODO]
    A --> D[context.WithCancel]
    A --> E[context.WithTimeout]
    A --> F[context.WithDeadline]
    A --> G[context.WithValue]
    
    B --> B1[Root context<br/>Non-cancelable<br/>No deadline]
    C --> C1[Placeholder<br/>Use when unsure<br/>Replace later]
    D --> D1[Manual cancellation<br/>Returns cancel func<br/>For explicit control]
    E --> E1[Auto-cancel after duration<br/>Time-based<br/>Returns cancel func]
    F --> F1[Auto-cancel at time<br/>Absolute time<br/>Returns cancel func]
    G --> G1[Attach key-value pairs<br/>Pass request data<br/>Type-safe retrieval]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style D fill:#F44336,color:#fff
    style E fill:#9C27B0,color:#fff
    style F fill:#E91E63,color:#fff
    style G fill:#00BCD4,color:#fff
```

### 📝 Example: Creating a Context with Timeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    go performTask(ctx)

    select {
    case <-ctx.Done():
        fmt.Println("Task timed out")
    }
}

func performTask(ctx context.Context) {
    select {
    case <-time.After(5 * time.Second):
        fmt.Println("Task completed successfully")
    case <-ctx.Done():
        fmt.Println("Task cancelled:", ctx.Err())
    }
}
```

**Output:**
```
Task timed out
```

### 🎯 Context Creation Decision Flow

![alt text](image-3.png)
```mermaid
flowchart LR
    A[Start] --> B{Need cancellation?}
    B -->|Yes| C{Time-based?}
    B -->|No| D{Need to pass values?}
    
    C -->|Yes - Duration| E[WithTimeout]
    C -->|Yes - Absolute| F[WithDeadline]
    C -->|No - Manual| G[WithCancel]
    
    D -->|Yes| H[WithValue]
    D -->|No| I[Background/TODO]
    
    E --> J[Context with timeout]
    F --> J
    G --> J
    H --> J
    I --> J
    
    J --> K[Pass to functions/goroutines]
    
    style A fill:#4CAF50,color:#fff
    style J fill:#2196F3,color:#fff
    style K fill:#FF9800,color:#fff
```

---

## 3. Propagating Context

Once you have a context, you can propagate it to downstream functions or goroutines by passing it as an argument.

### 🌳 Context Propagation Tree

![alt text](image-4.png)

```mermaid
graph TD
    A[Parent Context] --> B[Function 1]
    A --> C[Function 2]
    A --> D[Goroutine 1]
    
    B --> E[Child Function 1.1]
    B --> F[Child Function 1.2]
    
    C --> G[Goroutine 2]
    
    D --> H[Database Operation]
    
    G --> I[HTTP Request]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#2196F3,color:#fff
    style D fill:#FF9800,color:#fff
    style E fill:#9C27B0,color:#fff
    style F fill:#9C27B0,color:#fff
    style G fill:#FF9800,color:#fff
    style H fill:#F44336,color:#fff
    style I fill:#F44336,color:#fff
```

### 📝 Example: Propagating Context to Goroutines

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx := context.Background()
    ctx = context.WithValue(ctx, "UserID", 123)

    go performTask(ctx)
    
    time.Sleep(1 * time.Second) // Wait for goroutine
}

func performTask(ctx context.Context) {
    userID := ctx.Value("UserID")
    fmt.Println("User ID:", userID)
    
    // Can propagate further
    go subTask(ctx)
    time.Sleep(500 * time.Millisecond)
}

func subTask(ctx context.Context) {
    userID := ctx.Value("UserID")
    fmt.Println("SubTask - User ID:", userID)
}
```

**Output:**
```
User ID: 123
SubTask - User ID: 123
```

### 🔄 Context Propagation Sequence
![alt text](image-5.png)

```mermaid
sequenceDiagram
    participant Main as main()
    participant Ctx as Context
    participant F1 as Function A
    participant F2 as Function B
    participant G1 as Goroutine 1
    participant G2 as Goroutine 2
    
    Main->>Ctx: Create context
    Main->>Ctx: Add values (UserID: 123)
    Main->>F1: Pass context
    
    F1->>Ctx: Read values
    F1->>F2: Pass context
    F1->>G1: Launch with context
    
    F2->>Ctx: Read values
    F2->>G2: Launch with context
    
    G1->>Ctx: Access values
    G2->>Ctx: Access values
    
    Note over Main,G2: All operations share same context
```

---

## 4. Retrieving Values from Context

Context allows you to retrieve values stored within it, enabling access to important data within the scope of a specific goroutine or function.

### 🔑 Value Retrieval Flow

![alt text](image-6.png)

```mermaid
graph LR
    A[Context Creation] --> B[WithValue key: UserID]
    B --> C[Store value: 123]
    C --> D[Pass to Function]
    D --> E[Retrieve: ctx.Value UserID]
    E --> F[Type Assertion]
    F --> G[Use Value]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style D fill:#9C27B0,color:#fff
    style E fill:#F44336,color:#fff
    style F fill:#E91E63,color:#fff
    style G fill:#00BCD4,color:#fff
```

### 📝 Example: Retrieving User Information from Context

```go
package main

import (
    "context"
    "fmt"
)

// Define custom type for keys to avoid collisions
type contextKey string

const userIDKey contextKey = "UserID"

func main() {
    ctx := context.WithValue(context.Background(), userIDKey, 123)
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    if userID, ok := ctx.Value(userIDKey).(int); ok {
        fmt.Println("Processing request for User ID:", userID)
        
        // Use the user ID for further operations
        fetchUserData(ctx, userID)
    } else {
        fmt.Println("User ID not found in context")
    }
}

func fetchUserData(ctx context.Context, userID int) {
    // Verify we can still access the value
    if id := ctx.Value(userIDKey); id != nil {
        fmt.Printf("Fetching data for user: %d\n", userID)
    }
}
```

**Output:**
```
Processing request for User ID: 123
Fetching data for user: 123
```

### 🔍 Value Storage and Retrieval Pattern

![alt text](image-7.png)

```mermaid
flowchart TD
    A[Start] --> B[Create context.Background]
    B --> C[Add value: context.WithValue parent, key, value]
    C --> D[Pass context to function]
    D --> E[Retrieve: ctx.Value key]
    E --> F{Value exists?}
    
    F -->|Yes| G[Perform type assertion]
    F -->|No| H[Handle nil case]
    
    G --> I{Type valid?}
    I -->|Yes| J[Use value]
    I -->|No| K[Handle type error]
    
    J --> L[Continue processing]
    H --> L
    K --> L
    
    style A fill:#4CAF50,color:#fff
    style C fill:#2196F3,color:#fff
    style E fill:#FF9800,color:#fff
    style J fill:#00BCD4,color:#fff
```

---

## 5. Cancelling Context

Cancellation is an essential aspect of context management. It allows you to gracefully terminate operations and propagate cancellation signals to related goroutines.

### ⛔ Cancellation Propagation

![alt text](image-8.png)

```mermaid
graph TB
    A[Create WithCancel Context] --> B[Launch Goroutines]
    B --> C[Goroutine 1: Monitor ctx.Done]
    B --> D[Goroutine 2: Monitor ctx.Done]
    B --> E[Goroutine 3: Monitor ctx.Done]
    
    F[Main Thread] --> G[Call cancel]
    
    G --> H[Signal sent to ctx.Done]
    H --> I[All goroutines receive signal]
    
    I --> C
    I --> D
    I --> E
    
    C --> J[Cleanup & Exit]
    D --> J
    E --> J
    
    style A fill:#4CAF50,color:#fff
    style G fill:#F44336,color:#fff
    style H fill:#FF9800,color:#fff
    style J fill:#2196F3,color:#fff
```

### 📝 Example: Cancelling Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go performTask(ctx, "Worker-1")
    go performTask(ctx, "Worker-2")
    go performTask(ctx, "Worker-3")

    time.Sleep(2 * time.Second)
    fmt.Println("\nCancelling all workers...")
    cancel()
    
    time.Sleep(1 * time.Second)
    fmt.Println("Main exiting")
}

func performTask(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: Task cancelled - %v\n", name, ctx.Err())
            return
        default:
            fmt.Printf("%s: Performing task...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

**Output:**
```
Worker-1: Performing task...
Worker-2: Performing task...
Worker-3: Performing task...
Worker-1: Performing task...
Worker-2: Performing task...
Worker-3: Performing task...
Worker-1: Performing task...
Worker-2: Performing task...
Worker-3: Performing task...

Cancelling all workers...
Worker-1: Task cancelled - context canceled
Worker-2: Task cancelled - context canceled
Worker-3: Task cancelled - context canceled
Main exiting
```

### 🔄 Cancellation Flow Sequence


![alt text](image-9.png)
```mermaid
sequenceDiagram
    participant Main
    participant Ctx as Context
    participant G1 as Goroutine 1
    participant G2 as Goroutine 2
    participant G3 as Goroutine 3
    
    Main->>Ctx: WithCancel(Background)
    Main->>G1: Launch with context
    Main->>G2: Launch with context
    Main->>G3: Launch with context
    
    loop Work Loop
        G1->>G1: Perform task
        G2->>G2: Perform task
        G3->>G3: Perform task
    end
    
    Main->>Ctx: cancel()
    Ctx-->>G1: ctx.Done() closed
    Ctx-->>G2: ctx.Done() closed
    Ctx-->>G3: ctx.Done() closed
    
    G1->>G1: Cleanup & Exit
    G2->>G2: Cleanup & Exit
    G3->>G3: Cleanup & Exit
```

---

## 6. Timeouts and Deadlines

Setting timeouts and deadlines is crucial when working with context. It ensures that operations complete within a specified timeframe.

### ⏱️ Timeout Management
![alt text](image-10.png)

```mermaid
graph TD
    A[Timeout Management] --> B[WithTimeout]
    A --> C[WithDeadline]
    
    B --> D[Relative Duration<br/>e.g., 5 seconds from now]
    C --> E[Absolute Time<br/>e.g., 2:00 PM today]
    
    D --> F[Auto-cancel after duration]
    E --> G[Auto-cancel at specific time]
    
    F --> H{Operation completes?}
    G --> H
    
    H -->|Before timeout| I[Success]
    H -->|After timeout| J[context deadline exceeded]
    
    I --> K[Continue normally]
    J --> L[Handle timeout error]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style J fill:#F44336,color:#fff
```

### 📝 Example: Setting a Deadline for Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Example 1: WithTimeout
    fmt.Println("=== WithTimeout Example ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel1()
    
    go performTask(ctx1, "Task-1", 5*time.Second)
    
    <-ctx1.Done()
    fmt.Println("Task-1 result:", ctx1.Err())
    
    time.Sleep(1 * time.Second)
    
    // Example 2: WithDeadline
    fmt.Println("\n=== WithDeadline Example ===")
    deadline := time.Now().Add(2 * time.Second)
    ctx2, cancel2 := context.WithDeadline(context.Background(), deadline)
    defer cancel2()
    
    go performTask(ctx2, "Task-2", 1*time.Second)
    
    <-ctx2.Done()
    fmt.Println("Task-2 result:", ctx2.Err())
}

func performTask(ctx context.Context, name string, duration time.Duration) {
    select {
    case <-time.After(duration):
        fmt.Printf("%s completed successfully\n", name)
    case <-ctx.Done():
        fmt.Printf("%s cancelled: %v\n", name, ctx.Err())
    }
}
```

**Output:**
```
=== WithTimeout Example ===
Task-1 cancelled: context deadline exceeded
Task-1 result: context deadline exceeded

=== WithDeadline Example ===
Task-2 completed successfully
Task-2 result: context deadline exceeded
```

### 📊 Timeout Execution Timeline

![alt text](image-11.png)

```mermaid
gantt
    title Context Timeout Scenario
    dateFormat ss
    axisFormat %Ss
    
    section Context
    Create Context (2s timeout)    :milestone, m1, 00, 0s
    Context Active                  :active, t1, 00, 2s
    Context Deadline Exceeded       :crit, t2, 02, 1s
    
    section Task
    Task Started                    :milestone, m2, 00, 0s
    Task Execution (5s required)    :task1, 00, 5s
    Task Cancelled                  :crit, milestone, m3, 02, 0s
    
    section Result
    Timeout Triggered               :milestone, m4, 02, 0s
```

---

## 7. Context in HTTP Requests

Context plays a vital role in managing HTTP requests in Go, allowing you to control request cancellation, timeouts, and pass important values.

### 🌐 HTTP Request with Context
![alt text](image-12.png)
```mermaid
graph TB
    A[HTTP Request] --> B[Create Context with Timeout]
    B --> C[NewRequestWithContext]
    C --> D[Send Request]
    
    D --> E{Response received?}
    
    E -->|Before timeout| F[Process Response]
    E -->|After timeout| G[Context Cancelled]
    
    F --> H[Success]
    G --> I[Handle Error]
    
    I --> J[Cleanup]
    H --> J
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style G fill:#F44336,color:#fff
    style H fill:#00BCD4,color:#fff
```

### 📝 Example: Using Context in HTTP Requests

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
    // Create a context with 2 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Make multiple requests with the same context
    urls := []string{
        "https://httpbin.org/delay/1",  // 1 second delay
        "https://httpbin.org/delay/3",  // 3 second delay (will timeout)
    }

    for _, url := range urls {
        makeRequest(ctx, url)
    }
}

func makeRequest(ctx context.Context, url string) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        fmt.Printf("Error creating request for %s: %v\n", url, err)
        return
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("Error making request to %s: %v\n", url, err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response from %s: Status=%d, Body length=%d bytes\n", 
        url, resp.StatusCode, len(body))
}
```

**Output:**
```
Response from https://httpbin.org/delay/1: Status=200, Body length=XXX bytes
Error making request to https://httpbin.org/delay/3: context deadline exceeded
```

### 🔄 HTTP Request Flow with Context
![alt text](image-13.png)
```mermaid
sequenceDiagram
    participant Client
    participant Ctx as Context (2s timeout)
    participant HTTP as HTTP Request
    participant Server
    
    Client->>Ctx: WithTimeout(2s)
    Client->>HTTP: NewRequestWithContext(ctx)
    HTTP->>Server: Send Request
    
    alt Fast Response (< 2s)
        Server-->>HTTP: Response (200 OK)
        HTTP-->>Client: Return response
        Client->>Client: Process data
    else Slow Response (> 2s)
        Ctx-->>HTTP: Timeout signal
        HTTP->>HTTP: Cancel request
        HTTP-->>Client: Error: context deadline exceeded
        Client->>Client: Handle timeout
    end
```

---

## 8. Context in Database Operations

Context is useful when dealing with database operations, allowing you to manage query cancellations, timeouts, and pass relevant data.

### 🗄️ Database Operation Flow
![alt text](image-14.png)
```mermaid
graph TD
    A[Database Query] --> B[Create Context with Timeout]
    B --> C[db.QueryContext ctx, query]
    C --> D[Execute Query]
    
    D --> E{Query completes?}
    
    E -->|Before timeout| F[Process Rows]
    E -->|After timeout| G[Query Cancelled]
    
    F --> H[Return Results]
    G --> I[Rollback/Cleanup]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style G fill:#F44336,color:#fff
```

### 📝 Example: Using Context in Database Operations

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq"
)

type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    // Connect to database
    db, err := sql.Open("postgres", 
        "postgres://username:password@localhost/mydatabase?sslmode=disable")
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return
    }
    defer db.Close()

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Execute query with context
    users, err := getUsersWithContext(ctx, db)
    if err != nil {
        fmt.Println("Error executing query:", err)
        return
    }

    // Process results
    for _, user := range users {
        fmt.Printf("User: ID=%d, Name=%s, Email=%s\n", 
            user.ID, user.Name, user.Email)
    }
}

func getUsersWithContext(ctx context.Context, db *sql.DB) ([]User, error) {
    query := "SELECT id, name, email FROM users WHERE active = true"
    
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration failed: %w", err)
    }

    return users, nil
}
```

### 🔄 Database Query Execution Flow
![alt text](image-15.png)

```mermaid
flowchart LR
    A[Start] --> B[Create Context<br/>Timeout: 2s]
    B --> C[Open DB Connection]
    C --> D[Begin Transaction]
    D --> E[Execute Query<br/>with context]
    
    E --> F{Timeout?}
    
    F -->|No| G[Fetch Results]
    F -->|Yes| H[Cancel Query]
    
    G --> I[Process Data]
    H --> J[Rollback Transaction]
    
    I --> K[Commit Transaction]
    
    K --> L[Close Connection]
    J --> L
    
    L --> M[End]
    
    style B fill:#4CAF50,color:#fff
    style E fill:#2196F3,color:#fff
    style H fill:#F44336,color:#fff
    style K fill:#00BCD4,color:#fff
```

---

## 9. Best Practices

When working with context in Golang, follow these best practices to ensure efficient and reliable concurrency management.

### 💡 Best Practices Mind Map

```mermaid
mindmap
    root((Context<br/>Best Practices))
        Pass Explicitly
            Function arguments
            No globals
            Clear lifecycle
        Use TODO wisely
            Temporary placeholder
            Replace later
            Document intent
        Avoid Background
            Create specific context
            Add cancellation
            Prevent leaks
        Prefer Cancel
            Explicit control
            Clear intention
            Better than timeout
        Keep Size Small
            Minimal data
            Request-scoped only
            Avoid bloat
        No Chaining
            Single propagation
            Simple hierarchy
            Easy management
        Prevent Leaks
            Close goroutines
            Check Done channel
            Proper cleanup
```

### ✅ Best Practice Guidelines

1. **Pass Context Explicitly**: Always pass the context as the first argument to functions
2. **Use context.TODO()**: Only as a temporary placeholder when refactoring
3. **Avoid Using context.Background()**: Create specific contexts with cancellation/timeout
4. **Prefer Cancel Over Timeout**: Use explicit cancellation when possible
5. **Keep Context Size Small**: Only include necessary request-scoped data
6. **Avoid Chaining Contexts**: Propagate a single context throughout
7. **Be Mindful of Goroutine Leaks**: Ensure proper cleanup

### 📝 Example: Implementing Best Practices

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// ✅ GOOD: Context as first parameter
func processRequest(ctx context.Context, userID int) error {
    // Create child context with timeout for this specific operation
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel() // ✅ GOOD: Always call cancel

    // ✅ GOOD: Check context before expensive operations
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Simulate work
    return fetchUserData(ctx, userID)
}

func fetchUserData(ctx context.Context, userID int) error {
    // ✅ GOOD: Pass context to downstream operations
    resultChan := make(chan error, 1)
    
    go func() {
        // Simulate database query
        time.Sleep(2 * time.Second)
        resultChan <- nil
    }()

    // ✅ GOOD: Monitor context while waiting
    select {
    case <-ctx.Done():
        return ctx.Err()
    case err := <-resultChan:
        return err
    }
}

// Custom type for context keys (✅ GOOD: Prevents collisions)
type contextKey string

const userIDKey contextKey = "userID"

func main() {
    // ✅ GOOD: Create context with specific timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // ✅ GOOD: Use custom type for context keys
    ctx = context.WithValue(ctx, userIDKey, 123)

    if err := processRequest(ctx, 123); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Request processed successfully")
    }
}
```

### 🎯 Best Practice Decision Tree

```mermaid
flowchart TD
    A[Need Context?] --> B{What for?}
    
    B -->|Cancellation| C{Time-based?}
    B -->|Pass values| D[Use WithValue]
    B -->|Unsure| E[Use TODO temporarily]
    
    C -->|Yes| F[WithTimeout/Deadline]
    C -->|No| G[WithCancel]
    
    F --> H[Set reasonable timeout]
    G --> I[Call cancel when done]
    D --> J[Keep values small]
    
    H --> K[Pass explicitly]
    I --> K
    J --> K
    E --> K
    
    K --> L[Check ctx.Done in loops]
    L --> M{Goroutines?}
    
    M -->|Yes| N[Ensure cleanup]
    M -->|No| O[Complete]
    
    N --> O
    
    style A fill:#4CAF50,color:#fff
    style K fill:#2196F3,color:#fff
    style O fill:#00BCD4,color:#fff
```

---

## 10. Common Pitfalls to Avoid

### ⚠️ Common Pitfalls Overview

```mermaid
graph TB
    A[Common Pitfalls] --> B[Not Propagating Context]
    A --> C[Forgetting cancel]
    A --> D[Leaking Goroutines]
    A --> E[Using Background everywhere]
    A --> F[Passing nil]
    A --> G[Checking context too early]
    A --> H[Using blocking calls]
    A --> I[Overusing contexts]
    
    B --> B1[❌ Child functions can't cancel]
    C --> C1[❌ Resources not released]
    D --> D1[❌ Memory leaks]
    E --> E1[❌ No timeout/cancellation]
    F --> F1[❌ Panics occur]
    G --> G1[❌ Work never starts]
    H --> H1[❌ Can't cancel]
    I --> I1[❌ Wrong tool for job]
    
    style A fill:#F44336,color:#fff
    style B fill:#FF9800,color:#000
    style C fill:#FF9800,color:#000
    style D fill:#FF9800,color:#000
    style E fill:#FF9800,color:#000
```

### 📝 Examples of Common Pitfalls

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// ❌ BAD: Not propagating context
func badExample1() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // ❌ Creating new context instead of propagating
    doWork(context.Background()) // Wrong! Ignores parent timeout
}

// ✅ GOOD: Propagating context
func goodExample1() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // ✅ Pass the parent context
    doWork(ctx)
}

// ❌ BAD: Forgetting to call cancel
func badExample2() {
    ctx, _ := context.WithCancel(context.Background())
    // ❌ Cancel function is ignored - resource leak!
    
    doWork(ctx)
}

// ✅ GOOD: Always call cancel
func goodExample2() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // ✅ Ensures resources are released
    
    doWork(ctx)
}

// ❌ BAD: Goroutine doesn't check context
func badExample3(ctx context.Context) {
    go func() {
        for {
            // ❌ Never checks ctx.Done() - goroutine leak!
            fmt.Println("Working...")
            time.Sleep(100 * time.Millisecond)
        }
    }()
}

// ✅ GOOD: Goroutine monitors context
func goodExample3(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine exiting")
                return // ✅ Properly exits
            default:
                fmt.Println("Working...")
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
}

// ❌ BAD: Using Background for everything
func badExample4() {
    // ❌ No timeout, no cancellation
    doWork(context.Background())
}

// ✅ GOOD: Use specific context type
func goodExample4() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // ✅ Has timeout protection
    doWork(ctx)
}

// ❌ BAD: Passing nil context
func badExample5() {
    doWork(nil) // ❌ Will panic!
}

// ✅ GOOD: Always pass valid context
func goodExample5() {
    doWork(context.Background()) // ✅ At minimum, use Background
}

func doWork(ctx context.Context) {
    // Implementation
    fmt.Println("Doing work...")
}

func main() {
    fmt.Println("=== Good Examples ===")
    goodExample1()
    goodExample2()
    goodExample3(context.Background())
    goodExample4()
    goodExample5()
}
```

### 🔍 Pitfall Prevention Checklist

```mermaid
flowchart LR
    A[Creating Context] --> B{Remember to:}
    
    B --> C[✓ Pass to children]
    B --> D[✓ Call cancel]
    B --> E[✓ Check Done]
    B --> F[✓ Add timeout if needed]
    
    C --> G{In Functions}
    D --> G
    E --> G
    F --> G
    
    G --> H[✓ Never pass nil]
    G --> I[✓ Don't check too early]
    G --> J[✓ Wrap blocking calls]
    
    H --> K[Context Used Correctly]
    I --> K
    J --> K
    
    style A fill:#4CAF50,color:#fff
    style K fill:#00BCD4,color:#fff
```

---

## 11. Context and Goroutine Leaks

Goroutine leaks occur when goroutines are started but never properly terminated, leading to memory leaks.

### 🚨 Leak vs No Leak Comparison

```mermaid
graph TD
    subgraph "❌ LEAK: Improper Handling"
        A1[Main] --> B1[Create ctx]
        B1 --> C1[Start goroutine]
        C1 --> D1[Goroutine runs]
        D1 --> E1[Cancel context]
        E1 --> F1[Goroutine continues!]
        F1 --> G1[Memory leak]
    end
    
    subgraph "✅ CORRECT: Proper Handling"
        A2[Main] --> B2[Create ctx with cancel]
        B2 --> C2[Start goroutine]
        C2 --> D2[Monitor ctx.Done]
        D2 --> E2[Cancel context]
        E2 --> F2[Goroutine exits]
        F2 --> G2[No leak]
    end
    
    style G1 fill:#F44336,color:#fff
    style G2 fill:#4CAF50,color:#fff
```

### 📝 Example: Goroutine Leak vs Proper Handling

```go
package main

import (
    "context"
    "fmt"
    "runtime"
    "time"
)

// ❌ BAD: Leaking goroutine
func leakingExample() {
    ctx := context.Background()
    
    go func(ctx context.Context) {
        for {
            // ❌ No check for ctx.Done() - goroutine never exits!
            fmt.Println("Leaking goroutine working...")
            time.Sleep(500 * time.Millisecond)
        }
    }(ctx)
    
    time.Sleep(2 * time.Second)
    // Goroutine continues running even after function returns!
}

// ✅ GOOD: Proper goroutine cleanup
func properExample() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func(ctx context.Context) {
        for {
            select {
            case <-ctx.Done():
                // ✅ Properly handle cancellation
                fmt.Println("Goroutine exiting cleanly")
                return
            default:
                fmt.Println("Proper goroutine working...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }(ctx)
    
    time.Sleep(2 * time.Second)
    cancel() // ✅ Signal goroutine to stop
    time.Sleep(100 * time.Millisecond) // Give it time to exit
}

// ✅ BETTER: Using timeout for automatic cleanup
func betterExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    go func(ctx context.Context) {
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop() // ✅ Clean up ticker
        
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine exiting due to:", ctx.Err())
                return
            case <-ticker.C:
                fmt.Println("Better goroutine working...")
            }
        }
    }(ctx)
    
    time.Sleep(3 * time.Second) // Wait longer than timeout
}

func printGoroutineCount(label string) {
    fmt.Printf("%s - Goroutine count: %d\n", label, runtime.NumGoroutine())
}

func main() {
    printGoroutineCount("Start")
    
    fmt.Println("\n=== Leaking Example ===")
    leakingExample()
    printGoroutineCount("After leaking example")
    
    time.Sleep(1 * time.Second)
    
    fmt.Println("\n=== Proper Example ===")
    properExample()
    printGoroutineCount("After proper example")
    
    time.Sleep(1 * time.Second)
    
    fmt.Println("\n=== Better Example (with timeout) ===")
    betterExample()
    printGoroutineCount("After better example")
}
```

**Output:**
```
Start - Goroutine count: 1

=== Leaking Example ===
Leaking goroutine working...
Leaking goroutine working...
Leaking goroutine working...
After leaking example - Goroutine count: 2

=== Proper Example ===
Proper goroutine working...
Proper goroutine working...
Proper goroutine working...
Goroutine exiting cleanly
After proper example - Goroutine count: 2

=== Better Example (with timeout) ===
Better goroutine working...
Better goroutine working...
Better goroutine working...
Goroutine exiting due to: context deadline exceeded
After better example - Goroutine count: 1
```

### 🔄 Leak vs No Leak Sequence

```mermaid
sequenceDiagram
    participant M as Main
    participant C as Context
    participant GL as Goroutine (Leaking)
    participant GC as Goroutine (Correct)
    
    Note over M,GC: ❌ Leaking Pattern
    M->>C: Create context
    M->>GL: Start goroutine
    loop Forever
        GL->>GL: Do work (no ctx check)
    end
    M->>C: cancel()
    Note over GL: Goroutine never exits!
    
    Note over M,GC: ✅ Correct Pattern
    M->>C: WithCancel()
    M->>GC: Start goroutine
    loop Until cancelled
        GC->>C: Check ctx.Done()
        GC->>GC: Do work
    end
    M->>C: cancel()
    C-->>GC: Done signal
    GC->>GC: Exit cleanly
```

---

## 12. Using Context with Third-Party Libraries

Many third-party libraries don't natively support context. Here's how to integrate context properly.

### 🔌 Third-Party Integration Pattern

```mermaid
graph TB
    A[Third-Party API<br/>No Context Support] --> B[Create Wrapper Function]
    B --> C[Accept Context Parameter]
    C --> D{Before API Call}
    
    D --> E[Check ctx.Done]
    E --> F{Cancelled?}
    
    F -->|Yes| G[Return Early]
    F -->|No| H[Call API in Goroutine]
    
    H --> I[API Processing]
    I --> J{After API Call}
    
    J --> K[Check ctx.Done again]
    K --> L{Cancelled?}
    
    L -->|Yes| M[Discard Result]
    L -->|No| N[Return Result]
    
    style A fill:#FF9800,color:#000
    style B fill:#4CAF50,color:#fff
    style C fill:#2196F3,color:#fff
    style N fill:#00BCD4,color:#fff
```

### 📝 Example: Wrapping Third-Party APIs

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// Simulated third-party library that doesn't support context
type ThirdPartyAPI struct{}

func (api *ThirdPartyAPI) FetchData(query string) (string, error) {
    // Simulates a slow API call
    time.Sleep(3 * time.Second)
    return fmt.Sprintf("Data for: %s", query), nil
}

// ✅ GOOD: Wrapper that adds context support
func FetchDataWithContext(ctx context.Context, api *ThirdPartyAPI, query string) (string, error) {
    // Check if context is already cancelled before starting
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    default:
    }

    // Create a channel to receive the result
    type result struct {
        data string
        err  error
    }
    resultChan := make(chan result, 1)

    // Call the third-party API in a goroutine
    go func() {
        data, err := api.FetchData(query)
        resultChan <- result{data: data, err: err}
    }()

    // Wait for either the result or context cancellation
    select {
    case <-ctx.Done():
        // Context cancelled - return immediately
        return "", ctx.Err()
    case res := <-resultChan:
        // Check one more time if context was cancelled
        select {
        case <-ctx.Done():
            return "", ctx.Err()
        default:
            return res.data, res.err
        }
    }
}

// ✅ BETTER: Wrapper with timeout defaults
func FetchDataWithContextAndDefaults(ctx context.Context, api *ThirdPartyAPI, query string) (string, error) {
    // Add a reasonable default timeout if none exists
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }

    return FetchDataWithContext(ctx, api, query)
}

// Example of batch operations with context
func FetchMultipleWithContext(ctx context.Context, api *ThirdPartyAPI, queries []string) ([]string, error) {
    results := make([]string, len(queries))
    errChan := make(chan error, len(queries))
    
    for i, query := range queries {
        go func(index int, q string) {
            data, err := FetchDataWithContext(ctx, api, q)
            if err != nil {
                errChan <- err
                return
            }
            results[index] = data
            errChan <- nil
        }(i, query)
    }

    // Collect all results or errors
    for range queries {
        if err := <-errChan; err != nil {
            return nil, err
        }
    }

    return results, nil
}

func main() {
    api := &ThirdPartyAPI{}

    // Example 1: Normal operation
    fmt.Println("=== Example 1: Normal Operation ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel1()

    data, err := FetchDataWithContext(ctx1, api, "users")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Success:", data)
    }

    // Example 2: Timeout scenario
    fmt.Println("\n=== Example 2: Timeout Scenario ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel2()

    data, err = FetchDataWithContext(ctx2, api, "products")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Success:", data)
    }

    // Example 3: Manual cancellation
    fmt.Println("\n=== Example 3: Manual Cancellation ===")
    ctx3, cancel3 := context.WithCancel(context.Background())

    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Cancelling context...")
        cancel3()
    }()

    data, err = FetchDataWithContext(ctx3, api, "orders")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Success:", data)
    }

    // Example 4: Batch operations
    fmt.Println("\n=== Example 4: Batch Operations ===")
    ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel4()

    queries := []string{"users", "products", "orders"}
    results, err := FetchMultipleWithContext(ctx4, api, queries)
    if err != nil {
        fmt.Println("Batch error:", err)
    } else {
        for i, result := range results {
            fmt.Printf("Result %d: %s\n", i+1, result)
        }
    }
}
```

### 🔄 Third-Party Integration Sequence

```mermaid
sequenceDiagram
    participant App as Application
    participant Wrapper as Wrapper Function
    participant Ctx as Context
    participant API as Third-Party API
    
    App->>Wrapper: Call with context
    Wrapper->>Ctx: Check Done()
    
    alt Context already cancelled
        Wrapper-->>App: Return immediately
    else Context active
        Wrapper->>API: Launch API call
        
        par API Processing
            API->>API: Execute operation
        and Context Monitoring
            Wrapper->>Ctx: Monitor Done()
        end
        
        alt Context cancelled during call
            Ctx-->>Wrapper: Cancel signal
            Wrapper->>Wrapper: Discard result
            Wrapper-->>App: Return error
        else API completes first
            API-->>Wrapper: Return result
            Wrapper->>Ctx: Check Done()
            Wrapper-->>App: Return result
        end
    end
```

---

## 13. New Features in Go 1.21

Go 1.21 introduced several new context-related functions that enhance error handling and cleanup management.

### 🆕 Go 1.21 Features Overview

```mermaid
quadrantChart
    title Go 1.21 Context Features
    x-axis Low Complexity --> High Complexity
    y-axis Low Impact --> High Impact
    
    AfterFunc: [0.6, 0.7]
    WithDeadlineCause: [0.5, 0.8]
    WithTimeoutCause: [0.5, 0.85]
    WithoutCancel: [0.4, 0.6]
```

### 1️⃣ AfterFunc

Schedules a function to run after a context ends, enabling deferred cleanup routines.

```mermaid
graph LR
    A[Create Context] --> B[context.AfterFunc]
    B --> C[Register cleanup function]
    C --> D[Context runs normally]
    D --> E{Context Done?}
    
    E -->|Yes| F[Cleanup function executes]
    E -->|No| G[Continue]
    
    F --> H[Deferred work completed]
    
    I[stop function] -.->|Can cancel cleanup| F
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style F fill:#FF9800,color:#fff
```

#### Example: AfterFunc

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // Schedule cleanup to run after context is cancelled
    stop := context.AfterFunc(ctx, func() {
        fmt.Println("Cleanup: Processing remaining items...")
        fmt.Println("Cleanup: Closing connections...")
        fmt.Println("Cleanup: Complete!")
    })
    
    // Simulate work
    go func() {
        for i := 1; i <= 3; i++ {
            select {
            case <-ctx.Done():
                return
            default:
                fmt.Printf("Working... step %d\n", i)
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()
    
    time.Sleep(2 * time.Second)
    
    // Cancel context - triggers AfterFunc
    fmt.Println("Cancelling context...")
    cancel()
    
    time.Sleep(1 * time.Second)
    
    // Optionally, we could have called stop() to prevent cleanup
    _ = stop
}
```

**Output:**
```
Working... step 1
Working... step 2
Working... step 3
Cancelling context...
Cleanup: Processing remaining items...
Cleanup: Closing connections...
Cleanup: Complete!
```

### 2️⃣ WithDeadlineCause

Associates a custom error cause with a context's deadline.

```mermaid
flowchart TD
    A[Start] --> B[WithDeadlineCause ctx, deadline, cause]
    B --> C[Set custom error message]
    C --> D[Operation runs]
    
    D --> E{Deadline exceeded?}
    
    E -->|No| F[Complete successfully]
    E -->|Yes| G[Return error with cause]
    
    G --> H[Error: context deadline exceeded:<br/>custom cause message]
    
    H --> I[Better debugging info]
    
    style B fill:#4CAF50,color:#fff
    style G fill:#F44336,color:#fff
    style I fill:#2196F3,color:#fff
```

#### Example: WithDeadlineCause

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
)

func main() {
    deadline := time.Now().Add(100 * time.Millisecond)
    ctx, cancel := context.WithDeadlineCause(
        context.Background(),
        deadline,
        errors.New("RPC timeout: backend service unavailable"),
    )
    defer cancel()

    // Simulate work that takes too long
    time.Sleep(200 * time.Millisecond)

    // Check the error
    fmt.Println("Error:", ctx.Err())
    fmt.Println("Cause:", context.Cause(ctx))
}
```

**Output:**
```
Error: context deadline exceeded
Cause: RPC timeout: backend service unavailable
```

### 3️⃣ WithTimeoutCause

Associates a custom error cause with a context's timeout duration.

```mermaid
graph TB
    A[Create Context] --> B[WithTimeoutCause<br/>duration + cause]
    B --> C[Set timeout: 100ms<br/>cause: Backend RPC timeout]
    C --> D[Execute Operation]
    
    D --> E{Timeout reached?}
    
    E -->|No| F[Success]
    E -->|Yes| G[Timeout Error]
    
    G --> H[Error message includes:<br/>context deadline exceeded:<br/>Backend RPC timeout]
    
    H --> I[Improved debugging]
    
    F --> J[Continue]
    I --> K[Better error tracking]
    
    style B fill:#4CAF50,color:#fff
    style G fill:#F44336,color:#fff
    style I fill:#2196F3,color:#fff
```

#### Example: WithTimeoutCause

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
)

func main() {
    // Example 1: Database query timeout
    fmt.Println("=== Database Query Example ===")
    ctx1, cancel1 := context.WithTimeoutCause(
        context.Background(),
        100*time.Millisecond,
        errors.New("database query timeout: query too complex"),
    )
    defer cancel1()

    queryDatabase(ctx1)

    time.Sleep(200 * time.Millisecond)

    // Example 2: API call timeout
    fmt.Println("\n=== API Call Example ===")
    ctx2, cancel2 := context.WithTimeoutCause(
        context.Background(),
        100*time.Millisecond,
        errors.New("API timeout: external service not responding"),
    )
    defer cancel2()

    callAPI(ctx2)

    time.Sleep(200 * time.Millisecond)
}

func queryDatabase(ctx context.Context) {
    select {
    case <-time.After(150 * time.Millisecond):
        fmt.Println("Query completed")
    case <-ctx.Done():
        fmt.Println("Query failed")
        fmt.Println("  Error:", ctx.Err())
        fmt.Println("  Cause:", context.Cause(ctx))
    }
}

func callAPI(ctx context.Context) {
    select {
    case <-time.After(150 * time.Millisecond):
        fmt.Println("API call completed")
    case <-ctx.Done():
        fmt.Println("API call failed")
        fmt.Println("  Error:", ctx.Err())
        fmt.Println("  Cause:", context.Cause(ctx))
    }
}
```

**Output:**
```
=== Database Query Example ===
Query failed
  Error: context deadline exceeded
  Cause: database query timeout: query too complex

=== API Call Example ===
API call failed
  Error: context deadline exceeded
  Cause: API timeout: external service not responding
```

### 4️⃣ WithoutCancel

Creates a child context that is detached from the parent's cancellation.

```mermaid
graph TD
    A[Parent Context] --> B[WithoutCancel parent]
    B --> C[Detached Child Context]
    
    D[Parent Cancelled] --> E{Children behavior}
    
    E --> F[Normal children:<br/>Also cancelled]
    E --> G[WithoutCancel children:<br/>Continue running]
    
    A --> H[Normal Child 1]
    A --> I[Normal Child 2]
    C --> J[Detached Child 1]
    C --> K[Detached Child 2]
    
    D -.-> H
    D -.-> I
    D -.X J
    D -.X K
    
    H --> L[Terminates]
    I --> L
    J --> M[Continues independently]
    K --> M
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style M fill:#00BCD4,color:#fff
```

#### Example: WithoutCancel

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    parentCtx, cancel := context.WithCancel(context.Background())
    
    // Normal child - will be cancelled with parent
    normalChild := parentCtx
    
    // Detached child - won't be cancelled with parent
    detachedChild := context.WithoutCancel(parentCtx)
    
    // Start workers
    go worker(normalChild, "Normal Worker")
    go worker(detachedChild, "Detached Worker")
    
    time.Sleep(2 * time.Second)
    
    fmt.Println("\n=== Cancelling Parent Context ===\n")
    cancel()
    
    time.Sleep(3 * time.Second)
    fmt.Println("\n=== Main Exiting ===")
}

func worker(ctx context.Context, name string) {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: Context cancelled - exiting\n", name)
            return
        case <-ticker.C:
            fmt.Printf("%s: Working...\n", name)
        }
    }
}
```

**Output:**
```
Normal Worker: Working...
Detached Worker: Working...
Normal Worker: Working...
Detached Worker: Working...
Normal Worker: Working...
Detached Worker: Working...

=== Cancelling Parent Context ===

Normal Worker: Context cancelled - exiting
Detached Worker: Working...
Detached Worker: Working...
Detached Worker: Working...
Detached Worker: Working...

=== Main Exiting ===
```

---

## 14. FAQs

### Q1: Can I pass custom values in the context?

**A:** Yes, you can pass custom values using `context.WithValue()`. However, it's recommended to use custom types for keys to avoid collisions.

```go
type contextKey string
const userIDKey contextKey = "userID"

ctx := context.WithValue(context.Background(), userIDKey, 123)
```

### Q2: How can I handle context cancellation errors?

**A:** Use `ctx.Err()` to get the cancellation reason:

```go
select {
case <-ctx.Done():
    if err := ctx.Err(); err == context.Canceled {
        fmt.Println("Context was cancelled")
    } else if err == context.DeadlineExceeded {
        fmt.Println("Context deadline exceeded")
    }
}
```

### Q3: Can I create a hierarchy of contexts?

**A:** Yes, contexts naturally form a hierarchy. Child contexts inherit from their parents:

```go
parent := context.Background()
child1 := context.WithValue(parent, "key1", "value1")
child2 := context.WithTimeout(child1, 5*time.Second)
```

### Q4: Is it possible to pass context through HTTP middleware?

**A:** Yes, middleware can enrich or modify the request context:

```go
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), "requestID", generateID())
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Q5: How does context help with graceful shutdowns?

**A:** By canceling a context, you signal all goroutines to complete their work:

```go
ctx, cancel := context.WithCancel(context.Background())

// On shutdown signal
signal.Notify(shutdownChan, os.Interrupt)
<-shutdownChan
cancel() // All goroutines monitoring ctx.Done() will exit
```

---

## 📚 Complete Context Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Created: context.Background()<br/>or context.TODO()
    
    Created --> WithValue: Add values
    Created --> WithCancel: Add cancellation
    Created --> WithTimeout: Add timeout
    Created --> WithDeadline: Add deadline
    
    WithValue --> Active
    WithCancel --> Active
    WithTimeout --> Active
    WithDeadline --> Active
    
    Active --> Cancelled: cancel() called
    Active --> TimedOut: Deadline/Timeout exceeded
    Active --> Completed: Work finished normally
    
    Cancelled --> Done
    TimedOut --> Done
    Completed --> Done
    
    Done --> [*]
    
    note right of Active
        Goroutines monitor ctx.Done()
        Operations check context state
    end note
    
    note right of Done
        Resources cleaned up
        Goroutines terminated
        Error available via ctx.Err()
    end note
```

---

## 🎯 Context Patterns Summary

```mermaid
graph TB
    A[Context Patterns] --> B[Creation]
    A --> C[Propagation]
    A --> D[Cancellation]
    A --> E[Value Passing]
    
    B --> B1[Background/TODO]
    B --> B2[WithCancel]
    B --> B3[WithTimeout/Deadline]
    B --> B4[WithValue]
    
    C --> C1[Pass as first parameter]
    C --> C2[Never store in structs]
    C --> C3[Propagate to all children]
    
    D --> D1[Manual: cancel]
    D --> D2[Automatic: timeout]
    D --> D3[Monitor: ctx.Done]
    
    E --> E1[Request-scoped data]
    E --> E2[Type-safe keys]
    E --> E3[Avoid large values]
    
    style A fill:#4CAF50,color:#fff
    style B fill:#2196F3,color:#fff
    style C fill:#FF9800,color:#fff
    style D fill:#F44336,color:#fff
    style E fill:#9C27B0,color:#fff
```

---

## 🔄 End-to-End Context Usage

```mermaid
sequenceDiagram
    participant Main
    participant Ctx as Context
    participant API as API Layer
    participant DB as Database Layer
    participant Cache as Cache Layer
    
    Main->>Ctx: WithTimeout(5s)
    Main->>Ctx: WithValue("userID", 123)
    Main->>API: HandleRequest(ctx)
    
    API->>Ctx: Get userID
    API->>Cache: CheckCache(ctx, userID)
    
    alt Cache Hit
        Cache-->>API: Return cached data
    else Cache Miss
        API->>DB: QueryDatabase(ctx, userID)
        
        alt Query succeeds
            DB-->>API: Return data
            API->>Cache: Store(ctx, data)
        else Timeout
            Ctx-->>DB: Cancel signal
            DB-->>API: Error: timeout
        end
    end
    
    API-->>Main: Response or Error
    Main->>Ctx: Cleanup (defer cancel)
```

---

## 📖 Conclusion

Context is a fundamental part of Go programming that enables robust concurrent operations. By understanding and properly implementing context:

- ✅ You can gracefully handle cancellations
- ✅ You can prevent resource leaks
- ✅ You can set appropriate timeouts
- ✅ You can pass request-scoped values
- ✅ You can build scalable and responsive applications

Remember to:
1. Always pass context as the first parameter
2. Call `cancel()` using `defer` to prevent leaks
3. Monitor `ctx.Done()` in long-running operations
4. Use specific context types instead of `Background()`
5. Keep context values small and request-scoped

---

## 📝 License

This guide is provided for educational purposes. Feel free to use and share!

## 🤝 Contributing

If you find any issues or have suggestions for improvements, please feel free to contribute!

---

**Happy Coding with Go Context! 🚀**