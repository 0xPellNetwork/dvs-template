package dvse2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TriggerAnvilNewBlocks triggers new blocks on the Anvil blockchain by transferring Ether from one account to another.
// It sends the specified number of transactions with a given private key and receiver address.
func (per *PellDVSE2ERunner) TriggerAnvilNewBlocks(times int, privateKey, receiverAddress string) error {
	for i := 0; i < times; i++ {
		err := per.transferEther(privateKey, receiverAddress)
		if err != nil {
			per.logger.Error("Failed to trigger new block", "err", err)
		}
	}
	return nil
}

// transferEther transfers Ether from senderPrivteKey to receiverAddress.
func (per *PellDVSE2ERunner) transferEther(senderPrivteKey, receiverAddress string) error {
	client, err := ethclient.Dial(per.ethRPCURL)
	if err != nil {
		return err
	}
	defer client.Close()

	// check if the client is connected
	currentBblockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	per.logger.Info("Current block number", "currentBblockNumber", currentBblockNumber)

	senderPrivteKey = strings.TrimPrefix(senderPrivteKey, "0x")
	privateKey, err := crypto.HexToECDSA(senderPrivteKey)
	if err != nil {
		return err
	}

	// get the public key from the private key
	publicKey := privateKey.Public()
	per.logger.Info("Public key", "publicKey", publicKey)

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	// get the address from the public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	per.logger.Info("Address", "address", address)

	// get the balance of the address
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return err
	}

	amount := big.NewInt(0).SetUint64(1000000000000000000) // 1 ETH in Wei
	per.logger.Info("Balance", "balance", balance)
	// check if the balance is greater than 1 ETH
	if balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	// get the nonce of the address
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return err
	}
	per.logger.Info("Nonce", "nonce", nonce)
	// get the gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	per.logger.Info("Gas price", "gasPrice", gasPrice)

	// get the gas limit
	gasLimit := uint64(21000) // in units
	per.logger.Info("Gas limit", "gasLimit", gasLimit)

	// get the chain id
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	per.logger.Info("Chain ID", "chainID", chainID)

	toAddress := common.HexToAddress(receiverAddress)
	per.logger.Info("Transfer amount", "amount", amount)

	// create a transaction
	transaction := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     []byte{},
	})

	// sign the transaction
	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	per.logger.Info("Signed transaction", "signedTx", signedTx)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	per.logger.Info("Transaction sent", "txHash", signedTx.Hash())
	return nil
}
