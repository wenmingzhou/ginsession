package ginsession

import "fmt"

//内存版Session服务

//SessionData 支持的操作

//Get 根据key获取值
func (s *SessionData) Get(key string) (value interface{}, err error) {
	//获取读锁
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	value, ok := s.Data[key]
	if !ok {
		return nil, fmt.Errorf("invalid key")
	}
	return value, nil
}

//Set 根据key设置值
func (s *SessionData) Set(key string, value interface{}) {
	//获取写锁
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.Data[key] = value
	return
}

//Del 根据key删除值
func (s *SessionData) Del(key string) {
	//获取写锁
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	delete(s.Data, key)
	return
}
