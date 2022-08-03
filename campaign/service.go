package campaign

type Service interface {
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
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
