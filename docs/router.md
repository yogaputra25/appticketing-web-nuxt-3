## Router Patterns (Chi)

### No-Fallthrough Rule

Chi uses a **radix trie** for route matching. When a request enters a subroute
(e.g. `/admin/events`), chi matches **only** handlers declared inside that
subroute. It does **not** fall through to sibling subroutes on miss.

**Rule:** Event-scoped admin routes (paths starting with `/admin/events/{id}/...`)
MUST be declared inside `r.Route("/admin/events", ...)`, NOT in the outer
`r.Route("/admin", ...)` block.

```go
// ✅ CORRECT — inside /admin/events subroute
r.Route("/admin/events", func(r chi.Router) {
    r.Use(jwtMgr.Authenticator, auth.RequireAdmin)
    r.Get("/", eventH.ListAdmin)
    r.Post("/", eventH.Create)
    r.Get("/{id}", eventH.DetailAdmin)
    r.Put("/{id}", eventH.Update)
    r.Delete("/{id}", eventH.Delete)
    r.Post("/{id}/publish", eventH.Publish)
    r.Post("/{eventId}/categories", catH.Create)  // ← event-scoped
})

// ❌ WRONG — will 404 because chi never reaches this sibling subroute
r.Route("/admin", func(r chi.Router) {
    r.Use(jwtMgr.Authenticator, auth.RequireAdmin)
    r.Put("/categories/{id}", catH.Update)
    r.Post("/events/{eventId}/categories", catH.Create)  // ← NEVER MATCHED
})
```

### Smoke Test

Run the smoke test to verify all routes are registered correctly:

```bash
cd apps/api && make smoke
```

Requires a running API (default: `http://localhost:8080`). Override with:

```bash
API_BASE=http://localhost:8081 make smoke
```
