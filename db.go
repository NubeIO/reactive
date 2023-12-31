package reactive

import (
	"gorm.io/gorm"
)

func (n *BaseNode) AddDB(db *gorm.DB) {
	n.db = db
}
func (n *BaseNode) GetDB() *gorm.DB {
	return n.db
}
