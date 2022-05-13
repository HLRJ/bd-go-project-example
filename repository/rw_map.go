// 以下实现参考《Go 并发之三种线程安全的 map - AFreeCoder的文章 - 知乎》
// 链接: https://zhuanlan.zhihu.com/p/356739568
package repository

import "sync"

//go的内建的map对象不是线程(goroutine)安全的，并发读写的时候运行时会有检测，遇到并发问题就会导致panic。sync包的中sync.map是并发安全的
//加了读写锁实现线程安全的map[int]int类型
type RWTopicMap struct {
	sync.RWMutex // 加了读写锁的Topicmap
	m            map[int64]*Topic
}

func NewRWTopicMap(n int) *RWTopicMap { //创建一个读写锁的map
	return &RWTopicMap{
		m: make(map[int64]*Topic), //map对象必须在使用之前初始化，如果不初始化会出现panic异常
	}
}

func (m *RWTopicMap) Get(k int64) (*Topic, bool) { //查询map的值
	m.RLock()            //加个读锁
	defer m.RUnlock()    //解读锁
	v, existed := m.m[k] //在锁的保护下从map读取，map[key]返回的结果可以是一个，也可以是两个值
	return v, existed    // 如果读到值的话，existed就是true，返回的是真正的值，true，如果没有值的话，就返回0，false,因为有的值就是0

}

func (m *RWTopicMap) Set(k int64, v *Topic) { //更新map的值
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

type RWPostMap struct { //加了读写锁的PostMap
	sync.RWMutex
	m map[int64][]*Post
}

func NewRWPostMap(n int) *RWPostMap {
	return &RWPostMap{
		m: make(map[int64][]*Post),
	}
}

func (m *RWPostMap) Get(k int64) ([]*Post, bool) {
	m.RLock()
	defer m.RUnlock()
	v, existed := m.m[k]
	return v, existed
}

func (m *RWPostMap) Set(k int64, v []*Post) {
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}
