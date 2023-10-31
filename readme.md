# baton
**baton** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli). Baton is configured to operate using mutlihop IBC (ics-033). As time goes on, more and more features will be added. Please see section "Configuring a Multihop Testnet" for instructions on how to run a local testnet with multihop configured.

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Ignite CLI has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/baton@latest! | sudo bash
```
`username/baton` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Configuring a Multihop Testnet

This blockchain has a local version of the Go relayer stored within *external/*. Furthermore, the specification for multihop IBC has been implemented in the x/ibcx-go directory. Both the modified Go relayer and ibc-go implementation were provided by PolymerLabs. In order to build this testnet, you will need to have **Docker Compose** installed on your machine.

### Running the blockchains
The testnet is configured with three blockchains and a single relayer. The setup can be visualized as follows:
```
||chain0||     ||chain1|| ||chain2||
    ^              ^          ^
    |              |          |
    ---------------------------
                   |
              ||relayer||
```

The blockchains as well as the relayer run in their own Docker containers. To start the docker containers in detached mode, run the following command in the root directory:

```Bash
docker compose up -d --build
```

### Configuring the Relayer
Once the four docker containers are running, start the relayer's shell by running:

```Bash
docker container exec -it relayer /bin/bash
```

The working directory within the relayer container should be */go/*. Now, run the following commands:

```Bash
rly config init
```

This initializes the relayer configuration file.

```Bash
rly chains add --file chain0-config.json baton0 && rly chains add --file chain1-config.json baton1 && rly chains add --file chain2-config.json baton2 
```

This adds the chains' data to the relayer's configuration file.

```Bash
rly keys add baton-0 k0 > k0.txt && rly keys add baton-1 k1 > k1.txt && rly keys add baton-2 k2 > k2.txt
```

This creates three keys. These keys are for accounts that will eventually be used by the relayer, however they currently do not have any tokens. To fund these accouts, you will need to execute BASH for each of the blockchain containers and send "10000000stake" tokens from the "alice" account to the corresponding account for that blockchain. You can find each account's address by reading the files k0.txt, k1.txt and k2.txt. To fund the key "k0" for example, you would run

```Bash
batond tx bank send alice <k0 account address> 10000000stake
```

from within chain0. Alternatively, you may run

```Bash
docker container exec -it chain0 batond tx bank send alice <k0 account address> 10000000stake -y
```

from your local machine.

Now that every account is funded, you must tell the relayer which keys to use for which blockchains. Run the following within the relayer container:

```Bash
rly keys use baton-0 k0 && rly keys use baton-1 k1 && rly keys use baton-2 k2
```

Now you must setup the paths between each of the blockchains. Run the following commands:

```Bash
rly paths new baton-0 baton-1 path1 && rly tx clients path1 && rly tx connection path1
```

This establishes a connection between chain0 and chain1. For the connection between chain1 and chain2, run:

```Bash
rly paths new baton-1 baton-2 path2 && rly tx clients path2 && rly tx connection path2
```

Now we must create a path from chain0 to chain2 that uses chain1 as an intermediate hop. This path will piggyback off of the two connections we have already established. Run:

```Bash
rly paths new baton-0 baton-2 path3 baton-1 && \
rly tx channel path3 --src-port transfer --dst-port transfer --order unordered --version ics20-1
```

You can check that the channel endpoints are correctly configured on baton-0 and baton-2 (chain0 and chain2) by running:

```Bash
rly q channels baton-0 && rly q channels baton-2
```

The channel states should be open. Finally, you may start the relayer by running:

```Bash
rly start
```

Now the relayer is listening for messages to send between chain0 and chain2!

### Simple Multihop Test

To test the multihop connection, first enter into chain0's bash. Now, run the following command:

```Bash
batond tx ibc-transfer transfer transfer channel-0 <account on chain2> 100000stake --from alice --fees 4000stake --packet-timeout-height 0-0
```

**Important**: As of the time of writing, it is important that you set the timeout height to "0-0". This effectively ignores packet timeouts. The timeout feature is not fully functional with multihop IBC yet.

You can then query the account balance on chain2 using the following command (make sure you are now in chain2's container):

```Bash
batond query bank balances <account on chain2>
```

It may take some time for the balance to show up, as new blocks need to be committed, but you should eventually see the IBC transfer.

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
