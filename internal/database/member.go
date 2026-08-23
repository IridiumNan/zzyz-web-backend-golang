package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

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
	// use bcrypt to hash raw password
	RawPasswd string `json:"passwd" binding:"required"`

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
func checkPasswdAndHash(passwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))

	return err == nil
}

// parseMembersFromRows Convert rows into member slices
// Below attributes will be select
// id
// power
// nick
// email
// is_delete
func parseMembersFromRows(rows *sql.Rows) (members []Member, err error) {
	members = []Member{}
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

// InsertMember : insert the new member into the database
// it will not use member.ID, the ID should handle by sqlite itself (auto increment)
// As it's a write operation on database, so the request will handled by [RunDBWriteTasker]
func (mp *MemberSQLProducer) InsertMember(member Member) error {
	passwd, err := hashPasswd(member.RawPasswd)
	if err != nil {
		return fmt.Errorf("InsertMember: error when hash passwd: %w", err)
	}

	isDeleteInt := 0
	if member.IsDelete {
		isDeleteInt = 1
	}
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

	pushWriteTaskToChan(execFunc, sqlStr)

	return nil
}

// ListMember : exec for SELECT * (without passwd) from TableName
// return all member struct slice and error
// This function is read database so doesn't push the task to sqlWriteTaskChan
// If the len(members) == 0. it will still return, you should check
func (mp *MemberSQLProducer) ListMember() ([]Member, error) {
	queryStr := fmt.Sprintf("SELECT id, power, nick, email, is_delete FROM %s", mp.TableName)

	rows, err := globalDB.Query(queryStr)
	if err != nil {
		return nil, fmt.Errorf("ListMember: error when query: %w", err)
	}

	defer rows.Close()

	members, err := parseMembersFromRows(rows)

	return members, err
}

// QueryMember The function will query with sqlStr as below
// SELECT id, power, nick, email, is_delete FROM TableName WHERE attribute = value
// Ensure you convert all type as the data on sql before calling this func
// If the length of members is 0, it will still return, remember to check it
// the param like if True, will query with WHERE attribute like value
func (mp *MemberSQLProducer) QueryMember(attribute string, value string, like bool) (members []Member, err error) {
	operator := "="
	if like {
		operator = "like"
	}

	sqlStr := fmt.Sprintf("SELECT id, power, nick, email, is_delete FROM %s WHERE %s %s %s", mp.TableName, attribute, operator, value)

	utils.TextLogger.Info("exec sql query", "sqlStr", sqlStr)

	query := fmt.Sprintf("SELECT id, power, nick, email, is_delete FROM %s WHERE %s %s ?", mp.TableName, attribute, operator)

	var rows *sql.Rows
	switch attribute {
	case "id", "power", "is_delete":
		var valueInt int
		valueInt, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("error when convert into integer type %s, err: %w", value, err)
		}
		rows, err = globalDB.Query(query, valueInt)
	default:

		rows, err = globalDB.Query(query, value)
	}
	if err != nil {
		return nil, fmt.Errorf("error when Query Member: %w", err)
	}

	defer rows.Close()

	members, err = parseMembersFromRows(rows)

	return
}

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
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("error when hard delete member: id = %d", id)
		}

		return nil
	}

	pushWriteTaskToChan(execFunc, sqlStr)

	return nil
}

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

	pushWriteTaskToChan(execFunc, sqlStr)

	return nil
}

// UpdateMember The function wil update members table as sqlStr below
// UPDATE TableName SET attribute = value WHERE id = ?
// If the attribute = passwd, it will hash it first then udpate
func (mp *MemberSQLProducer) UpdateMember(id int, attribute string, value string) (err error) {
	// If attribute is passwd, you should hash it first
	if attribute == "passwd" {
		value, err = hashPasswd(value)
		if err != nil {
			return fmt.Errorf("error when hash the passwd: %w", err)
		}
		utils.TextLogger.Warn("updating passwd", "hashed_passwd", value, "id", id)
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s = %s WHERE id = %d", mp.TableName, attribute, value, id)

	execFunc := func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", mp.TableName, attribute)
		_, err := tx.ExecContext(ctx, query, value, id)
		if err != nil {
			return fmt.Errorf("error when update member: sqlStr: %s, err: %w", sqlStr, err)
		}

		return nil
	}

	pushWriteTaskToChan(execFunc, sqlStr)

	return nil
}

// FetchMemberPower Return the power of this
// return -1 if the passwd is wrong else return exact power
func (mp *MemberSQLProducer) FetchMemberPower(nick string, rawPasswd string) (power int, err error) {
	errPower := -1
	sqlStr := fmt.Sprintf("SELECT power, passwd FROM %s WHERE nick = ?", mp.TableName)

	utils.TextLogger.Info(fmt.Sprintf("exec sql query: %s", sqlStr), "nick", nick)
	rows, err := globalDB.Query(sqlStr, nick)
	if err != nil {
		return errPower, fmt.Errorf("error when get passwd of %s, err: %s", nick, err)
	}

	// NOTE: even the nick is unique on this table
	// Use slices for get all passwd

	type AuthPower struct {
		Power  int
		Passwd string
	}
	var authList []AuthPower
	for rows.Next() {
		var power int
		var passwd string

		err = rows.Scan(&power, &passwd)
		if err != nil {
			utils.TextLogger.Error("error when scan passwd for member", "nick", nick, "err", err)
			continue
		}

		authList = append(authList, AuthPower{
			Power:  power,
			Passwd: passwd,
		})
	}

	if len(authList) != 1 {
		utils.TextLogger.Warn("the passwd for a specific member is more than one", "nick", nick)
	}

	// Return the first match power else errPower
	for idx := range authList {
		if checkPasswdAndHash(rawPasswd, authList[idx].Passwd) {
			return authList[idx].Power, nil
		}
	}

	return errPower, fmt.Errorf("the passwd for nick %s not match", nick)
}
