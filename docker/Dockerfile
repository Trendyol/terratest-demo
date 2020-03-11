FROM golang:1.14.0-alpine
ARG K8S_VERSION=v1.14.3
RUN set -x                  && \
    apk --update upgrade    && \
    apk add ca-certificates build-base && \
    rm -rf /var/cache/apk/* && \
    wget https://storage.googleapis.com/kubernetes-release/release/$K8S_VERSION/bin/linux/amd64/kubectl && \
    chmod +x kubectl &&  \
    mv kubectl /usr/local/bin/kubectl

ENTRYPOINT ["/kubectl"]