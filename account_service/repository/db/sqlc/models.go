// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type GptKeyType string

const (
	GptKeyTypeT3 GptKeyType = "t3"
	GptKeyTypeT4 GptKeyType = "t4"
)

func (e *GptKeyType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = GptKeyType(s)
	case string:
		*e = GptKeyType(s)
	default:
		return fmt.Errorf("unsupported scan type for GptKeyType: %T", src)
	}
	return nil
}

type NullGptKeyType struct {
	GptKeyType GptKeyType `json:"gpt_key_type"`
	Valid      bool       `json:"valid"` // Valid is true if GptKeyType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullGptKeyType) Scan(value interface{}) error {
	if value == nil {
		ns.GptKeyType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.GptKeyType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullGptKeyType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.GptKeyType), nil
}

type AccountKey struct {
	UserID    pgtype.UUID        `json:"user_id"`
	KeyID     pgtype.UUID        `json:"key_id"`
	ExpiredAt time.Time          `json:"expired_at"`
	CrDate    time.Time          `json:"cr_date"`
	UpDate    pgtype.Timestamptz `json:"up_date"`
}

type GptKey struct {
	KeyID            pgtype.UUID    `json:"key_id"`
	Type             GptKeyType     `json:"type"`
	ExpiredAt        time.Time      `json:"expired_at"`
	AssoicateAccount pgtype.UUID    `json:"assoicate_account"`
	MaxUsage         pgtype.Numeric `json:"max_usage"`
	CurrentUsage     pgtype.Numeric `json:"current_usage"`
	MaxShare         pgtype.Numeric `json:"max_share"`
	CurrentShare     pgtype.Numeric `json:"current_share"`
}

type Msg struct {
	MsgID        int64       `json:"msg_id"`
	MsgSessionID pgtype.UUID `json:"msg_session_id"`
	UserMsg      string      `json:"user_msg"`
	Response     pgtype.Text `json:"response"`
	CrDate       time.Time   `json:"cr_date"`
}

type MsgSession struct {
	MsgSessionID pgtype.UUID `json:"msg_session_id"`
	UserID       pgtype.UUID `json:"user_id"`
	CrDate       time.Time   `json:"cr_date"`
}

type Role struct {
	RoleID   pgtype.UUID        `json:"role_id"`
	RoleName string             `json:"role_name"`
	IsEnable bool               `json:"is_enable"`
	CrDate   time.Time          `json:"cr_date"`
	UpDate   pgtype.Timestamptz `json:"up_date"`
	CrUser   string             `json:"cr_user"`
	UpUser   pgtype.Text        `json:"up_user"`
}

type Session struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       pgtype.UUID        `json:"user_id"`
	RefreshToken string             `json:"refresh_token"`
	UserAgent    string             `json:"user_agent"`
	ClientIp     string             `json:"client_ip"`
	IsBlocked    bool               `json:"is_blocked"`
	CrDate       time.Time          `json:"cr_date"`
	ExpiredAt    pgtype.Timestamptz `json:"expired_at"`
}

type User struct {
	UserID            pgtype.UUID        `json:"user_id"`
	UserName          string             `json:"user_name"`
	Email             string             `json:"email"`
	IsEmailVerified   bool               `json:"is_email_verified"`
	HashedPassword    string             `json:"hashed_password"`
	PasswordChangedAt time.Time          `json:"password_changed_at"`
	SsoIdentifer      pgtype.Text        `json:"sso_identifer"`
	IsInternal        bool               `json:"is_internal"`
	CrDate            time.Time          `json:"cr_date"`
	UpDate            pgtype.Timestamptz `json:"up_date"`
	CrUser            string             `json:"cr_user"`
	UpUser            pgtype.Text        `json:"up_user"`
}

type UserRole struct {
	UserID pgtype.UUID        `json:"user_id"`
	RoleID pgtype.UUID        `json:"role_id"`
	CrDate time.Time          `json:"cr_date"`
	UpDate pgtype.Timestamptz `json:"up_date"`
	CrUser string             `json:"cr_user"`
	UpUser pgtype.Text        `json:"up_user"`
}

type UserRoleView struct {
	UserID            pgtype.UUID        `json:"user_id"`
	UserName          string             `json:"user_name"`
	Email             string             `json:"email"`
	IsEmailVerified   bool               `json:"is_email_verified"`
	HashedPassword    string             `json:"hashed_password"`
	PasswordChangedAt time.Time          `json:"password_changed_at"`
	SsoIdentifer      pgtype.Text        `json:"sso_identifer"`
	IsInternal        bool               `json:"is_internal"`
	CrDate            time.Time          `json:"cr_date"`
	UpDate            pgtype.Timestamptz `json:"up_date"`
	CrUser            string             `json:"cr_user"`
	UpUser            pgtype.Text        `json:"up_user"`
	RoleID            pgtype.UUID        `json:"role_id"`
	RoleName          pgtype.Text        `json:"role_name"`
}

type VertifyEmail struct {
	ID          int64              `json:"id"`
	UserID      pgtype.UUID        `json:"user_id"`
	Email       string             `json:"email"`
	SecretCode  string             `json:"secret_code"`
	IsUsed      bool               `json:"is_used"`
	IsValidated bool               `json:"is_validated"`
	CrDate      time.Time          `json:"cr_date"`
	UsedDate    pgtype.Timestamptz `json:"used_date"`
	ExpiredAt   time.Time          `json:"expired_at"`
}
