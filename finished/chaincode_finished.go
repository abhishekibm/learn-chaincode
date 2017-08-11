/*
My first chaincode :  Money payment
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//args[0] = "{\"zabi\":6.13,\"abhi\":23}"
	err := stub.PutState("moneyWorld", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	var dat map[string]int
	if err := json.Unmarshal([]byte(args[0]), &dat); err != nil {
		panic(err)
	}

	for k, v := range dat {
		err := stub.PutState(k, []byte(strconv.Itoa(v)))
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "transferMoney" {
		return t.transferMoney(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Pff...Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) transferMoney(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var sender, recipient, jsonResp string
	var amount int
	var err error
	fmt.Println("running transferMoney()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. Sender, recipient, amount")
	}

	sender = args[0]    //sender
	recipient = args[1] // recipient
	amount, _ = strconv.Atoi(args[2])

	//first get current state
	currState, err := stub.GetState("moneyWorld")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "moneyWorld" + "\"}"
		return nil, errors.New(jsonResp)
	}

	var dat map[string]int
	if err := json.Unmarshal([]byte(currState), &dat); err != nil {
		panic(err)
	}
	dat[sender] = dat[sender] - amount
	dat[recipient] = dat[recipient] + amount
	strB, _ := json.Marshal(dat)
	err1 := stub.PutState("moneyWorld", []byte(strB))
	if err1 != nil {
		return nil, err1
	}

	for k, v := range dat {
		err := stub.PutState(k, []byte(strconv.Itoa(v)))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	fmt.Println("Updated")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
