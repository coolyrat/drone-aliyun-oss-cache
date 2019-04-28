FROM luischan/alpine:3.9
COPY ./oa-suite /oa-suite
ENTRYPOINT [ "/oa-suite" ]