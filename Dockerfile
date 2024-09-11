# Dockerfile
FROM envoyproxy/envoy:v1.18.3

COPY ./envoy.yaml /etc/envoy/envoy.yaml
COPY ./image.binpb /tmp/envoy/proto.binpb

EXPOSE 5001

CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml