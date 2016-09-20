FROM scratch
MAINTAINER Stephane Jourdan <stephane.jourdan@squarescale.com>
COPY bin/sqsc /
ENTRYPOINT ["/sqsc"]
