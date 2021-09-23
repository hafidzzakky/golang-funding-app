package transaction

import (
	"time"
)

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	var transactionFormatter []CampaignTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}

type UserTransactionFormatter struct {
	ID        int                              `json:"id"`
	Amount    int                              `json:"amount"`
	Status    string                           `json:"status"`
	CreatedAt time.Time                        `json:"created_at"`
	Campaign  CampaignUserTransactionFormatter `json:"campaign"`
}

type CampaignUserTransactionFormatter struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	trans := UserTransactionFormatter{}
	trans.ID = transaction.ID
	trans.Amount = transaction.Amount
	trans.Status = transaction.Status
	trans.CreatedAt = transaction.CreatedAt

	// Campaign
	campaignFormatter := CampaignUserTransactionFormatter{}
	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageUrl = ""
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageUrl = transaction.Campaign.CampaignImages[0].FileName
	}
	trans.Campaign = campaignFormatter
	return trans
}

func FormatUserCampaignTransactions(transactions []Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	var userTransactionFormatter []UserTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		userTransactionFormatter = append(userTransactionFormatter, formatter)
	}

	return userTransactionFormatter
}
