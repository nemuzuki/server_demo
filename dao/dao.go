package dao

import (
	"bufio"
	"encoding/json"
	"os"
	"server_demo/models"
)

var (
	postIndexMap  map[int32][]*models.Post
	topicIndexMap map[int32]*models.Topic
)

//初始化topic map
func InitTopicIndexMap() error {
	open, err := os.Open("./data/topic") //相对路径，从根目录开始算
	if err != nil {
		return err
	}
	topicTmpMap := make(map[int32]*models.Topic)
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		text := scanner.Text()
		var topic models.Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			return err
		}
		topicTmpMap[int32(topic.ID)] = &topic
	}
	topicIndexMap = topicTmpMap
	return nil
}

//初始化post map
func InitPostIndexMap() error {
	open, err := os.Open("./data/post")
	if err != nil {
		return err
	}
	postTmpMap := make(map[int32][]*models.Post)
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		text := scanner.Text()
		var post models.Post
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			return err
		}
		parentId := int32(post.ParentID)
		_, ok := postTmpMap[parentId]
		if !ok {
			postTmpMap[parentId] = []*models.Post{}
		}
		postTmpMap[parentId] = append(postTmpMap[parentId], &post)

	}
	postIndexMap = postTmpMap
	return nil
}

//根据topic id查询topic。注意函数首字母大写才能被别的模块调用
func QueryTopicById(id int32) *models.Topic {
	return topicIndexMap[id]
}

//根据请求的topic id查询post
func QueryPostByParentId(parentID int32) []*models.Post {
	return postIndexMap[parentID]
}
