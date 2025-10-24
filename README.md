# go-modules

This repository contains small example Go modules and notes on how to create and publish Go modules, both at the repository root and inside subdirectories.

The document below shows the common workflows, commands, and examples for:

- Creating a module at the repository root
- Creating a module inside a subdirectory (a nested module)
- Tagging and publishing module versions
- Importing and updating modules with `go` commands

---

## Module at repository root

Typical steps when you want the repository itself to be a single module:

1. Create a repository on GitHub (for example `github.com/username/app`).
2. Initialize the module in the project root using the repository path as the module path:

```bash
go mod init github.com/username/app
cat go.mod
# module github.com/username/app
```

Project layout example:

```bash
.
├── go.mod
└── main.go
```

3. Commit and push to GitHub. Tag releases with semantic version tags (vMAJOR.MINOR.PATCH).

```bash
git add -A
git commit -m "initial commit"
git tag -a v1.0.0 -m "Version 1.0.0"
git push origin master
git push origin v1.0.0
```

4. Follow Semantic Versioning for future releases.

---

## Module in a subdirectory

You can also create a separate module inside a subdirectory. This is useful when a single repository contains multiple independent modules.

Example repository layout:

```bash
app2/
└── sub/
    └── main.go
```

From the `sub` directory initialize a module using the full import path including the subdirectory:

```bash
cd app2/sub
go mod init github.com/username/app2/sub
cat go.mod
# module github.com/username/app2/sub

# go 1.22.1
```

When tagging a release for a nested module, include a tag name that reflects the submodule if you prefer (for example `sub/v1.0.0`), but you can also tag at repository level depending on your release strategy.

```bash
git add -A
git commit -m "release submodule"
git tag -a sub/v1.0.0 -m "sub module v1.0.0"
git push origin master
git push origin sub/v1.0.0
```

---

## Importing modules

When importing, use the full module path as declared in the `go.mod` file. Examples:

```go
import (
    "github.com/username/app"
    "github.com/username/app2/sub"
)
```

Run `go mod tidy` to add and prune dependencies for your module. If a dependency is a repo on GitHub, Go will by default resolve it to a version based on tags; otherwise it may use the latest available commit.

Updating a dependency to a specific version can be done with `go get`:

```bash
go get github.com/username/app@v1.0.1
go get github.com/username/app2/sub@v1.0.1
```

Use `go get -u` to update to the latest minor/patch version allowed by your go.mod constraints.

---

## Notes

- Prefer semantic versioning for published modules.
- When designing multi-module repositories, decide on a tagging/release strategy that fits your workflow (per-submodule tags vs repository-level tags).
- This repository contains small examples; adapt the commands above for your project and CI workflow.

---
