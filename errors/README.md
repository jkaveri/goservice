# errors

Composable error helpers for services: wrap causes for `errors.Is` / `errors.As`, attach client-facing text, optional **stable codes** (domain business identifiers and transport mapping), metadata, and stack capture for debugging. The API follows patterns familiar from `github.com/pkg/errors`, extended with message, code, metadata, and stack types.

Import as a named package to avoid clashing with the standard library:

```go
import errors "github.com/jkaveri/goservice/errors"
```

Use `errors.Is`, `errors.As`, and `errors.Unwrap` from this package when you want the same behavior as the standard library `errors` package with documentation aligned to this module.

---

## Leaf errors: `New`, `Errorf`

**What:** A simple error value with a message string. `Unwrap()` is nil.

**When:** Creating a root failure inside your package (validation failed, invariant broken, downstream returned nothing to wrap).

```go
return errors.New("tenant not found")
return errors.Errorf("limit must be positive, got %d", n)
```

Prefer these over `errors.New` / `fmt.Errorf` from stdlib when you want errors typed by this package and consistent formatting with the rest of the stack.

---

## Context wrapping: `Wrap`, `Wrapf`

**What:** Adds a short operator-facing message and sets `Unwrap()` to the inner error so `errors.Is` / `errors.As` still see the cause. `Error()` becomes `message: cause.Error()`. When `err` is non-nil and `HasStack(err)` is false, **`Wrap` captures a stack trace at this call site** before adding the message (`Wrapf` same).

**When:** Returning an error up the call stack with an operator breadcrumb (“what we were doing”). Not the same as `WithMessage`: `Wrap` does **not** implement `MessageError`. Stacks are attached on the **first** `Wrap` when the chain has none—no separate stack API.

```go
if err != nil {
    return errors.Wrap(err, "load tenant")
}
```

**Nil behavior:** `Wrap(nil, "msg")` is still a non-nil error with `Unwrap() == nil` and `Error() == "msg"`. Check `err == nil` before wrapping if you need stdlib-style nil propagation.

---

## Client-facing text: `WithMessage`, `WithMessagef`, `Message`

**What:** `WithMessage` adds a **short segment** meant for HTTP/gRPC bodies or UI. `Error()` is `segment: underlying.Error()`; `Message()` on the wrapper returns only that segment. `Message(err)` walks the unwrap chain, collects every `MessageError`, joins outer→inner with `": "`, and **skips** layers that are not `MessageError` (including `Wrap` / `New` / stdlib errors).

**When:** Boundaries that must return a stable, user-safe string while logs keep the full chain. Typical pattern: inner detail with `New` / `Wrap`, outer layers with `WithMessage` for what the user should read.

```go
err := errors.WithMessage(
    errors.Wrap(dbErr, "query tenants"),
    "could not load workspace",
)
// logs / Error(): full chain
// API body: errors.Message(err) → "could not load workspace"
```

Use `WithMessagef` when the segment is formatted from values.

---

## Stable codes (domain + boundary): `WithCode`, `Code`, `ContainsCode`

**What:** Wraps an error with a **string code**. `Error()` becomes `[code] underlying`. The value implements `CodeError`. `Code(err)` returns the first code in the chain (via `As`). `ContainsCode` walks the full chain (including `WalkErrorChain`) to see if any layer carries a matching code.

**When — domain (replace sentinels):** Each bounded context or package can define **its own code constants** (plain strings) that name a **business problem**, not only HTTP/gRPC outcomes. Callers identify failures with `errors.ContainsCode(err, billing.InvoiceAlreadyPosted)` or `errors.Code(err)` instead of comparing `err == ErrInvoiceAlreadyPosted` or brittle `errors.Is` against package-global sentinel variables. Keep codes **stable** across releases so upstream services and clients can branch on them.

**When — boundary:** The same codes can be copied into API responses, logs, or metrics, and mapped to HTTP status or gRPC status in an interceptor. For a **shared** catalog (cross-service “not found”, “rate limited”, etc.), **`errorcode`** in this module provides constructors and helpers; you can still attach a domain-specific wrapper or use this package’s `WithCode` for package-local codes.

When `HasStack(err)` is false, **`WithCode` also captures a stack** at this call site (same policy as `Wrap`).

```go
// In package billing (example)
const CodeInvoiceAlreadyPosted = "BILLING_INVOICE_ALREADY_POSTED"

func PostInvoice(...) error {
    if posted {
        return errors.WithCode(errors.New("invoice locked"), CodeInvoiceAlreadyPosted)
    }
    // ...
}

// Caller
if errors.ContainsCode(err, billing.CodeInvoiceAlreadyPosted) {
    // handle business case
}
```

Use **constants** for codes (avoid scattering magic strings). Prefer `ContainsCode` when the code might be wrapped under other layers; use `Code(err)` when you only care about the first tagged layer.

---

## Debugging context: `WithMetadata`, `Metadata`

**What:** Attaches a shallow `map[string]any` for operators and tools. **`Error()` is unchanged** (still the inner error’s text). `Metadata(err)` returns the map from the first `MetadataError` in the chain. Verbose formatting (`%+v`) prints the wrapped error’s `%+v` output, then a **`metadata:`** block with **sorted keys** and tab-separated `key: value` lines.

**When:** You want correlation fields (request ID, tenant ID, feature flag) in logs or structured handlers without putting them in the user-visible `Error()` string. Good for support dashboards and gateways that copy metadata into details.

```go
err = errors.WithMetadata(err, map[string]any{
    "request_id": reqID,
    "tenant_id":  tenantID,
})
```

---

## Stack capture: `StackError`, `HasStack`

**What:** Stack frames are stored on an inner `withStack` layer. `Wrap`, `WithMessage`, and `WithCode` call **`ensureStack`** when the chain has no stack yet; otherwise the error is unchanged. Wrappers cache **`HasStack`** as `stacked != nil` (no chain walk). `fmt.Sprintf("%+v", err)` prints the logical chain, then stack frames after the innermost `withStack`.

**When:** Use normal `Wrap` / `WithMessage` / `WithCode`—do not add stacks manually. `New` / `Errorf` leaves have no stack until the first wrap. `WithMetadata` does **not** attach a stack.

```go
return errors.Wrap(dbErr, "insert tenant") // stack + message in one step
return errors.Wrap(err, "outer")           // HasStack true → no second capture
```

Use `errors.As(err, &stackErr)` with `StackError` for programmatic `StackTrace()`. Use `errors.HasStack(err)` to branch without walking the chain.

---

## Utilities

| API | Role |
|-----|------|
| `HasStack` | O(1) on package wrappers: whether a stack was already captured below (see `stackMarker`). |
| `WalkErrorChain` | Depth-first walk over `Unwrap()` (single or `[]error`), stop when the visitor returns true. |
| `Join` | Same as standard library `errors.Join`—combine multiple independent errors. |
| `Is` / `As` | Same semantics as the standard library, re-exported for one import path. |

---

## Formatting cheat sheet

| Construct | `%s` / `%v` | `%q` | `%+v` |
|-----------|-------------|------|--------|
| `New` / `Errorf` | message | quoted **message only** | message |
| `Wrap` / `Wrapf` | `msg: cause` | quoted **msg only** | message line, indented cause chain, then stack (first wrap only) |
| `WithMessage` | full `Error()` | quoted full `Error()` | outer message line, `%+v` of unwrap, stack if first attach |
| `WithCode` | `[code] inner` | quoted full `Error()` | `[code]` line, tab-indented unwrap, stack if first attach |
| `WithMetadata` | inner `Error()` | quoted full `Error()` | inner `%+v`, then `metadata:` block (sorted keys) |

Prefer **`%+v` in logs** when you want causes, metadata, and stacks where implemented.

---

## Composition tips

1. **Order:** Build inside-out: leaf (`New` / stdlib / `errorcode`) → `Wrap` for internal steps (stack on first wrap) → `WithMessage` at boundaries for client text → `WithCode` for stable business codes → `WithMetadata` for correlation → map to HTTP/gRPC at the transport layer.
2. **Do not duplicate:** Avoid both `Wrap` and `WithMessage` with the same English sentence; use `Wrap` for operator breadcrumbs and `WithMessage` for client-safe wording.
3. **Typed inspection:** Prefer `errors.Is` / `errors.As` for sentinel or type targets; use `errors.ContainsCode` / `errors.Code` for **string business codes**; use `errorcode.Is*` when the chain uses **`errorcode`** values. Avoid parsing `Error()` strings.

---

## Mocks

Generated mocks for `CodeError`, `MetadataError`, `StackError`, and related shapes live under `errors/mock/` for tests in other packages.
