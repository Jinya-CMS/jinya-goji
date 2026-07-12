package database

import "encoding/json"

type FormItem struct {
	Id            int      `db:"id,primarykey,autoincrement"`
	Position      int      `db:"position"`
	Type          string   `db:"type"`
	Options       []string `db:"options"`
	SpamFilter    []string `db:"spam_filter"`
	Label         string   `db:"label"`
	HelpText      string   `db:"help_text"`
	FormId        int      `db:"form_id"`
	IsFromAddress bool     `db:"is_from_address"`
	IsSubject     bool     `db:"is_subject"`
	IsRequired    bool     `db:"is_required"`
	Placeholder   *string  `db:"placeholder"`
}

func (item *FormItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id            int      `json:"id"`
		Position      int      `json:"position"`
		Type          string   `json:"type"`
		Options       []string `json:"options"`
		SpamFilter    []string `json:"spamFilter"`
		Label         string   `json:"label"`
		HelpText      string   `json:"helpText"`
		IsFromAddress bool     `json:"isFromAddress"`
		IsSubject     bool     `json:"isSubject"`
		IsRequired    bool     `json:"isRequired"`
		Placeholder   *string  `json:"placeholder"`
	}{
		Id:            item.Id,
		Position:      item.Position,
		Type:          item.Type,
		Options:       item.Options,
		SpamFilter:    item.SpamFilter,
		Label:         item.Label,
		HelpText:      item.HelpText,
		IsFromAddress: item.IsFromAddress,
		IsSubject:     item.IsSubject,
		IsRequired:    item.IsRequired,
		Placeholder:   item.Placeholder,
	})
}

func (item *FormItem) Create() error {
	return Insert(item)
}

func (item *FormItem) Update() error {
	return Update(item)
}

func (item *FormItem) Delete() error {
	return Delete(item)
}

func FindFormItemById(id int) (*FormItem, error) {
	return Get[FormItem](id)
}

func FindFormItemsByForm(formId int) ([]FormItem, error) {
	return Select[FormItem](`
SELECT
    id,
    position,
    type,
    options,
    spam_filter,
    label,
    help_text,
    form_id,
    is_from_address,
    is_subject,
    is_required,
    placeholder
FROM
    form_item
WHERE
    form_id = ?
ORDER BY
    position
`, formId)
}

func (item *FormItem) GetForm() (*Form, error) {
	return FindFormById(item.FormId)
}
