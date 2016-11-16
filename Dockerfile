FROM golang:1.7
MAINTAINER Logmatic.io Support team, @gpolaert


# Install tools required
RUN go get github.com/Masterminds/glide

# Add sources
COPY .  $GOPATH/src/github.com/logmatic/beats-forwarder
WORKDIR $GOPATH/src/github.com/logmatic/beats-forwarder


# Install deps and build
RUN glide install && go build


# Ports
EXPOSE 5044
ENTRYPOINT ["./beats-forwarder"]