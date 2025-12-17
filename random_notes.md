# Formbricks Hub work trial

## Task 1: Plan and document tech-stack and basic architecture decisions (3 hours)
Please write up a planning document outlining the important implementation decisions. Please ask Matti all the questions you have in order to come up with this concept.

- What tech stack & programming language would you recommend and why?
- What third party libraries or additional packages are needed and recommended (if any) (e.g. frameworks, ORM, webserver, …)
- What’s the base project structure you would recommend?

## Scope:
- A webserver with PostgreSQL database and full “experience” CRUD functionality as outlined in the documentation
- Same experience data model as the current Hub implementation
- Paginated API with sorting and filtering implemented (depth of filtering capabilities depends on time available and is lower priority)
- API authentication using API key (environment variable)

## Nice to have:
- AI enrichment
- Semantic search
- Webhook functionality
- Tests


---

## Key decisions

### Language choice

- Go
  - Fast
  - Simple (and easy to learn)
  - Lots of support and tooling

- Python
  - Slow
  - Simple (not as simple as go actually)
  - Lots of support and tooling
  - Lots of nice frameworks (Django, FastAPI...)

- Node.js
  - Slow'ish
  - Lots of support and stuff
  - Some nice frameworks (Express...)

- Rust
  - Very fast
  - Huge learning curve, great for high performance systems
  - Nice frameworks like Tokyo and stuff

*Conclusion:* Go typicially gives the best balance of having an easy to use ecosystem, performance, and being a super simple language to learn (https://learnxinyminutes.com/go/). Hub's Bottlenecks should are on data access, performance should come from architecture.


### Framework

- Just raw Go:
  - Kind of simple, kind of not
  - More control
  - Slower to develop (worse documentation, in the future)

- Huma:
  - Simple
  - API first design
    - Single source of truth for OpenAPI (documentation generated from code instead of comments that need to be kept in sync)
  - Router agnostic

- Gin:
  - Very fast

*Conclusion:* I heard someone once say that "code is a liability", and that you should have as little code as possible, so in that sense it's better to use some Framework.
From personal experience OpenAPI is a great tool, and we should leverage it, however it's not always easy to keep it synced with the code, so Huma's approach sounds good... It's also inspired by FastAPI (Python) which I like

### Database

- Postgres
  - SQL support
  - Fast
  - Vector search
  - Lots of extensions
  - Everybody likes Postgres (I'm sure some guys on youtube don't)...

- Clickhouse:
  - Clickhouse has no updates, only adds data
    - Updates are super slow
  - No vector search
  - No data integrity guarantees
  - Good for immutable, logging kind of data

- TimescaleDB:
  - It's PostgresSQL with time-series optimizations
  - Better than PostgreSQL for analytical queries (aggregate, summarize, and analyze large amounts of data)
    - Faster for big (>10M records) tables
  - Easy to port from PostgresSQL in the future (just converting tables)
  - https://www.tigerdata.com/timescaledb

*Conclusion:* Use **PostgreSQL** for now due to lack of familiarity with TimescaleDB and time budget concerns, but consider **TimescaleDB** in the future (migration should be trivial)
**Redis** should also be introduced for caching and job queues.


### Datamodels

Keeping data model form the existing implementation as the scope asks:

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

Additionally we need tables for registering webhooks:

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

...and API keys, but I'll skip that for now.

## Project structure

The following structure is backend focused, so it could go on a "backend" folder if we have a flat repo.
Go projects typicially follow the following structure:

- `/cmd` Multiple entry points (API server, migrations, workers)
- `/internal` Private application code organized by responsibility
  - `/api` HTTP Layer (Handlers, middleware, routing)
  - `/service` Service Layer (Business logic)
  - `/repository` Repository Layer (Data access)
- `/pkg` Reusable utilities that could be extracted later, typically independent from the project... what most developers just call "utils"
- `/tests` Integration tests and test utilities (unit tests should be next to their respective code)

Additionally it should have:
- `makefile` Clear Makefile targets for common tasks, build API, run tests, generate documentation / OpenAPI stuff
- docker compose for local development dependencies (PostgreSQL and Redis)
