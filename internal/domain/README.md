# Domain layer

Purpose: define business entities, value objects and domain-specific logic. This layer has no external dependencies and exposes interfaces used by higher layers. Keep it free from framework or transport-specific concerns.

Contains:
- `project` entity and repository interface

Guidelines:
- No imports from `internal/usecase` or `internal/interfaces` or `internal/infrastructure`.
- Keep domain logic independent and pure for easy testing.
