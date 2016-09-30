FROM scratch
MAINTAINER Stephane Jourdan <stephane.jourdan@squarescale.com>
COPY ./sqsc-docker /sqsc
ENTRYPOINT ["/sqsc"]
