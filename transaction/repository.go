package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignId(campaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetByCampaignId(campaignID int) ([]Transaction, error) {
	var trans []Transaction
	err := r.db.Preload("User").Where("campaign_id = ? ", campaignID).Order("id desc").Find(&trans).Error
	if err != nil {
		return trans, err
	}
	return trans, nil
}

func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var trans []Transaction
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ? ", userID).Order("id desc").Find(&trans).Error
	if err != nil {
		return trans, err
	}
	return trans, nil
}
