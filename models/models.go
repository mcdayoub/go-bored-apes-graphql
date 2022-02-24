package models

// TransferInput holds the info that will be written to the graphQL server.
type TransferInput struct {
	Transaction string
	Sender      string
	Receiver    string
	TokenID     int
}

// TransferData holds the info for a Transfer event from the ethereum network.
type TransferData struct {
	Transaction string
	Transfer    string
	Sender      string
	Receiver    string
	TokenID     string
}

// ApprovalForAllData holds the info for an ApprovalForAll event from the ethereum network.
type ApprovalForAllData struct {
	Owner    string
	Operator string
	Approved string
}

// Approval hex is just for the purposes of testing
var Approval = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
