#!/bin/bash

rly keys use baton-1 k1 && rly keys use baton-2 k2

rly paths new baton-1 baton-2 path2
sleep 1
rly tx clients path2
sleep 4
rly tx connection path2
sleep 10
rly tx channel path2 --src-port transfer --dst-port transfer --order unordered --version "ics20-1/validators:1"