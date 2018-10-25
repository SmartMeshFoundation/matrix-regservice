# matrix-regservice
matrix-regservice 作为matrix的第三方 application service 使用,主要作用为限制SmartRaiden所使用的Matrix服务器只能接受有效的Raiden用户注册.
即注册在此Matrix服务器的用户只能是SmartRaiden节点,这些节点拥有统一的名字,并且可以保证这些用户必须掌握对应的私钥才能注册.


## Matrix 服务器的部署
Matrix的安装部署参考[matrix](https://github.com/matrix-org/synapse).

### 修改配置文件

1. enable_registration 为False, 保证用户不能通过正常的接口来注册,只能通过第三方application service注册用户
2. search_all_users 修改为True,保证可以检索用户
3. expire_access_token 修改为True, 保证用户会自动退出,防止第三方重放攻击
4. port 8008 修改为8007
5. 移除port 8007下的webclient,禁止通过web方式登录
6. app_service_config_files 修改为[ registration.yaml]
7. trusted_third_party_id_servers 全部注释

## matrix-regservice 部署安装

### 安装matrix-regservice
```bash
go get github.com/SmartMeshFoundation/matrix-regservice
```
然后将matrix-regservice复制到$PATH路径下

### 生成配置文件
切换到matrix工作目录(homeserver.yaml 所在目录)
```bash
matrix-regservice --matrixdomain yourdomain.com genconfig
```
可以看到在matrix工作目录下生成了registration.yaml和run.sh
我的样例
registration.yaml
```yaml
id: Q7PM2E53RE-transport02.smartmesh.cn
hs_token: RNI4CGEDTKC4WJTB4RZWRK4NOKA7M4PREUW6F2GZ
as_token: LODE52N2CKVXMOURUMAWLEEEXMWB4DKIKRI246XD
url: http://127.0.0.1:8009/regapp/1
sender_localpart: transport02.smartmesh.cn
protocols:
- regapp.transport02.smartmesh.cn
namespaces:
  users:
  - exclusive: false
    regex: '@.*'
  aliases: []
  rooms: []
```
run.sh
```bash
matrix-regservice --astoken LODE52N2CKVXMOURUMAWLEEEXMWB4DKIKRI246XD --hstoken RNI4CGEDTKC4WJTB4RZWRK4NOKA7M4PREUW6F2GZ --matrixurl http://127.0.0.1:8008/_matrix/client/api/v1/createUser --host 127.0.0.1 --port 8009 --datapath .matrix --matrixdomain transport02.smartmesh.cn --verbosity 5
```

### nginx 反向代理配置
the following is a example configuration file for nginx.
```conf
  server {
    listen 8008;
    server_name localhost;

  location /_matrix {
    proxy_pass http://127.0.0.1:8007;
    proxy_max_temp_file_size 0;
    proxy_connect_timeout 30;
    }

  location /regapp/1 {
    proxy_pass http://127.0.0.1:8009;
    proxy_max_temp_file_size 0;
    proxy_connect_timeout 30;
    }
  }
```