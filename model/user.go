package model

import "fmt"

type User struct {
	ID         int64  `gorm:primary key;not_nil`
	Login      string `gorm:not_nil`
	Pass       string `gorm:not_nil`
	WorkNumber int32
}

func (u *User) Get(login string, pass string) error {
	query := fmt.Sprintf("login = '%s' and pass = '%s'", login, pass)
	return DBConn.Where(query).First(&u).Error
}

func (u *User) Save() error {
	return DBConn.Save(&u).Error
}

func (u *User) UpdatePass(login string, pass string, newPass string) error {
  query := fmt.Sprintf("login = '%s' and pass = '%s'", login, pass)
  return DBConn.Model(&u).Where(query).Update("pass", newPass).Error
}
