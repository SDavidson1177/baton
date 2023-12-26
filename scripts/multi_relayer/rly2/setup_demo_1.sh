#!/bin/bash

# rly config init

# rly chains add --file chain0-config.json baton-0
# rly chains add --file chain1-config.json baton-1
# rly chains add --file chain2-config.json baton-2

mkdir /root/.relayer2
cp -r /root/.relayer/* /root/.relayer2/

rly keys restore baton-0 k3 "frown nature exile renew begin avocado situate connect unique grocery tribe volume foster month gauge expect street work kind keen festival better believe above" --home /root/.relayer2
rly keys restore baton-1 k4 "mother museum under glue broom race violin knife phone invest essay lonely illness churn minute dance cradle uncle alter night tag bounce film slogan" --home /root/.relayer2
rly keys restore baton-2 k5 "slice urban chaos pudding cat castle hint charge follow music volume grunt dynamic dismiss derive capital skate evolve lounge together garden reflect comfort wisdom" --home /root/.relayer2
