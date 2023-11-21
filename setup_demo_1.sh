#!/bin/bash

rly config init

rly chains add --file chain0-config.json baton-0
rly chains add --file chain1-config.json baton-1
rly chains add --file chain2-config.json baton-2

rly keys add baton-0 k0 > k0.txt
rly keys add baton-1 k1 > k1.txt
rly keys add baton-2 k2 > k2.txt

cat k*