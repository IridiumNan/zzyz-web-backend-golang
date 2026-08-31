package test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/database"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
)

// dumpPostConfig
// return true if success
// return false if fail, you should terminate the test
func dumpPostConfig(fileName string) bool {
	content := `
# 成员 id (由后端添加成员的时候生成)
member_id = 2

# 作者名称 (允许昵称)
author = "张三"

# 邮箱 可选, 没有填写放空即可
email = ""

# 主文件类型
# TODO: 可选项 [ md, html, docx ]
# WARN: 当前支持 [ md, html ]
# 0 md
# 1 html
# 2 docx
index_format = 0

# 文章标题
title = "如何做好一个golang后端"

# 文章概要(多行使用三个双引号)
overview = """
这是一篇关于如何做好golang后端的文章.
主要涉及接口设计， 系统分析，流程图制作等
"""


# NOTE: 如果是全新的文章则不需要填写
# WARN: 如果是退回的或者更新的文章需要填写
# TODO: 通过简易的 html 页面查看文章列表
post_id = -1

# 直接写标签, 如果不存在则自动创建新的标签
# 每个文章只支持 3 个标签以内
# WARN: 如果系统的总标签数量超过 20 则会拒绝创建
# 查看网页前端获取现在已有的标签
tags = ["golang", "CS"]
	`

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		slog.Error("when crate file ", "file_name", fileName, "err", err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		slog.Error("when write file", "file_name", fileName, "err", err)
		return false
	}
	return true
}

// checkResult return false if the attr not match
func checkResult(postConfig *models.PostConfig, expectedConfig *models.PostConfig) bool {
	reportIfNotMatch := func(attrName string, postVal any, expectedVal any) bool {
		if postVal != expectedVal {
			slog.Error("Found not match attr", "attr", attrName, "provided_val", postVal, "expected_val", expectedVal)
			return false
		}
		return true
	}

	var match bool
	if match = reportIfNotMatch("member_id", postConfig.MemberID, expectedConfig.MemberID); !match {
		return false
	}
	if match = reportIfNotMatch("author", postConfig.Author, expectedConfig.Author); !match {
		return false
	}
	if match = reportIfNotMatch("email", postConfig.Email, expectedConfig.Email); !match {
		return false
	}
	if match = reportIfNotMatch("index_format", postConfig.IndexFormat, expectedConfig.IndexFormat); !match {
		return false
	}
	if match = reportIfNotMatch("title", postConfig.Title, expectedConfig.Title); !match {
		return false
	}
	if match = reportIfNotMatch("overview", postConfig.Overview, expectedConfig.Overview); !match {
		return false
	}
	if match = reportIfNotMatch("post_id", postConfig.PostID, expectedConfig.PostID); !match {
		return false
	}

	if len(postConfig.Tags) != len(expectedConfig.Tags) {
		return false
	}

	for idx := range postConfig.Tags {
		if postConfig.Tags[idx] != expectedConfig.Tags[idx] {
			return false
		}
	}

	return true
}

func TestLoadPostConfig(t *testing.T) {
	fileName := "test.toml"
	success := dumpPostConfig(fileName)
	if !success {
		t.Fatal("fail to dump test configuration, exit the test")
	}

	postConfig, err := database.LoadPostConfigFromToml(fileName)
	if err != nil {
		slog.Error("error when load config from toml file", "err", err)
	}

	match := checkResult(postConfig, &models.PostConfig{
		MemberID:    2,
		Author:      "张三",
		Email:       "",
		IndexFormat: models.FormatMD,
		Title:       "如何做好一个golang后端",
		Overview: `这是一篇关于如何做好golang后端的文章.
主要涉及接口设计， 系统分析，流程图制作等
`,
		PostID: -1,
		Tags:   []string{"golang", "CS"},
	})

	if !match {
		t.Fatal("test failed")
	}

	err = os.Remove(fileName)
	if err != nil {
		slog.Error("error when remove the test config file, you should remove it manually", "fileName", fileName)
	}
}
