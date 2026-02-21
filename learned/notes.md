The architectural flow is 
client->routes(url+http methods)->business layer->service layer->repository layer(DB logic)->db

folder structire followed at enterprise level is 

    ├── cmd/
    │   └── api/
    │       └── main.go                 # Application entry point
    ├── internal/
    │   ├── domain/                      # Core domain
    │   │   ├── entity/                   # Domain entities
    │   │   ├── repository/                # Repository interfaces
    │   │   ├── service/                    # Service interfaces
    │   │   └── valueobject/                 # Value objects
    │   ├── app/                          # Application layer
    │   │   ├── service/                   # Service implementations
    │   │   └── dto/                        # Data Transfer Objects
    │   ├── infra/                        # Infrastructure layer
    │   │   ├── repository/                 # Repository implementations (PostgreSQL)
    │   │   ├── db/                           # Database connection, migrations
    │   │   ├── middleware/                   # Custom middleware
    │   │   ├── validator/                     # Validation setup
    │   │   └── logger/                         # Logging configuration
    │   ├── api/                           # Delivery layer (HTTP)
    │   │   ├── handler/                       # Request handlers
    │   │   ├── router/                         # Route registration
    │   │   └── response/                         # Standard response helpers
    │   └── pkg/                               # Internal shared utilities
    ├── migrations/                         # SQL migration scripts
    ├── configs/                             # Configuration files (if any)
    ├── scripts/                              # Build/deploy scripts
    ├── test/                                  # Integration/e2e tests
    ├── api/                                    # OpenAPI specifications
    ├── go.mod
    └── Makefile

Why this structure is solid

Domain is independent → No DB, no HTTP, no framework

Application layer coordinates logic → Use cases live here

Infrastructure is replaceable → PostgreSQL today, Mongo tomorrow

API layer is thin → Just HTTP mapping & validation

The flow: 

HTTP Request
   ↓
Router
   ↓
Handler
   ↓
Application Service
   ↓
Domain Entity
   ↓
Repository Interface
   ↓
Repository Implementation (Postgres)
   ↓
Database

The flow with folder:

cmd/api/main.go
   ↓
internal/api/router
   ↓
internal/api/handler
   ↓
internal/app/service
   ↓
internal/domain/entity
   ↓
internal/domain/repository (interface)
   ↓
internal/infra/repository (Postgres)
   ↓
PostgreSQL

    routes->only URL and handler mapping
    handler-> HTTP logic(req,res)
    models->DB schema representation
    repository->SQL Queries
    config-> DB and env setup


🧱 Folder-wise Responsibility (One-line Rules)
cmd/

Application entry point

cmd/api/

➡️ Starts the app, loads config, wires dependencies, starts HTTP server
👉 No business logic

internal/

All core application code (cannot be imported by other projects)

🔷 DOMAIN LAYER (Business rules – pure & protected)
internal/domain/entity/

Business models (NOT DB models)

➡️ Represents real-world concepts (User, Order, Payment)
➡️ No JSON / SQL / ORM tags
➡️ Contains business rules

entity → core business objects
internal/domain/valueobject/

Strongly typed domain values

➡️ Email, UserID, Money, PhoneNumber
➡️ Self-validated, immutable

valueobject → validated domain data
internal/domain/repository/

Repository interfaces (contracts)

➡️ Defines what data operations are needed
➡️ Does NOT contain SQL
➡️ Used by services

repository → data access contracts
internal/domain/service/

Domain service interfaces

➡️ Complex business logic spanning multiple entities
➡️ Pure business rules

service → business rules across entities
🔷 APPLICATION LAYER (Use-cases)
internal/app/service/

Use-case implementations

➡️ Orchestrates domain logic
➡️ Calls repositories
➡️ Enforces application rules

service → application workflows
internal/app/dto/

Data Transfer Objects

➡️ Request / response structures
➡️ Converts HTTP input → domain objects

dto → request/response shapes
🔷 INFRASTRUCTURE LAYER (Technical details)
internal/infra/repository/

Repository implementations

➡️ Actual SQL queries
➡️ PostgreSQL / MySQL / Mongo code

repository → SQL & DB logic
internal/infra/db/

Database setup

➡️ DB connection
➡️ Migrations
➡️ Pooling

db → database connection & migrations
internal/infra/middleware/

HTTP middleware

➡️ Auth
➡️ Logging
➡️ CORS
➡️ Rate limiting

middleware → request interception
internal/infra/validator/

Validation setup

➡️ Central validation logic
➡️ Struct validation rules

validator → input validation
internal/infra/logger/

Logging configuration

➡️ Zap / Logrus setup
➡️ Log formats & levels

logger → logging system
🔷 DELIVERY LAYER (HTTP)
internal/api/router/

Routes only

➡️ URL → handler mapping
➡️ No logic

router → URL mapping only
internal/api/handler/

HTTP logic

➡️ Parse request
➡️ Call service
➡️ Return response

handler → req/res handling
internal/api/response/

Standard response helpers

➡️ Success / error formatting
➡️ HTTP status consistency

response → uniform API responses
🔷 SHARED INTERNAL CODE
internal/pkg/

Shared utilities

➡️ Helpers
➡️ Common tools
➡️ Reusable internal code

pkg → shared helpers
🔷 OTHER IMPORTANT FOLDERS
migrations/

SQL migration files

migrations → DB versioning
configs/

Configuration files

➡️ Environment configs
➡️ YAML / JSON / ENV

configs → env & app config
scripts/

Automation scripts

➡️ Build
➡️ Deploy
➡️ CI/CD helpers

scripts → automation
test/

Integration / E2E tests

test → black-box testing
api/

OpenAPI / Swagger specs

api → API documentation



step2:
Mod initalization in  folder
go mod init <foldername>

install required packages like fiber and postgres driver pgx.
go get github.com/gofiber/fiber/v2
go get github.com/jackc/pgx/v5

question: why pgx?
1. faster than database/sql
2. to get native postgres support
3.postgres pgx is prodution level adoption

step3: postgres connection
1.go to config file as it contains configs regarding application
2.so create a folder database.go
3.

pgxpool.Pool: This is a concurrency-safe connection pool that manages a collection of reusable database connections. In contrast to using a single pgx.Conn, a pool allows multiple goroutines (such as those handling web requests) to share and reuse connections efficiently without creating a new connection for every single operation.
* (pointer): The * indicates that DB is a pointer to a pgxpool.Pool instance.
Global/Package-Level Variable: Defining DB outside of any function makes it accessible throughout the entire package (or application, if defined in main). This allows different parts of the application (e.g., different handler functions) to access the same single pool instance

the usage pattern:
The typical pattern involves initializing the DB variable once during application startup and keeping it open for the lifetime of the application
so, var DB *pgxpool.Pool 
pgxpool → a package from github.com/jackc/pgx/v5/pgxpool

Pool → a struct type inside that package

*Pool → a pointer to that struct
DB holds a pointer to a PostgreSQL connection pool

Why a pointer (*) is used?

Because:

A connection pool is a large, stateful object

You don’t want to copy it around

You want all parts of the app to share the same pool

Name starts with	Meaning
Capital letter	Exported (public)
Small letter	Unexported (private)
so here DB is in capital letter
DB is capitalized, it can be used in other packages like import "learned/config"
you can keep the db in samll letter but only if you acess the db inside the same package
package config

var db *pgxpool.Pool

// so in order to pass the access of db connection from this package to another the other package must call this below function to access the conneection
func GetDB() *pgxpool.Pool {
	return db
}
now other package will access it like
config.GetDB().Query(...)

This is actually better design in many cases 👍
Because:

DB cannot be modified directly

You control access

#DSN
DSN stands for Data Source Name.
dsn := "postgres://user:password@localhost:5432/mydb"
It is a single connection string that tells your program:
Where is the database and how should I connect to it?
Breaking the DSN into parts:
| Part          | Meaning                      |
| ------------- | ---------------------------- |
| `postgres://` | Database **type / protocol** |
| `user`        | Database **username**        |
| `password`    | Database **password**        |
| `localhost`   | Database **host**            |
| `5432`        | PostgreSQL **port**          |
| `mydb`        | Database **name**            |

Think of DSN like a home address + key.
Why do we need DSN?
When you call:pgxpool.New(context.Background(), dsn)
The library uses the DSN to:

Locate the database server

Authenticate the user

Select the database

Set up connections

rules:
1.donot use hardcoded dsn in production as code is goin to be pushed to github
2.always  Using environment variables
e.g:dsn := os.Getenv("DATABASE_URL")
e.g2:export DATABASE_URL="postgres://user:password@localhost:5432/mydb"
this makes:
✔ Secure
✔ Configurable
✔ Works across dev / staging / prod

# Other DSN formats (PostgreSQL)
Key-value style
host=localhost port=5432 user=user password=password dbname=mydb
With SSL
postgres://user:password@localhost:5432/mydb?sslmode=disable

# What happens internally when DSN is used?

Parse the DSN

Validate credentials

Open TCP connection

Authenticate with Postgres

Add connection to pool

Ready for queries 🚀



# What does pgxpool.New(...) do?

First, the function signature:

func New(ctx context.Context, connString string) (*pgxpool.Pool, error)

So it returns two things:

*pgxpool.Pool → a PostgreSQL connection pool

error → nil if successful, otherwise an error


# When this line runs:
db, err = pgxpool.New(context.Background(), dsn)
Internally:

dsn is parsed

PostgreSQL server address is resolved

Initial connections are created

Pool configuration is applied

Pool object is allocated in memory

Pointer to pool → assigned to db

If anything fails → err is set

# Why context.Background()?
What is context.Context?

It is used to:

control cancellation

enforce timeouts

pass request-scoped values

context.Background() means:

“No timeout, no cancellation, root context.”

Later you can do:
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

db, err = pgxpool.New(ctx, dsn)



ConnectDB() → initializes the database (one-time setup)

DB() → provides access to the database (used everywhere else)

What does func DB() *pgxpool.Pool mean?

DB() → a getter function

*pgxpool.Pool → it returns a pointer to the DB connection pool


❌ Bad idea
func DB() *pgxpool.Pool {
	return db
}

Then somewhere else:

config.DB().Query(...) // 💥 nil pointer panic (hard to debug)

You’ll get:

panic: runtime error: invalid memory address
No clue why it happened.

With explicit panic
panic("ConnectDB must be called before DB()")

You immediately know:

What went wrong

Where to fix it

This is called fail fast.


Think of db as a car engine:

ConnectDB() → starts the engine

DB() → lets you drive

If engine is off and you try to drive:
PANIC: Start the engine first
Better than crashing silently on the road 


after configuring the db also check for alive
pgxpool.New() validates config
It may not fully verify DB reachability

if err := db.Ping(context.Background()); err != nil {
	return err
}


Why do we even need CloseDB()?

pgxpool.Pool:

Manages multiple open DB connections

Holds sockets, memory, goroutines

Calling:

db.Close()

Closes all connections

Frees resources

Stops background workers

If you don’t:

OS will clean up eventually

BUT graceful shutdown is broken

Long-running apps may leak resources


Where should CloseDB() be called?
✅ BEST PLACE: main()
func main() {
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	defer config.CloseDB()

	startServer()
}
# Why defer?

Runs no matter how main() exits

Clean, idiomatic Go

Guaranteed execution on normal exit

3️⃣ What does NOT trigger CloseDB()?
Scenario	Closed?
HTTP request ends	❌ NO
Handler finishes	❌ NO
Goroutine exits	❌ NO
App keeps running	❌ NO



Graceful shutdown (REAL WORLD)

In real backend services, apps don’t just exit.
They receive OS signals like:

SIGINT → Ctrl+C

SIGTERM → Docker/K8s stop

func main() {
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	defer config.CloseDB()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go startServer()

	<-stop
	log.Println("shutting down...")
}

When signal arrives:

server stops accepting requests

DB closes gracefully

app exits cleanly


When NOT to call CloseDB()

❌ Inside handlers:

func handler() {
	db := config.GetDB()
	defer db.Close() // ❌ VERY WRONG
}

Why?

Pool shared across app

Closing it breaks other goroutines



defer config.GetDB().Close()

This line is very important.

What happens:

config.GetDB()

Checks if DB is initialized

Returns the pool pointer

.Close()

Schedules pool cleanup

Runs when main() exits

So effectively:

“When the application shuts down, close all database connections.”


# Note: always follow this while designing backend system

1.config

2.repository

3.service

4.handler

Each layer depends only on the layer above it, never below it.


Think of your backend as a restaurant:

Layer	Role
config	Kitchen setup (gas, water, tools)
repository	Raw ingredients access
service	Cooking & recipes
handler	Serving customers

You cannot serve food before cooking
You cannot cook without ingredients
You cannot get ingredients without a kitchen


1️⃣ config — Infrastructure Setup (FOUNDATION)
Responsibility

DB connection

Env variables

External clients (Redis, Kafka, etc.)

Why FIRST?

Because everything else needs infrastructure.

config.ConnectDB()

Nothing works without this.

❌ If skipped or mixed

DB logic leaks everywhere

Impossible to change infra later

2️⃣ repository — Data Access Layer
Responsibility. Move ALL database access out of the rest of the app.
After this step:

❌ No SQL in main

❌ No SQL in handlers

❌ No SQL in services

✅ SQL lives ONLY in repositories

This is the single most important structural step.

Why Repository Layer comes NOW (and not service/handler)

Because:

Service logic depends on data

Handler logic depends on business rules

Repository is the first consumer of your DB

If you skip repository:

SQL spreads everywhere

You lose control immediately

So repository must be built first.

SQL

DB schema knowledge

Persistence logic

userRepo.GetByID(ctx, id)
Why SECOND?

Because:

Business logic should NOT know SQL

Storage is a lower-level concern

❌ If skipped
handler -> SQL -> business logic

You get:

duplicated queries

untestable code

schema changes breaking everything

3️⃣ service — Business Rules (THE CORE)
Responsibility

Validation

Business decisions

Orchestration

userService.RegisterUser(...)
Why THIRD?

Because:

Services depend on repositories

Services define “what the app DOES”

❌ If skipped

Business logic goes into handlers

HTTP becomes tightly coupled to logic

Cannot reuse logic (CLI, cron, jobs)

4️⃣ handler — Delivery Layer (EDGE)
Responsibility

HTTP / gRPC / CLI

Parsing input

Formatting output

POST /users
Why LAST?

Because:

It is just a transport mechanism

Should be replaceable

Today: HTTP
Tomorrow: gRPC
Later: Message queue
Business logic stays untouched.


# Now Repository layer at 2nd level
🎯 What is the goal of this step?

Move ALL database access out of the rest of the app.

After this step:

❌ No SQL in main

❌ No SQL in handlers

❌ No SQL in services

✅ SQL lives ONLY in repositories

This is the single most important structural step.

🧠 Why Repository Layer comes NOW (and not service/handler)

Because:

Service logic depends on data

Handler logic depends on business rules

Repository is the first consumer of your DB

If you skip repository:

SQL spreads everywhere

You lose control immediately

So repository must be built first.

We always start with one domain, not everything.
repository/
 └── user_repository.go

 Context is passed from above
ctx context.Context

This allows:

request cancellation

timeouts

tracing


Why these imports?

pgxpool → to execute queries

context → to control cancellation, timeouts, tracing

⚠️ Repository methods must accept context.Context
This allows:

request cancellation

graceful shutdown

observability



type UserRepository struct { ... }
type UserRepository struct {
	db *pgxpool.Pool
}
What is this?

This is a struct that wraps the database pool.

Think of it as:

“A specialized object whose job is to fetch/store Users.”

Why store db inside a struct?

❌ BAD approach:

func GetUserByID(...) {
	config.GetDB().Query(...)
}

Problems:

Hard to test

Tight coupling

SQL everywhere

✅ GOOD approach:
type UserRepository struct {
	db *pgxpool.Pool
}

Benefits:

DB access is encapsulated

Easy to mock

Easy to replace DB

Clear ownership

4️⃣ NewUserRepository(...) — the MOST IMPORTANT PART
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

This is a constructor function.

What is happening step-by-step?

Function receives a DB pool

Creates a new UserRepository

Stores the DB inside it

Returns a pointer to it

5️⃣ Why is this called Dependency Injection?

Because:

Repository does NOT create the DB

DB is injected from outside

main()
  ↓
config.ConnectDB()
  ↓
config.GetDB()
  ↓
NewUserRepository(db)

This means:

Repository depends on an abstraction (DB)


Not on how DB is created



Why domain exists (important)

Domain models must not belong to infra (DB) or transport (HTTP).

If User lives in:

repository ❌ → tied to DB

handler ❌ → tied to HTTP

service ❌ → tied to orchestration

So we give it a neutral home.
Create domain/user.go here user is an entity
package domain

type User struct {
	ID    int
	Name  string
	Email string
}
That’s it.
No DB tags.
No JSON tags.
No framework imports.

👉 Pure business concept.

Before (what you might have)
// ❌ User defined inside repository
type User struct {
	ID    int
	Name  string
	Email string
}

This tightly couples repository with domain.

🧠 What just happened (this is important)
Repository:

Knows how to fetch data

Maps DB rows → domain.User

Domain:

Knows what a User is

Service:

Will later decide what to do with User

Handler:

Will later decide how to expose User

# Dependency direction (VERY IMPORTANT)
domain        ← depends on nothing
repository    → depends on domain
service       → depends on domain + repository(interface)
handler       → depends on domain + service
config        → depends on nothing


# What exactly belongs in user_repository
✔ Repository methods answer questions like:

How do I fetch a user?

How do I insert a user?

How do I update a user row?

How do I delete a user?

What SQL should be used?

How to map DB rows → domain.User?

❌ What does NOT belong in repository

Validation rules

Business decisions

Authorization

HTTP logic

JSON parsing

Those belong to service or handler.

# Typical UserRepository interface (clean & complete)
type UserRepository interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

This interface defines WHAT can be done, not HOW.

but uou can implement without interface with concrete implementation

That includes:

GetByID

Create

Update

Delete

List

Patch (partial update)

# create User
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
	).Scan(&user.ID)
}
Notice:

Repository does SQL

Repository mutates domain.User

Repository does no validation


# 3️⃣ func NewUserRepository(...) UserRepository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}
Create a concrete Postgres-backed repository and return it as a UserRepository.”

So although the real object is *userRepository,
the type exposed to the caller is UserRepository.

This line is VERY IMPORTANT.

What it does step-by-step

Takes a DB pool (dependency)

Creates a concrete userRepository

Returns it as a UserRepository interface

What does it take as input?
db *pgxpool.Pool

This is dependency injection.

Meaning:

Repository does NOT create the DB

DB is given from outside

Repository only uses it

This keeps layers decoupled.
3️⃣ What does it return?
UserRepository -> it's an interface that mean “I promise you something that behaves like a UserRepository." 
It does NOT return:*userRepository

return &userRepository{db: db}
This creates:
A concrete struct: userRepository
With its db field set
So in memory:
&userRepository
 ├── db -> *pgxpool.Pool


# Universal method template in your head

Every repository method follows this shape:

1. Accept context
2. Build SQL
3. Execute SQL
4. Map result → domain model
5. Return error or result




1️⃣ Why ctx context.Context exists

(and how you were writing code earlier without knowing it)

First: what problem does context solve?

Control + cancellation + deadlines across function boundaries

Backend systems are not linear programs:

HTTP requests can be cancelled

Timeouts must stop DB queries

Servers shut down while work is running

Tracing IDs must flow everywhere

Without context, you cannot control this.

What context.Context actually is

Conceptually:

context = request lifecycle controller

It carries:

cancellation signal

deadline / timeout

request-scoped values (trace IDs, auth info)

Why it appears in repository methods
func (r *userRepository) GetByID(ctx context.Context, id int)

This means:

“This DB operation belongs to some request or process.
If that request dies, stop the DB work.”

What happens internally when context is cancelled

Example:

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

repo.GetByID(ctx, 1)

Internally:

pgx listens to ctx.Done()

If timeout hits:

DB query is aborted

connection is returned to pool

error is returned

⚠️ Without context:

DB keeps running

Resources leak

System becomes unstable

Why you didn’t see context earlier

Because:

Small scripts don’t need cancellation

Tutorials hide complexity

Real backend systems cannot afford to ignore it

Context appears when code becomes production-grade.

2️⃣ Why user *domain.User (pointer, domain)

Let’s break this into two questions.

A) Why domain.User?

Because repository returns business entities, not DB rows.

Repository’s job:

DB rows → domain entity

Not:

DB rows → HTTP JSON

This keeps layers clean.

B) Why pointer *domain.User?
Reason 1: Efficiency

Passing structs by value copies memory.

Reason 2: Mutability (IMPORTANT)

Repository modifies the user:

user.ID = generatedID

That change must be visible to caller.

Reason 3: Nil semantics

You can return:

nil, err

You can’t do that with values.

Rule you should remember

Repositories return pointers to domain entities.

This is idiomatic Go.

3️⃣ What is r.db?
r.db

This is:

*pgxpool.Pool

Meaning:

A connection pool

Not a single DB connection

Safe for concurrent use

Internally:

Maintains multiple open connections

Hands out one connection per query

Returns it automatically

You never manage connections manually.

4️⃣ What is QueryRow, Exec, Query REALLY doing?

This is the part most people never explain properly.

Exec — fire-and-forget
_, err := r.db.Exec(ctx, query, args...)

Internally:

Borrow a connection from pool

Send SQL to Postgres

Execute command

Return rows affected

Release connection

Used for:

UPDATE

DELETE

INSERT (when no return needed)

QueryRow — exactly ONE row expected
row := r.db.QueryRow(ctx, query, args...)
row.Scan(...)

Internally:

Borrow connection

Execute SQL

Expect one row

Hold result cursor

Scan pulls columns into variables

Release connection after scan

Used when:

SELECT ... WHERE id = ?

INSERT ... RETURNING id

⚠️ If zero rows → error
⚠️ If multiple rows → only first is used

Query — multiple rows
rows, err := r.db.Query(ctx, query)

Internally:

Borrow connection

Execute SQL

Return a cursor

You iterate row-by-row

You must rows.Close() → connection returned

Used for:

SELECT * FROM users

Why you didn’t see these before

Because:

ORMs hide them

Tutorials abstract them away

Raw SQL drivers expose them

You are now writing lower-level, more powerful code.

5️⃣ How all this works TOGETHER (big picture)

When this runs:

user, err := repo.GetByID(ctx, 1)

What actually happens:

1. HTTP request comes in
2. Handler creates ctx
3. ctx passed → service
4. service passes ctx → repository
5. repository calls QueryRow(ctx)
6. pgx checks ctx:
   - cancelled? → stop
   - alive? → proceed
7. DB executes query
8. Row scanned into domain.User
9. Connection returned to pool
10. Control returns upward

This is request-scoped execution.

6️⃣ The mental model you should keep

Whenever you see this:

ctx context.Context

Think:

“Who owns this work, and when should it stop?”

Whenever you see this:

user *domain.User

Think:

“This is a business entity, not a DB row.”

Whenever you see this:

r.db.QueryRow