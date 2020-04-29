package ginsession

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/satori/go.uuid"
	"sync"
)


const (
	sessionCookieName  ="session_id"  //session_id 在cookie中对应的key
	sessionContextName ="session"    //session data在gin上下文对应的key
)

var (
	MgrObj *Mgr  //全局的Session管理对象(大仓库)
)
//session服务
//表示一个具体的用户session数据
type SessionData struct {
	ID string
	Data map[string]interface{}
	rwLock sync.RWMutex //读写锁,锁的时上面的data
}

//NewSessionData 构造函数
func NewSessionData(id string)*SessionData  {
	return &SessionData{
		ID:     id,
		Data:   make(map[string]interface{},8),
	}
}

//是一个全局的session管理
type Mgr struct {
	Session map[string]*SessionData
	rwLock sync.RWMutex
}

func InitMgr(){
	MgrObj = &Mgr{
		Session: make(map[string]*SessionData,1024), //初始化1024用来存取用户的session data
	}
}
//根据传进来的session id 找到对应的session data
func (m *Mgr)GetSessionData(sessionID string)(sd *SessionData,err error)  {
	//取之前加锁
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	sd,ok :=m.Session[sessionID]
	if !ok{
		err =fmt.Errorf("invalid session id")
		return
	}
	return
}

//创建一条Session记录
func (m *Mgr)CreateSession()(sd *SessionData){
	// 1.造一个sessionID
	//uuidObj :=uuid.  // ?
	uuidObj :=uuid.NewV4()
	// 2.造一个和它对应的sessionData
	sd =NewSessionData(uuidObj.String())
	// 3.返回SessionData
	return
}


//实现一个gin框架的中间件
//所有流经我这个中间件的请求，它的上下文中肯定会有一个session ->session data
func SessionMiddleware(mgrObj *Mgr) gin.HandlerFunc{
	if mgrObj ==nil{
		panic("must call InitMgr before use it")
	}
	return func(c *gin.Context{
		// 1 从请求的Cookie中获取Session ID
		var sd *SessionData //session data
		session_id ,err :=c.Cookie(sessionCookieName)
		if err !=nil{
			// 1.1取不到session_id ->给这个新用户创建一个新的Session data 同时分配一个session id
			sd =mgrObj.CreateSession()
		}
		// 1.2取到session_id
		sd,err =mgrObj.GetSessionData(session_id)
		// 2根据session_id 去session大仓库中取到对应的session  data
		if err !=nil {
			//2.1根据用户传过来的session_id 在大仓库中根本取不到session data
			sd :=mgrObj.CreateSession()
			//2.2更新用户cookie中保存的哪个session_id
			session_id =sd.ID

		}
		// 3如何实现让后续的所有的处理请求的方法都能拿到session data
		// 3利用gin的c.set("session",session data)
		c.Set(sessionContextName,sd)
		//在gin框架中，要回写cookie必须在处理请求的函数返回之前
		c.SetCookie(sessionCookieName,session_id,3600,"/","127.0.0.1",false,true)
		c.Next() //执行后续的请求
	})
}
