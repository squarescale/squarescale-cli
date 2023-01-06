# squarescale-cli/sqsc

Squarescale CLI

## Description

## Usage

## Install

To install, use `go get`:

```bash
$> go get -d github.com/squarescale/squarescale-cli
```
### Troubleshooting

When installing squarescale-cli, if you have an error `Unable to read Username for 'https://github.com'` but you are all set to connect to Github using ssh, just run the following command and try again:
`git config --global url.ssh://git@github.com/.insteadOf https://github.com/`

## Publish a release on Github

Github action is triggered on tag creation named with pattern 'v*'.

To release, just tag the version and push this tag:


```bash
$> git tag -m v1.x.x 1.x.x
$> git push --tags
```

Under the hood, *goreleaser* is used to build binaries and distribution packages.
Refer to the [documentation](https://goreleaser.com/intro/) to further details.

To build locally into *dist/* directory, run the following command:

```bash
$> goreleaser release --snapshot --rm-dist
```

## Contribution

1. Fork ([https://github.com/squarescale/squarescale-cli/fork](https://github.com/squarescale/squarescale-cli/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create a new Pull Request

## Author

[squarescale](https://github.com/squarescale)
