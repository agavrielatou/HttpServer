# Configure IPv6 support inside container
DOCKER_SYSCTL="--sysctl net.ipv6.bindv6only=1 --sysctl net.ipv6.conf.all.disable_ipv6=0 --sysctl net.ipv6.conf.default.disable_ipv6=0"
docker run $DOCKER_SYSCTL --net="host" -it --publish 8080:8080 hellocloudreachmain:1.0
