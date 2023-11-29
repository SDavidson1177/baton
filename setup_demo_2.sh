#!/bin/bash

rly keys use baton-0 k0 && rly keys use baton-1 k1 && rly keys use baton-2 k2

rly paths new baton-0 baton-1 path1
sleep 1
rly tx clients path1
sleep 4
rly tx connection path1 --version "ics20-1-ORDER_ORDERED"
sleep 10
rly tx channel path1 --src-port transfer --dst-port transfer --version "ics20-1-ORDER_ORDERED"
sleep 10

# rly paths new baton-1 baton-2 path2
# sleep 1
# rly tx clients path2
# sleep 4
# rly tx connection path2
# sleep 10

# rly paths new baton-0 baton-2 path3 baton-1 
# sleep 1
# rly tx channel path3 --src-port transfer --dst-port transfer --order unordered --version ics20-1