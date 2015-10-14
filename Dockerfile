FROM google/golang

WORKDIR /gopath/src/github.com/maleck13/locals-api
ADD . /gopath/src/github.com/maleck13/locals-api
RUN go get github.com/maleck13/locals-api
RUN ls -al /gopath/bin
EXPOSE 9005
CMD []
ENTRYPOINT ["/gopath/bin/locals-api"]