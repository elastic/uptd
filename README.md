# Up To Date

`uptd` (shorthand for _up to date_) is a Go package which provides primitives
and mechanisms to check if a version is up to date according to the latest one
published in a remote artifact repository.

## Installation

To download the source, run: `go get github.com/elastic/uptd`.

## Usage

```go

import (
    "fmt"

    "github.com/elastic/uptd"
)

func main() {
    var githubToken = "your personal github token"

    provider, err := uptd.NewGithubProvider("elastic", "go-licenser", githubToken)
    if err != nil {
        panic(err)
    }

    uptodate, err := uptd.New(provider, version)
    if err != nil {
        panic(err)
    }

    res, err := uptodate.Check()
    if err != nil {
        panic(err)
    }

    if res.NeedsUpdate {
        fmt.Printf(
            "new version %s available, release URL is %s",
            res.Latest.Version.String(), res.Latest.URL,
        )
    }
}
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).
