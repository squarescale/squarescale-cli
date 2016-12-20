FROM scratch
LABEL name="squarescale-cli"
LABEL version=0.1.0
MAINTAINER SquareScale Engineering <engineering@squarescale.com>
COPY ./sqsc-docker /sqsc
ENTRYPOINT ["/sqsc"]
CMD ["--help"]
