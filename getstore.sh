#!/bin/bash
set -x
ordererfile="orderer-storage.log"
if [ -f $ordererfile ]; then
    rm $ordererfile
fi
for i in $(seq 1 3)
do
    echo "orderer$i.example.com" >> $ordererfile
    docker exec -ti orderer$i.example.com du -h /var/hyperledger/production >> $ordererfile
done

org0file="org0-storage.log"
if [ -f $org0file ]; then
    rm $org0file
fi
for i in $(seq 1 2)
do
    echo "peer$i.org0.example.com" >> $org0file
    docker exec -ti peer$i.org0.example.com du -h /var/hyperledger/production >> $org0file
done

org1file="org1-storage.log"
if [ -f $org1file ]; then
    rm $org1file
fi
for i in $(seq 1 4)
do
    echo "peer$i.org1.example.com" >> $org1file
    docker exec -ti peer$i.org1.example.com du -h /var/hyperledger/production >> $org1file
done