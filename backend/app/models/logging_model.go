package models

type Logging struct {
	ID            int          `db:"id" json:"id"`
	CreatedAt     string       `db:"created_at" json:"created_at"`
	User_id       int          `db:"user_id" json:"user_id"`
	Slug          string       `db:"slug" json:"slug"`
	Action_status int          `db:"action_status" json:"action_status"`
	Logging_attrs LoggingAttrs `db:"logging_attrs" json:"logging_attrs"`
}

type LoggingAttrs struct {
	Ip_address string `db:"ip_address" json:"ip_address"`
	Device     string `db:"device" json:"device"`
	User_agent string `db:"user_agent" json:"user_agent"`
}
