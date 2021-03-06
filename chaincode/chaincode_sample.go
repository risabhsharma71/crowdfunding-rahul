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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var userIndexStr = "_userindex"
var campaignIndexStr = "_campaignindex"
var transactionIndexStr = "_transactionindex"

type User struct {
	User_Type        string `json:"usertype"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            int    `json:"phone"`
	Password         string `json:"password"`
	ReTypePassword   string `json:"retypepassword"`
	Documenttype     string `json:"documenttype"`
	Organisationtype string `json:"organisationtype"`
	Facebook         string `json:"facebook"`
	Blog             string `json:"blog"`
	Websiteurl       string `json:"websiteurl"`
	Youtube          string `json:"youtube"`
	Designation      string `json:"designation"`
}
type AllUsers struct {
	Userlist []User `json:"userlist"`
}
type SessionAunthentication struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
type Session struct {
	StoreSession []SessionAunthentication `json:"session"`
}

type CreateCampaign struct {
	Campaign_Id            int       `json:"campaignid"`
	Campaign_Title         string    `json:"campaigntitle"`
	Campaign_Category      string    `json:"campaigncategory"`
	Campaign_Story         string    `json:"campaignstory"`
	Campaign_Description   string    `json:"campaigndescription"`
	Campaign_Image         string    `json:"campaignimage"`
	Campaign_Video         string    `json:"campaignvideo"`
	Campaign_Goal_amount   int       `json:"campaigngoalamount"`
	Campaign_Creation_Date string    `json:"campaigncreationdate`
	Project_Start_Date     string    `json:"projectstartdate"`
	Initial_amount         int       `json:"initialamount"`
	Campaign_End_Date      string    `json:"campaignenddate"`
	Rewards                []Rewards `json:"rewards"`
}
type CampaignList struct {
	Campaignlist []CreateCampaign `json:"campaignlist"`
}
type Rewards struct {
	Reward_Title        string `json:"reawardtitle"`
	Reward_Offer_Amount int    `json:"reawardofferamount"`
	Reward_Description  string `json:"reawardescription"`
}

type Login struct {
	Emailid      string `json:"emailid"`
	UserPassword string `json:"userpassword"`
}

type Contribution struct {
	UserId            string `json:"userid"`
	AmountContributed int    `json:"amountcontributed"`
	TransactionStatus string `json:"transactionstatus"`
	Campaign_Id       int    `json:"campaignid"`
}
type ContributionList struct {
	Contributionlist []Contribution `json:"contributionlist"`
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

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//_, args := stub.GetFunctionAndParameters()
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(userIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is ur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)

	} else if function == "registerUser" {
		return t.registerUser(stub, args)

	} else if function == "Delete" {
		return t.Delete(stub, args)

	} else if function == "SaveSession" {
		return t.SaveSession(stub, args)

	} else if function == "CreateCampaign" {
		return t.CreateCampaign(stub, args)

	} else if function == "AddReward" {
		return t.AddReward(stub, args)

	} else if function == "Contribute" {
		return t.Contribute(stub, args)

	} else if function == "UpdateTxStatus" {
		return t.UpdateTxStatus(stub, args)

	}

	/*else if function == "create_campaign" {
		return t.create_campaign(stub, args)

	}else if function == "transactions" {
		return t.transactions(stub, args)

	}*/

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "readuser" { //read a variable
		return t.readuser(stub, args)
	} else if function == "login" {
		return t.login(stub, args)

	} else if function == "auntheticatetoken" {
		return t.SetUserForSession(stub, args)

	}
	/*else if function == "readcampaign" {
			return t.readcampaign(stub, args)
	    }else if function == "readtxdetails" {
			return t.readtxdetails(stub, args)
	}*/
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair

func (t *SimpleChaincode) readuser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error
	//var campaign_title,jsonResp string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

func (t *SimpleChaincode) registerUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//input sanitation
	fmt.Println("- start user register")

	user := User{}
	user.User_Type = args[0]
	user.Name = args[1]
	user.Email = args[2]
	user.Phone, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Failed to get phone as cannot convert it to int")
	}
	user.Password = args[4]
	user.ReTypePassword = args[5]
	user.Documenttype = args[6]
	user.Facebook = args[7]
	user.Blog = args[8]
	user.Websiteurl = args[9]
	user.Youtube = args[10]
	user.Organisationtype = args[11]
	user.Designation = args[12]

	UserAsBytes, err := stub.GetState("getcfusers")
	if err != nil {
		return nil, errors.New("Failed to get users")
	}
	var allusers AllUsers
	json.Unmarshal(UserAsBytes, &allusers) //un stringify it aka JSON.parse()

	allusers.Userlist = append(allusers.Userlist, user)
	fmt.Println("allusers", allusers.Userlist) //append to allusers
	fmt.Println("! appended user to allusers")
	jsonAsBytes, _ := json.Marshal(allusers)
	fmt.Println("json", jsonAsBytes)
	err = stub.PutState("getcfusers", jsonAsBytes) //rewrite allusers
	if err != nil {
		return nil, err
	}
	fmt.Println("- end user_register")
	return nil, nil
}
func (t *SimpleChaincode) login(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	//input sanitation
	fmt.Println("- login")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	emailid := args[0]

	password := args[1]

	UserAsBytes, err := stub.GetState("getcfusers")
	if err != nil {
		return nil, errors.New("Failed to get users")
	}
	var allusers AllUsers
	json.Unmarshal(UserAsBytes, &allusers) //un stringify it aka JSON.parse()

	for i := 0; i < len(allusers.Userlist); i++ {

		if allusers.Userlist[i].Email == emailid && allusers.Userlist[i].Password == password {

			return []byte(allusers.Userlist[i].Email), nil
		}
	}
	return nil, nil
}

func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	err := stub.DelState(name) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the marble index
	userAsBytes, err := stub.GetState(userIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get array index")
	}
	var userIndex []string
	json.Unmarshal(userAsBytes, &userIndex) //un stringify it aka JSON.parse()

	//remove marble from index
	for i, val := range userIndex {
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name { //find the correct marble

			userIndex = append(userIndex[:i], userIndex[i+1:]...) //remove it
			for x := range userIndex {                            //debug prints...
				fmt.Println(string(x) + " - " + userIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(userIndex) //save new index
	err = stub.PutState(userIndexStr, jsonAsBytes)
	return nil, nil
}
func (t *SimpleChaincode) SaveSession(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	authsession := SessionAunthentication{}
	authsession.Token = args[0]
	authsession.Email = args[1]
	UserAsBytes, err := stub.GetState("savesession")
	if err != nil {
		return nil, errors.New("Failed to get users")
	}
	var session Session
	json.Unmarshal(UserAsBytes, &session) //un stringify it aka JSON.parse()

	session.StoreSession = append(session.StoreSession, authsession)
	fmt.Println("allsessions", session.StoreSession) //append to allusers
	fmt.Println("! appended user to allsessions")
	jsonAsBytes, _ := json.Marshal(session)
	fmt.Println("json", jsonAsBytes)
	err = stub.PutState("session", jsonAsBytes) //rewrite allusers
	if err != nil {
		return nil, err
	}
	fmt.Println("- end save session")
	return nil, nil
}
func (t *SimpleChaincode) SetUserForSession(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var token string
	var err error
	fmt.Println("running write()")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}
	token = args[0]

	UserAsBytes, err := stub.GetState("session")
	if err != nil {
		return nil, errors.New("failed to get sessions")
	}
	var session Session
	json.Unmarshal(UserAsBytes, &session)
	for i := 0; i < len(session.StoreSession); i++ {
		if session.StoreSession[i].Token == token {

			return []byte(session.StoreSession[i].Email), nil
		}
	}
	return nil, nil
}

func (t *SimpleChaincode) CreateCampaign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 11 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}

	//input sanitation
	fmt.Println("- start registration")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[9]) <= 0 {
		return nil, errors.New("6th argument must be a non-empty string")
	}
	if len(args[10]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}

	cuser := CreateCampaign{}

	cuser.Campaign_Id, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Failed to get loanamount as cannot convert it to int")
	}
	cuser.Campaign_Title = args[1]

	cuser.Campaign_Category = args[2]
	cuser.Campaign_Story = args[3]
	cuser.Campaign_Description = args[4]
	cuser.Campaign_Image = args[5]
	cuser.Campaign_Video = args[6]
	cuser.Campaign_Goal_amount, err = strconv.Atoi(args[7])
	if err != nil {
		return nil, errors.New("Failed to get loanamount as cannot convert it to int")
	}
	cuser.Campaign_Creation_Date = makeTimestamp()
	cuser.Project_Start_Date = args[8]
	cuser.Initial_amount, err = strconv.Atoi(args[9])
	if err != nil {
		return nil, errors.New("Failed to get interest as cannot convert it to int")
	}
	cuser.Campaign_End_Date = args[10]

	fmt.Println("cuser", cuser)

	UserAsBytes, err := stub.GetState("getcausers")
	if err != nil {
		return nil, errors.New("Failed to get users")
	}
	var campaignlist CampaignList
	json.Unmarshal(UserAsBytes, &campaignlist) //un stringify it aka JSON.parse()

	campaignlist.Campaignlist = append(campaignlist.Campaignlist, cuser)
	fmt.Println("campaignallusers", campaignlist.Campaignlist) //append to allusers
	fmt.Println("! appended cuser to campaignallusers")
	jsonAsBytes, _ := json.Marshal(campaignlist)
	fmt.Println("json", jsonAsBytes)
	err = stub.PutState("getcausers", jsonAsBytes) //rewrite allusers
	if err != nil {
		return nil, err
	}
	fmt.Println("- end campaignlist")
	return nil, nil
}

func (t *SimpleChaincode) AddReward(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation
	fmt.Println("- start addreward")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}

	reward := Rewards{}
	reward.Reward_Title = args[0]
	reward.Reward_Offer_Amount, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Failed to get id as cannot convert it to int")
	}

	reward.Reward_Description = args[2]
	Campaign_Id, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Failed to get id as cannot convert it to int")
	}

	UserAsBytes, err := stub.GetState("getcausers")
	if err != nil {
		return nil, errors.New("Failed to get users")
	}

	var campaignlist CampaignList
	json.Unmarshal(UserAsBytes, &campaignlist)

	for i := 0; i < len(campaignlist.Campaignlist); i++ {

		if campaignlist.Campaignlist[i].Campaign_Id == Campaign_Id {

			campaignlist.Campaignlist[i].Rewards = append(campaignlist.Campaignlist[i].Rewards, reward)

		}

		jsonAsBytes, _ := json.Marshal(campaignlist)
		fmt.Println("json", jsonAsBytes)
		err = stub.PutState("getcausers", jsonAsBytes) //rewrite allusers
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("- end addreward")
	return nil, nil
} //un stringify it aka JSON.parse()

func (t *SimpleChaincode) Contribute(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//input sanitation
	fmt.Println("- start contribute")

	contribute := Contribution{}
	contribute.UserId = args[0]
	contribute.AmountContributed, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Failed to get amountcontributed as cannot convert it to int")
	}
	contribute.TransactionStatus = args[2]
	contribute.Campaign_Id, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Failed to get campaign_id as cannot convert it to int")
	}

	UserAsBytes, err := stub.GetState("getcontribute")
	if err != nil {
		return nil, errors.New("Failed to get contribute")
	}
	var contributionlist ContributionList
	json.Unmarshal(UserAsBytes, &contributionlist) //un stringify it aka JSON.parse()

	contributionlist.Contributionlist = append(contributionlist.Contributionlist, contribute)
	fmt.Println("contributionlist", contributionlist.Contributionlist) //append to allusers
	fmt.Println("! appended contribute to contributionlist")
	jsonAsBytes, _ := json.Marshal(contributionlist)
	fmt.Println("json", jsonAsBytes)
	err = stub.PutState("getcontribute", jsonAsBytes) //rewrite allusers
	if err != nil {
		return nil, err
	}
	fmt.Println("- end Contributionlist")
	return nil, nil
}

func (t *SimpleChaincode) UpdateTxStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var TransactionStatus string

	//input sanitation
	fmt.Println("- start updatetxstatus")

	PaymentStatus := args[0]
	Userid := args[1]
	TransactionStatus = "payment clear"
	UserAsBytes, err := stub.GetState("getcontribute")
	if err != nil {
		return nil, errors.New("Failed to get contributelist")
	}

	var contributionlist ContributionList
	json.Unmarshal(UserAsBytes, &contributionlist)

	for i := 0; i < len(contributionlist.Contributionlist); i++ {
		if PaymentStatus == "Payment Successful" {
			if contributionlist.Contributionlist[i].TransactionStatus == "pending" && contributionlist.Contributionlist[i].UserId == Userid {
				contributionlist.Contributionlist[i].TransactionStatus = TransactionStatus
			}
		}

		jsonAsBytes, _ := json.Marshal(contributionlist)
		fmt.Println("json", jsonAsBytes)
		err = stub.PutState("getcontribute", jsonAsBytes) //rewrite allusers
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("- end updatepaymentstatus")
	return nil, nil
}

func makeTimestamp() string {
	t := time.Now()

	return t.Format(("2006-01-02T15:04:05.999999-07:00"))
	//return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
