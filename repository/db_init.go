package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

var (
	topicIndexMap map[int64]*Topic
	postIndexMap  map[int64][]*Post
	rwMutex       sync.RWMutex //读写锁
)

func Init(filePath string) error {
	if err := initTopicIndexMap(filePath); err != nil {
		return err
	}
	if err := initPostIndexMap(filePath); err != nil {
		return err
	}
	return nil
}

func initTopicIndexMap(filePath string) error {
	open, err := os.Open(filePath + "topic") //打开topic文件
	if err != nil {
		return err
	}
	defer func(open *os.File) { //在闭包中进行关闭文件和处理错误
		err := open.Close()
		if err != nil {

		}
	}(open)
	scanner := bufio.NewScanner(open)     //*file实现了read方法，即实现了reader接口类型
	topicTmpMap := make(map[int64]*Topic) //临时map
	for scanner.Scan() {                  //判断是否读完 循环读取
		text := scanner.Text() // 将读到的一行转化为string
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil { //由于数据文件为json格式，所以反序列化即可
			return err
		}
		topicTmpMap[topic.Id] = &topic //存到临时map里面
	}
	topicIndexMap = topicTmpMap //指向临时map
	return nil
}

func initPostIndexMap(filePath string) error {
	open, err := os.Open(filePath + "post")
	if err != nil {
		return err
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)
	scanner := bufio.NewScanner(open)
	postTmpMap := make(map[int64][]*Post)
	for scanner.Scan() {
		text := scanner.Text()
		var post Post
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			return err
		}
		posts, ok := postTmpMap[post.ParentId]
		if !ok { //没查到 就新建一个键值对
			postTmpMap[post.ParentId] = []*Post{&post} //一个主题可以对应多个回帖，以主题的id为主键
			continue
		}
		posts = append(posts, &post)      //查到的话 就切片增加一个
		postTmpMap[post.ParentId] = posts //更新value
	}
	postIndexMap = postTmpMap //将临时赋给永久的
	return nil
}
