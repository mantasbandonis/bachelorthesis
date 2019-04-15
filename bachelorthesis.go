package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Information struct {
	Id       string  `json:"Id"`
	Calories float64 `json:"Calories"`
	Device   string  `json:"Device"`
}

type PrivateInformation struct {
	Id    string  `json:"Id"`
	Name  string  `json:"Name"`
	Email string  `json:"Email"`
	Socnr float64 `json:"Socnr"`
}

// MainChaincode Defined for Chaincode Interface
type MainChaincode struct {
	// define to pmplement CC interface
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(MainChaincode))

	if err != nil {
		fmt.Printf("Error starting the MainChaincode Contract: %s", err)
	}
}

// Init is initalized and defined for Chaincode
func (t *MainChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("Chaincode Initialized")

	return shim.Success(nil)
}

// Invoke Chaincode
func (t *MainChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Start Invoke")
	defer fmt.Println("Stop Invoke")
	// Get function name and args
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "createInformation":
		return t.createInformation(stub, args)
	case "getInformation":
		return t.getInformation(stub, args)
	case "getInformationPriv":
		return t.getInformationPriv(stub, args)
	default:
		return shim.Error("Invalid Invoke")
	}
}

// CREATE Information
func (t *MainChaincode) createInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	id := args[0]
	calories, err1 := strconv.ParseFloat(args[1], 32)
	device := args[2]
	name := args[3]
	email := args[4]
	socnr, err2 := strconv.ParseFloat(args[5], 32)

	if err1 != nil || err2 != nil {
		return shim.Error("Error parsing the values")
	}

	information := &Information{id, calories, device}
	informationBytes, err3 := json.Marshal(information)

	if err3 != nil {
		return shim.Error(err2.Error())
	}

	informationPriv := &PrivateInformation{id, name, email, socnr}
	// Change struct into bytes for the Blockchain
	informationPrivBytes, err4 := json.Marshal(informationPriv)

	if err4 != nil {
		return shim.Error(err1.Error())
	}

	// We store the No Privacy Data with our open Data Collection
	err5 := stub.PutPrivateData("collectionThesis", id, informationBytes)
	if err5 != nil {
		return shim.Error(err5.Error())
	}

	err6 := stub.PutPrivateData("collectionPrivate", id, informationPrivBytes)
	if err6 != nil {
		return shim.Error(err6.Error())
	}

	jsonInformation, err7 := json.Marshal(information)
	if err7 != nil {
		return shim.Error(err7.Error())
	}
	return shim.Success(jsonInformation)
}

// GET TRANSACTION
func (t *MainChaincode) getInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	id := args[0]
	Information := Information{}

	informationBytes, err1 := stub.GetPrivateData("collectionThesis", id)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err2 := json.Unmarshal(informationBytes, &Information)
	if err2 != nil {
		fmt.Println("Error changing object with ID: " + id + "back into JSON ")
		return shim.Error(err2.Error())
	}

	jsonInformation, err3 := json.Marshal(Information)
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	return shim.Success(jsonInformation)
}

func (t *MainChaincode) getInformationPriv(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	id := args[0]
	PrivateInformation := PrivateInformation{}

	informationPrivBytes, err1 := stub.GetPrivateData("collectionPrivate", id)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err2 := json.Unmarshal(informationPrivBytes, &PrivateInformation)
	if err2 != nil {
		fmt.Println("Error changing object with ID: " + id + "back into Json")
		return shim.Error(err2.Error())
	}

	jsonPrivateInformation, err3 := json.Marshal(PrivateInformation)
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	return shim.Success(jsonPrivateInformation)
}
