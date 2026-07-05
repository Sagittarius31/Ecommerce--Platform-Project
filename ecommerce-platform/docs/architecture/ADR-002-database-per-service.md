# ADR-002: Database Per Service

**Date:** 2024-01-01 | **Status:** Accepted

## Context
Services sharing a database create schema coupling and performance interference.

## Decision
Each service owns one PostgreSQL database. Cross-service data goes through APIs only.

## Consequences
**Positive:** Schema independence, performance isolation.
**Negative:** No cross-service JOINs, eventual consistency, more infrastructure.
