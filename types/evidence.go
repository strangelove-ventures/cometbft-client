package types

import (
	"fmt"
	"time"

	"github.com/strangelove-ventures/cometbft-client/crypto/merkle"
)

// Evidence represents any provable malicious activity by a validator.
// Verification logic for each evidence is part of the evidence module.
type Evidence interface {
	Bytes() []byte        // bytes which comprise the evidence
	Hash() []byte         // hash of the evidence
	Height() int64        // height of the infraction
	String() string       // string format of the evidence
	Time() time.Time      // time of the infraction
	ValidateBasic() error // basic consistency check
}

// EvidenceList is a list of Evidence. Evidences is not a word.
type EvidenceList []Evidence

// Hash returns the simple merkle root hash of the EvidenceList.
func (evl EvidenceList) Hash() []byte {
	// These allocations are required because Evidence is not of type Bytes, and
	// golang slices can't be typed cast. This shouldn't be a performance problem since
	// the Evidence size is capped.
	evidenceBzs := make([][]byte, len(evl))
	for i := 0; i < len(evl); i++ {
		// TODO: We should change this to the hash. Using bytes contains some unexported data that
		// may cause different hashes
		evidenceBzs[i] = evl[i].Bytes()
	}
	return merkle.HashFromByteSlices(evidenceBzs)
}

func (evl EvidenceList) String() string {
	s := ""
	for _, e := range evl {
		s += fmt.Sprintf("%s\t\t", e)
	}
	return s
}
