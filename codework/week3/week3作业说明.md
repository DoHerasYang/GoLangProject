# Week3作业说明

## 1. 环境准备

### 1. 删除所有的Deployment

<img src="./img/01.png" alt="03" style="zoom:60%;" />

<img src="./img/02.png" alt="03" style="zoom:60%;" />

## 2. 将 webook 的端口更改为8081

- 修改原始代码 - `config.types.go` 新增结构体 `GinConfig`

  <img src="./img/03.png" alt="03" style="zoom:40%;" />

- 修改 `k8s.go`  和  `dev.go` 两个配置文件

  <img src="./img/04.png" alt="03" style="zoom:40%;" />

- 更改 `main.go`

  <img src="./img/05.png" alt="03" style="zoom:50%;" />

- 修改webook的Web服务Kubernetes - `deployment` / `service` 将所有的端口更改，并修改为两个`replicas` 

  <img src="./img/06.png" alt="03" style="zoom:50%;" />

---

## 3. 修改Redis端口为6380

- 修改 `config.k8s.go` 的 redis 连接

  <img src="./img/07.png" alt="03" style="zoom:50%;" />

- 修改 `webook-redis-deployment.yaml`  以及  `webook-redis-service.yaml`

  <img src="./img/08.png" alt="03" style="zoom:50%;" />

---

## 4. MySQL 修改端口为3308

<img src="./img/09.png" alt="03" style="zoom:40%;" />

---

## 5. 运行截图

- 运行获得所有`pods`的运行状态

<img src="./img/10.png" alt="03" style="zoom:40%;" />

- 运行所有`services`并获得

<img src="./img/11.png" alt="03" style="zoom:40%;" />

- 可以能够自由注册

  ![12](./img/12.png)

---

## 6. 思考题  - 修改MySQL的默认端口

- 修改MySQL配置文件`mysql-deployment.yaml ` / `mysql-service.yaml` 中的 `containerPort` / `targetPort` 需要重新定制`MySQL`Docker镜像；

- 创建 `my.cnf`

  ```
  [mysqld]
  user=mysql
  pid-file=/var/run/mysqld/mysqld.pid
  socket=/var/run/mysqld/mysqld.sock
  port=3308
  datadir=/var/lib/mysql
  ```

- 创建 - `Dockerfile`

  ```dockerfile
  FROM mysql:8.0
  COPY my.cnf /etc/mysql/conf.d/
  EXPOSE 3308
  ```

- 打包新的镜像

  ```
  docker build -t doheras/mysql:8.0-customport .
  ```

- 引用新的镜像并修改配置文件

  <img src="./img/13.png" alt="03" style="zoom:40%;" />









