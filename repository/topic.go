package repository

import (
	"sync"
)

type Topic struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type TopicDao struct {
}

var (
	topicDao  *TopicDao //空结构体指针
	topicOnce sync.Once //单例模式，使用sync.Once来保证只被调用一次
)

//返回一个空结构体的指针  ，Do里面的内容只会执行一次
func NewTopicDaoInstance() *TopicDao {
	topicOnce.Do(
		func() {
			topicDao = &TopicDao{}
		})
	return topicDao
}

//通过id查询 主题
func (*TopicDao) QueryTopicById(id int64) *Topic {
	return topicIndexMap[id] //topicIndexMap是db_init.go中定义的map[int64]*Topic类型变量
}
