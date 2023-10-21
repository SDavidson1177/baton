FROM golang:1.19.1

RUN mkdir /baton
RUN apt update
RUN apt install curl -y
RUN curl https://get.ignite.com/cli@v0.27.1! | bash
RUN mkdir ./baton
COPY . ./baton
WORKDIR ./baton

EXPOSE 26656 26657 1317 9090 8545 8546

CMD ["tail", "-f", "/dev/null"]
