package ComprehensiveSmartContract

import (
	"encoding/json"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	UserId       string `json:"userId"`
	DocumentId   string `json:"documentId"`
	DocumentData string `json:"documentData"`
}

func (s *SmartContract) CreateAssert(ctx contractapi.TransactionContextInterface, userId string, documentId string, documentData string) string {
	newAsset := Asset{
		UserId:       userId,
		DocumentId:   documentId,
		DocumentData: documentData,
	}

	assetJson, err := json.Marshal(newAsset)
	if err != nil {
		errorString := "Error Parsing asset..."
		log.Fatal(errorString)
		return errorString
	}

	err = ctx.GetStub().PutState(documentId, assetJson)
	if err != nil {
		errorString := "Error storing asset onto chain..."
		log.Fatal(errorString)
		return errorString
	}

	return "Successfully stored data onto chain..."
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, documentId string) (*Asset, string) {
	assetJson, err := ctx.GetStub().GetState(documentId)
	if err != nil {
		errorString := "Error fetching asset from chain..."
		log.Fatal(errorString)
		return nil, errorString
	}
	if assetJson != nil {
		errorString := "Asset with Id: " + documentId + "Not found..."
		log.Fatal(errorString)
		return nil, errorString
	}

	var asset Asset
	err = json.Unmarshal(assetJson, &asset)
	if err != nil {
		errorString := "Error Parsing asset..."
		log.Fatal(errorString)
		return nil, errorString
	}

	return &asset, "Successfully retrived asset from chain."
}

func (s *SmartContract) RealAllAsset(ctx contractapi.TransactionContextInterface) ([]*Asset, string) {
	resultsIterator, err := ctx.GetStateByRange("", "")
	if err != nil {
		errorString := "Error fetching asset from chain..."
		log.Fatal(errorString)
		return nil, errorString
	}
	defer resultsIterator.close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, "Successfully retrived asset from chain."
}

