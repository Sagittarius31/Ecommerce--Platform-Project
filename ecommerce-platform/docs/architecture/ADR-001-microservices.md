# ADR-001: Microservices Architecture

**Date:** 2024-01-01 | **Status:** Accepted

## Context
Need a platform that scales independently per domain and allows independent deployment.

## Decision
6 independent Go microservices, each owning its data and deploying independently.

## Consequences
**Positive:** Independent scaling, fault isolation, technology flexibility.
**Negative:** Network latency, distributed debugging, eventual consistency.
