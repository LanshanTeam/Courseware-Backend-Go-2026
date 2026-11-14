## 蓝山工作室Go组第六节课——关系型数据库（MySQL & PostgreSQL）

### 什么是数据库

在前面的课程中，我们学习了 Go 语言的基础语法、并发编程模型，并用 **Gin / Hertz** 搭建起了属于自己的 HTTP Web 服务。现在的你已经可以写出用户注册、登录、发布文章、提交表单的 API 接口了。如果我们将发送过来的数据存储在 Go 的 `map` 里，一旦服务重启后任务就没了。有的同学可能会想到[第三节课](https://github.com/LanshanTeam/Courseware-Backend-Go-2026/tree/main/lesson03)讲的，可以将数据写入文件之中，但同时也带来了一些问题，如果数据量过大，每次加载文件和查询数据时会异常缓慢。因此我们这节课学习通过数据库管理数据。

为了彻底解决数据的**安全持久化**、**高效并发访问**与**快速条件检索**，现代后端开发的核心基石——**数据库（Database）**应运而生。

数据库是一个专门用于组织、存储和管理数据的系统软件。它就像一个高度智能化的“数据仓库”，支持多人同时安全操作，并通过优化的索引结构保证海量数据毫秒级精准查询。

> **配套视频推荐**：  
> [数据库其实就那么回事 - B站科普视频](https://www.bilibili.com/video/BV1Hh411C7Up/)（前段为核心概念科普）  
> [什么是MySQL - B站视频](https://www.bilibili.com/video/BV1p5qhYsE4f/)  

---

### 数据库的分类

在后端体系中，数据库通常按照数据组织形式与一致性保证划分为两大阵营：

#### 1. 关系型数据库（RDBMS）
- **核心特征**：数据以严谨的**二维数据表（行与列）**形式组织，拥有固定的表结构定义（Schema）；
- **查询方式**：使用国际标准的 **SQL（结构化查询语言）**，支持复杂的跨表关联查询（JOIN）；
- **数据安全**：强力支持 **ACID 事务**，保证数据分毫不差；
- **适用场景**：用户账号、订单支付、账户资金、权限系统等核心商业数据；
- **工业代表**：***MySQL、PostgreSQL***。

#### 2. 非关系型数据库（NoSQL）
- **核心特征**：数据结构更加灵活自由（键值对、文档、列族、图结构等），通常不要求固定表结构；
- **优势与取舍**：放宽强事务约束，换取极致的吞吐性能与横向扩展能力，多数遵循最终一致性；
- **适用场景**：高性能缓存、热点排行榜、海量日志收集、临时会话（Session）；
- **工业代表**：***Redis、MongoDB***（下节课会讲）。

---

### MySQL 数据库

#### 部署

在本地开发中，我们**强烈推荐使用 Docker 容器化部署**，不仅环境隔离干净，而且随用随起、无需繁琐的系统服务配置。

- **使用 Docker 部署（推荐）**：

拉取官方镜像并一键后台启动：
```bash
docker pull mysql:8.0
docker run -d -p 3306:3306 --name mysql-dev -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
```

进入容器内部交互终端：
```bash
docker exec -it mysql-dev /bin/bash
```

使用内置命令行客户端登录：
```bash
mysql -u root -p
# 随后输入设置的密码 123456 即可成功进入
```

- **传统本地安装**：如需直接安装在操作系统中，可前往 [MySQL 官方下载页面](https://dev.mysql.com/downloads/mysql/)，参考 [Windows 安装配置指南](https://juejin.cn/post/7337186171210170383)。

#### 数据表结构

关系型数据库的数据存放在一张张二维表中。下面是一个典型的用户表（`users`）结构：

```sql
+--------------+-----------------+------+-----+---------+----------------+
| Field        | Type            | Null | Key | Default | Extra          |
+--------------+-----------------+------+-----+---------+----------------+
| id           | bigint unsigned | NO   | PRI | NULL    | auto_increment |
| created_at   | datetime(3)     | YES  |     | NULL    |                |
| updated_at   | datetime(3)     | YES  |     | NULL    |                |
| nick_name    | varchar(42)     | YES  |     | NULL    |                |
| gender       | varchar(12)     | YES  |     | NULL    |                |
| password     | varchar(258)    | YES  |     | NULL    |                |
| avatar       | varchar(256)    | YES  |     | NULL    |                |
| introduction | longtext        | YES  |     | NULL    |                |
| email        | varchar(128)    | YES  |     | NULL    |                |
| qq           | bigint          | YES  |     | NULL    |                |
| tel          | varchar(18)     | YES  |     | NULL    |                |
| birthday     | varchar(128)    | YES  |     | NULL    |                |
| role         | tinyint         | YES  |     | 1       |                |
| user_name    | varchar(42)     | YES  |     | NULL    |                |
+--------------+-----------------+------+-----+---------+----------------+
```

每一列的属性含义如下：
1. **Field（列名/字段名）**：当前列的标识名称；
2. **Type（数据类型）**：规定该字段存放什么格式的数据：详细内容可参考下表。
   
   | **类别**         | **数据类型**          | **说明**                             | **范围/长度**                                                |
   | ---------------- | --------------------- | ------------------------------------ | ------------------------------------------------------------ |
   | **数字类型**     | `TINYINT`             | 存储小范围的整数                     | 有符号：-128 到 127；无符号：0 到 255                        |
   |                  | `SMALLINT`            | 存储较小范围的整数                   | 有符号：-32,768 到 32,767；无符号：0 到 65,535               |
   |                  | `MEDIUMINT`           | 存储中等范围的整数                   | 有符号：-8,388,608 到 8,388,607；无符号：0 到 16,777,215     |
   |                  | `INT` / `INTEGER`     | 存储普通范围的整数                   | 有符号：-2,147,483,648 到 2,147,483,647；无符号：0 到 4,294,967,295 |
   |                  | `BIGINT`              | 存储大范围的整数                     | 有符号：-9,223,372,036,854,775,808 到 9,223,372,036,854,775,807；无符号：0 到 18,446,744,073,709,551,615 |
   |                  | `FLOAT`               | 存储单精度浮点数                     | 4 字节，精度通常为 7 位数字                                  |
   |                  | `DOUBLE`              | 存储双精度浮点数                     | 8 字节，精度通常为 15 位数字                                 |
   |                  | `DECIMAL` / `NUMERIC` | 存储定点数（高精度）                 | 定义时指定精度和小数位数，例如 `DECIMAL(10,2)`，表示最多 10 位，2 位小数 |
   | **字符串类型**   | `CHAR`                | 固定长度字符串                       | 1 到 255 字符                                                |
   |                  | `VARCHAR`             | 可变长度字符串                       | 1 到 65,535 字符                                             |
   |                  | `TEXT`                | 变长文本数据                         | 最多 65,535 字符                                             |
   |                  | `TINYTEXT`            | 变长文本数据                         | 最多 255 字符                                                |
   |                  | `MEDIUMTEXT`          | 变长文本数据                         | 最多 16,777,215 字符                                         |
   |                  | `LONGTEXT`            | 变长文本数据                         | 最多 4,294,967,295 字符                                      |
   | **日期时间类型** | `DATE`                | 存储日期                             | `YYYY-MM-DD`，范围：1000-01-01 到 9999-12-31                 |
   |                  | `DATETIME`            | 存储日期和时间                       | `YYYY-MM-DD HH:MM:SS`，范围：1000-01-01 00:00:00 到 9999-12-31 23:59:59 |
   |                  | `TIMESTAMP`           | 存储时间戳                           | `YYYY-MM-DD HH:MM:SS`，范围：1970-01-01 00:00:01 到 2038-01-19 03:14:07 (UTC) |
   |                  | `TIME`                | 存储时间                             | `HH:MM:SS`，范围：`-838:59:59` 到 `838:59:59`                |
   |                  | `YEAR`                | 存储年份                             | `YYYY`，范围：1901 到 2155                                   |
   | **布尔类型**     | `BOOLEAN` / `BOOL`    | 布尔值                               | 0（`FALSE`）或 1（`TRUE`）                                   |
   | **二进制类型**   | `BINARY`              | 固定长度的二进制数据                 | 1 到 255 字节                                                |
   |                  | `VARBINARY`           | 可变长度的二进制数据                 | 1 到 65,535 字节                                             |
   |                  | `BLOB`                | 二进制大对象                         | 最多 65,535 字节                                             |
   |                  | `TINYBLOB`            | 小型二进制大对象                     | 最多 255 字节                                                |
   |                  | `MEDIUMBLOB`          | 中型二进制大对象                     | 最多 16,777,215 字节                                         |
   |                  | `LONGBLOB`            | 大型二进制大对象                     | 最多 4,294,967,295 字节                                      |
   | **JSON 类型**    | `JSON`                | 存储 JSON 格式数据                   | 最多 4GB 数据                                                |
   | **集合类型**     | `ENUM`                | 枚举类型（限制为一组预定义值之一）   | 1 到 65,535 个预定义值                                       |
   |                  | `SET`                 | 集合类型（可存储多个预定义值的组合） | 1 到 64 个预定义值的组合                                     |
3. **Null**：`YES` 代表此列允许为空，`NO` 代表必填；
4. **Key（索引属性）**：
   - `PRI`：主键索引（Primary Key），整行数据的唯一凭据，全局非空且唯一；
   - `UNI`：唯一索引（Unique Key），该列数值不允许重复（如手机号、唯一账号）；
   - `MUL`：普通索引（Multiple），允许重复，专为加速条件筛选而建立。
5. **Default（默认值）**：插入数据没有赋值时自动填入的值；
6. **Extra（附加特性）**：如 `auto_increment` 表示数值会自动自增递增。

**插入数据后的实际数据表外观：**
```sql
+----+------------+-----------+------------------------+------+---------------------+
| id | user_name  | nick_name | email                  | role | created_at          |
+----+------------+-----------+------------------------+------+---------------------+
| 1  | xiaoming   | 小明      | xiaoming@lanshan.email | 1    | 2026-10-01 10:00:00 |
| 2  | xiaohong   | 小红      | xiaohong@lanshan.email | 1    | 2026-11-01 10:15:20 |
| 3  | admin      | 超管      | admin@lanshan.email    | 2    | 2026-11-13 08:30:00 |
| 5  | airiseina  | 绵眠一夏   | airiseina@example.com  | 1    | 2026-11-18 21:33:59 |
+----+------------+-----------+------------------------+------+---------------------+
```

---

### PostgreSQL 数据库

在当今后端技术演进中，另一个不可忽视且备受推崇的关系型数据库是 **PostgreSQL**（通常简称为 **PG** 或 **Postgres**）。

PostgreSQL 自诞生起就以“世界上功能最强大的开源关系型数据库”著称。在海外科技大厂、现代全栈项目以及当下的 **AI / 向量应用** 浪潮中，PostgreSQL 的份额正迅速攀升。

#### 部署

同样使用 Docker 一键启动：
```bash
docker pull postgres:16
docker run -d -p 5432:5432 --name postgres-dev -e POSTGRES_PASSWORD=123456 postgres:16
```

进入容器并使用 `psql` 交互工具连接：
```bash
docker exec -it postgres-dev psql -U postgres
```

> **psql 常用快捷指令速查**：
> - `\l`：列出当前所有的数据库；
> - `\c <dbname>`：切换连接到指定数据库；
> - `\dt`：查看当前数据库下的所有数据表；
> - `\d <tablename>`：查看指定数据表的表结构定义；
> - `\q`：安全退出 `psql` 命令行终端。

与 MySQL 相比，PostgreSQL 具备许多非常亮眼的现代化特性：
1. **严格遵循国际标准**：PG 极度遵循 ANSI SQL 国际标准，类型校验严谨，SQL 优化器对复杂的多表关联与窗口函数表现出色；
2. **原生支持丰富高级数据类型**：
   - **`JSONB`**：不仅能存 JSON，还是以二进制格式紧凑存储，**支持直接针对 JSON 内部键建立 GIN 高性能索引**，让关系型数据库兼备文档存储的威力；
   - 原生数组（`INT[]`、`TEXT[]`）、网络地址（`INET`）、UUID 等；
3. **主键自增机制差异**：MySQL 习惯使用 `AUTO_INCREMENT`，而 PG 经典使用 `SERIAL`（底层序列），现代语法也原生支持标准的 `GENERATED ALWAYS AS IDENTITY`；
4. **丰富的插件生态（Extensions）**：
   - **`pgvector`**：当前在 **大模型知识库（RAG）** 中极为火热。通过安装该扩展，PG 可以直接存储文本 Embedding 向量，并执行高效的余弦相似度检索，一套数据库兼顾传统业务与 AI 向量检索！

---

### SQL 核心操作

SQL（Structured Query Language）是程序员与关系型数据库沟通的标准语言。语法非常直观，贴近英文表达：

借用去年一个学长的图![](E:\Document\Github\Courseware-Backend-Go-2026\lesson06\img.png)

大家还是要花时间去好好学习sql啊！！！😭

> **推荐教程**：  
> [SQL 教程 - 菜鸟教程](https://www.runoob.com/sql/sql-tutorial.html)  
> [MySQL 基础与进阶 - 黑马程序员 (B站)](https://www.bilibili.com/video/BV1Kr4y1i7ru/)  

下面以学生表（`students`）为例，演示核心操作语法：

##### 1. 创建数据库与数据表（DDL）

```sql
-- 创建数据库
CREATE DATABASE lanshan_school;

-- MySQL 创建表 (自增使用 AUTO_INCREMENT)
CREATE TABLE students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    age INT DEFAULT 18,
    major VARCHAR(50),
    score INT
);

-- PostgreSQL 创建表 (自增使用 SERIAL)
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT DEFAULT 18,
    major VARCHAR(50),
    score INT
);
```

##### 2. 插入数据（INSERT）

```sql
-- 单条插入
INSERT INTO students (name, age, major, score)
VALUES ('张三', 20, 'Computer Science', 88);

-- 批量插入多条数据
INSERT INTO students (name, age, major, score)
VALUES 
    ('李四', 19, 'Software Engineering', 92),
    ('王五', 21, 'Mathematics', 79),
    ('赵六', 20, 'Physics', 85);
```

> **小技巧**：PostgreSQL 支持 `INSERT INTO ... RETURNING id;`，可以在插入成功的瞬间直接返回新生成的自增 ID，省去二次查询。

##### 3. 查询数据（SELECT）

```sql
-- 投影查询所有字段
SELECT * FROM students;

-- 条件筛选查询 (WHERE)
SELECT name, score FROM students WHERE age >= 20;

-- 复合逻辑与模糊匹配 (LIKE: % 表示匹配任意字符)
SELECT * FROM students WHERE major = 'Computer Science' AND score > 80;
SELECT * FROM students WHERE name LIKE '张%';

-- 排序与分页查询 (ORDER BY 排序, LIMIT 限制数量, OFFSET 偏移跳过行数)
-- 场景：按成绩由高到低排序，获取第 2 页的数据 (每页 2 条)
SELECT * FROM students 
ORDER BY score DESC 
LIMIT 2 OFFSET 2;
```

##### 4. 修改与删除（UPDATE / DELETE）

> ⚠️ **生产高危预警**：执行 `UPDATE` 和 `DELETE` 语句时，**务必先检查 `WHERE` 条件**！一旦遗漏 `WHERE`，整张表的所有数据都会被修改或清空！

```sql
-- 修改指定学生的成绩
UPDATE students 
SET score = 95 
WHERE name = '张三';

-- 删除指定记录
DELETE FROM students WHERE name = '赵六';

-- 清空整张表的所有数据 (TRUNCATE 比逐行 DELETE 效率高很多)
TRUNCATE TABLE students;
```

---

### GORM 框架实战

平时写代码时，我们时常会使用 **ORM（对象关系映射，Object-Relational Mapping）**。

ORM 将 Go 语言中的**结构体（Struct）**映射为数据库中的**数据表（Table）**，将结构体实例映射为表中的**行（Row）**。这不仅极大提升了开发效率，还能通过预编译占位符从根本上杜绝 SQL 注入风险。

Go 社区中最主流的 ORM 库便是 **GORM**（[官方文档中文站](https://gorm.io/zh_CN/docs/index.html)）。

> **配套视频推荐**：  
> [基于 Gin 和 GORM 的 Web 开发实战 - 七米出品 (B站)](https://www.bilibili.com/video/BV1gJ411p7xC/)  
> [【GORM教学】手把手带你入门GORM](https://www.bilibili.com/video/BV1E64y1472a/?spm_id_from=333.337.search-card.all.click&vd_source=c1ea60555dbef494273824f133c56fd9)  

#### 1. 安装依赖

```bash
# 安装 GORM 核心包
go get -u gorm.io/gorm

# 根据所选数据库安装驱动（支持同时安装）
go get -u gorm.io/driver/mysql
go get -u gorm.io/driver/postgres
```

#### 2. 初始化连接与连接池

```go
package main

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 连接 MySQL
func ConnectMySQL() (*gorm.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lanshan_school?charset=utf8mb4&parseTime=True&loc=Local"
	return openDB(mysql.Open(dsn))
}

// 连接 PostgreSQL
func ConnectPostgres() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123456 dbname=lanshan_school port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	return openDB(postgres.Open(dsn))
}

// 统一配置底层连接池
func openDB(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池调优：避免连接暴增或频繁握手
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大活跃连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最长时间

	return db, nil
}
```

#### 3. 业务模型定义（Model）

内嵌 `gorm.Model` 会自动赋予模型 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt`（支持软删除）：

```go
type Member struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	Password string `gorm:"size:255;not null"`
    Department string `gorm:"size:64;index;comment:部门名称"`
}
```

通过 `AutoMigrate` 自动建表：
```go
// 启动时自动将结构体映射同步到数据库中生成数据表
if err := db.AutoMigrate(&Article{}); err != nil {
	log.Fatal("数据库表迁移失败:", err)
}
```

#### 4. 数据操作实战演练（增删改查）

```go
// 准备数据
member1 := Member{Name: "thm", Age: 19, Department: "Backend-Go"}
member2 := &Member{Name: "whh", Age: 20, Department: "Backend-Go"}
members := []*Member{
    {Name: "xzy", Age: 19, Department: "Backend-Go"},
    {Name: "lmx", Age: 19, Department: "Backend-Go"},
    {Name: "chr", Age: 18, Department: "Backend-Go"},
    {Name: "sxf", Age: 19, Department: "Backend-Go"},
    {Name: "dhy", Age: 20, Department: "Backend-Go"},
    {Name: "ll", Age: 19, Department: "Backend-Go"},
    {Name: "lqz", Age: 20, Department: "Backend-Go"},
    {Name: "cxl", Age: 19, Department: "Backend-Java"},
    {Name: "lkg", Age: 19, Department: "Backend-Java"},
}

// =======================
//       增（Create）
// =======================
db.Create(&member1)
db.Create(member2)
db.Create(members)
log.Println("数据创建完成")

// =======================
//       查（Read）
// =======================

// 1. 查第一条记录
var m1 Member
db.First(&m1)
fmt.Println("First 查询:", m1)

// 2. 根据 ID 查询
var m2 Member
db.First(&m2, 3) // 查 ID = 3 的记录
fmt.Println("根据 ID 查询:", m2)

// 3. 条件查询（查 Backend-Go 的成员）
var goMembers []Member
db.Where("department = ?", "Backend-Go").Find(&goMembers)
fmt.Println("Backend-Go 成员:")
for _, m := range goMembers {
    fmt.Println("成员:", m.Name)
}
var age19plusMembers []Member
db.Where("age >= ?", 20).Find(&age19plusMembers)
fmt.Println("年龄大于等于20的成员:")
for _, m := range age19plusMembers {
    fmt.Println("成员:", m.Name)
}

// =======================
//       改（Update）
// =======================

// 1. 更新单个字段
db.Model(&m1).Update("Age", 20)

// 2. 更新多个字段
db.Model(&m1).Updates(Member{Name: "airiseina", Department: "Backend-Java"})

fmt.Println("更新完成:", m1)

// =======================
//       删（Delete）
// =======================

// 1. 根据 id 删除
db.Delete(&Member{}, 2) // 删除 ID = 2 的成员

// 2. 按条件删除
db.Where("department = ?", "Backend-Java").Delete(&Member{})

var javaMember Member
err = db.Where("department = ?", "Backend-Java").First(&javaMember).Error
if err != nil {
    fmt.Println(err.Error())
}

log.Println("删除操作完成")
}
```

#### 5. 事务（Transaction）

事务能够确保一系列数据操作具备 **ACID 特性**（原子性、一致性、隔离性、持久性），即“要么全部成功提交，要么全部回滚撤销”。

GORM 提供了闭包风格的事务 API，极简且安全：函数内部只要返回任意 `err != nil`，底层自动触发 `Rollback`；若返回 `nil`，则自动提交 `Commit`：

```go
func transactionExample(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 插入数据 1
		member1 := Member{Name: "md", Age: 20, Department: "Go", ID: 50001}
		if err := tx.Create(&member1).Error; err != nil {
			return err
		}

		// 插入数据 2（这里模拟失败）
		member2 := Member{Name: "md", Age: 21, Department: "Python", ID: 50001} // 假设唯一键冲突
		if err := tx.Create(&member2).Error; err != nil {
			return err
		}

		return nil // 没有错误则提交事务
	})

	if err != nil {
		log.Println("事务回滚：", err)
	} else {
		log.Println("事务提交成功")
	}
}
```

---

### 其他

#### 数据库管理工具

大家可以使用一些图形化客户端配合学习开发：
- **DBeaver**：[DBeaver 官网](https://dbeaver.io/)，跨平台，同时原生支持 MySQL、PostgreSQL、SQLite 等所有主流数据库；
- **DataGrip**：[DataGrip 官网](https://www.jetbrains.com/zh-cn/datagrip/)（JetBrains 家族专业工具，高校学生可免费认证）；
- **IDE 自带工具**：GoLand 右侧的 Database 窗口，或 VS Code 安装 Database Client 扩展。

---

## 课后作业

#### lv0

1. 动手练习手写 SQL（建表、增删改查、排序分页）

#### lv1

**完成一个通讯录 **

*基本要求*：使用`gorm`或者`sql语句`实现通讯录的的增删改查

*进阶要求*：

1. 结合上一节课的 Gin / Hertz 框架，将增删改查操作封装为标准的 RESTful HTTP API；
2. 支持用户注册与登录，注册时将用户信息存入数据库（**密码尽量通过哈希加密存储**）；
3. 用jwt进行鉴权，用户只能修改自己的通讯录；

#### lv2

1. 课程只是简单的展示了一下，可以去了解一下GORM 中模型之间的`一对一`、`一对多`和`多对多`关系。

---

要求不高，希望同学们还是去尝试完成以上的小项目，基本要求只要看懂这节课的示例代码就能写出来，进阶要求也尽量去完成一下。写好了代码可以提交到`tanghaomin@lanshan.email`，有不懂的地方可以来私聊我。
