// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen

type TranferInput struct {
	Transaction string `json:"transaction"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	TokenID     int    `json:"token_id"`
	Read        bool   `json:"read"`
}
