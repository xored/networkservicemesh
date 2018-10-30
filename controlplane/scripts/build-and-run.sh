#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
CONTAINERTXT=/tmp/container.txt

docker build -t nsmd/nsmd -f "${DIR}/../build/Dockerfile.nsmd" "${DIR}/../../"
docker build -t nsmd/nse -f "${DIR}/../build/Dockerfile.nse" "${DIR}/../../"
docker build -t nsmd/nsc -f "${DIR}/../build/Dockerfile.nsc" "${DIR}/../../"

echo "Starting nsmd..."
docker run -d -v "/var/lib/networkservicemesh:/var/lib/networkservicemesh" nsmd/nsmd > "${CONTAINERTXT}"
echo "Starting nse..."
docker run -d -v "/var/lib/networkservicemesh:/var/lib/networkservicemesh" nsmd/nse >> "${CONTAINERTXT}"
echo "Starting vpp..."
docker run --network=host --privileged=true --volume=/var/run:/var/run --volume=/var/lib:/var/lib --volume=/lib/modules:/lib/modules --ipc=host --pid=host -d networkservicemesh/vpp >> "${CONTAINERTXT}"
echo "Starting vpp-daemon..."
docker run --network=host --privileged=true --volume=/var/run:/var/run --volume=/var/lib:/var/lib --volume=/lib/modules:/lib/modules --volume=/var/lib/networkservicemesh:/var/lib/networkservicemesh/ --ipc=host --pid=host -d networkservicemesh/vpp-daemon >> "${CONTAINERTXT}"

echo "vpp-daemon takes time (unnecessarily) to register with nsmd, so let it connect. waiting 60 seconds..."
sleep 60

echo "Running nsm client..."
docker run -d -v "/var/lib/networkservicemesh:/var/lib/networkservicemesh" nsmd/nsc >> "${CONTAINERTXT}"

echo "waiting 20 seconds..."
sleep 20

echo "Showing nsmd logs..."
docker logs "$(sed '1q;d' ${CONTAINERTXT})"

echo "Showing nse logs..."
docker logs "$(sed '2q;d' ${CONTAINERTXT})"

echo "Showing vpp-daemon logs..."
docker logs "$(sed '4q;d' ${CONTAINERTXT})"

echo "Showing nsc logs..."
docker logs "$(sed '5q;d' ${CONTAINERTXT})"

echo "Showing nse interfaces..."
docker exec "$(sed '2q;d' ${CONTAINERTXT})" ifconfig -a

echo "Showing nsc interfaces..."
docker exec "$(sed '5q;d' ${CONTAINERTXT})" ifconfig -a

echo "Ping nse from nsc interfaces..."
docker exec "$(sed '5q;d' ${CONTAINERTXT})" ping -c 5 2.2.2.3

echo "Kill containers..."
xargs docker kill < ${CONTAINERTXT}
