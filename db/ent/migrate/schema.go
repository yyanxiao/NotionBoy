// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AccountsColumns holds the columns for the "accounts" table.
	AccountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted", Type: field.TypeBool, Default: false},
		{Name: "user_id", Type: field.TypeString},
		{Name: "user_type", Type: field.TypeEnum, Nullable: true, Enums: []string{"wechat", "telegram"}, Default: "wechat"},
		{Name: "database_id", Type: field.TypeString, Nullable: true},
		{Name: "access_token", Type: field.TypeString, Nullable: true},
		{Name: "notion_user_id", Type: field.TypeString, Nullable: true},
		{Name: "notion_user_email", Type: field.TypeString, Nullable: true},
		{Name: "is_latest_schema", Type: field.TypeBool, Default: false},
		{Name: "is_openai_api_user", Type: field.TypeBool, Default: false},
	}
	// AccountsTable holds the schema information for the "accounts" table.
	AccountsTable = &schema.Table{
		Name:       "accounts",
		Columns:    AccountsColumns,
		PrimaryKey: []*schema.Column{AccountsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "account_user_id_user_type",
				Unique:  true,
				Columns: []*schema.Column{AccountsColumns[4], AccountsColumns[5]},
			},
		},
	}
	// ChatHistoriesColumns holds the columns for the "chat_histories" table.
	ChatHistoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted", Type: field.TypeBool, Default: false},
		{Name: "user_id", Type: field.TypeInt},
		{Name: "conversation_idx", Type: field.TypeInt},
		{Name: "conversation_id", Type: field.TypeUUID},
		{Name: "message_id", Type: field.TypeString, Nullable: true},
		{Name: "message_idx", Type: field.TypeInt, Nullable: true},
		{Name: "request", Type: field.TypeString, Nullable: true},
		{Name: "response", Type: field.TypeString, Nullable: true},
	}
	// ChatHistoriesTable holds the schema information for the "chat_histories" table.
	ChatHistoriesTable = &schema.Table{
		Name:       "chat_histories",
		Columns:    ChatHistoriesColumns,
		PrimaryKey: []*schema.Column{ChatHistoriesColumns[0]},
	}
	// QuotaColumns holds the columns for the "quota" table.
	QuotaColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted", Type: field.TypeBool, Default: false},
		{Name: "user_id", Type: field.TypeInt},
		{Name: "category", Type: field.TypeEnum, Enums: []string{"chatgpt"}},
		{Name: "daily", Type: field.TypeInt, Nullable: true},
		{Name: "monthly", Type: field.TypeInt, Nullable: true},
		{Name: "yearly", Type: field.TypeInt, Nullable: true},
		{Name: "daily_used", Type: field.TypeInt, Nullable: true},
		{Name: "monthly_used", Type: field.TypeInt, Nullable: true},
		{Name: "yearly_used", Type: field.TypeInt, Nullable: true},
	}
	// QuotaTable holds the schema information for the "quota" table.
	QuotaTable = &schema.Table{
		Name:       "quota",
		Columns:    QuotaColumns,
		PrimaryKey: []*schema.Column{QuotaColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "quota_user_id_category",
				Unique:  true,
				Columns: []*schema.Column{QuotaColumns[4], QuotaColumns[5]},
			},
		},
	}
	// WechatSessionColumns holds the columns for the "wechat_session" table.
	WechatSessionColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted", Type: field.TypeBool, Default: false},
		{Name: "session", Type: field.TypeBytes},
		{Name: "dummy_user_id", Type: field.TypeString, Unique: true},
	}
	// WechatSessionTable holds the schema information for the "wechat_session" table.
	WechatSessionTable = &schema.Table{
		Name:       "wechat_session",
		Columns:    WechatSessionColumns,
		PrimaryKey: []*schema.Column{WechatSessionColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AccountsTable,
		ChatHistoriesTable,
		QuotaTable,
		WechatSessionTable,
	}
)

func init() {
	WechatSessionTable.Annotation = &entsql.Annotation{
		Table: "wechat_session",
	}
}
