---
name: pre-commit-review
description: Reviews staged or unstaged git changes for bugs, security issues, and performance problems in this Go + Gin SMS/WhatsApp messaging API. Generates conventional commit messages when no issues are found. Use when the user says "Review commit", "Check before commit", or asks to review changes before committing.
---

# Pre-Commit Review (MessageApi — Go + Gin)

## Trigger

When the user says "Review commit" or "Check before commit":

1. Run `git diff --cached` to get staged changes.
2. If output is empty, run `git diff` for unstaged changes.
3. If both are empty, respond: "No changes to review. Stage files with `git add` or make edits first."

## Project Context

This is a Go REST API built with **Gin** that sends SMS and WhatsApp messages through multiple providers.

- Module path: `messageapi.e-vrit.co.il/...`
- Framework: `github.com/gin-gonic/gin`
- Config: `github.com/joho/godotenv` with `.env.<env>` files (e.g. `.env.qa`, `.env.prod`)
- Tests: `github.com/stretchr/testify/assert`, table-driven

### Architecture

```
routes/routes.go          → registers Gin routes
services/<domain>/httpHandler.go  → Gin handler, parses request, auth, dispatches to provider
services/<domain>/iService.go     → provider-agnostic interface (e.g. ISms)
services/<domain>/<provider>/     → concrete providers (e.g. inforu/, flashy/)
services/<domain>/models/         → shared request/response models
enums/                            → shared constants (e.g. SendingType)
middleware/                       → Gin/HTTP middleware
```

The current messaging flow (reference for new work):

- `routes.SmsRoutes` → `smsmessage.SendSms` → `ISms.SendSms` / `ISms.SendWhatsApp`
- Provider selection: `sms.SendingType` (enum `enums.Inforu` / `enums.Flashy`)
- Channel selection: query `channel=sms|whatsapp`

## Review Process

Act as a senior Go developer familiar with Gin and external HTTP integrations. Analyze changes for:

### Logical Bugs

- Unhandled errors — never ignore `err`; check at call site (watch for `_, _ = client.Do(req)` style).
- **Panic-prone header access** like `ctx.Request.Header["Token"][0]` — use `ctx.GetHeader("Token")` to avoid index-out-of-range when the header is missing.
- **`log.Fatalf` / `os.Exit(1)` inside request-handling code paths** (e.g. inside `NewInforuModel`, `SendSms`). These crash the server on a single bad request and must be replaced with returned errors.
- Nil pointer dereference after `resp, err := client.Do(req)` — always check `err` before `defer resp.Body.Close()`.
- Missing `defer resp.Body.Close()` after a successful HTTP response.
- Type assertions on decoded JSON without the `, ok` form (e.g. `response["StatusId"].(float64)` without the comma-ok).
- Missing edge cases: empty `Message`, empty `ReciverPhoneNumber`, `TemplateId == 0` for WhatsApp, empty `Recipients` slice.
- Incorrect provider dispatch — new `SendingType` values must be added to the `switch` in `httpHandler.go` **and** wired through the `ISms` interface.
- Off-by-one / unbounded loops when iterating `RecipientCustomFields` or `TemplateParameters`.

### Security (Go + Gin + external HTTP specific)

- **Auth bypass**: every protected endpoint must validate the `Token` header against `os.Getenv("TOKEN")` using `ctx.GetHeader` and a constant-time compare (`subtle.ConstantTimeCompare`) when feasible. New routes in `routes/` without this check are a critical finding.
- **Hardcoded secrets**: API keys (`FLASHY_KEY`, `INFORU_AUT`), `TOKEN`, `Authorization: Basic ...` headers must come from env vars, never string literals in `.go` files. Test files should also avoid real credentials — reference env vars or mock them.
- **Sensitive data in logs**: never log the full `Authorization` header, `x-api-key`, `Token`, or full recipient phone numbers / message bodies at info level. Prefer redacting.
- **SSRF / URL injection**: provider URLs (`InforuUrl`, `FlashyUrl`) must come from env, never from request input.
- **Input validation before outbound HTTP**: validate `ReciverPhoneNumber` (digits, length), `SenderName`, and `Message` length before calling the provider — both to prevent abuse and to catch user errors early.
- **Missing request body size limits** on Gin routes that could receive large payloads (`router.MaxMultipartMemory`, `ctx.Request.Body = http.MaxBytesReader(...)`).
- **Error leakage**: don't return raw provider error bodies to the client via `gin.H{"error": err.Error()}` if they may contain provider-internal details. Map to a generic message and log the detail.

### Performance

- Creating `http.Client{}` per request is acceptable here but **`time.Second * 30`** is long for SMS; consider a tighter timeout and document it.
- N+1 outbound calls in loops (e.g. sending to a recipient list one-by-one without batching when the provider supports batches).
- `godotenv.Load` called on **every** request (as in `NewInforuModel`) — should be loaded once at startup in `main.go`, not per-request.
- Unnecessary allocations: building `map[string]interface{}` when a typed struct with `json` tags would do.
- Blocking outbound HTTP inside the Gin handler without a timeout-aware `ctx` (`http.NewRequestWithContext(ctx.Request.Context(), ...)`).

### Go & Project Conventions

- **Error handling**: wrap with context using `fmt.Errorf("inforu send sms: %w", err)`. Avoid `fmt.Printf("...: %s\n", err)` as the sole error reporting.
- **Layering**: `httpHandler.go` should parse/validate/dispatch only. Business logic and outbound HTTP belong in the provider package under `services/smsmessage/<provider>/`.
- **New provider**: must implement the full `ISms` interface (`SendSms`, `SendWhatsApp`). Add a matching `enums.SendingType` constant and a `case` in `httpHandler.go`'s switch. Add a `New<Provider>Model()` constructor that reads config from env.
- **New endpoint**: `routes/routes.go` → new `httpHandler.go` function → (if new domain) new `iService.go` interface → provider(s). Keep the `services/<domain>/...` layout.
- **Config**: read from env via `os.Getenv`; `.env.<env>` files loaded in `main.go`. Do not call `godotenv.Load` inside handlers or constructors.
- **Naming**: exported types end with `Model` for request/response payloads and provider configs (e.g. `SmsModel`, `InforuModel`, `FlashySmsModel`) — follow existing convention unless refactoring the whole package.
- **Tests**: table-driven, `*_test.go` alongside implementation, package `<pkg>_test`. Use `testify/assert`. Token is loaded via `GetToken()` from `.env` in `smsmessages_test.go` — mirror that pattern.
- **Gin idioms**: prefer `ctx.GetHeader("Token")` over `ctx.Request.Header["Token"][0]`; prefer `ctx.ShouldBindJSON` (already used); return via `ctx.JSON(status, gin.H{...})` and `return`.
- **Build tags / env flavors**: watch for changes in `main_qa.go` / `main_prod.go` / `main_stage.go` — each must stay consistent (build tag header + correct env value).
- **Commented-out code**: avoid committing large commented blocks (e.g. leftover `json.Marshal` alternatives in `flashy/service.go`). Flag them.

## Output Format

### If issues are found

For each issue:

- **Severity**: Critical / Suggestion
- **File**: `path/to/file.go:line`
- **Issue**: Clear explanation
- **Fix**: Exact Go code snippet to apply

### If no issues found

1. Confirm: "Code looks good."
2. Generate a conventional commit message:

```
<type>(<scope>): <short description>

[optional body]
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `build`, `ci`

**Scope examples** (match the package/domain touched):
`sms`, `inforu`, `flashy`, `whatsapp`, `routes`, `middleware`, `enums`, `models`, `main`, `deps`, `docker`, `ci`

**Examples:**

```
feat(whatsapp): add WhatsApp channel via InfoRU

Dispatch channel=whatsapp in httpHandler and implement
SendWhatsApp on InforuModel with TemplateId support.
```

```
fix(sms): use ctx.GetHeader to avoid panic on missing Token

Replaces ctx.Request.Header["Token"][0] which panicked when
the Token header was absent.
```

```
refactor(inforu): move godotenv.Load out of NewInforuModel

Load .env once in main.go; constructor now only reads env vars.
```

```
perf(flashy): reuse http.Client across requests

Avoid per-request client allocation and tighten timeout to 10s.
```

```
test(sms): add table-driven cases for WhatsApp template params
```

```
chore(deps): bump gin to v1.10.0
```

```
ci(github): add messages-api-flow workflow
```
