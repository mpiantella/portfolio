# Interfaces (adapters) layer

Purpose: contain adapters that translate between the external world (HTTP/CLI/GRPC) and the application's use cases. Handlers should be thin, focus on validation, context/timeouts, and mapping transport-level concerns to use case input/output.

Contains:
- HTTP handlers for web UI and JSON APIs

Guidelines:
- Keep handlers small and testable; delegate business logic to use cases.
- Implement request/response mapping and status codes here.
