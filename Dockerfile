FROM scratch

LABEL name="squarescale-cli"
LABEL version="0.1.0"
LABEL maintainer="SquareScale Engineering <engineering@squarescale.com>"

COPY ./sqsc-docker /sqsc

ENTRYPOINT ["/sqsc"]
CMD ["--help"]
