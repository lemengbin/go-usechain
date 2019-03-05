package manager

import (
	"github.com/usechain/go-usechain/core/state"
	"github.com/usechain/go-usechain/common"
	"github.com/usechain/go-usechain/crypto/sha3"
	"github.com/usechain/go-usechain/common/hexutil"
	"strings"
)

const (
	ManagerContract 		= "0xfffffffffffffffffffffffffffffffff0000003"
	ABI 					= "[{\"constant\": true,\"inputs\": [{\"name\": \"_index\",\"type\": \"uint256\"}],\"name\": \"getCommittee\",\"outputs\": [{\"name\": \"\",\"type\": \"address\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"whichRound\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"mode\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"getCommitteeIndex\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [{\"name\": \"index\",\"type\": \"uint256\"}],\"name\": \"getCommitteeAsymkey\",\"outputs\": [{\"name\": \"\",\"type\": \"string\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [{\"name\": \"_index\",\"type\": \"uint256\"}],\"name\": \"getCandidate\",\"outputs\": [{\"name\": \"\",\"type\": \"address\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"getBlockNumber\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"IsCommittee\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"isEntireConfirmed\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"Requirement\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": false,\"inputs\": [{\"name\": \"_candidate\",\"type\": \"address\"}],\"name\": \"vote\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"getCommitteePubkey\",\"outputs\": [{\"name\": \"\",\"type\": \"string\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"getCandidateLen\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"name\": \"rounds\",\"outputs\": [{\"name\": \"selected\",\"type\": \"bool\"},{\"name\": \"committeePublicKey\",\"type\": \"string\"},{\"name\": \"committeePublicKey_candidate\",\"type\": \"string\"},{\"name\": \"confirmCount\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [{\"name\": \"_candidate\",\"type\": \"address\"}],\"name\": \"getVotes\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": false,\"inputs\": [{\"name\": \"_pubkey\",\"type\": \"string\"}],\"name\": \"confirmCommitteePubkey\",\"outputs\": [],\"payable\": false,\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"MAX_COMMITTEEMAN_COUNT\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": false,\"inputs\": [{\"name\": \"_pubkey\",\"type\": \"string\"}],\"name\": \"uploadCommitteePubkey\",\"outputs\": [],\"payable\": false,\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"getCommitteeConfirmStat\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": false,\"inputs\": [{\"name\": \"_flag\",\"type\": \"bool\"}],\"name\": \"controlVote\",\"outputs\": [],\"payable\": false,\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"vote_enabled\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": false,\"inputs\": [],\"name\": \"confirmVoting\",\"outputs\": [],\"payable\": false,\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"constant\": false,\"inputs\": [{\"name\": \"_asymPubkey\",\"type\": \"string\"}],\"name\": \"confirmAndKeyUpload\",\"outputs\": [],\"payable\": false,\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"Election_duration\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"},{\"constant\": true,\"inputs\": [],\"name\": \"Election_cycle\",\"outputs\": [{\"name\": \"\",\"type\": \"uint256\"}],\"payable\": false,\"stateMutability\": \"view\",\"type\": \"function\"}]"
)

// Get the committee max count based on state db reader
// Normal mode: 5 committees
// Single mode: 1 committees
func GetCommitteeCount(statedb *state.StateDB) uint {
	// detect the contract running mode
	res := statedb.GetState(common.HexToAddress(ManagerContract), common.BigToHash(common.Big0))

	// if running in single mode, just one committee
	if res.Big().Cmp(common.Big1) == 0 {
		return 1
	}
	// if normal mode, should be 5 committees
	return 5
}


//Check the address whether be a committee based on state db reader
func IsCommittee(statedb *state.StateDB, addr common.Address) bool {
	// detect the contract running mode
	res := statedb.GetState(common.HexToAddress(ManagerContract), common.BigToHash(common.Big0))

	// if running in single mode, anyone could be committee
	if res.Big().Cmp(common.Big1) == 0 {
		return true
	}
	// if normal mode, need strict check based on state db
	///TODO:change the committeeOnDuty to map struct for decline db queries
	hash := sha3.NewKeccak256()
	hash.Write(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000002"))
	var keyIndex []byte
	keyIndex = hash.Sum(keyIndex)
	for i := int64(0); i < int64(common.MaxCommitteemanCount); i++ {
		res := statedb.GetState(common.HexToAddress(ManagerContract), common.HexToHash(common.IncreaseHexByNum(keyIndex, i)))
		if strings.EqualFold(addr.Hash().String(), res.String()) {
			return true
		}
	}
	return false
}
