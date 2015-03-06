FROM google/debian:wheezy

COPY gearlogger /usr/bin/gearlogger

CMD ["gearlogger"]
