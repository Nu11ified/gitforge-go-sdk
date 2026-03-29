# gitforge-go

Go SDK for [GitForge](https://git-forge.dev) — Git infrastructure for developers who build on Git.

## Install

```bash
go get github.com/gitforge/sdk-go
```

Requires Go 1.21+.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    gitforge "github.com/gitforge/sdk-go"
)

func main() {
    client := gitforge.New(
        "https://api.git-forge.dev",
        "gf_your_token_here",
    )
    ctx := context.Background()

    // Create a repo
    repo, err := client.Repos.Create(ctx, gitforge.CreateRepoOpts{
        Name:       "my-repo",
        Visibility: "private",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("Created:", repo.Name)

    // List repos
    repos, err := client.Repos.List(ctx, gitforge.ListOpts{Limit: 10})
    if err != nil {
        panic(err)
    }
    for _, r := range repos.Data {
        fmt.Println(r.Name)
    }

    // Create a commit
    err = client.Commits.Create(ctx, repo.ID, gitforge.CreateCommitOpts{
        Branch:      "main",
        Message:     "initial commit",
        AuthorName:  "Your Name",
        AuthorEmail: "you@example.com",
        Files: []gitforge.CommitFile{
            {Path: "README.md", Content: "# My Project"},
            {Path: "main.go", Content: "package main\n\nfunc main() {}"},
        },
    })
    if err != nil {
        panic(err)
    }
}
```

## Resources

| Resource | Methods |
|----------|---------|
| `Repos` | `Create`, `List`, `Get`, `Update`, `Delete` |
| `Branches` | `List`, `Create`, `Delete`, `Promote` |
| `Tags` | `List`, `Create`, `Delete` |
| `Commits` | `Create`, `List`, `Get` |
| `Files` | `GetBlob`, `GetTree` |
| `Search` | `Code` |
| `Tokens` | `Create`, `List`, `Revoke` |
| `Mirrors` | `List`, `Create`, `Sync`, `Delete` |
| `Webhooks` | `Create`, `List`, `Update`, `Delete`, `Test` |

## Webhook Validation

```go
import gitforge "github.com/gitforge/sdk-go"

isValid := gitforge.ValidateWebhook(gitforge.ValidateOpts{
    Payload:   rawBody,
    Signature: r.Header.Get("X-GitForge-Signature"),
    Secret:    "your_webhook_secret",
    Timestamp: r.Header.Get("X-GitForge-Timestamp"),
    Tolerance: 300,
})
```

## Error Handling

```go
repo, err := client.Repos.Get(ctx, "nonexistent")
if err != nil {
    var gfErr *gitforge.Error
    if errors.As(err, &gfErr) {
        fmt.Println(gfErr.Status)  // 404
        fmt.Println(gfErr.Code)    // "NOT_FOUND"
    }
}
```

## Contributing

This SDK is developed inside the [GitForge monorepo](https://github.com/Nu11ified/GitForge) at `sdks/go/` and published to this repo via git subtree.

To contribute:

1. Clone the monorepo: `git clone https://github.com/Nu11ified/GitForge.git`
2. Make changes in `sdks/go/`
3. Run tests: `cd sdks/go && go test ./...`
4. Submit a PR to the monorepo

## License

MIT
