package db

// APIToken in database.
type APIToken struct {
	ID        int64  `json:"id,omitempty" db:"id,type=INT,auto_increment,primary,key"`
	Type      int32  `json:"type,omitempty" db:"type,type=INT"`
	Token     string `json:"token,omitempty" db:"token,type=INT"`
	UserEmail string `json:"user_email,omitempty" db:"user_email,type=VARCHAR(128)"`

	Name              string `json:"name,omitempty" db:"name,type=VARCHAR(64)"`
	Message           string `json:"message,omitempty" db:"omitempty,type=CARCHAR(128)"`
	DeadlineTimestamp int64  `json:"deadline_timestamp,omitempty" db:"deadline_timestamp,type=INT(64)"`
	// BitMap.
	AllowedActions int32 `json:"allowed_actions,omitempty" db:"allowed_actions,type=INT(32)"`
}

// APITokenRawModel return APIToken's rawmodel for sqlm.
func APITokenRawModel() interface{} {
	return &APIToken{}
}
