### Go语言工程实践——实现简易论坛服务端

这篇文章是在我完成极简版抖音项目之后撰写的，相比刚听网课的时候的懵懵懂懂，书写过一些go代码后对这个作业的理解深入了很多。这个作业是实现一个服务端go程序，下面从零开始开发。

#### 准备工作

首先建立一个名为server_demo的空目录作为项目根目录，然后使用`go mod init server_demo`初始化一个server_demo模块，此时会发现自动生成一个go.mod文件，该文件是模块的配置文件，此时里面只有模块名和go版本

```
 module server_demo
 
 go 1.12
```

之后，我们需要导入go的web框架gin，执行`go get github.com/gin-gonic/gin`，go.mod会多出一行，表示该模块依赖的模块（包）及其版本

```
 require github.com/gin-gonic/gin v1.9.1 // indirect
```

还会生成一个go.sum文件，包含了每个包的名称、版本和校验和信息

```
 github.com/bytedance/sonic v1.5.0/go.mod h1:ED5hyg4y6t3/9Ku1R6dU/4KyJ48DZ4jPhfY1O2AihPM=
 github.com/bytedance/sonic v1.9.1 h1:6iJ6NqdoxCDr6mbY8h18oSO+cShGSMRGCEo7F2h0x8s=
```

#### 项目描述

项目背景：实现一个论坛，论坛里有很多Topic（话题），每个话题下有若干Post（帖子）。要实现的业务：根据topic id查询一个话题的名称，根据topic id查询一个话题下的所有帖子

根据项目分层原则，可以将代码分成如下目录：

-   models目录下定义Topic、Post以及数据在内存中存放的数据结构
-   dao目录下实现业务所需的原子操作，主要是对数据库的初始化和增删改查
-   controller目录下实现服务接口
-   routes目录下定义接口路由
-   主目录下main.go中写主函数

项目结构如下

```
 .                  
 ├── controller     
 │   ├── bmi.go     
 │   └── query.go   
 ├── dao            
 │   └── dao.go     
 ├── data           
 │   ├── post       
 │   └── topic      
 ├── go.mod         
 ├── go.sum         
 ├── main.go        
 ├── models         
 │   ├── post.go    
 │   └── topic.go   
 ├── routes         
 │   ├── api.go     
 │   └── route.go   
 ├── server_demo.exe
```

下面将分层阐述具体实现方法

#### models

首先根据data中的数据格式，使用<https://oktools.net/json2go>工具将json转为go结构体

```
 //Post
 {"id":1,"parent_id":1,"content":"小姐姐快来1","create_time":1650437616,"user_id":1}
 {"id":2,"parent_id":1,"content":"小姐姐快来2","create_time":1650437617,"user_id":2}
 {"id":3,"parent_id":1,"content":"小姐姐快来3","create_time":1650437618,"user_id":13}
 {"id":5,"parent_id":1,"content":"测试内容嗨嗨嗨嗨","create_time":1652073758,"user_id":2}
 
 //Topic
 {"id":1,"title":"青训营来啦!","content":"小姐姐，快到碗里来~","create_time":1650437625}
 {"id":2,"title":"青训营来啦!","content":"小哥哥，快到碗里来~","create_time":1650437640}
```

在Post.go定义Post结构体，其中parent id就是topic id

```
 package models
 
 type Post struct{
     ID int `json:"id"`
     ParentID int `json:"parent_id"`
     Content string `json:"content"`
     CreateTime int `json:"create_time"`
     UserID int `json:"user_id"`
 }
```

在Topic.go定义Topic结构体，每个topic有一个标题title和一个解释内容content

```
 package models
 
 type Topic struct{
     ID int `json:"id"`
     Title string `json:"title"`
     Content string `json:"content"`
     CreateTime int `json:"create_time"`
 }
```

#### dao

dao.go定义数据存放的结构，存储Post的数据结构postIndexMap，存储Topic的数据结构topicIndexMap。postIndexMap是一个将topic id映射到Post列表的Map，每个列表中是一个话题下的帖子。topicIndexMap是一个将topic id映射到Topic的Map

```
 var (
     postIndexMap  map[int32][]*models.Post
     topicIndexMap map[int32]*models.Topic
 )
```

根据data内容对这两个存储结构初始化。InitTopicIndexMap读取data中的数据存入topicIndexMap中，InitPostIndexMap将数据存入postTmpMap中，要注意的是**这两个函数必须首字母大写，这样才能被外部模块调用**（go特性）

```
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
```

之后定义从两个Map中读取数据的函数

```
 //根据topic id查询topic。注意函数首字母大写才能被别的模块调用
 func QueryTopicById(id int32) *models.Topic {
     return topicIndexMap[id]
 }
 
 //根据请求的topic id查询post
 func QueryPostByParentId(parentID int32) []*models.Post {
     return postIndexMap[parentID]
 }
```

#### controller

query.go中编写接口，即收到请求的处理过程。QueryTopicById根据请求的topic id获取topic内容，首先调用dao.QueryTopicById查找到Topic，然后读取里面的各个属性，使用c.JSON来写响应体。QueryPostById也是同理，响应一个topic下的所有post信息，返回一个列表

```
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
 
```

#### routes

api.go中定义上面实现的两个接口的路由，比如我们要调用controller.QueryTopicById接口，那么URL就是127.0.0.1:8080/server_demo/getTopic/

```
 func ApiRoutes(r *gin.Engine) {
 
     apiRouter := r.Group("server_demo")
 
     apiRouter.GET("/getTopic/", controller.QueryTopicById)
     apiRouter.GET("/getPost/", controller.QueryPostById)
 }
```

route.go中对路由初始化，包括三个过程

-   注册全局中间件：日志、发生panic时恢复程序运行两个功能
-   注册 API 路由
-   配置 404 请求

```
 // SetupRoute 路由初始化
 func SetupRoute(router *gin.Engine) {
 
     // 注册全局中间件
     registerGlobalMiddleWare(router)
 
     // 注册 API 路由
     ApiRoutes(router)
 
     // 配置 404 请求
     setup404Handler(router)
 }
 
 func registerGlobalMiddleWare(router *gin.Engine) {
     router.Use(
         gin.Logger(),//日志中间件
         gin.Recovery(),//发生panic时恢复程序运行
     )
 }
 
 // 处理404请求
 func setup404Handler(router *gin.Engine) {
     router.NoRoute(func(c *gin.Context) {
         c.JSON(http.StatusNotFound, gin.H{
             "error_code":    404,
             "error_message": "路由未定义，请确认 url 和请求方法是否正确。",
         })
     })
 }
```

#### main

main.go

-   初始化数据库
-   配置路由
-   启动整个服务，绑定在8080端口上

```
 func main() {
     //初始化数据库
     dao.InitPostIndexMap()
     dao.InitTopicIndexMap()
     //配置路由
     router := gin.New()
     routes.SetupRoute(router)
     //启动
     err := router.Run(":8080")
     if err != nil {
         fmt.Println(err.Error())
     }
 }
```

#### 测试

运行方法

```sh
go build
./server_demo.exe
```

可以使用Postman来测试接口，可以看到两个接口都可以正常工作


![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9cdc563c0451492f99e0a17efa2eee1c~tplv-k3u1fbpfcp-watermark.image?)


![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/287c8558739f4d20b21dfd1fdc01fdb3~tplv-k3u1fbpfcp-watermark.image?)