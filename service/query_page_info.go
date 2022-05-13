package service

import (
	"errors"
	"github.com/Moonlight-Zhao/go-project-example/repository"
	"sync"
)

type PageInfo struct { //对外提供服务
	Topic    *repository.Topic //进一步封装数据
	PostList []*repository.Post
}

func QueryPageInfo(topicId int64) (*PageInfo, error) { //PageInfo
	return NewQueryPageInfoFlow(topicId).Do()
}

func NewQueryPageInfoFlow(topId int64) *QueryPageInfoFlow {
	return &QueryPageInfoFlow{
		topicId: topId,
	}
}

type QueryPageInfoFlow struct { //包内使用
	topicId  int64
	pageInfo *PageInfo //方便给对外服务的结构体赋值

	topic *repository.Topic
	posts []*repository.Post
}

func (f *QueryPageInfoFlow) Do() (*PageInfo, error) { //方法 返回一个PageInfo指针   ,作用是判断和验证信息
	if err := f.checkParam(); err != nil {
		return nil, err
	}
	if err := f.prepareInfo(); err != nil { //封装到小写
		return nil, err
	}
	if err := f.packPageInfo(); err != nil { //封装到大写
		return nil, err
	}
	return f.pageInfo, nil
}

func (f *QueryPageInfoFlow) checkParam() error { //检查参数
	if f.topicId <= 0 {
		return errors.New("topic id must be larger than 0")
	}
	return nil
}

func (f *QueryPageInfoFlow) prepareInfo() error { //将信息封装到QueryPageInfoFlow结构体
	//获取topic信息
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { //获取话题和回复信息程序中使用了go协程来并行获取提高相应速度
		defer wg.Done()
		topic := repository.NewTopicDaoInstance().QueryTopicById(f.topicId)
		f.topic = topic
	}()
	//获取post列表
	go func() {
		defer wg.Done()
		posts := repository.NewPostDaoInstance().QueryPostsByParentId(f.topicId)
		f.posts = posts
	}()
	wg.Wait()
	return nil
}

func (f *QueryPageInfoFlow) packPageInfo() error { //将小写的数据封装为大写的  对包外提供服务
	f.pageInfo = &PageInfo{
		Topic:    f.topic,
		PostList: f.posts,
	}
	return nil
}
