# go_redis
## TCP服务器框架
* 实现了系统信号获取，以便于通知客户端服务关闭
* 使用atomic包完成了布尔类型锁，避免handler的同步问题
* 封装WaitGroup方法，实现可延时退出
