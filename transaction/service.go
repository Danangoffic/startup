package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

type Service interface {
	GetCampaignTransactions(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetUserTransactions(userId int) ([]Transaction, error)
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
