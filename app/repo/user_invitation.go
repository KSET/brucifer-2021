package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
)

type userInvitationRepo struct {
}

type UserInvitationCreateOptions struct {
	Inviter models.User
	Comment string
}

func UserInvitation() userInvitationRepo {
	return userInvitationRepo{}
}

func (r userInvitationRepo) Create(data UserInvitationCreateOptions) (invitation models.UserInvitation, err error) {
	db := database.DatabaseProvider().Client()

	invitation = models.UserInvitation{
		CreatorID: data.Inviter.ID,
		Code:      uuid.New().String(),
		Comment:   data.Comment,
	}

	err = db.Save(&invitation).Error

	return
}

func (r userInvitationRepo) List() (invitations *[]models.UserInvitationUsed, err error) {
	db := database.DatabaseProvider().Client()

	invitations = &[]models.UserInvitationUsed{}
	err = db.Table("user_invitations").Model(models.UserInvitationUsed{}).Preload("UsedBy").Preload("Creator").Find(invitations).Error

	return
}

func (r userInvitationRepo) GetByCode(code string) (invitation *models.UserInvitation, err error) {
	db := database.DatabaseProvider().Client()

	invitation = &models.UserInvitation{}
	err = db.Where("code = ?", code).First(invitation).Error

	return
}

func (r userInvitationRepo) Invalidate(
	invitation *models.UserInvitation,
	claimer *models.User,
) (*models.UserInvitation, error) {
	db := database.DatabaseProvider().Client()

	invitation.UsedByID = claimer.ID

	err := db.Preload(clause.Associations).Save(invitation).Error

	return invitation, err
}
