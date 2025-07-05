package myipld

import (
	// "crypto/sha256"
	"encoding/json"
	"fmt"
)

type MyLink struct {
	Name string
	Cid  MyCID
}

/*{comment}

MyNode is representing a simplified ipld node
it holds arbitrary data and a list of links to other nodes

{/comment}*/

type MyNode struct {
	Data    json.RawMessage
	Links   []MyLink
	Cid     MyCID
	rawData []byte
}

func NewMyNode(data interface{}) (*MyNode, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to json : %w", err)
	}

	node := &MyNode{
		Data: dataBytes,
	}

	/* {comment}
	compute the cid based on the serialized content (data + links)
	for a new node, links are initially empty, but the cid still depends on
	the data
	*/

	if err != nil {
		return nil, fmt.Errorf("failed to compute CID for new new node, %w", err)
	}

	return node, nil
}

func (n *MyNode) AddLink(name string, targetCID MyCID) error {
	n.Links = append(n.Links, MyLink{Name: name, Cid: targetCID})
	return n.recomputeCID()
}

func (n *MyNode) recomputeCID() error {
	serializableNode := struct {
		Data  json.RawMessage `json:"data"`
		Links []MyLink        `json:"links"`
	}{
		Data:  n.Data,
		Links: n.Links,
	}

	rawBytes, err := json.Marshal(serializableNode)
	if err != nil {
		return fmt.Errorf("failed to marshal node for CID computations : %w", err)
	}

	// cache raw bytes for ToBytes method
	n.rawData = rawBytes

	cid, err := ComputeSHA256(rawBytes)

	if err != nil {
		return fmt.Errorf("failed to compute sha256 hash : %w", err)
	}

	n.Cid = cid
	return nil
}

func (n *MyNode) ToBytes() ([]byte, error) {
	if n.rawData == nil {
		/* {comment}
		this should ideally not happen if newMyNode or AddLink were used,
		but as a safeguard
		{/comment}*/

		if err := n.recomputeCID(); err != nil {
			return nil, err
		}
	}

	return n.rawData, nil
}

func FromBytes(data []byte) (*MyNode, error) {
	var serializableNode struct {
		Data  json.RawMessage `json:"data"`
		Links []MyLink        `json:"links"`
	}

	if err := json.Unmarshal(data, &serializableNode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes to MyNode : %w", err)
	}

	node := &MyNode {
		Data : serializableNode.Data,
		Links: serializableNode.Links,
	}

	if err := node.recomputeCID(); err != nil {
		return nil, fmt.Errorf("failed to recompute CID after deserialization : %w", err)
	}

	return node, nil

}
