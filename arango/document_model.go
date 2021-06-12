package arango

import "time"

type DocumentModel struct {
	ArangoInterface `json:"-"`
	Id              string    `json:"_id,omitempty"`
	Key             string    `json:"_key,omitempty"`
	Rev             string    `json:"_rev,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type EdgeModel struct {
	DocumentModel
	From string `json:"_from"`
	To   string `json:"_to"`
}

type ArangoInterface interface {
	GetId() string
	GetKey() string
	Set(Id, Key, Rev string)
	InitializeTimestamp()
	UpdateTimestamp()
}

func (d *DocumentModel) GetId() string {
	return d.Id
}

func (d *DocumentModel) GetKey() string {
	return d.Key
}

func (d *DocumentModel) InitializeTimestamp() {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
}

func (d *DocumentModel) Set(Id, Key, Rev string) {
	d.Id = Id
	d.Key = Key
	d.Rev = Rev
}

func (d *DocumentModel) UpdateTimestamp() {
	d.UpdatedAt = time.Now()
}
