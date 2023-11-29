docker exec -it relayer ./setup_demo_1.sh
echo "Transfering funds..."
# docker container exec -it chain0 batond tx bank send alice cosmos1m873lan7tvuvzahndzyl9lg9932azmelgz9yjv 10000000stake -y
# docker container exec -it chain1 batond tx bank send alice cosmos1xk5nnhhq6842rnpqx80hjteq2ydkxyfjc884xg 10000000stake -y
# docker container exec -it chain2 batond tx bank send alice cosmos1unh6ecqlkaa09t65dt5wyxt7yrv9ffkq9z7nnd 10000000stake -y
echo "Setting up multihop channel baton-0 <--> baton-1 <--> baton-2"
docker exec -it relayer ./setup_demo_2.sh
