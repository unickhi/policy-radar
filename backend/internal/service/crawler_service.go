package service

import (
	"fmt"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
	"time"
)

type CrawlerService struct {
	logRepo       *repository.CrawlerLogRepo
	nationalRepo  *repository.NationalStandardRepo
	industryRepo  *repository.IndustryStandardRepo
	localRepo     *repository.LocalStandardRepo
}

func NewCrawlerService(
	logRepo *repository.CrawlerLogRepo,
	nationalRepo *repository.NationalStandardRepo,
	industryRepo *repository.IndustryStandardRepo,
	localRepo *repository.LocalStandardRepo,
) *CrawlerService {
	return &CrawlerService{
		logRepo:      logRepo,
		nationalRepo: nationalRepo,
		industryRepo: industryRepo,
		localRepo:    localRepo,
	}
}

func (s *CrawlerService) ExecuteScript(script, query string) (*model.CrawlerLog, error) {
	log := &model.CrawlerLog{
		Script:    script,
		Query:     query,
		Status:    "success",
		CreatedAt: time.Now(),
	}
	return log, nil
}

func (s *CrawlerService) ImportData(data []map[string]interface{}, targetType string) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	count := 0
	for _, item := range data {
		switch targetType {
		case "national":
			std := mapToNationalStandard(item)
			if std.StandardNo == "" {
				continue
			}
			std.Source = "crawl"
			if err := s.nationalRepo.Upsert(std); err != nil {
				fmt.Printf("导入国标 %s 失败: %v\n", std.StandardNo, err)
				continue
			}
			count++
		case "industry":
			std := mapToIndustryStandard(item)
			if std.StandardNo == "" {
				continue
			}
			std.Source = "crawl"
			if err := s.industryRepo.Upsert(std); err != nil {
				fmt.Printf("导入行标 %s 失败: %v\n", std.StandardNo, err)
				continue
			}
			count++
		case "local":
			std := mapToLocalStandard(item)
			if std.StandardNo == "" {
				continue
			}
			std.Source = "crawl"
			if err := s.localRepo.Upsert(std); err != nil {
				fmt.Printf("导入地标 %s 失败: %v\n", std.StandardNo, err)
				continue
			}
			count++
		}
	}
	return count, nil
}

func mapToNationalStandard(item map[string]interface{}) *model.NationalStandard {
	std := &model.NationalStandard{}
	if v, ok := item["standard_no"].(string); ok {
		std.StandardNo = v
	}
	if v, ok := item["standard_name"].(string); ok {
		std.StandardName = v
	}
	if v, ok := item["english_name"].(string); ok {
		std.EnglishName = v
	}
	if v, ok := item["publish_date"].(string); ok {
		std.PublishDate = v
	}
	if v, ok := item["implement_date"].(string); ok {
		std.ImplementDate = v
	}
	if v, ok := item["status"].(string); ok {
		std.Status = v
	}
	if v, ok := item["nature"].(string); ok {
		std.Nature = v
	}
	if v, ok := item["category"].(string); ok {
		std.Category = v
	}
	if v, ok := item["is_adopted"].(string); ok {
		std.IsAdopted = v
	}
	if v, ok := item["link1"].(string); ok {
		std.Link1 = v
	}
	if v, ok := item["link2"].(string); ok {
		std.Link2 = v
	}
	if v, ok := item["ccs_code"].(string); ok {
		std.CCSCode = v
	}
	if v, ok := item["ics_code"].(string); ok {
		std.ICSCode = v
	}
	if v, ok := item["department"].(string); ok {
		std.Department = v
	}
	if v, ok := item["technical_dept"].(string); ok {
		std.TechnicalDept = v
	}
	if v, ok := item["publisher"].(string); ok {
		std.Publisher = v
	}
	if v, ok := item["description"].(string); ok {
		std.Description = v
	}
	if v, ok := item["download_url"].(string); ok {
		std.DownloadURL = v
	}
	std.StandardType = "国标"
	std.CheckStatus = 0
	return std
}

func mapToIndustryStandard(item map[string]interface{}) *model.IndustryStandard {
	std := &model.IndustryStandard{}
	if v, ok := item["standard_no"].(string); ok {
		std.StandardNo = v
	}
	if v, ok := item["standard_name"].(string); ok {
		std.StandardName = v
	}
	if v, ok := item["publish_date"].(string); ok {
		std.PublishDate = v
	}
	if v, ok := item["implement_date"].(string); ok {
		std.ImplementDate = v
	}
	if v, ok := item["status"].(string); ok {
		std.Status = v
	}
	if v, ok := item["revision_type"].(string); ok {
		std.RevisionType = v
	}
	if v, ok := item["technical_owner"].(string); ok {
		std.TechnicalOwner = v
	}
	if v, ok := item["approve_dept"].(string); ok {
		std.ApproveDept = v
	}
	if v, ok := item["industry_class"].(string); ok {
		std.IndustryClass = v
	}
	if v, ok := item["standard_class"].(string); ok {
		std.StandardClass = v
	}
	if v, ok := item["replace_standard"].(string); ok {
		std.ReplaceStandard = v
	}
	if v, ok := item["detail_link"].(string); ok {
		std.DetailLink = v
	}
	if v, ok := item["download_url"].(string); ok {
		std.DownloadURL = v
	}
	if v, ok := item["ccs_code"].(string); ok {
		std.CCSCode = v
	}
	if v, ok := item["ics_code"].(string); ok {
		std.ICSCode = v
	}
	std.StandardType = "行标"
	std.CheckStatus = 0
	return std
}

func mapToLocalStandard(item map[string]interface{}) *model.LocalStandard {
	std := &model.LocalStandard{}
	if v, ok := item["standard_no"].(string); ok {
		std.StandardNo = v
	}
	if v, ok := item["standard_name"].(string); ok {
		std.StandardName = v
	}
	if v, ok := item["publish_date"].(string); ok {
		std.PublishDate = v
	}
	if v, ok := item["implement_date"].(string); ok {
		std.ImplementDate = v
	}
	if v, ok := item["status"].(string); ok {
		std.Status = v
	}
	if v, ok := item["nature"].(string); ok {
		std.Nature = v
	}
	if v, ok := item["department"].(string); ok {
		std.Department = v
	}
	if v, ok := item["publisher"].(string); ok {
		std.Publisher = v
	}
	if v, ok := item["detail_link"].(string); ok {
		std.DetailLink = v
	}
	if v, ok := item["download_url"].(string); ok {
		std.DownloadURL = v
	}
	if v, ok := item["ccs_code"].(string); ok {
		std.CCSCode = v
	}
	if v, ok := item["ics_code"].(string); ok {
		std.ICSCode = v
	}
	std.StandardType = "地标"
	std.CheckStatus = 0
	return std
}