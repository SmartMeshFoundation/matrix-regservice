package models

//User register user for this service
type User struct {
	LocalPart   string `gorm:"primary_key"` //@someone:matrix.org someone is localpoart,matrix.org is domain
	DisplayName string
	Password    string
}

//IsUserAlreadyExists return true when this user already registered
func IsUserAlreadyExists(localPart string) bool {
	var u User
	u.LocalPart = localPart
	err := db.Where(&u).Find(&u).Error
	return err == nil
}

//NewUser create a new user in db
func NewUser(localPart, displayName, password string) (err error) {
	u := &User{
		LocalPart:   localPart,
		DisplayName: displayName,
		Password:    password,
	}
	err = db.Save(u).Error
	return
}
