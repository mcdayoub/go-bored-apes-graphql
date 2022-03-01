package listener

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mcdayoub/go-bored-apes-graphql/models"
	"github.com/sirupsen/logrus"
)

// Listener listens for new ethereum blocks
type Listener struct {
	EthereumURL     string
	Client          *ethclient.Client
	ContractAddress string

	// The listener could listen to other methods if desired
	Method string

	// The Send chan is specifically for TransferData, but this could be configured for other methods
	Send *chan models.TransferData
}

// NewListener inits a new listener with specific attributes to listen for
func NewListener(ethereumURL, contractAddress, method string, send *chan models.TransferData) *Listener {
	listener := Listener{
		EthereumURL:     ethereumURL,
		ContractAddress: contractAddress,
		Method:          method,
		Send:            send,
	}

	// Start the infura client
	// Use a plan that has Archive Data add on ($250/month) to get past transactions
	// Archive data includes data outside of the most recent 128 blocks
	client, err := ethclient.Dial(listener.EthereumURL)
	if err != nil {
		logrus.Fatal(err)
	}

	listener.Client = client

	return &listener
}

// Listen gets the newest ethereum block headers
func (listener *Listener) Listen() error {
	headers := make(chan *types.Header)

	// Create a subscription object
	sub, err := listener.Client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return err
	}

	// Listen for new messages
	for {
		select {
		case err := <-sub.Err():
			logrus.Error(err)
		case header := <-headers:
			block, err := listener.Client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				// Infura sometimes reports "got null header for uncle 0 of block XXX..."
				logrus.Error(err)
				break
			}
			for _, tx := range block.Transactions() {
				// tx.To() is nil if it is a contract-creation tx
				if tx.To() != nil {
					to := tx.To().Hex()

					// Check if the tx is interacting with the address of the listener's desired smart contract
					if to == listener.ContractAddress {
						listener.GetReceipt(tx)
					}
				}
			}
		}
	}
}

// GetReceipt gets the receipt of the transaction
func (listener *Listener) GetReceipt(tx *types.Transaction) {
	receipt, err := listener.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		logrus.Error(err)
		return
	}

	logs := receipt.Logs

	listener.HandleLogs(logs)
}

// HandleLogs iterates through the provided logs.
// If the 0 index matches the listener's method
// then it drops the transfer data into the output channel
func (listener *Listener) HandleLogs(logs []*types.Log) {
	for i := range logs {
		// Check if the logs show that the method is a Transfer
		if logs[i].Topics[0].Hex() == listener.Method {
			// Build the transfer data object
			t := toTransferData(logs[i].Topics, logs[i].TxHash)

			// Helpful log to see that the listener has found a BAYC transfer
			logrus.Info("Listener found Transfer: ", t)

			// Drop the transfer object into the send channel
			*listener.Send <- t
		}

		// Approval logs are useful for the sanity check that the listener is listening
		// During development, BAYC methods were mostly this type
		// For now these events are just logged
		if logs[i].Topics[0].Hex() == models.Approval {
			t := toApprovalData(logs[i].Topics)
			logrus.Info("Listener found Approval: ", t)
		}
	}
}

// toTransferData builds the object to hold the transfer data
func toTransferData(topics []common.Hash, txHash common.Hash) models.TransferData {
	transferData := models.TransferData{
		Transaction: txHash.Hex(),
		Transfer:    topics[0].Hex(),
		Sender:      topics[1].Hex(),
		Receiver:    topics[2].Hex(),
		TokenID:     topics[3].Hex(),
	}

	return transferData

}

// toApprovalData builds the object to hold the approval data
// This was primarily used for testing purposes
func toApprovalData(topics []common.Hash) models.ApprovalForAllData {
	approvalData := models.ApprovalForAllData{
		Owner:    topics[0].Hex(),
		Operator: topics[1].Hex(),
		Approved: topics[2].Hex(),
	}

	return approvalData
}

// PastBlocks could be used to get transfers from a range of blocks
// Useful in case you need to get data from past blocks
func (listener *Listener) PastBlocks(fromBlock, toBlock *big.Int) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock.Int64()),
		ToBlock:   big.NewInt(toBlock.Int64()),
		Addresses: []common.Address{
			common.HexToAddress(listener.ContractAddress),
		},
	}

	logs, err := listener.Client.FilterLogs(context.Background(), query)
	if err != nil {
		logrus.Error(err)
	}

	// Convert the logs to pointers
	logPointers := make([]*types.Log, len(logs))
	for i := 0; i < len(logs); i++ {
		logPointers = append(logPointers, &logs[i])
	}

	listener.HandleLogs(logPointers)
}
