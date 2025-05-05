FROM golang:1.23-bullseye AS build

ARG GITHUB_TOKEN
RUN if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN is not set" && exit 1; fi
RUN git config --global url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork"

ARG HTTP_PROXY
ARG HTTPS_PROXY

## Setup env
WORKDIR /app/squaring-e2e


COPY go.mod go.sum ./
RUN --mount=type=cache,target="/go/pkg/mod" \
    go mod download

COPY . .
########## Build dvse2e ##########
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    make build-dvse2e

########## Setup runtime env ##########
FROM golang:1.23-bullseye AS runtime
RUN apt-get update -yqq && apt-get install -yqq openssh-server curl jq less
RUN mkdir -p /run/sshd && chmod 0755 /run/sshd

RUN go install github.com/mikefarah/yq/v4@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY --from=build /app/squaring-e2e/bin/dvse2e /usr/bin/dvse2e

COPY ./docker/ssh /root/.ssh
RUN chmod 600 /root/.ssh/*

WORKDIR /root

# expose port
EXPOSE 22

CMD ["bash"]
