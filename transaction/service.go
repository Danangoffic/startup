package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

type Service interface {
	GetCampaignTransactions(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetUserTransactions(userId int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

func (s *service) GetCampaignTransactions(input GetCampaignTransactionsInput) ([]Transaction, error) {

	// get campaign
	// check campaign userId != user_id yg request
	campaign, err := s.campaignRepository.FindById(input.ID)

	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UsersID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the Campaign")
	}

	transactions, err := s.repository.FindByCampaignId(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	return transactions, nil
}

func (s *service) GetUserTransactions(userId int) ([]Transaction, error) {
	transactions, err := s.repository.FindByUserId(userId)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignId
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	// transaction.Code = ""

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{}
	paymentTransaction.Amount = newTransaction.Amount
	paymentTransaction.ID = newTransaction.ID

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
