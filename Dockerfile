FROM scratch
MAINTAINER Stephane Jourdan <stephane.jourdan@squarescale.com>
COPY bin/linux_amd64/sqsc /
ENTRYPOINT ["/sqsc"]
