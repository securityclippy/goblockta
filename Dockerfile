FROM ubuntu:latest

RUN apt update
RUN apt-get install -y ca-certificates

ADD goblockta /


CMD ["./goblockta"]