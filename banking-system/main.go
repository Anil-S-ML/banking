package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Account struct {
	Number       string
	Balance      float64
	Transactions []Transaction
}

type Transaction struct {
	Type       string
	Amount     float64
	NewBalance float64
	Timestamp  time.Time
}

type BankingService struct {
	accounts map[string]*Account
}

func NewBankingService() *BankingService {
	bs := &BankingService{
		accounts: make(map[string]*Account),
	}
	bs.loadAccounts()
	return bs
}

func (bs *BankingService) loadAccounts() {
	file, err := os.ReadFile("accounts.json")
	if err == nil {
		json.Unmarshal(file, &bs.accounts)
	}
}

func (bs *BankingService) saveAccounts() {
	data, _ := json.MarshalIndent(bs.accounts, "", "  ")
	_ = os.WriteFile("accounts.json", data, 0644)
}

func (bs *BankingService) CreateAccount() string {
	accountNumber := generateAccountNumber()
	bs.accounts[accountNumber] = &Account{
		Number:       accountNumber,
		Balance:      0,
		Transactions: make([]Transaction, 0),
	}
	bs.saveAccounts()
	return accountNumber
}

func (bs *BankingService) Deposit(accountNumber string, amount float64) (float64, error) {
	account, exists := bs.accounts[accountNumber]
	if !exists {
		return 0, errors.New("account not found")
	}

	account.Balance += amount
	transaction := Transaction{
		Type:       "CREDIT",
		Amount:     amount,
		NewBalance: account.Balance,
		Timestamp:  time.Now(),
	}
	account.Transactions = append(account.Transactions, transaction)
	bs.saveAccounts()
	return account.Balance, nil
}

func (bs *BankingService) Withdraw(accountNumber string, amount float64) (float64, error) {
	account, exists := bs.accounts[accountNumber]
	if !exists {
		return 0, errors.New("account not found")
	}

	if account.Balance < amount {
		return 0, errors.New("insufficient funds")
	}

	account.Balance -= amount
	transaction := Transaction{
		Type:       "DEBIT",
		Amount:     amount,
		NewBalance: account.Balance,
		Timestamp:  time.Now(),
	}
	account.Transactions = append(account.Transactions, transaction)
	bs.saveAccounts()
	return account.Balance, nil
}

func (bs *BankingService) CheckBalance(accountNumber string) (float64, error) {
	account, exists := bs.accounts[accountNumber]
	if !exists {
		return 0, errors.New("account not found")
	}
	return account.Balance, nil
}

func (bs *BankingService) ViewTransactions(accountNumber string) ([]Transaction, error) {
	account, exists := bs.accounts[accountNumber]
	if !exists {
		return nil, errors.New("account not found")
	}
	return account.Transactions, nil
}

func loadAccountCounter() int {
	data, err := os.ReadFile("counter.txt")
	if err != nil {
		return 0
	}
	num, err := strconv.Atoi(string(data))
	if err != nil {
		return 0
	}
	return num
}

func saveAccountCounter(counter int) {
	_ = os.WriteFile("counter.txt", []byte(strconv.Itoa(counter)), 0644)
}

var accountCounter int = loadAccountCounter()

func generateAccountNumber() string {
	accountCounter++
	saveAccountCounter(accountCounter)
	return fmt.Sprintf("%d", accountCounter)
}

var service = NewBankingService()

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./app [command] [arguments]")
		return
	}

	command := os.Args[1]

	switch command {
	case "create-account":
		accountNumber := service.CreateAccount()
		fmt.Println(accountNumber)

	case "deposit":
		if len(os.Args) != 4 {
			fmt.Println("Usage: ./bank deposit [account-number] [amount]")
			return
		}
		amount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Invalid amount")
			return
		}
		balance, err := service.Deposit(os.Args[2], amount)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%.2f\n", balance)

	case "withdraw":
		if len(os.Args) != 4 {
			fmt.Println("Usage: ./bank withdraw [account-number] [amount]")
			return
		}
		amount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Invalid amount")
			return
		}
		balance, err := service.Withdraw(os.Args[2], amount)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%.2f\n", balance)

	case "check-balance":
		if len(os.Args) != 3 {
			fmt.Println("Usage: ./bank check-balance [account-number]")
			return
		}
		balance, err := service.CheckBalance(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%.2f\n", balance)

	case "view-transactions":
		if len(os.Args) != 3 {
			fmt.Println("Usage: ./bank view-transactions [account-number]")
			return
		}
		transactions, err := service.ViewTransactions(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, t := range transactions {
			fmt.Printf("%s;%.2f;%.2f\n", t.Type, t.Amount, t.NewBalance)
		}

	default:
		fmt.Println("Unknown command")
	}
}
