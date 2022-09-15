FROM golang:1.18-alpine3.16

LABEL "com.github.actions.name"="squarescale-cli"
LABEL "com.github.actions.description"="Squarescale CLI"
LABEL "com.github.actions.icon"="refresh-cw"
LABEL "com.github.actions.color"="green"

LABEL name="squarescale-cli"
LABEL version="1.0.9"
LABEL maintainer="SquareScale Engineering <engineering@squarescale.com>"
LABEL repository="https://github.com/squarescale/squarescale-cli"

COPY ./dist/sqsc-alpine-amd64-gh-actions /sqsc
COPY ./gh-actions/gh-actions /gh-actions

ENTRYPOINT ["/gh-actions"]
CMD ["/sqsc version"]
