package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server_demo/dao"
	"server_demo/models"
	"strconv"
)

//根据请求的topic id获取topic内容
func QueryTopicById(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 8)
	topic := dao.QueryTopicById(int32(id))
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"status_msg":  "查询标题成功",
		"comment": models.Topic{
			ID:         int(id),
			Title:      topic.Title,
			Content:    topic.Content,
			CreateTime: topic.CreateTime,
		},
	})
}

//根据请求的topic id和post id获取post内容
func QueryPostById(c *gin.Context) {
	parentId, _ := strconv.ParseInt(c.Query("parent_id"), 10, 8)
	postList := dao.QueryPostByParentId(int32(parentId))
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"status_msg":  "查询帖子成功",
		"post_list":   postList,
	})
}
