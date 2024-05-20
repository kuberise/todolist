FROM golang:1.22-alpine as builder 

WORKDIR /opt 

# https://megamorf.gitlab.io/2019/09/08/alpine-go-builds-with-cgo-enabled/
RUN apk add build-base

COPY . .  

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -mod=vendor -o app ./main.go 


FROM alpine as runner 

COPY --from=builder /opt/app /opt/

CMD ["/opt/app"]