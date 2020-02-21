FROM ubuntu:latest
WORKDIR /app

RUN apt-get update  && apt-get -y install cron && apt-get -y install wget
                   && apt-get install -y git &&  apt-get install -y build-essential

RUN wget https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz  &&
    tar -C /usr/local -xzf go1.13.8.linux-amd64.tar.gz
RUN rm -rf go1.13.8.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin

RUN mkdir -p /root/go/src/github.com/mrturkmen06

RUN cd /root/go/src/github.com/mrturkmen06 && git clone https://github.com/mrturkmen06/scheduler.git
COPY cron /etc/cron.d/cron
RUN chmod 0644 /etc/cron.d/cron
RUN crontab /etc/cron.d/cron
RUN touch /var/log/cron.log
RUN service cron start
SHELL ["/bin/bash", "-c"]
CMD ["cron","-f"]
