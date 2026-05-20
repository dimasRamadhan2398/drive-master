package models

// TableStatus represents the status of a database table
type TableStatus struct {
	TableName string `json:"tableName"`
	RowCount  int64  `json:"rowCount"`
	Exists    bool   `json:"exists"`
}