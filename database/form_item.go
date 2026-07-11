package database

type FormItem struct {
	Id            int      `db:"id,primarykey,autoincrement" json:"id"`
	Position      int      `db:"position" json:"position"`
	Type          string   `db:"type" json:"type"`
	Options       []string `db:"options" json:"options"`
	SpamFilter    []string `db:"spam_filter" json:"spamFilter"`
	Label         string   `db:"label" json:"label"`
	HelpText      string   `db:"help_text" json:"helpText"`
	FormId        int      `db:"form_id" json:"-"`
	IsFromAddress bool     `db:"is_from_address" json:"isFromAddress"`
	IsSubject     bool     `db:"is_subject" json:"isSubject"`
	IsRequired    bool     `db:"is_required" json:"isRequired"`
	Placeholder   *string  `db:"placeholder" json:"placeholder"`
}
