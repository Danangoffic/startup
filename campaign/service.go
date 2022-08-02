package campaign

type Service interface {
	GetCampaigns(userId int) ([]Campaign, error)
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
