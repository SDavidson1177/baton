docker exec -it relayer1 ./setup_demo_1.sh

# echo "chain 2..."
# sleep 2

echo "Transfering funds..."

# Relayer 1
docker container exec -it chain0 batond tx bank send alice cosmos1m873lan7tvuvzahndzyl9lg9932azmelgz9yjv 10000000stake -y
docker container exec -it chain1 batond tx bank send alice cosmos1xk5nnhhq6842rnpqx80hjteq2ydkxyfjc884xg 10000000stake -y
docker container exec -it chain2 batond tx bank send alice cosmos1unh6ecqlkaa09t65dt5wyxt7yrv9ffkq9z7nnd 10000000stake -y

sleep 3 # wait for blocks to be processed

# Relayer 2
docker container exec -it chain0 batond tx bank send alice cosmos1aqm5m56tnf4ya2ukcwzmvx35lp0hzk8cl0grhp 10000000stake -y
docker container exec -it chain1 batond tx bank send alice cosmos153da2sfhlgst9r6qu6cdteqtpgs0xjpddrk6ej 10000000stake -y
docker container exec -it chain2 batond tx bank send alice cosmos1d9aadrju5rmz6gmcuashhsptx23c97j3kknqnz 10000000stake -y

echo "Setting up multihop channel baton-0 <--> baton-1 <--> baton-2"
docker exec -it relayer1 ./setup_demo_2.sh

echo "chain 2..."
# sleep 2

docker exec -it relayer2 ./setup_demo_1.sh

sleep 2

docker exec -it relayer2 ./setup_demo_2.sh
