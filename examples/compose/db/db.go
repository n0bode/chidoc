package db

import "time"

type DB struct {
	lastID int64
	tmp    map[string]UserOrm
}

func NewDB() *DB {
	return &DB{
		tmp:    make(map[string]UserOrm),
		lastID: 1,
	}
}

func (db *DB) Filter(filter func(user UserOrm) (UserOrm, bool)) []UserOrm {
	var users []UserOrm
	pass := false

	for _, user := range db.tmp {
		if filter != nil {
			user, pass = filter(user)
			// ignore user
			if pass {
				continue
			}
		}
		users = append(users, user)
	}
	return users
}

func (db *DB) Update(filter func(user UserOrm) (UserOrm, bool)) (user UserOrm, updated bool) {
	for key, user := range db.tmp {
		found, updated := filter(user)
		db.tmp[key] = found
		// ignore user
		if updated {
			return found, updated
		}
	}
	return user, updated
}

func (db *DB) UserByID(id int64) (user UserOrm, exists bool) {
	for _, user := range db.tmp {
		if user.ID == id {
			return user, true
		}
	}
	return user, false
}

func (db *DB) AddUser(username, password, name, email string) UserOrm {
	now := time.Now().UTC()

	orm := UserOrm{
		ID:        db.lastID,
		Username:  username,
		Password:  password,
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdateAt:  now,
	}

	db.lastID++
	db.tmp[username] = orm
	return orm
}
