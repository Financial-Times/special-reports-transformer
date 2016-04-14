FROM alpine:3.3

ADD . /special-reports-transformer/

RUN apk add --update bash \
  && apk --update add git bzr \
  && apk --update add go \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/special-reports-transformer" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && cp -r special-reports-transformer/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && go build \
  && mv special-reports-transformer /app \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*
CMD [ "/app" ]