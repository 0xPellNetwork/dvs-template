FROM golang:1.23-bullseye AS build

ARG HTTP_PROXY
ARG HTTPS_PROXY

## Setup env
WORKDIR /app

ARG GITHUB_TOKEN
ARG PELLDVS_VERSION
RUN if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN is not set" && exit 1; fi
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork"
RUN git clone https://github.com/0xPellNetwork/pelldvs.git -b $PELLDVS_VERSION --depth 1

# Build pelldvs
WORKDIR /app/pelldvs
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    make build

########## Build squaringd ##########
## Download dependencies before copying the app to leverage docker cache
WORKDIR /app/squaringd
COPY go.mod go.sum ./
RUN --mount=type=cache,target="/go/pkg/mod" go mod download

## Build the app
COPY ./ /app/squaringd/
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    make build

########## Setup runtime env ##########
FROM golang:1.23-bullseye AS runtime
RUN apt-get update -yqq && apt-get install -yqq openssh-server curl jq less
RUN mkdir -p /run/sshd && chmod 0755 /run/sshd

RUN curl -L https://foundry.paradigm.xyz | bash
ENV PATH="$PATH:/root/.foundry/bin"
RUN foundryup

RUN go install github.com/mikefarah/yq/v4@latest
COPY --from=build /app/squaringd/bin/squaringd /usr/bin/squaringd
COPY --from=build /app/pelldvs/build/pelldvs /usr/bin/pelldvs
COPY ./docker/operator/scripts /root/scripts
COPY ./docker/ssh /root/.ssh

RUN chmod 600 /root/.ssh/*

WORKDIR /root

# expose port
EXPOSE 22