# 实现一个review-job完成数据的处理

## 仿照kratos的internal下的文件自定义一个job文件夹
### 1. 在main函数里面，仿照gs,hs的实现方式，实现一个js使得job层可以一直运行
### 2. 定义自己的conf.yaml和conf.prot文件，实现参数的定义
### 3. 各结构体
#### 3.1 关于Msg结构体。是因为针对Ela的插入还是更新操作，需要拿到数据的字段，因此可以通过postman发送一个插入的操作，拿到其在elasticsearch界面的数据，然后复制到一个```json.cn```的网站去观察，注意一定要匹配
#### 3.2 关于Ela结构体，是因为原本的Ela结构体没有index字段，因此可以通过再包一层得到index
### 4. 要完成这个服务必须启动几个组件：canal\kafka\elasticsearch， 这些都可以在李文周的博客找到