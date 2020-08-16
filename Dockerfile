FROM golang

COPY . .

ENV GOPATH=$PWD

RUN chmod +x start.sh

RUN ./start.sh

ENTRYPOINT /bin/http-server

EXPOSE 8080
