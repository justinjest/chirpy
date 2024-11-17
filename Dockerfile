FROM debian:stable-slim

# COPY source destination
COPY out /bin/goserver

CMD ["/bin/goserver"]