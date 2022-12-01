# go_redis
## TCP服务器框架
* 实现了系统信号获取，以便于通知客户端服务关闭
* 使用atomic包完成了布尔类型锁，避免handler的同步问题
* 封装WaitGroup方法，实现可延时退出
## Redis DB核心
* 实现了 RESP 协议的封装与解析
* 实现了基于sync/map的线程安全字典类型交互
  * `TYPE $KEY`
  * `KEYS $KEY`
  * `FLUSHDB`
  * `RENAME $KEY1 $NEWKEY1`
  * `RENAMENX $KEY1 $NEWKEY1`
  * `DEL $KEY1 [$KEY2 ...]`
  * `EXIST $KEY1 [$KEY2 ...]`
* 实现了基于string的类型交互
  * `GET $KEY`
  * `SET $KEY1 $VALUE [$KEY2 $VALUE2 ...]`
  * `SETNX $KEY1 $VALUE`
  * `GETSET $KEY1 $VALUE`
  * `STRLEN $KEY`