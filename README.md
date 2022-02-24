# go-bored-apes-graqphql

### Go / GraphQL / Bored Ape Yacht Club

This is a GraphQL server written in Go that listens to the ethereum transfer events for Bored Ape Yacht Club (BAYC).

These events are found at:
https://etherscan.io/address/0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d

### Prereqs to running the app
* Clone this repo
* Download postgres. On mac I use https://postgresapp.com/

### Running the app
* Create the db: in the terminal run `createdb transfers`
* From the `/graphql` directory run `psql -d transfers -a -f schema.sql`
* Run `./api` from the root of this directory
* Go to `http://localhost:8080` to play in the GraphQL playground
* Wait for some BAYC transactions to happen

If you'd like to run the app yourself with Go, make sure to have go 1.17 installed.
You can find downloads here: https://go.dev/dl/

### GraphQL Playground
In the playground you can find the queries and mutations used to:
* query for `transfers` by `transaction`, `token_id`, `sender`, `receiver`
* read `transfers` by `transaction` 
* create `transfers`

Here are a few examples:
```
mutation {
  createTransfer(
    input: { transaction: "0xTransaction", sender: "0xSender", receiver: "0xReceiver", token_id: 101, read:false }
  ) {
    transaction
    sender
    receiver
    read
    token_id
  }
}
```

```
query {
  unreadTransfers {
    transaction
    sender
    receiver
    read
  }
}
```

```
mutation {
  readTransfer(transaction: "0xTransaction") {
    transaction
    sender
    receiver
    read
    token_id
  }
}
```

### Tools and References
* Infura for ethereum network information https://infura.io/
* `go-ethereum` https://github.com/ethereum/go-ethereum
* `gqlgen` https://github.com/99designs/gqlgen
* This blog post for a quick briefing on how to gen a lot of the GraphQL + Postgres code https://w11i.me/graphql-server-go-part1


### Takeaways (If this were in prod I would instead do...)
* I looked into ent: https://entgo.io/ and it seems suitable for building large apps, but not suitable for this small project.
* The postgres setup could be done better with `docker-compose`
* There are a lot of hard coded values that would live better as a kube secret or an environment variable.
* Tests: I want to add tests to this app. A lot of the testing so far is just sanity checks via
  * GraphQL queries and mutations in the playground
  * Log messages while the app is running
* If I were to make tests I would want:
  * Unit tests for each of the methods that handle data conversions (eth hex -> transfer data)
  * Service tests for GraphQL + Postgres
* BAYC kinda ugly tho

![image](https://user-images.githubusercontent.com/38268139/155456881-fc5d954e-80b6-4ed1-9e1e-c1d3cd7d17ca.png)
