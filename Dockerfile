# this is a multistage docker image! this helps the final image be much smaller

# build stage
FROM docker.io/golang:1.21-alpine AS build-env

WORKDIR /go/src/github.com/DanielOaks/codingtest-birdhouse-mock
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build birdhouse-mock.go

CMD ["./birdhouse-mock"]

# run stage
FROM docker.io/alpine:3.13

LABEL maintainer="Daniel Oaks <daniel@danieloaks.net>" \
      description="Mock server for the birdhouse-admin frontend"

EXPOSE 5031/tcp

COPY --from=build-env /go/src/github.com/DanielOaks/codingtest-birdhouse-mock/birdhouse-mock \
                      /bh-bin/

ENTRYPOINT ["/bh-bin/birdhouse-mock"]

# # uncomment to debug
# RUN apk add --no-cache bash
# RUN apk add --no-cache vim
# CMD /bin/bash
