package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// jsonMarshalWithTime 辅助函数，添加格式化的创建时间
func jsonMarshalWithTime(v interface{}, createdAt time.Time) ([]byte, error) {
	type TimeWrapper struct {
		CreatedAtStr string `json:"created_at"`
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	if !createdAt.IsZero() {
		m["created_at"] = createdAt.Format("2006-01-02 15:04:05")
	}

	return json.Marshal(m)
}

// 国家标准政策表
type NationalStandard struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Link1          string         `gorm:"size:500" json:"link1"`
	Link2          string         `gorm:"size:500" json:"link2"`
	StandardNo     string         `gorm:"size:100;index" json:"standard_no"`
	StandardName   string         `gorm:"size:500;index" json:"standard_name"`
	EnglishName    string         `gorm:"size:500" json:"english_name"`
	PublishDate    string         `gorm:"size:20" json:"publish_date"`
	ImplementDate  string         `gorm:"size:20" json:"implement_date"`
	Status         string         `gorm:"size:20" json:"status"`
	Nature         string         `gorm:"size:20" json:"nature"`
	Category       string         `gorm:"size:20" json:"category"`
	IsAdopted      string         `gorm:"size:10" json:"is_adopted"`
	CCSCode        string         `gorm:"size:20" json:"ccs_code"`
	ICSCode        string         `gorm:"size:50" json:"ics_code"`
	Department     string         `gorm:"size:100" json:"department"`
	TechnicalDept  string         `gorm:"size:100" json:"technical_dept"`
	Publisher      string         `gorm:"size:200" json:"publisher"`
	Description    string         `gorm:"type:text" json:"description"`
	StandardType   string         `gorm:"size:20;default:'国标'" json:"standard_type"`
	DownloadURL    string         `gorm:"size:500" json:"download_url"`
	CheckStatus    int            `gorm:"default:0" json:"check_status"` // 0待核验 1已核验 2不通过
	CategoryID     uint           `gorm:"index" json:"category_id"`
	Source         string         `gorm:"size:20;default:'manual'" json:"source"` // crawl/爬取 manual/人工
	CreatedAt      time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// MarshalJSON 自定义JSON输出，格式化创建时间
func (n NationalStandard) MarshalJSON() ([]byte, error) {
	type Alias NationalStandard
	return jsonMarshalWithTime(Alias(n), n.CreatedAt)
}

func (NationalStandard) TableName() string {
	return "national_standard"
}

// 行业标准政策表
type IndustryStandard struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	DetailLink      string         `gorm:"size:500" json:"detail_link"`
	StandardNo      string         `gorm:"size:100;uniqueIndex" json:"standard_no"`
	StandardName    string         `gorm:"size:500" json:"standard_name"`
	PublishDate     string         `gorm:"size:20" json:"publish_date"`
	ImplementDate   string         `gorm:"size:20" json:"implement_date"`
	RevisionType    string         `gorm:"size:20" json:"revision_type"`
	CCSCode         string         `gorm:"size:20" json:"ccs_code"`
	ICSCode         string         `gorm:"size:50" json:"ics_code"`
	TechnicalOwner  string         `gorm:"size:200" json:"technical_owner"`
	ApproveDept     string         `gorm:"size:100" json:"approve_dept"`
	IndustryClass   string         `gorm:"size:100" json:"industry_class"`
	StandardClass   string         `gorm:"size:50" json:"standard_class"`
	Status          string         `gorm:"size:20" json:"status"`
	ReplaceStandard string         `gorm:"size:200" json:"replace_standard"`
	StandardType    string         `gorm:"size:20;default:'行标'" json:"standard_type"`
	DownloadURL     string         `gorm:"size:500" json:"download_url"`
	CheckStatus     int            `gorm:"default:0" json:"check_status"`
	CategoryID      uint           `gorm:"index" json:"category_id"`
	Source          string         `gorm:"size:20;default:'manual'" json:"source"` // crawl/爬取 manual/人工
	CreatedAt       time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i IndustryStandard) MarshalJSON() ([]byte, error) {
	type Alias IndustryStandard
	return jsonMarshalWithTime(Alias(i), i.CreatedAt)
}

func (IndustryStandard) TableName() string {
	return "industry_standard"
}

// 地方标准政策表
type LocalStandard struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DetailLink    string         `gorm:"size:500" json:"detail_link"`
	StandardNo    string         `gorm:"size:100;uniqueIndex" json:"standard_no"`
	StandardName  string         `gorm:"size:500" json:"standard_name"`
	PublishDate   string         `gorm:"size:20" json:"publish_date"`
	ImplementDate string         `gorm:"size:20" json:"implement_date"`
	Status        string         `gorm:"size:20" json:"status"`
	Nature        string         `gorm:"size:20" json:"nature"`
	CCSCode       string         `gorm:"size:20" json:"ccs_code"`
	ICSCode       string         `gorm:"size:50" json:"ics_code"`
	Department    string         `gorm:"size:100" json:"department"`
	Publisher     string         `gorm:"size:200" json:"publisher"`
	Description   string         `gorm:"type:text" json:"description"`
	StandardType  string         `gorm:"size:20;default:'地标'" json:"standard_type"`
	DownloadURL   string         `gorm:"size:500" json:"download_url"`
	CheckStatus   int            `gorm:"default:0" json:"check_status"`
	CategoryID    uint           `gorm:"index" json:"category_id"`
	Source        string         `gorm:"size:20;default:'manual'" json:"source"` // crawl/爬取 manual/人工
	CreatedAt     time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l LocalStandard) MarshalJSON() ([]byte, error) {
	type Alias LocalStandard
	return jsonMarshalWithTime(Alias(l), l.CreatedAt)
}

func (LocalStandard) TableName() string {
	return "local_standard"
}

// 政策分类表
type PolicyCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex" json:"code"`
	Description string         `gorm:"size:500" json:"description"`
	CreatedAt   time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PolicyCategory) TableName() string {
	return "policy_category"
}

// 推荐政策解析表
type PolicyRecommend struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PolicyID   uint           `gorm:"not null;index" json:"policy_id"`
	PolicyType string         `gorm:"size:20;not null" json:"policy_type"` // 国标/行标/地标
	PolicyName string         `gorm:"size:500" json:"policy_name"`         // 关联政策名称（冗余存储）
	Title      string         `gorm:"size:200" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	Sort       int            `gorm:"default:0" json:"sort"`
	CreatedAt  time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PolicyRecommend) TableName() string {
	return "policy_recommend"
}

// 爬虫日志表
type CrawlerLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Script    string    `gorm:"type:text" json:"script"`
	Query     string    `gorm:"size:200" json:"query"`
	Result    string    `gorm:"type:text" json:"result"`
	Status    string    `gorm:"size:20" json:"status"` // success/failed
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}

func (CrawlerLog) TableName() string {
	return "crawler_log"
}

// 自动迁移
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&NationalStandard{},
		&IndustryStandard{},
		&LocalStandard{},
		&PolicyCategory{},
		&PolicyRecommend{},
		&CrawlerLog{},
	)
}