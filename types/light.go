package types

// LightBlock is a SignedHeader and a ValidatorSet.
// It is the basis of the light client
type LightBlock struct {
	*SignedHeader `json:"signed_header"`
	ValidatorSet  *ValidatorSet `json:"validator_set"`
}

// SignedHeader is a header along with the commits that prove it.
type SignedHeader struct {
	*Header `json:"header"`

	Commit *Commit `json:"commit"`
}
