package models

type Member struct {
	// The primary key ID
	// It is auto increment on the sqlite
	// You should not provide the value
	ID int `json:"id"`

	// NOTE: The power of this member
	// for now 1 has all power and 0 has no power
	// use the binary 000, 001
	Power int `json:"power" binding:"required"`

	// The name for member assigne
	Nick string `json:"nick" binding:"required"`

	// NOTE: It's optional
	Email string `json:"email"`

	// NOTE: This is original passwd and system will not hold it, sql store the hash value of this
	// RawPasswd is required
	// use bcrypt to hash raw password
	RawPasswd string `json:"passwd" binding:"required"`

	// There is no need to provide the create_time sql will use current time
	// CreateTime time.Time `json:"create_time"`

	// WARN: On sql is INTERGER type so you must convert it as int before you insert it into database
	IsDelete bool `json:"is_delete" binding:"required"`
}
