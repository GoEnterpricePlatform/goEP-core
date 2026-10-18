# Builder registration and manifest validation

## 1.NewCompiledRegistry

Registers the available constructors: **Identity**, **Posts**, **Billing**, **Catalog** and **AI**.

In this record, "available" means that your code is included in the executable.

## 2. Initialize

Looks up each constructor by name and calls it with the typed configuration that module needs.

For example, you construct **Identity** and then pass your authentication middleware to **Posts**.

## 3. Registry.Validate

Compare the registry with the manifests:

- Detects if there is a manifest without compiled code.
- Detects if there is a factory without a manifest.

That validation is tested, but still **is not executed when the application starts**.

## Tour

```text
NewCompiledRegistry
        ↓
Initialize searches each factory by name
        ↓
Call the constructors in the current order
        ↓
Modules record their routes as before
```