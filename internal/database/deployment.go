package database

import (
	"awesomeProject/internal/model"
)

func MigrateDeploymentScheme() error {
	return DB.AutoMigrate(
		model.Deployment{},
		model.ApplicationTemplate{},
	)
}

// Sorry for the old WordPress app :d
func DummyDeploymentData() error {
	return DB.Create(&model.ApplicationTemplate{
		Name:        "WordPress",
		Description: "WordPress (also known as WP or WordPress.org) is a web content management system. It was originally created as a tool to publish blogs but has evolved to support publishing other web content, including more traditional websites, mailing lists and Internet forum, media galleries, membership sites, learning management systems and online stores. Available as free and open-source software, WordPress is among the most popular content management systems",
		ExecCommand: "docker run --name some-wordpress -p 8080:80 -d wordpress",
	}).Error
}
