# Formbricks Hub – Technical Design Proposal

This document describes the tech stack, libraries, and base project structure I’d recommend for building **Formbricks Hub**. The goal is to keep things simple, pragmatic, and easy to iterate on, while still leaving room to scale later.

---

## Key decisions

### Language choice

I looked mainly at Go, Python, Node.js, and Rust.

**Go**

* Fast
* Simple (small language, easy to learn)
* Great standard library and tooling
* Very well suited for backend APIs and services

**Python**

* Simple and expressive
* Huge ecosystem
* Slower runtime
* Async / performance can get tricky at scale
* Great frameworks (Django, FastAPI)

**Node.js**

* Large ecosystem
* Familiar to many devs
* Performance is okay but not amazing
* Dependency and tooling complexity can grow fast
* NPM sucks (new hacks every week)

**Rust**

* Very fast and memory-safe
* Steep learning curve
* Best suited for low-level or very performance-critical systems
* Probably overkill here

**Conclusion:**
Go usually hits the best balance between simplicity, performance, and developer experience. It’s easy to onboard new developers, and the performance characteristics are more than good enough for this kind of system. Any real bottlenecks in Hub are likely to be around data access and architecture rather than raw CPU speed.

---

### Framework / HTTP layer

**Raw Go (`net/http`)**

* Maximum control
* Minimal abstractions
* But also more boilerplate
* Slower to develop and maintain over time

**Gin**

* Very fast
* Mature and well-known
* Mostly focused on routing, not API design

**Huma**

* Simple and lightweight
* API-first design
* OpenAPI is generated directly from code (single source of truth)
* Avoids the usual problem of docs getting out of sync
* Inspired by FastAPI, which has proven to work very well

**Conclusion:**
There’s a saying that *“code is a liability”*, and I generally agree with that. Using a framework that removes boilerplate and enforces good patterns is usually worth it.

Huma’s OpenAPI-first approach is a big win here. OpenAPI is extremely useful, but keeping specs and code in sync manually is painful. Having this built into the framework feels like the right trade-off.

After a brief discussion I ended up going with the RAW go approach.

---

### Database

**PostgreSQL**

* SQL
* Fast and reliable
* JSONB, vector search, and lots of extensions
* Very well understood and widely used

**ClickHouse**

* Optimized for append-only, analytical workloads
* Updates are slow
* No vector search
* Weak data integrity guarantees
* Not a great fit for this use case

**TimescaleDB**

* PostgreSQL-compatible
* Optimized for time-series and analytical queries
* Much better performance on very large datasets (>10M rows)
* Easy to migrate to from Postgres later

**Conclusion:**
Start with **PostgreSQL** to reduce risk and development time. It already fits the use case well and keeps things simple.

If analytics or data volume become a problem later, **TimescaleDB** is an easy upgrade path.

In addition, **Redis** should be used for caching and background jobs / queues.

---

### Data models

The existing data model is kept as requested by the scope.

```sql
CREATE TABLE experience_data (
  id UUID PRIMARY KEY DEFAULT gen_uuidv7(),

  collected_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  source_type VARCHAR NOT NULL,
  source_id VARCHAR,
  source_name VARCHAR,

  field_id VARCHAR NOT NULL,
  field_label VARCHAR,
  field_type VARCHAR NOT NULL,

  value_text TEXT,
  value_number DOUBLE PRECISION,
  value_boolean BOOLEAN,
  value_date TIMESTAMP,
  value_json JSONB,

  metadata JSONB,
  language VARCHAR(10),
  user_identifier VARCHAR
);
```

Additional tables are needed for registering webhooks:

```sql
CREATE TABLE webhooks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  url TEXT NOT NULL,
  name VARCHAR(255),
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE webhook_events (
  webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
  event_type VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (webhook_id, event_type)
);
```

(API key tables are omitted for brevity.)

---

## Third-party libraries / packages

Recommended additions on top of Go itself:

* **Huma** – API framework + OpenAPI generation
* **pgx / sqlc** – PostgreSQL driver and type-safe queries
* **go-redis** – Redis client

---

## Project structure

This is a backend-focused structure. In a monorepo it could live under a `/backend` folder.

Go projects typically follow a layout like this:

```
/cmd
  /api        # API server entrypoint
  /worker     # Background jobs / queues
  /migrate    # Migration runner

/internal
  /api        # HTTP handlers, middleware, routing
  /service    # Business logic
  /repository # Data access layer

/pkg          # Shared utilities (potentially extractable later)
/tests        # Integration tests and test helpers
```

Additional recommended files:

* **Makefile**

  * Build binaries
  * Run tests
  * Run migrations
  * Generate OpenAPI docs

* **docker-compose.yml**

  * PostgreSQL
  * Redis

## System design diagram

                            ┌──────────┐
                            │   User   │
                            └────┬─────┘
                                 │ HTTP / API
                                 ▼
                    ┌─────────────────────────┐
                    │      API Gateway        │
                    │     (Go + Huma)         │
                    └───────┬─────────┬───────┘
                            │         │
                  Fast reads │         │ Writes / queries
                   & cache   │         │
                            ▼         ▼
                   ┌────────────────┐  ┌────────────────────┐
                   │     Redis       │  │    PostgreSQL       │
                   │  - cache        │  │  - experience data │
                   │  - rate limits  │  │  - webhooks        │
                   │  - job queues   │  │  - API keys        │
                   └────────┬───────┘  └─────────┬──────────┘
                            │                    │
                            │ async jobs         │
                            ▼                    │
                 ┌────────────────────┐          │
                 │   Worker Service    │◀─────────┘
                 │ - webhooks delivery │
                 │ - background jobs  │
                 └────────────────────┘

                         (future, optional swap)
                          ┌──────────────────┐
                          │   TimescaleDB    │
                          │ (Postgres-based  │
                          │  analytics DB)   │
                          └──────────────────┘


---

## Final notes

This setup aims to stay boring and predictable. It avoids premature optimization while still being production-ready.
