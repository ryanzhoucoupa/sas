FROM alpine:3.2
RUN apk add --update ca-certificates

ENV SAS_ENVIRONMENT production
ENV SAS_DB_ADAPTER mysql
ENV SAS_DB_CONFIG root:my-secret-pw@tcp(0.0.0.0:3306)/sftp

ADD sas /bin/sas
ADD conf.json /etc/conf.d/sas.json
ENTRYPOINT ["/bin/sas"]
EXPOSE 8080
CMD ["-configuration-path=/etc/conf.d/sas.json"]
