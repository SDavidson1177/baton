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

The blockchains as well as the relayer run in their own Docker containers. First, you will need to build the most recent version of the blockchain. To do that, run the following command from the root directory:

```Bash
ignite chain build
```

To start the docker containers in detached mode, run the following command in the root directory:

```Bash
docker compose up -d --build
```

You **must** wait for the three blockchains to finish being initialized before proceeding to the next steps. You can check the logs for each chain by running the command:

```Bash
docker logs (chain0|chain1|chain2)
```

You will know that the chain is running when you see a bunch of verbose logging.

### Configuring the Demo

Once all three chains are running, execute the *./run_demo.sh* script from your root directory. This set up a channel from baton-0 to baton-2 with baton-1 as an intermediate hop. Once the script has finished running, enter into the relayer's bash by running:

```Bash
docker exec -it relayer /bin/bash
```

Start the relayer by running

```Bash
rly start
```

### Simple Multihop Test

To test the multihop connection, first enter into chain0's bash (from a seperate terminal). Now, run the following command:

```Bash
batond tx ibc-transfer transfer transfer channel-0 <account on chain2> 100000stake --from alice --fees 4000stake --packet-timeout-height 0-0
```

**Important**: As of the time of writing, it is important that you set the timeout height to "0-0". This effectively ignores packet timeouts. The timeout feature is not fully functional with multihop IBC yet.

You can then query the account balance on chain2 using the following command (make sure you are now in chain2's bash):

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
