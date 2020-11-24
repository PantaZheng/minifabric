module github.com/fabricgoapp

replace github.com/fabricgoapp => ../

require (
	github.com/google/uuid v1.1.2
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta3
	github.com/sirupsen/logrus v1.3.0
)

go 1.13
