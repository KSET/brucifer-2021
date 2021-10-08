package models

type UserInvitation struct {
	Base
	CreatorID uint   `json:"-"`
	Creator   User   `json:"creator" gorm:"foreignKey:CreatorID"`
	UsedByID  uint   `json:"used" gorm:"index"`
	Comment   string `json:"comment"`
	Code      string `json:"code" gorm:"unique;size:127"`
}

type UserInvitationUsed struct {
	UserInvitation
	UsedBy *User `json:"usedUpBy"`
}
