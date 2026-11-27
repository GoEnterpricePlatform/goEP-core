# UI interaction guidelines

This project uses Templ for server-rendered HTML, Datastar for reactive page state and server requests, and TypeScript for richer browser-side interactions. Choose the layer based on where the state and work belong.

## Use TypeScript for client-owned forms

Use TypeScript when a form needs immediate validation, dependent fields, or interactions that should update without a server round trip. Keep the form's client-side state and rendering together, then submit the result to a Go handler when the user saves. The Go handler must still validate the submitted data before calling a service.

For simple required fields or basic show/hide behavior, Datastar expressions may be enough; do not add TypeScript just because a form has several fields.

## Use Go UI handlers for complex or data-dependent updates

Use a Go UI handler when an update depends on services, database data, domain rules, or HTML rendered by Templ. A handler can read Datastar signals, perform the calculation or operation, then patch signals and HTML with SSE.

For example, Paddle plan variation selections call `GeneratePaddlePlanItems`. The handler uses the selected option IDs to calculate combinations, preserves item data for combinations that remain selected, and asks Templ to render the Items table.

If an editor has dynamic rows, keep its state mutations and rendering in one place. When Templ renders those rows, let the Go handler update them. If the editor later becomes fully client-owned and needs immediate interactions, move its state and rendering together to TypeScript instead of splitting them between TypeScript and server handlers.

## Use Datastar for small UI state and requests

Use Datastar signals and attributes for UI state such as opening or closing a popup, showing loading or error messages, and triggering requests to Go handlers.

For example, the Variations page opens its form modal by setting `$modal.domain` and `$modal.action`. Its Create/Update button then calls the matching Go route with `@post` or `@put`. Datastar controls the popup's visibility; the Go handler validates and persists the variation.

## Rule of thumb

- **TypeScript:** client-owned form behavior and real-time validation.
- **Go UI handlers + Templ:** complex updates that need backend data or server-rendered dynamic HTML.
- **Datastar:** lightweight reactive state, popup visibility, loading/error feedback, and requests to Go.

Keep backend validation and persistence in Go regardless of which layer provides the client-side interaction.
