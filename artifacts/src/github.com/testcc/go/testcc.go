package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SampleContract struct {
}

type User struct {
	ID                   string   `json:"id"`
	Email                string   `json:"email"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
}

func (s *SampleContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SampleContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()
	switch function {
	case "createUser":
		return createUser(APIstub, args)
	case "getUser":
		return getUser(APIstub, args)
	default:
		return shim.Error(fmt.Sprintf("%s", "No such methods..."))
	}
}

func createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var obj = args[0]
	var user User
	json.Unmarshal([]byte(obj), &user)
	fmt.Printf("obj:%s", obj)
	userKey, _ := stub.CreateCompositeKey("user", []string{user.ID})
	userAsBytes, _ := stub.GetState(userKey)
	if userAsBytes != nil {
		return shim.Error("Asset already exist")
	}
	if err := stub.PutState(userKey, []byte(obj)); err != nil {
		fmt.Printf("Error in storing user information", err)
		return shim.Error(fmt.Sprintf("%s", err))
	}
	return shim.Success(nil)
}

func getUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userKey, _ := stub.CreateCompositeKey("user", []string{args[0]})
	val, err := stub.GetState(userKey)
	fmt.Printf("\nID: %s\n", args[0])
	if err != nil {
		fmt.Printf("[ERROR] cannot get state, because of %s\n", err)
		return shim.Error(fmt.Sprintf("%s", err))
	} else {
		user := User{}
		json.Unmarshal(val, &user)
		fmt.Printf("Data available..", user)
		return shim.Success(val)
	}
}

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SampleContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
