# Golang Interview: Struct - Complete Guide

> A comprehensive guide to Go struct interview questions with code examples and key concepts.

---

## Table of Contents

1. [What is a Struct?](#1-what-is-a-struct)
2. [What is an Anonymous Struct?](#2-what-is-an-anonymous-struct)
3. [Struct Embedding in Go](#3-struct-embedding-in-go)
4. [Pointer to a Struct](#4-pointer-to-a-struct)
5. [Methods for a Struct](#5-methods-for-a-struct)
6. [Comparing Structs with == Operator](#6-comparing-structs-with--operator)
7. [Zero Value of a Struct](#7-zero-value-of-a-struct)
8. [Struct vs Map](#8-struct-vs-map)
9. [Struct Tags](#9-struct-tags)
10. [Deep Copy vs Shallow Copy](#10-deep-copy-vs-shallow-copy)
11. [Unexported Fields](#11-unexported-fields)
12. [Constructor-like Functions](#12-constructor-like-functions)
13. [Immutable Structs](#13-immutable-structs)
14. [Size of a Struct](#14-size-of-a-struct)

---

## 1. What is a Struct?

**Question:** What is a struct in Go and how do you define one?

**Key Description:**
- A struct is a composite data type that groups variables (fields) under a single name
- Represents a collection of related data fields
- Fundamental building block for organizing and modeling data in Go programs

**Example Code:**

```go
// Define a struct named 'Person' with fields 'Name' and 'Age'
type Person struct {
    Name string
    Age  int
}

// Create an instance of the 'Person' struct
personInstance := Person{
    Name: "John Doe",
    Age:  25,
}

// Accessing and modifying struct fields
fmt.Println(personInstance.Name) // Output: John Doe
fmt.Println(personInstance.Age)  // Output: 25

// Modifying a struct field
personInstance.Age = 26
fmt.Println(personInstance.Age)  // Output: 26
```

---

## 2. What is an Anonymous Struct?

**Question:** What is an anonymous struct, and in what situations might you use one?

**Key Description:**
- A struct without a predefined name, created on the fly
- Used for short-lived or temporary data structures
- Defined and used in a single place without being assigned a name

**Situations to Use:**
- Temporary data structures
- Inline initialization
- Function parameters
- Marshaling and unmarshaling JSON

**Example Code:**

```go
// Anonymous struct definition
var person = struct {
    Name string
    Age  int
}{
    Name: "John Doe",
    Age:  25,
}

// Temporary data structure
data := struct {
    Status  string
    Message string
}{
    Status:  "OK",
    Message: "Data retrieved successfully",
}

// Inline initialization
result := struct{ Code int }{Code: 200}

// Function parameter with anonymous struct
func processData(data struct{ ID int; Name string }) {
    fmt.Println("Processing data:", data)
}

func main() {
    processData(struct{ ID int; Name string }{ID: 1, Name: "Alice"})
}

// JSON unmarshaling
var result struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}
jsonStr := `{"status": "OK", "message": "Success"}`
json.Unmarshal([]byte(jsonStr), &result)
```

---

## 3. Struct Embedding in Go

**Question:** Explain the concept of struct embedding in Go. Provide an example.

**Key Description:**
- Achieves composition by including one or more fields of one struct type within another
- Embedded struct's fields and methods become part of the outer struct
- Allows the outer struct to "inherit" behavior (composition over inheritance)
- Sometimes called "anonymous fields"

**Key Points:**
- Fields from embedded struct are directly accessible from outer struct instances
- If conflicts exist, field from outer struct takes precedence
- Promotes code reuse and structuring of related types

**Example Code:**

```go
package main

import "fmt"

// Animal struct representing a basic animal
type Animal struct {
    Name string
}

// Dog struct embeds the Animal struct
type Dog struct {
    Animal        // Embedding the Animal struct
    Breed  string // Additional field specific to Dog
}

// Cat struct embeds the Animal struct
type Cat struct {
    Animal        // Embedding the Animal struct
    Color string  // Additional field specific to Cat
}

func main() {
    // Creating instances of Dog and Cat
    dog := Dog{
        Animal: Animal{Name: "Buddy"},
        Breed:  "Labrador",
    }

    cat := Cat{
        Animal: Animal{Name: "Whiskers"},
        Color:  "Gray",
    }

    // Accessing fields from embedded struct
    fmt.Printf("Dog's name: %s\n", dog.Name) // Accessing Animal's field directly
    fmt.Printf("Cat's name: %s\n", cat.Name) // Accessing Animal's field directly

    // Accessing fields specific to each struct
    fmt.Printf("Dog's breed: %s\n", dog.Breed)
    fmt.Printf("Cat's color: %s\n", cat.Color)
}
```

---

## 4. Pointer to a Struct

**Question:** How do you create a pointer to a struct? What benefits does using a pointer to a struct offer?

**Key Description:**
- Use `&` operator followed by struct instance to obtain its memory address
- Allows efficient memory management and modification of original struct

**Benefits:**
1. **Efficient Memory Management** - More memory-efficient for large structs
2. **Modify Original Struct** - Changes through pointer affect underlying struct
3. **Avoiding Copying** - No need to copy entire struct when passing to functions
4. **Nil Pointer Check** - Can check for uninitialized pointers
5. **Sharing Data** - Allows sharing and modifying same data between functions

**Example Code:**

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    // Create an instance of the 'Person' struct
    personInstance := Person{
        Name: "John Doe",
        Age:  25,
    }

    // Create a pointer to the 'Person' struct
    personPointer := &personInstance

    // Accessing struct fields through the pointer
    fmt.Println("Name:", personPointer.Name)
    fmt.Println("Age:", personPointer.Age)

    // Modify the original struct through the pointer
    personPointer.Age = 26
    fmt.Println("Updated Age:", personInstance.Age) // Output: 26

    // Check if the pointer is nil before dereferencing
    if personPointer != nil {
        fmt.Println("Name:", personPointer.Name)
    }
}

// Initializing a struct with 'new'
newPerson := new(Person)
newPerson.Name = "Bob"
newPerson.Age = 32
```

---

## 5. Methods for a Struct

**Question:** Can you define methods for a struct in Go? How are they different from regular functions?

**Key Description:**
- Methods are functions associated with a particular type
- Declared with a receiver (similar to "this" or "self" in other languages)
- Allow associating behavior with a specific type

**Key Points:**
- **Receiver:** Parameter specifying the type the method operates on
- **Value Receiver:** Operates on a copy of the struct
- **Pointer Receiver:** Operates on the original struct
- Methods are called using `instance.method()` syntax

**Example Code:**

```go
package main

import "fmt"

// Define a struct named 'Rectangle'
type Rectangle struct {
    Width  float64
    Height float64
}

// Method with value receiver
func (r Rectangle) area() float64 {
    return r.Width * r.Height
}

// Method with value receiver
func (r Rectangle) perimeter() float64 {
    return 2*r.Width + 2*r.Height
}

// Method with pointer receiver (can modify the struct)
func (r *Rectangle) scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rectangle := Rectangle{Width: 5, Height: 3}

    fmt.Println("Area:", rectangle.area())           // Output: 15
    fmt.Println("Perimeter:", rectangle.perimeter()) // Output: 16

    rectangle.scale(2)
    fmt.Println("Scaled Area:", rectangle.area())    // Output: 60
}
```

---

## 6. Comparing Structs with == Operator

**Question:** Is it possible to compare two structs in Go using the `==` operator? Why or why not?

**Key Description:**
- You CAN compare structs if all fields are of comparable types
- Comparison will result in compilation error if struct contains non-comparable types (slices, maps, functions)

**Example Code:**

```go
package main

import "fmt"

// Comparable struct
type Person struct {
    Name string
    Age  int
}

func main() {
    person1 := Person{Name: "John", Age: 30}
    person2 := Person{Name: "John", Age: 30}

    // This works - all fields are comparable
    if person1 == person2 {
        fmt.Println("The structs are equal.")
    } else {
        fmt.Println("The structs are not equal.")
    }
}

// Non-comparable struct (will cause compilation error)
type ExampleStruct struct {
    Numbers []int  // Slices are NOT comparable
}

func failingComparison() {
    struct1 := ExampleStruct{Numbers: []int{1, 2, 3}}
    struct2 := ExampleStruct{Numbers: []int{1, 2, 3}}

    // This will NOT compile!
    // if struct1 == struct2 { } // Error: invalid operation
}
```

---

## 7. Zero Value of a Struct

**Question:** What is the zero value of a struct in Go? Does it depend on the values of its fields?

**Key Description:**
- Zero value is a struct with ALL fields set to their respective zero values
- Zero value for each field is determined by its data type
- Not dependent on values assigned during initialization

**Zero Values by Type:**
| Type | Zero Value |
|------|------------|
| Numeric (int, float64) | `0` |
| String | `""` (empty string) |
| Boolean | `false` |
| Pointer | `nil` |
| Struct | Each field set to its zero value |

**Example Code:**

```go
package main

import "fmt"

type ExampleStruct struct {
    Integer      int
    Float        float64
    Text         string
    IsActive     bool
    Pointer      *int
}

func main() {
    // Create instance with no explicit values
    var zeroStruct ExampleStruct

    // All fields have their zero values
    fmt.Println("Integer:", zeroStruct.Integer)     // Output: 0
    fmt.Println("Float:", zeroStruct.Float)         // Output: 0
    fmt.Println("Text:", zeroStruct.Text)           // Output: ""
    fmt.Println("IsActive:", zeroStruct.IsActive)   // Output: false
    fmt.Println("Pointer:", zeroStruct.Pointer)     // Output: <nil>
}
```

---

## 8. Struct vs Map

**Question:** When would you choose to use a struct over a map in Go, and vice versa?

**Key Description:**

### Structs

| Advantages | Disadvantages |
|------------|---------------|
| Type safety at compile-time | Fixed structure |
| Faster field access | Needs initialization with all fields |
| More readable and self-documenting | |

### Maps

| Advantages | Disadvantages |
|------------|---------------|
| Dynamic structure | No compile-time type safety |
| Flexible - easily extended | Slower lookups |
| Concise for many fields | Less readable |

**Key Differences:**

| Aspect | Struct | Map |
|--------|--------|-----|
| Syntax | Fixed with named fields | Dynamic key-value pairs |
| Initialization | `var` or composite literal | `make` function |
| Access | Dot notation (`myStruct.field`) | Brackets (`myMap[key]`) |
| Type Safety | Compile-time | Runtime |

**Guidelines:**
- Use **structs** when data has fixed structure and type safety is essential
- Use **maps** when dealing with dynamic/unknown data structures

---

## 9. Struct Tags

**Question:** Explain the purpose of struct tags in Go. Provide an example where struct tags are commonly used.

**Key Description:**
- Metadata associated with struct fields
- Specified as string literal in backticks after field declaration
- Used for serialization, database mapping, and validation

**Common Uses:**
1. **Serialization/Deserialization** - JSON, XML encoding
2. **Database Mapping** - ORM column mapping
3. **Validation** - Field constraints and rules

**Example Code:**

```go
package main

import (
    "encoding/json"
    "fmt"
)

// Struct with tags for JSON serialization
type Person struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email,omitempty"` // omitempty: skip if empty
}

func main() {
    person := Person{
        Name:  "John Doe",
        Age:   30,
        Email: "john@example.com",
    }

    // Serialize to JSON
    jsonData, err := json.Marshal(person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(string(jsonData))
    // Output: {"name":"John Doe","age":30,"email":"john@example.com"}
}
```

---

## 10. Deep Copy vs Shallow Copy

**Question:** When you assign one struct variable to another in Go, is it a deep copy or a shallow copy?

**Key Description:**
- Assignment performs a **shallow copy**
- Individual fields receive copies of values
- Reference types (slices, maps, pointers) - only references are copied, not underlying data

**Example Code:**

```go
package main

import "fmt"

// Shallow Copy Example
type Person struct {
    Name string
    Age  int
}

func shallowCopyExample() {
    person1 := Person{Name: "Alice", Age: 25}
    person2 := person1  // Shallow copy

    person1.Name = "Bob"

    fmt.Println("Person 1:", person1) // Output: {Bob 25}
    fmt.Println("Person 2:", person2) // Output: {Alice 25}
}

// Deep Copy Example (with slices)
type PersonWithHobbies struct {
    Name    string
    Age     int
    Hobbies []string
}

func deepCopyExample() {
    person1 := PersonWithHobbies{
        Name:    "Alice",
        Age:     25,
        Hobbies: []string{"Reading", "Traveling"},
    }

    // Perform deep copy manually
    person2 := PersonWithHobbies{
        Name:    person1.Name,
        Age:     person1.Age,
        Hobbies: make([]string, len(person1.Hobbies)),
    }
    copy(person2.Hobbies, person1.Hobbies)

    // Modify original
    person1.Name = "Bob"
    person1.Hobbies[0] = "Swimming"

    fmt.Println("Person 1:", person1) // Modified
    fmt.Println("Person 2:", person2) // Unchanged
}
```

---

## 11. Unexported Fields

**Question:** Explain the concept of unexported fields in a struct. How are they accessed within and outside of the package?

**Key Description:**
- **Exported (Public):** Fields starting with uppercase letter - accessible outside package
- **Unexported (Private):** Fields starting with lowercase letter - only accessible within same package

**Example Code:**

```go
// In package mypackage
package mypackage

import "fmt"

type Person struct {
    FirstName string // Exported field
    lastName  string // Unexported field
    Age       int    // Exported field
}

// Within the same package - both fields accessible
func Example() {
    person := Person{
        FirstName: "John",
        lastName:  "Doe",
        Age:       30,
    }

    fmt.Println("First Name:", person.FirstName)
    fmt.Println("Last Name:", person.lastName)  // Works!
    
    person.lastName = "Smith"  // Can modify
}

// In package main (different package)
package main

import "mypackage"

func main() {
    person := mypackage.Person{
        FirstName: "Alice",
        // lastName:  "Doe", // ERROR: cannot refer to unexported field
        Age: 25,
    }
}
```

---

## 12. Constructor-like Functions

**Question:** Go doesn't have traditional constructors. How can you create a constructor-like function for a struct?

**Key Description:**
- Go uses factory functions (conventionally named `NewTypeName`)
- Returns a pointer to a new instance with initialized fields
- Encapsulates struct creation details

**Example Code:**

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

// Constructor-like function
func NewPerson(name string, age int) *Person {
    return &Person{
        Name: name,
        Age:  age,
    }
}

func main() {
    // Use constructor-like function
    person := NewPerson("Alice", 25)

    fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
```

---

## 13. Immutable Structs

**Question:** Can you create an immutable struct in Go? Explain the concept and provide an example.

**Key Description:**
- Go structs are mutable by default
- Achieve immutability through convention:
  - Use unexported fields
  - Provide only getter methods (no setters)
  - Use constructor function for initialization

**Example Code:**

```go
package main

import "fmt"

// "Immutable" struct with unexported fields
type Person struct {
    firstName string
    lastName  string
    age       int
}

// Constructor
func NewPerson(firstName, lastName string, age int) *Person {
    return &Person{
        firstName: firstName,
        lastName:  lastName,
        age:       age,
    }
}

// Getter methods only (no setters)
func (p Person) FirstName() string {
    return p.firstName
}

func (p Person) LastName() string {
    return p.lastName
}

func (p Person) Age() int {
    return p.age
}

func main() {
    person := NewPerson("John", "Doe", 30)

    // Read-only access
    fmt.Println("First Name:", person.FirstName())
    fmt.Println("Last Name:", person.LastName())
    fmt.Println("Age:", person.Age())

    // Cannot modify - compilation error:
    // person.firstName = "Jane" // Error: cannot assign to unexported field
}
```

---

## 14. Size of a Struct

**Question:** How is the size of a struct determined in Go? Are there any factors that affect the size?

**Key Description:**
- Size = sum of sizes of all fields
- May include padding for memory alignment

**Factors Affecting Size:**
1. **Field Types** - Each type has specific size (int: 4/8 bytes, etc.)
2. **Alignment** - Compiler adds padding for memory access optimization
3. **Padding** - Added at end to ensure proper alignment
4. **Field Order** - Reordering fields may change total size due to padding

**Example Code:**

```go
package main

import (
    "fmt"
    "unsafe"
)

type ExampleStruct struct {
    BoolValue    bool    // 1 byte
    Int32Value   int32   // 4 bytes
    Float64Value float64 // 8 bytes
}

func main() {
    example := ExampleStruct{}

    // Calculate size using unsafe package
    size := unsafe.Sizeof(example)

    fmt.Printf("Size of ExampleStruct: %d bytes\n", size)
    // Note: Actual size may be larger due to padding
}
```

---

## Quick Reference Summary

| Topic | Key Concept |
|-------|-------------|
| Struct Definition | `type Name struct { fields }` |
| Anonymous Struct | Struct without name for temporary use |
| Struct Embedding | Composition by including another struct |
| Pointer to Struct | Use `&` for address, enables modification |
| Methods | Functions with receiver attached to types |
| Struct Comparison | Only works if all fields are comparable |
| Zero Value | All fields set to their type's zero value |
| Struct vs Map | Struct for fixed data, Map for dynamic |
| Struct Tags | Metadata in backticks for serialization |
| Deep vs Shallow Copy | Assignment = shallow copy |
| Unexported Fields | Lowercase = private to package |
| Constructor | `NewTypeName()` function pattern |
| Immutability | Unexported fields + getters only |
| Struct Size | Sum of fields + alignment padding |

---

## Author

Based on the article by **Kiran Adhikari** - [Original Medium Article](https://medium.com/@kiruu1238/golang-interview-struct-758db20ef4c9)

---

*Good luck with your Go interviews!* 🚀