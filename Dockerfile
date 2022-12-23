FROM golang:alpine as build

WORKDIR /build

RUN apk add --no-cache git \
    && git clone https://github.com/layou233/ZBProxy.git \
    && cd ZBProxy \
    && CGO_ENABLED=0 GOAMD64=v3 go build -v -ldflags="-s -w" -o ../ZBProxy-linux \
    && echo '{ \
    "Services": [ \
        {\
            "Name": "HypixelDefault",\
            "TargetAddress": "mc.hypixel.net",\
            "TargetPort": 25565,\
            "Listen": 25565,\
            "Flow": "auto",\
            "IPAccess": {\
                "Mode": ""\
            },\
            "Minecraft": {\
                "EnableHostnameRewrite": true,\
                "OnlineCount": {\
                    "Max": 114514,\
                    "Online": -1,\
                    "EnableMaxLimit": false\
                },\
                "NameAccess": {\
                    "Mode": ""\
                },\
                "AnyDestSettings": {},\
                "MotdFavicon": "{DEFAULT_MOTD}",\
                "MotdDescription": "§d{NAME}§e service is working on §a§o{INFO}§r§c§lProxy for §6§n{HOST}:{PORT}§r" \
            },\
            "TLSSniffing": {\
                "RejectNonTLS": false\
            },\
            "Outbound": {\
                "Type": ""\
            }\
        }\
    ],\
    "Lists": {}\
}'\
> ../ZBProxy.json

FROM gcr.io/distroless/static-debian11:latest

COPY --from=build /build/ZBProxy-linux /
COPY --from=build /build/ZBProxy.json /

CMD [ "/ZBProxy-linux" ]
