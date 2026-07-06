# ADR-002: Database Per Service
**Status:** Accepted
## Decision
Each service owns one PostgreSQL database. Cross-service data goes through APIs only.
## Consequences
+ Schema independence, performance isolation
- No cross-service JOINs, eventual consistency
