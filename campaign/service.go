package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
	GetCampaignImageById(ID int) (CampaignImage, error)
	DeleteCampaignImage(input GetCampaignImageDetailInput) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userId int) ([]Campaign, error) {
	if userId != 0 {
		campaign, err := s.repository.FindByUserId(userId)
		if err != nil {
			return campaign, err
		}
		return campaign, err
	}

	campaign, err := s.repository.FindAll()
	if err != nil {
		return campaign, err
	}
	return campaign, err
}

func (s *service) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindById(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	campaign.UsersID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)
	// create a slug

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindById(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UsersID != inputData.User.ID {
		return campaign, errors.New("Not the campaign's owner!")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	// find campaign data by id from input campaign id
	campaign, err := s.repository.FindById(input.CampaignId)
	if err != nil {
		return CampaignImage{}, err
	}

	// validate is the user were a campaign owner or not
	if campaign.UsersID != input.User.ID {
		return CampaignImage{}, errors.New("Not the campaign's owner!")
	}

	// marking all images if input is primary is true
	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignId)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	// set campaignImage struct to pass to repository create image
	campaignImage := CampaignImage{}
	campaignImage.CampaignId = input.CampaignId
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}
	return newCampaignImage, nil
}

func (s *service) GetCampaignImageById(ID int) (CampaignImage, error) {
	campaignImage, err := s.repository.FindCampaignImageById(ID)
	if err != nil {
		return CampaignImage{}, err
	}
	return campaignImage, nil
}

func (s *service) DeleteCampaignImage(input GetCampaignImageDetailInput) (bool, error) {

	campaignImage, err := s.repository.FindCampaignImageById(input.ID)
	if err != nil {
		return false, err
	}

	campaign, _ := s.repository.FindById(campaignImage.CampaignId)

	_, err = s.repository.DeleteCampaignImage(campaignImage)
	if err != nil {
		return false, err
	}

	isDeletedPrimary := campaignImage.IsPrimary

	if isDeletedPrimary == 1 {
		var i = 0
		for _, v := range campaign.CampaignImages {
			if i == 0 {
				_, err2 := s.repository.MarkAllImagesAsNonPrimary(v.CampaignId)
				if err2 != nil {
					return false, err2
				}
			}
			i++
		}
	}

	return true, nil
}
