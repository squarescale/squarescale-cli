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

There is a Makefile rule for that ! But you need to create a Github OAuth2 token before you can publish anything:
- Go to your profile settings on Github
- Create a new personal access token (enter appropriate description)
- Check the `repository` box level to give your token the permissions to create new releases
- remember this token in a super safe place !

Publishing a new version of the CLI on Github is super simple:

```bash
$> GITHUB_USER_TOKEN=[USER_NAME]:[OAUTH2_TOKEN] make publish
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
