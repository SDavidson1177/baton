#!/bin/bash

for i in {1..25}; do
	docker exec -it chain1 batond tx ibc-transfer transfer transfer channel-1 cosmos1unh6ecqlkaa09t65dt5wyxt7yrv9ffkq9z7nnd 100000stake --from alice --fees 4000stake --packet-timeout-height 0-0 -y
	sleep 2
done
