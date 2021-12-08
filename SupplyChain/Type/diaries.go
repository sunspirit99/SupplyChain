package Type

import (
	model "SuperBank/Model"
	"SuperBank/database"

	"gorm.io/gorm"
)

type DiariesORM struct {
	model.Transaction
	db *gorm.DB
}

func NewDiariesORM() *DiariesORM {
	return &DiariesORM{db: database.Connector}
}

func (dr DiariesORM) CreateDiary(action model.Diaries) (model.Diaries, error) {

	err := dr.db.Create(&action).Error
	return action, err
}
