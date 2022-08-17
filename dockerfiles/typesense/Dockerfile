FROM typesense/typesense:0.23.1

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

EXPOSE 8108

ENTRYPOINT ["/opt/typesense-server"]
