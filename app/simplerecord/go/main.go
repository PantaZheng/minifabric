package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	logger "github.com/sirupsen/logrus"
)

const (
	ccID          = "simplerecord"
	channelID     = "mychannel"
	orgName       = "org1.example.com"
	orgAdmin      = "Admin"
	publicTxName  = "AddPublicRecord"
	privateTxName = "AddPrivateRecord"
)

/*
To run this app, make sure that one of the wallet files such as Admin.id from
vars/profiles/vscode/wallets directory is copied onto ./wallets directory,
then this example code will use the wallet file and connection file to make
connections to Fabric network
*/

type Record struct {
	Timestamp   string  `json:"timestamp"`
	DeviceID    string  `json:"device_id"`
	Temperature float64 `json:"temperature"`
}

func addRecord(contract *gateway.Contract, txName string) {
	var seededRand = rand.New(rand.NewSource(time.Now().Unix()))
	var wg sync.WaitGroup
	startTime := time.Now()
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := contract.SubmitTransaction(txName,
				strconv.FormatInt(time.Now().UnixNano(), 10),
				strconv.FormatInt(int64(seededRand.Intn(20)), 10),
				strconv.FormatFloat(seededRand.Float64(), 'g', 4, 64)); err != nil {
				logger.Printf(" %v addRecord error: %v \n", txName, err.Error())
			}
		}()
	}
	wg.Wait()
	logger.Infof("%v  took is %v \n", txName, time.Now().Sub(startTime))
}

//func addPublicRecord(contract *gateway.Contract, record *Record) error {
//	_, err := contract.SubmitTransaction("AddPublicRecord", record.Timestamp, record.DeviceID, strconv.FormatFloat(record.Temperature, 'g', 4, 64))
//	if err != nil {
//		return fmt.Errorf("failed to submit transaction: %v", err.Error())
//	}
//	return nil
//}

func useWalletGateway() {
	wallet, err := gateway.NewFileSystemWallet("./profiles/wallets")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("Admin") {
		fmt.Println("Failed to get Admin from wallet")
		os.Exit(1)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("./profiles/connection.yaml")),
		gateway.WithIdentity(wallet, "Admin"),
	)
	if err != nil {
		logger.Errorf("Failed to connect: %v", err)
	}
	if gw == nil {
		logger.Error("Failed to create gateway")
	}

	network, err := gw.GetNetwork(channelID)
	if err != nil {
		logger.Errorf("Failed to get network: %v", err)
	}

	nowTime := time.Now().UnixNano()
	logger.Info("time is %v", nowTime)

	contract := network.GetContract(ccID)
	uuid.SetRand(nil)

	addRecord(contract, publicTxName)
	addRecord(contract, privateTxName)
}

func main() {
	file, err := os.OpenFile("golang.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	useWalletGateway()
}
