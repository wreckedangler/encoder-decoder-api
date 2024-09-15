FROM ubuntu:latest
LABEL authors="raphi"

ENTRYPOINT ["top", "-b"]