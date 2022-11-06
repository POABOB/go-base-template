FROM alpine
ADD base.output /base.output
COPY filebeat.yml /filebeat.yml
ENTRYPOINT [ "/base.output" ]
