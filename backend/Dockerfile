FROM postgres:12.3

RUN  set -ex \
    \
    && apt update \
    && apt install -y\
        build-essential \
        wget \
        postgresql-server-dev-12 \
    \
    && wget -O is_jsonb_valid.tar.gz "https://github.com/furstenheim/is_jsonb_valid/archive/master.tar.gz" \
    \
    && mkdir -p /usr/src/is_jsonb_valid \
    && tar \
        --extract \
        --file is_jsonb_valid.tar.gz \
        --directory /usr/src/is_jsonb_valid \
        --strip-components 1 \
    && rm is_jsonb_valid.tar.gz \
    && cd /usr/src/is_jsonb_valid \
    && make install