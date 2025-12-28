# UseCase (application) layer

Purpose: implement application-specific business rules and orchestrate domain entities. Use cases are triggered by interfaces (HTTP handlers, CLI) and call domain interfaces for persistence or external operations.

Contains:
- `project` use cases (listing projects)

Guidelines:
- Accept interfaces for dependencies to follow dependency inversion.
- Keep business rules here, not in handlers or infrastructure.
