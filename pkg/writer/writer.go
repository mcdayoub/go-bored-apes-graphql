package writer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mcdayoub/go-bored-apes-graphql/models"
	"github.com/sirupsen/logrus"
)

type Writer struct {
	Receive *chan models.TransferData
}

// NewWriter returns a Writer that has an receive channel
func NewWriter(receive *chan models.TransferData) Writer {
	return Writer{
		Receive: receive,
	}
}

// Start makes the writer look for messages from the listener on the channel
// and write to the GraphQL server.
func (writer *Writer) Start() error {
	for transferData := range *writer.Receive {
		tokenID, err := strconv.Atoi(transferData.TokenID)
		if err != nil {
			logrus.Error(err)
			return err
		}

		transferInput := models.TransferInput{
			Transaction: transferData.Transaction,
			Sender:      transferData.Sender,
			Receiver:    transferData.Sender,
			TokenID:     tokenID,
		}

		write(transferInput)
	}
	return nil
}

// write creates the graphQL query to write to the server.
func write(transfer models.TransferInput) {
	// Post to the graphQL URL
	url := "http://localhost:8080/query"

	// Here is what each attribute needs to be in the graphQL mutation
	// transaction := `\"0xTransaction\"`
	// sender := `\"0xSender\"`
	// receiver := `\"0xReceiver\"`
	// token_id := "123"
	// read := "false"

	// Add the escape characters to each of the graphQL arguments
	transaction := `\"` + transfer.Transaction + `\"`
	sender := `\"` + transfer.Sender + `\"`
	receiver := `\"` + transfer.Receiver + `\"`
	tokenID := strconv.Itoa(transfer.TokenID)
	read := "false"

	// Create the mutation
	// I did not find a great go graphQL library to do this without some pain.
	// If this was in prod I would dig deeper and find a better way to generate this.
	s := fmt.Sprintf("{\"query\":\"mutation { createTransfer(input: {transaction: %v, sender: %v, receiver: %v, token_id: %v, read: %v}) {transaction sender receiver read token_id}} \"}`)", transaction, sender, receiver, tokenID, read)

	var jsonStr = []byte(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logrus.Error(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()

	// For debugging purposes we want to look at the status of the responses from the graphQL server.
	logrus.Info("response Status:", resp.Status)
	logrus.Info("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Info("response Body:", string(body))
}
