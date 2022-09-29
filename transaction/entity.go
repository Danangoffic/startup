package transaction

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"time"

	"github.com/leekchan/accounting"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.User `gorm:"foreignKey:UserID"`
	Campaign   campaign.Campaign
}

func (t Transaction) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp ", Precision: 0, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(t.Amount)
}
