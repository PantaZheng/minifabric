/*
Copyright 2009-2019 SAP SE or an SAP affiliate company. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"
)

//TODO： 依赖不应那么多，需要根据最新的实现来修改代码写法

// ManagementChaincode serves functionalities to communicate channel updates and signatures between different channel members.
type ManagementChaincode struct {
	contractapi.Contract
}

// Proposal gathers all information of a proposed update, including all added signatures.
type Proposal struct {
	// Description describes the proposal.
	Description string `json:"description,omitempty"`

	// Creator contains the msp ID of the proposal creator.
	Creator string `json:"creator"`

	// ConfigUpdate contains the base64 string representation of the common.ConfigUpdate.
	ConfigUpdate string `json:"config_update"`

	// Signatures contains a map of signatures: mspID -> base64 string representation of common.ConfigSignature
	Signatures map[string]string `json:"signatures,omitempty"`
}

const (
	NewProposalEvent    = "newProposalEvent"
	DeleteProposalEvent = "deleteProposalEvent"
	NewSignatureEvent   = "newSignatureEvent"
)

// ErrProposalNotFound is returned when the requested object is not found.
var ErrProposalNotFound = fmt.Errorf("proposal not found")

func getMSPID(creator []byte) (mspID string, err error) {
	identity := &msp.SerializedIdentity{}
	if err = proto.Unmarshal(creator, identity); err != nil {
		return "", fmt.Errorf("error happened unmarshalling the creator: %v", err)
	}
	return identity.Mspid, err
}

// getProposal fetches and decodes the proposal with the given id from the state or returns an error.
func getProposal(ctx contractapi.TransactionContextInterface, id string) (proposal *Proposal, err error) {
	proposalJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("error happened reading proposal with id (%s): %v", id, err)
	}
	if len(proposalJSON) == 0 {
		return nil, ErrProposalNotFound
	}
	err = json.Unmarshal(proposalJSON, proposal)
	if err != nil {
		return nil, fmt.Errorf("error happened unmarshalling the proposal JSON representation to struct: %v", err)
	}
	return proposal, err
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (mcc *ManagementChaincode) Init(contractapi.TransactionContextInterface) (err error) {
	return err
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
//func (mcc *ManagementChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
//	function, args := stub.GetFunctionAndParameters()
//	switch function {
//	case "proposeUpdate":
//		return mcc.proposeUpdate(stub, args)
//	case "addSignature":
//		return mcc.addSignature(stub, args)
//	case "getProposals":
//		return mcc.getProposals(stub, args)
//	case "getProposal":
//		return mcc.getProposal(stub, args)
//	case "deleteProposal":
//		return mcc.deleteProposal(stub, args)
//	default:
//		return shim.Error("Invalid invoke function name. Expecting \"proposeUpdate\" \"addSignature\" \"getProposals\" \"getProposal\" \"deleteProposal\".")
//	}
//}

// ProposeUpdate creates a new proposal containing the given update and a description.
//
// Arguments:
//   0: proposalID  - the ID of the new proposal
//   1: update      - base64 encoded proto/common.ConfigUpdate
//   2: description - a short of the update
//
// Returns:
//   the ID of the created proposal
//
// Events:
//   name: newProposalEvent(<proposalID>)
//   payload: ID of the proposal
//
func (mcc *ManagementChaincode) ProposeUpdate(ctx contractapi.TransactionContextInterface, proposalID, configUpdate, description string) (msg interface{}, err error) {

	// check if the configUpdate is in the correct format: base64 encoded proto/common.ConfigUpdate
	update, err := base64.StdEncoding.DecodeString(configUpdate)
	if err != nil {
		return nil, fmt.Errorf("error happened decoding the configUpdate base64 string: %v", err)
	}
	err = proto.Unmarshal(update, &common.ConfigUpdate{})
	if err != nil {
		return nil, fmt.Errorf("error happened decoding common.ConfigUpdate: %v", err)
	}

	_, err = getProposal(ctx, proposalID)
	if err != ErrProposalNotFound {
		return nil, fmt.Errorf("proposalID already in use")
	}

	// create and store the proposal
	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return nil, fmt.Errorf("error happened reading the transaction creator: " + err.Error())
	}
	mspID, err := getMSPID(creator)
	if err != nil {
		return nil, err
	}
	proposal := Proposal{
		ConfigUpdate: configUpdate,
		Description:  description,
		Creator:      mspID,
	}
	proposalJSON, err := json.Marshal(proposal)
	if err != nil {
		return nil, fmt.Errorf("error happened marshalling the new proposal: " + err.Error())
	}
	err = ctx.GetStub().PutState(proposalID, proposalJSON)
	if err != nil {
		return nil, fmt.Errorf("error happened persisting the new proposal on the ledger: " + err.Error())
	}
	err = ctx.GetStub().SetEvent(fmt.Sprintf("%s(%s)", NewProposalEvent, proposalID), []byte(proposalID))
	if err != nil {
		return nil, fmt.Errorf("error happened emitting event: " + err.Error())
	}
	return json.Marshal(map[string]interface{}{"proposal_id": proposalID})
}

// addSignature adds (or updates) a signature of the calling organization to the proposal.
//
// Arguments:
//   0: proposalID - the ID of the proposal where the signature is added
//   1: signature  - base64 encoded proto/common.ConfigSignature
//
// Events:
//   name: newSignatureEvent(<proposalID>)
//   payload: ID of the proposal
//
func (mcc *ManagementChaincode) addSignature(ctx contractapi.TransactionContextInterface, proposalID, signature string) (err error) {
	// check if the signature is in the correct format: base64 encoded proto/common.ConfigSignature
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("error happened decoding the signature base64 string: %v", err.Error())
	}
	err = proto.Unmarshal(sig, &common.ConfigSignature{})
	if err != nil {
		return fmt.Errorf("error happened decoding common.ConfigSignature: %v", err.Error())
	}

	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return fmt.Errorf("error happened reading the transaction creator: %v", err.Error())
	}
	mspID, err := getMSPID(creator)
	if err != nil {
		return err
	}

	// fetch and update the state of the proposal
	proposal, err := getProposal(ctx, proposalID)
	if err != nil {
		return err
	}
	if proposal.Signatures == nil {
		proposal.Signatures = make(map[string]string)
	}
	proposal.Signatures[mspID] = signature

	// store the updated proposal
	proposalJSONUpdated, err := json.Marshal(proposal)
	if err != nil {
		return fmt.Errorf("error happened marshalling the updated proposal: %v", err.Error())
	}
	err = ctx.GetStub().PutState(proposalID, proposalJSONUpdated)
	if err != nil {
		return fmt.Errorf("error happened persisting the updated proposal on the ledger: %v", err.Error())
	}
	err = ctx.GetStub().SetEvent(fmt.Sprintf("%s(%s)", NewSignatureEvent, proposalID), []byte(proposalID))
	if err != nil {
		return fmt.Errorf("error happened emitting event: %v", err.Error())
	}
	return err
}

// deleteProposal deletes the proposal with the given ID from the state.
// This can only be called by the proposal creator.
//
// Arguments:
//   0: proposalID - the ID of the proposal where the signature is added
//
// Events:
//   name: deleteProposalEvent(<proposalID>)
//   payload: ID of the proposal
//
func (mcc *ManagementChaincode) deleteProposal(ctx contractapi.TransactionContextInterface, proposalID string) (err error) {

	// fetch proposal
	proposal, err := getProposal(ctx, proposalID)
	if err != nil {
		return fmt.Errorf("error happened fetch proposal: " + err.Error())
	}

	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return fmt.Errorf("error happened reading the transaction creator: " + err.Error())
	}
	mspID, err := getMSPID(creator)
	if err != nil {
		return fmt.Errorf("error happened get mspID: " + err.Error())
	}

	// check if calling organization is proposal creator
	if proposal.Creator != mspID {
		return fmt.Errorf("forbidden. only the proposal creator (%s) can delete the proposal", proposal.Creator)
	}

	// delete the proposal
	err = ctx.GetStub().DelState(proposalID)
	if err != nil {
		return fmt.Errorf("error happened deleting the state: %v", err)
	}
	err = ctx.GetStub().SetEvent(fmt.Sprintf("%s(%s)", DeleteProposalEvent, proposalID), []byte(proposalID))
	if err != nil {
		return fmt.Errorf("error happened emitting event: " + err.Error())
	}
	return err
}

// getProposals returns all proposals.
//
// Arguments: none
//
// Returns:
//   a map from proposalID to proposal
//
func (mcc *ManagementChaincode) getProposals(ctx contractapi.TransactionContextInterface) (proposalsJSON []byte, err error) {
	proposals := make(map[string]*Proposal)
	proposalIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("error happened reading keys from ledger: " + err.Error())
	}
	defer func() {
		_ = proposalIterator.Close()
	}()

	for proposalIterator.HasNext() {
		proposalJSON, err := proposalIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("error happened iterating over available proposals: " + err.Error())
		}
		proposal := &Proposal{}
		err = json.Unmarshal(proposalJSON.Value, proposal)
		if err != nil {
			return nil, fmt.Errorf("error happened unmarshalling a proposal JSON representation to struct: " + err.Error())
		}
		proposals[proposalJSON.Key] = proposal
	}

	proposalsJSON, err = json.Marshal(proposals)
	if err != nil {
		return proposalsJSON, fmt.Errorf("error happened marshalling the update proposals: " + err.Error())
	}
	return proposalsJSON, err
}

// getProposal returns the proposal with the given ID.
//
// Arguments:
//   0: proposalID - the ID of a proposal
//
// Returns:
//   the proposal with the given ID
//
func (mcc *ManagementChaincode) getProposal(ctx contractapi.TransactionContextInterface, proposalID string) (proposalJSON []byte, err error) {
	proposalJSON, err = ctx.GetStub().GetState(proposalID)
	if err != nil {
		return proposalJSON, fmt.Errorf("error happened reading proposal with id (%v): %v", proposalID, err)
	}
	if len(proposalJSON) == 0 {
		return proposalJSON, fmt.Errorf(fmt.Sprintf("proposal with id (%s) not found", proposalID))
	}
	return proposalJSON, err
}

func main() {
	var err error
	rand.Seed(time.Now().UTC().UnixNano())
	cc, err := contractapi.NewChaincode(new(ManagementChaincode))
	if err != nil {
		panic(err.Error())
	}
	err = cc.Start()
	if err != nil {
		fmt.Printf("Error starting management chaincode: %s", err)
	}
}
