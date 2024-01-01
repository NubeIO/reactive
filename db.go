package reactive

import (
	"gorm.io/gorm"
)

func (n *BaseNode) AddDB(db *gorm.DB) {
	if db == nil {
		panic("AddBB() db is empty")
	}
	n.db = db
}
func (n *BaseNode) GetDB() *gorm.DB {
	if n.db == nil {
		panic("GetDB() node db has not been added please add")
	}
	return n.db
}
