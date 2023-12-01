package types

// When specifying versions for channels, they should be of the following format:
// "version1:metric/version2:metric..."
// For example, the version that is required to signal the ibc module to check that
// each intermediate chain has at least 10 validators is "validators:10"

const (
	VERSION_VALIDATORS = "validators"
)
