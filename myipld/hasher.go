package myipld

import (
	"crypto/sha256"
	"fmt"
)

const HashSize = sha256.Size

type MyCID struct {
	Hash [HashSize]byte
}

func (c MyCID) String() string {
	/* {comment}
	string returns a hexadecimal representation of the MyCID
	showing only the first 8 bytes for brevity in output
	{/comment} */
	return fmt.Sprintf("my-cid-%x", c.Hash[:8])
}

func ComputeSHA256(data []byte)(MyCID, error){
	hash := sha256.Sum256(data)
	return MyCID{Hash : hash}, nil
}