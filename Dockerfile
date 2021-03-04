FROM debian:10-slim

LABEL "com.github.actions.name"="squarescale-cli"
LABEL "com.github.actions.description"="Squarescale CLI"
LABEL "com.github.actions.icon"="refresh-cw"
LABEL "com.github.actions.color"="green"

LABEL name="squarescale-cli"
LABEL version="0.1.0"
LABEL maintainer="SquareScale Engineering <engineering@squarescale.com>"
LABEL repository="https://github.com/squarescale/squarescale-cli"

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

COPY ./dist/sqsc-alpine-amd64 /sqsc
COPY ./gh-actions.sh /gh-actions.sh

ENTRYPOINT ["/gh-actions.sh"]
CMD ["--help"]
