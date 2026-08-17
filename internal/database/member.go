package database

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const memberTableName = "members"

type Member struct {
	// The primary key ID
	// It is auto increment on the sqlite
	// You should not provide the value
	// ID int `json:"id"`

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
	RawPasswd string `json:"password" binding:"required"`

	// There is no need to provide the create_time sql will use current time
	// CreateTime time.Time `json:"create_time"`

	// WARN: On sql is INTERGER type so you must convert it as int before you insert it into database
	IsDelete bool `json:"is_delete" binding:"required"`
}

type MemberSQLProducer struct {
	// Hold the table name for query
	TableName string
}

func NewMemberDB() *MemberSQLProducer {
	return &MemberSQLProducer{
		memberTableName,
	}
}

func hashPasswd(rawPasswd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPasswd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPasswdAndHash : Check if the user's login attemp match the hash stored
// Return true if hash found on db
func CheckPasswdAndHash(passwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))

	return err == nil
}

func (mp *MemberSQLProducer) InsertMember(member Member) error {
	passwd, err := hashPasswd(member.RawPasswd)
	if err != nil {
		return fmt.Errorf("InsertMember: error when hash passwd: %w", err)
	}

	isDeleteInt := 0
	if member.IsDelete {
		isDeleteInt = 1
	}
	// Construct the sqlStr just for log record
	sqlStr := fmt.Sprintf(
		"INSERT INTO %s (power, nick, email, passwd, is_delete)\nVALUES(%d, %s, %s, %s, %d)",
		mp.TableName,
		member.Power,
		member.Nick,
		member.Email,
		passwd,
		isDeleteInt,
	)

	// Construct the execFunc
	execFunc := func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("INSERT INTO %s (power, nick, email, passwd, is_delete) VALUES (?, ?, ?, ?, ?)", mp.TableName)
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("error when exec PrepareContext: %w", err)
		}

		defer stmt.Close()

		_, err = stmt.ExecContext(
			ctx,
			member.Power,
			member.Nick,
			member.Email,
			passwd,
			isDeleteInt,
		)
		if err != nil {
			return fmt.Errorf("error when exec context: %w", err)
		}

		// NOTE: DO NOT call tx.Commit() here, let worker handler it

		return nil
	}

	// push new write task to chan
	sqlWriteTaskChan <- sqlWriteTask{
		execFunc: execFunc,
		sqlStr:   sqlStr,
	}

	return nil
}
