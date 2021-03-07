package main

import (
	"encoding/json"
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
	ccID      = "recordcc"
	channelID = "mychannel"
	orgName   = "org1.example.com"
	orgAdmin  = "Admin"
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

func addRecord(contract *gateway.Contract, record *Record) error {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}
	logger.Println("add Record %v", record)
	transient := make(map[string][]byte)
	transient["record"] = recordJSON
	tx, err := contract.CreateTransaction("AddRecord", gateway.WithTransient(transient))
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err.Error())
	}
	_, err = tx.Submit()
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err.Error())
	}
	return nil
}

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
	var seededRand = rand.New(rand.NewSource(nowTime))

	contract := network.GetContract(ccID)
	uuid.SetRand(nil)

	var wg sync.WaitGroup
	start := time.Now()
	recordTime := start.Unix()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		record := &Record{
			Timestamp:   strconv.FormatInt(recordTime, 10),
			DeviceID:    strconv.FormatInt(int64(seededRand.Intn(20)), 10),
			Temperature: seededRand.Float64(),
		}
		func() {
			defer wg.Done()
			if err := addRecord(contract, record); err != nil {
				logger.Println("addRecord error", err.Error())
			}
		}()
		recordTime += 1
	}
	wg.Wait()
	logger.Info("The time took is ", time.Now().Sub(start))
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
