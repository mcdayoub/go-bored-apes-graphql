package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mcdayoub/go-bored-apes-graphql/graphql/gqlgen"
	"github.com/mcdayoub/go-bored-apes-graphql/graphql/pg"
	"github.com/mcdayoub/go-bored-apes-graphql/models"
	"github.com/mcdayoub/go-bored-apes-graphql/pkg/listener"
	"github.com/mcdayoub/go-bored-apes-graphql/pkg/writer"
	"github.com/sirupsen/logrus"
	tombv2 "gopkg.in/tomb.v2"
)

const (
	// infuraURL is how we get data from the ethereum network
	infuraURL = "wss://mainnet.infura.io/ws/v3/"

	// BoredApeContractAddress to track the BAYC transactions
	BoredApeContractAddress = "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"

	// TransferHex for ethereum transfers
	TransferHex = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

func run(infuraKey string) error {
	// Initialize tomb. This handles clean goroutine tracking and termination.
	tomb := &tombv2.Tomb{}

	// Run the API in a tomb
	tomb.Go(func() error {
		// initialize the db
		// This is hard coded but in prod it would be better to derive this with a config file.
		db, err := pg.Open("dbname=transfers sslmode=disable")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// initialize the repository
		repo := pg.NewRepository(db)

		// configure the server
		mux := http.NewServeMux()
		mux.Handle("/", gqlgen.NewPlaygroundHandler("/query"))
		mux.Handle("/query", gqlgen.NewHandler(repo))

		// run the server
		port := ":8080"
		fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
		fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))

		return nil
	})

	// Create the channel that the TransferData events are sent and received on
	transfers := make(chan models.TransferData, 100)

	// Start the listener for the bored ape transfers
	tomb.Go(func() error {
		logrus.Info("Starting Listener")

		// Combine the infuraURL and the infuraKey that was inputted when running the app.
		infuraURLAndKey := infuraURL + infuraKey

		// Here we set the attributes for the listener.
		// These are meant to be configurable in case we want to have a listener for different events and contracts.
		listener := listener.NewListener(infuraURLAndKey, BoredApeContractAddress, TransferHex, &transfers)
		err := listener.Listen()

		return err
	})

	// Start the writer that will write to our graphQL server
	tomb.Go(func() error {
		logrus.Info("Starting Writer")
		writer := writer.NewWriter(&transfers)
		err := writer.Start()

		return err
	})

	// Wait blocks until all goroutines have finished running.
	return tomb.Wait()
}

func main() {
	// infuraKey is the project id key from your infura account
	infuraKey := flag.String("key", "", "the infura project id")
	flag.Parse()

	if err := run(*infuraKey); err != nil {
		log.Fatal(err)
	}
}
