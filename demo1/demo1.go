package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("init is running ")

	// 		0 		   1 
	// "1377****023", "2.5"
	var userId, value string
	var money float64
	var err error
	
	if len(args) > 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 - 2")
	}
	
	fmt.Println("- start init user")
	
	userId = args[0]
	money = 0.0
	if len(args) == 2 {
		money, _ = strconv.ParseFloat(args[1], 64)
	}
	
	value = strconv.FormatFloat(money, 'f', 6, 64)
	err = stub.PutState(userId + "balance", []byte(value))
	if err != nil {
		return nil, err
	}
	
	err = stub.PutState(userId + "points", []byte("0"))
	if err != nil {
		return nil, err
	}
	
	fmt.Println("new user " + userId + " is created")
	fmt.Println("account balance for user " + userId + " is " + value)
	fmt.Println("account points for user " + userId + " is " + "0")
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke is running " + function)
	
	// handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}
	if function == "recharge" {
		return t.recharge(stub, args)
	}
	if function == "change" {
		return t.change(stub, args)
	}
	if function == "settle" {
		return t.settle(stub, args)
	}
	
	fmt.Println("Invoke did not find function: " + function)
	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) recharge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId, value string
	var val_recharge, val_previous, val_end float64
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	userId = args[0]
	
	val, err := stub.GetState(userId + "balance")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get account balance for user" + userId + "\"}"
		return nil, errors.New(jsonResp)
	}
	if val == nil{
		val_previous = 0.0;
	}else{
		val_previous, _ = strconv.ParseFloat(string(val), 64);
	}

	val_recharge, _ = strconv.ParseFloat(args[1], 64);
	val_end = val_recharge + val_previous
	value = strconv.FormatFloat(val_end, 'f', 6, 64)
	err = stub.PutState(userId + "balance", []byte(value))

	if err != nil {
		return nil, err
	}

	fmt.Println("account balance for user " + userId + " is " + value)
	return nil, nil
}

func (t *SimpleChaincode) settle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId, value string
	var val_settle, val_previous, val_end float64
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	userId = args[0]
	
	val, err := stub.GetState(userId + "balance")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get account balance for user" + userId + "\"}"
		return nil, errors.New(jsonResp)
	}
	if val == nil{
		val_previous = 0.0;
	}else{
		val_previous, _ = strconv.ParseFloat(string(val), 64);
	}

	val_settle, _ = strconv.ParseFloat(args[1], 64);
	val_end = val_previous + val_settle
	value = strconv.FormatFloat(val_end, 'f', 6, 64)
	err = stub.PutState(userId + "balance", []byte(value))
	if err != nil {
		return nil, err
	}
    
	fmt.Println("account balance for user " + userId + " is " + value)
	fmt.Println("invoke is running change")
	args1 := []string{userId, "1"}
	return t.change(stub, args1)
}

func (t *SimpleChaincode) change(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId, value string
	var val_change, val_previous, val_end int
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	userId = args[0]
	
	val, err := stub.GetState(userId + "points")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get account points for user" + userId + "\"}"
		return nil, errors.New(jsonResp)
	}
	if val == nil{
		val_previous = 0;
	}else{
		val_previous, _ = strconv.Atoi(string(val));
	}

	val_change, _ = strconv.Atoi(args[1])
	val_end = val_change + val_previous
	value = strconv.Itoa(val_end)
	err = stub.PutState(userId + "points", []byte(value))
	if err != nil {
		return nil, err
	}

	fmt.Println("account points for user " + userId + " is " + string(val_end))
	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query is running " + function)
	
	// handle different functions
	if function == "queryBalance" {
		return t.queryBalance(stub, args)
	}
	if function == "queryPoints" {
		return t.queryPoints(stub, args)
	}
	if function == "queryAll" {
		return t.queryAll(stub, args)
	}
	
	fmt.Println("Query did not find function: " + function)
	return nil, errors.New("Received unknown function query")
}

func (t *SimpleChaincode) queryBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting userID of the user to query")
	}
	
	userId = args[0]

	val, err := stub.GetState(userId + "balance")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for user " + userId + "\"}"
		return nil, errors.New(jsonResp)
	}

	if val == nil {
		val = []byte("0.0")
	}

	jsonResp := "{\"user\":\"" + userId + "\",\"'s account balance is \":\"" + string(val) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil
}

func (t *SimpleChaincode) queryPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting userID of the user to query")
	}
	
	userId = args[0]

	val, err := stub.GetState(userId + "points")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for user " + userId + "\"}"
		return nil, errors.New(jsonResp)
	}

	if val == nil {
		val = []byte("0")
	}

	jsonResp := "{\"user\":\"" + userId + "\",\"'s account points is \":\"" + string(val) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil
}

func (t *SimpleChaincode) queryAll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting userID of the user to query")
	}
	
	userId = args[0]

	val1, err := stub.GetState(userId + "balance")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for user " + userId + "\"}"
		return nil, errors.New(jsonResp)
	}

	if val1 == nil {
		val1 = []byte("0.0")
	}

	jsonResp1 := "{\"user\":\"" + userId + "\",\"'s account balance is \":\"" + string(val1) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp1)
	
	val2, err := stub.GetState(userId + "points")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for user " + userId + "\"}"
		return nil, errors.New(jsonResp)
	}

	if val2 == nil {
		val2 = []byte("0")
	}

	jsonResp2 := "{\"user\":\"" + userId + "\",\"'s account points is \":\"" + string(val2) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp2)
	
	return []byte(jsonResp1 + jsonResp2), nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
