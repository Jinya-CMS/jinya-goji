package database

type BlogCategory struct {
	Id             int     `db:"id,primarykey,autoincrement" json:"id"`
	Name           string  `db:"name" json:"name"`
	Description    *string `db:"description" json:"description"`
	ParentId       *int    `db:"parent_id" json:"-"`
	WebhookEnabled bool    `db:"webhook_enabled" json:"webhookEnabled"`
	WebhookUrl     *string `db:"webhook_url" json:"webhookUrl"`
}
