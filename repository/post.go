package repository

import (
	"encoding/json"
	"os"
	"sync"
)

type Post struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parent_id"` //一个主题可以对应多个回帖
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type PostDao struct {
} //大写的空结构体  空结构体作用：成为map的value，只是为了利用map的key不可重复性；
// 在channel里面，空结构体 与 channel 可谓是一个经典组合，有时候我们只是需要一个信号来控制程序的运行逻辑，并不在意其内容如何

var (
	postDao  *PostDao //空结构体的指针
	postOnce sync.Once
)

func NewPostDaoInstance() *PostDao {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}
func (*PostDao) QueryPostsByParentId(parentId int64) []*Post {
	return postIndexMap[parentId]
}

func (*PostDao) InsertPost(post *Post) error {
	f, err := os.OpenFile("./data/post", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	marshal, _ := json.Marshal(post)
	if _, err = f.WriteString(string(marshal) + "\n"); err != nil {
		return err
	}

	rwMutex.Lock()
	postList, ok := postIndexMap[post.ParentId]
	if !ok {
		postIndexMap[post.ParentId] = []*Post{post}
	} else {
		postList = append(postList, post)
		postIndexMap[post.ParentId] = postList
	}
	rwMutex.Unlock()
	return nil
}
