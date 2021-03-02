FROM golang:1.13.4 as builder

ARG APP_PATH="rest"
ARG COPY_LIST=""
RUN mkdir /app /build
COPY . /build/
RUN if [ "${COPY_LIST}" != "" ];then \
for item in ${COPY_LIST};do cp /build/${item} /app;done \
fi
RUN cd /build/$APP_PATH && CGO_ENABLED=0 go build -o /app/$APP_PATH

FROM alpine
COPY --from=builder /app /bin/
EXPOSE 8080
ENTRYPOINT [ "rest" ] 
