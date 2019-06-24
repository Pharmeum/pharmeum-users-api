package db

type User struct {
	ID          uint64 `db:"id"`
	Email       string `db:"pk,email"`
	Password    string `db:"password"`
	Name        string `db:"name"`
	Phone       string `db:"phone"`
	DateOfBirth string `db:"date_of_birth"`
}

func (u User) TableName() string {
	return "users"
}

func (d *DB) GetUser(email string) (*User, error) {
	var user User
	err := d.db.Select().Model(email, &user)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
	}

	return &user, err
}

func (d *DB) CreateUser(user *User) error {
	return d.db.Model(user).Insert()
}

func (d *DB) SetUserNewPassword(email string, password string) error {
	user := &User{Password: password}
	return d.db.Model(&user).Update("password")
}
