package database

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Node 代表一个客户端节点
type Node struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Hostname  string    `json:"hostname"`
	Version   string    `json:"version"`
	IP        string    `json:"ip"`
	Status    string    `json:"status"` // online, offline
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}

// Task 代表一个异步任务或操作记录
type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type"` // search, download, play
	NodeID    string    `json:"node_id" gorm:"index"`
	Payload   string    `json:"payload"` // JSON string
	Result    string    `json:"result"`  // JSON string
	Status    string    `json:"status"`  // pending, success, failed
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Setting 系统设置
type Setting struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
}

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(path string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移模式
	return DB.AutoMigrate(&Node{}, &Task{}, &Setting{})
}

// Node Operations
func UpsertNode(node *Node) error {
	return DB.Save(node).Error
}

func UpdateNodeStatus(id string, status string) error {
	return DB.Model(&Node{}).Where("id = ?", id).Update("status", status).Error
}

func GetNodes() ([]Node, error) {
	var nodes []Node
	err := DB.Order("last_seen desc").Find(&nodes).Error
	return nodes, err
}

func GetNode(id string) (*Node, error) {
	var node Node
	err := DB.First(&node, "id = ?", id).Error
	return &node, err
}

// Task Operations
func CreateTask(task *Task) error {
	return DB.Create(task).Error
}

func UpdateTaskResult(id uint, status string, result string, errStr string) error {
	updates := map[string]interface{}{
		"status": status,
		"result": result,
		"error":  errStr,
	}
	return DB.Model(&Task{}).Where("id = ?", id).Updates(updates).Error
}

func GetTasks(nodeID string, offset, limit int) ([]Task, int64, error) {
	var tasks []Task
	var count int64

	query := DB.Model(&Task{})
	if nodeID != "" {
		query = query.Where("node_id = ?", nodeID)
	}

	query.Count(&count)
	err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&tasks).Error
	return tasks, count, err
}
