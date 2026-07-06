# ADR-001: Microservices Architecture
**Status:** Accepted
## Decision
6 independent Go microservices, each owning its data.
## Consequences
+ Independent scaling, fault isolation, technology flexibility
- Network latency, distributed debugging, eventual consistency
