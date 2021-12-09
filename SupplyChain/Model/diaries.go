package Model

import (
	"gorm.io/gorm"
)

type DiariesORM struct {
	Diaries
	db *gorm.DB
}

func NewDiariesORM() *DiariesORM {
	return &DiariesORM{db: Connector}
}

func (dr DiariesORM) CreateDiary(action Diaries) (Diaries, error) {

	err := dr.db.Create(&action).Error
	return action, err
}

func (dr DiariesORM) GetDiaryById(id string) (Diaries, error) {
	var diary Diaries

	err := dr.db.Raw("SELECT * from diaries where id = ?", id).Scan(&diary).Error
	return diary, err
}

func (dr DiariesORM) GetDiaryByProduct(productID string) ([]Diaries, error) {
	var diary []Diaries

	err := dr.db.Raw("SELECT * from diaries where related_products = ?", productID).Scan(&diary).Error
	return diary, err
}
