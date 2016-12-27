/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
    "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Customer struct {
	custName	string  `json:"custName"`
	points 		int		`json:"points"`
}

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
	

	var custA, custB string //custName
    var err error 
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	custA = args[0]
	custB = args[1]

	// Write the state to the ledger
	err = stub.PutState(custA, []byte("1000"))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(custB, []byte("1000"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

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

func (t *SimpleChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error){
	
    fmt.Println("transfer is running ")

	if len(args) != 3{
		return nil, errors.New("Incorrect Number of arguments.Expecting 3 for transfer")
	}

	var A, B string //custName
	var Apoints, Bpoints int // custName points
	var X int // transfer amt

	A = args[0]
	B = args[1]


	Apointsbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state of " + A)
	}
	if Apointsbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Apoints, _ = strconv.Atoi(string(Apointsbytes))

	Bpointsbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state of " + B)
	}
	if Bpointsbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bpoints, _ = strconv.Atoi(string(Bpointsbytes))

	// Perform the transfer
	X, err = strconv.Atoi(args[2])
	Apoints = Apoints - X
	Bpoints = Bpoints + X
	fmt.Printf("Apoints = %d, Bpoints = %d\n", Apoints, Bpoints)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Apoints)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bpoints)))
	if err != nil {
		return nil, err
	}

    return nil, nil
}
