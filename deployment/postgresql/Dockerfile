FROM postgres:latest

ENV POSTGRES_DB=DMQS \
    POSTGRES_USER=ykarabul \
    POSTGRES_PASSWORD=yigitsh \
    POSTGRES_HOST=localhost \
    POSTGRES_PORT=5432

COPY ../../init-db.sh /tmp/init-db.sh

RUN chmod +x /tmp/init-db.sh

RUN ./tmp/init-db.sh

EXPOSE 5432
