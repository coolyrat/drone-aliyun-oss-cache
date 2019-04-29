FROM luischan/alpine:3.9
COPY ./drone-oss-cache /drone-oss-cache
ENTRYPOINT [ "/drone-oss-cache" ]