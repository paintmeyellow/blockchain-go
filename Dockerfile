FROM ubuntu:18.04

RUN apt update && apt upgrade -y
RUN apt install software-properties-common -y
RUN add-apt-repository ppa:openjdk-r/ppa && apt update
RUN apt install openjdk-11-jdk wget systemd -y

RUN mkdir $HOME/minima
RUN wget -qO $HOME/minima/minima.jar https://github.com/minima-global/Minima/raw/master/jar/minima.jar

EXPOSE 9002 9003 9004 9005

VOLUME $HOME/.minima

COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]
