package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	memberTableName = "members"
	logOnlyHint     = "[LOG ONLY]"
)

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

// NewMemberDB : return a new MemberSQLProducer with Default TableName [memberTableName]
func NewMemberDB() *MemberSQLProducer {
	return &MemberSQLProducer{
		memberTableName,
	}
}

// hashPasswd : the rawPasswd as input and output hashedpasswd which should be insert into database
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

// NOTE:
// Begin of the function for handling INSERT operation to database
// functions [InsertMember]

// InsertMember : insert the new member into the database
// it will not use member.ID, the ID should handle by sqlite itself (auto increment)
// As it's a write operation on database, so the request will handled by [RunDBWriteTasker]
// NOTE: sqlStr is just LOG ONLY, DON'T USE IT DIRECTLY
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
	// NOTE: LOG ONLY
	sqlStr := fmt.Sprintf(
		"%s INSERT INTO %s (power, nick, email, passwd, is_delete)\nVALUES(%d, %s, %s, %s, %d)",
		logOnlyHint,
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
		// WARN: DO NOT call tx.Commit() here, let worker handler it
		return nil
	}

	// push new write task to chan
	// exec on func [RunDBWriteTasker]
	sqlWriteTaskChan <- sqlWriteTask{
		execFunc: execFunc,
		sqlStr:   sqlStr,
	}

	return nil
}

// NOTE:
// Begin of the function for handling SELECT operation to database
// functions
// [ListMember]

// ListMember : exec for SELECT * (without passwd) from TableName
// return all member struct slice and error
// This function is read database so doesn't push the task to sqlWriteTaskChan
func (mp *MemberSQLProducer) ListMember() ([]Member, error) {
	queryStr := fmt.Sprintf("SELECT id, power, nick, email, is_delete FROM %s", mp.TableName)

	rows, err := globalDB.Query(queryStr)
	if err != nil {
		return nil, fmt.Errorf("ListMember: error when query: %w", err)
	}

	defer rows.Close()

	members := []Member{}
	// traverse the rows and push to members slice
	for rows.Next() {
		var id int
		var power int
		var nick string
		var email string
		var isDelete int

		err = rows.Scan(&id, &power, &nick, &email, &isDelete)
		if err != nil {
			utils.TextLogger.Error("error when scan member", "err", err, "id", id)
			continue
		}

		isDeleteBool := false
		if isDelete != 0 {
			isDeleteBool = true
		}
		members = append(members, Member{
			ID:        id,
			Power:     power,
			Nick:      nick,
			Email:     email,
			RawPasswd: "",
			IsDelete:  isDeleteBool,
		})
	}
	return members, nil
}

// TODO: function QeuryMember

// NOTE:
// Begin of the function for handling DELETE operation to database
// functions
// [DeleteMember]
// [hardDeleteMember]

// DeleteMember : delete the member by id, it will not throw error now
// It can be soft delete or hard delete
// See func [hardDeleteMember] and [softDeleteMember]
func (mp *MemberSQLProducer) DeleteMember(id int, softDelete bool) error {
	switch softDelete {
	case true:
		return mp.softDeleteMember(id)
	case false:
		return mp.hardDeleteMember(id)
	}

	return nil
}

// hardDeleteMember : hard delete remove the data from database file
func (mp *MemberSQLProducer) hardDeleteMember(id int) error {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE id = %d", mp.TableName, id)

	execFunc := func(ctx context.Context, tx *sql.Tx) error {
		query := sqlStr
		_, err := tx.Exec(query)
		if err != nil {
			return fmt.Errorf("error when hard delete member: id = %d", id)
		}

		return nil
	}

	// push new write task to chan
	// exec on func [RunDBWriteTasker]
	sqlWriteTaskChan <- sqlWriteTask{
		execFunc: execFunc,
		sqlStr:   sqlStr,
	}

	return nil
}

// NOTE:
// Begin of the function for handling UPDATE operation to database
// functions
// [softDeleteMember]

// softDeleteMember : soft delete just set the is_delete on database as 1
// It will not really delete the data
func (mp *MemberSQLProducer) softDeleteMember(id int) error {
	sqlStr := fmt.Sprintf("UPDATE %s SET is_delete = 1 WHERE id = %d", mp.TableName, id)

	execFunc := func(ctx context.Context, tx *sql.Tx) error {
		query := sqlStr
		_, err := tx.Exec(query)
		if err != nil {
			return fmt.Errorf("error when soft delete member: id = %d", id)
		}

		return nil
	}

	// push new write task to chan
	// exec on func [RunDBWriteTasker]
	sqlWriteTaskChan <- sqlWriteTask{
		execFunc: execFunc,
		sqlStr:   sqlStr,
	}

	return nil
}

// TODO: function UpdateMember
