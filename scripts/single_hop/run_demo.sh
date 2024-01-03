# Copy scripts for this demo into the relayer container
docker cp ./setup_demo_1.sh relayer1:/go/baton/external/setup_demo_1.sh
docker cp ./setup_demo_2.sh relayer1:/go/baton/external/setup_demo_2.sh

docker exec -it relayer1 ./setup_demo_1.sh
echo "Transfering funds..."

# Relayer 1
docker container exec -it chain1 batond tx bank send alice cosmos1xk5nnhhq6842rnpqx80hjteq2ydkxyfjc884xg 10000000stake -y
docker container exec -it chain2 batond tx bank send alice cosmos1unh6ecqlkaa09t65dt5wyxt7yrv9ffkq9z7nnd 10000000stake -y

echo "Setting up single hop channel baton-1 <--> baton-2"
docker exec -it relayer1 ./setup_demo_2.sh

