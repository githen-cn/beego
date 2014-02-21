package orm

import (
	"fmt"
	"os"
	"strings"
	"time"

	// _ "github.com/bylevel/pq"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Id       int `orm:"auto"`
	Boolean  bool
	Char     string `orm:"size(50)"`
	Text     string
	Date     time.Time `orm:"type(date)"`
	DateTime time.Time
	Byte     byte
	Rune     rune
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	Uint     uint
	Uint8    uint8
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Float32  float32
	Float64  float64
	Decimal  float64 `orm:"digits(8);decimals(4)"`
}

type DataNull struct {
	Id       int       `orm:"auto"`
	Boolean  bool      `orm:"null"`
	Char     string    `orm:"size(50);null"`
	Text     string    `orm:"null"`
	Date     time.Time `orm:"type(date);null"`
	DateTime time.Time `orm:"null"`
	Byte     byte      `orm:"null"`
	Rune     rune      `orm:"null"`
	Int      int       `orm:"null"`
	Int8     int8      `orm:"null"`
	Int16    int16     `orm:"null"`
	Int32    int32     `orm:"null"`
	Int64    int64     `orm:"null"`
	Uint     uint      `orm:"null"`
	Uint8    uint8     `orm:"null"`
	Uint16   uint16    `orm:"null"`
	Uint32   uint32    `orm:"null"`
	Uint64   uint64    `orm:"null"`
	Float32  float32   `orm:"null"`
	Float64  float64   `orm:"null"`
	Decimal  float64   `orm:"digits(8);decimals(4);null"`
}

type User struct {
	Id         int    `orm:"auto"`
	UserName   string `orm:"size(30);unique"`
	Email      string `orm:"size(100)"`
	Password   string `orm:"size(100)"`
	Status     int16
	IsStaff    bool
	IsActive   bool      `orm:"default(1)"`
	Created    time.Time `orm:"auto_now_add;type(date)"`
	Updated    time.Time `orm:"auto_now"`
	Profile    *Profile  `orm:"null;rel(one);on_delete(set_null)"`
	Posts      []*Post   `orm:"reverse(many)" json:"-"`
	ShouldSkip string    `orm:"-"`
}

func NewUser() *User {
	obj := new(User)
	return obj
}

type Profile struct {
	Id    int     `orm:"auto"`
	Age   int16   ``
	Money float64 ``
	User  *User   `orm:"reverse(one)" json:"-"`
}

func (u *Profile) TableName() string {
	return "user_profile"
}

func NewProfile() *Profile {
	obj := new(Profile)
	return obj
}

type Post struct {
	Id      int       `orm:"auto"`
	User    *User     `orm:"rel(fk)"` //
	Title   string    `orm:"size(60)"`
	Content string    ``
	Created time.Time `orm:"auto_now_add"`
	Updated time.Time `orm:"auto_now"`
	Tags    []*Tag    `orm:"rel(m2m)"`
}

func NewPost() *Post {
	obj := new(Post)
	return obj
}

type Tag struct {
	Id    int     `orm:"auto"`
	Name  string  `orm:"size(30)"`
	Posts []*Post `orm:"reverse(many)" json:"-"`
}

func NewTag() *Tag {
	obj := new(Tag)
	return obj
}

type Comment struct {
	Id      int       `orm:"auto"`
	Post    *Post     `orm:"rel(fk)"`
	Content string    ``
	Parent  *Comment  `orm:"null;rel(fk)"`
	Created time.Time `orm:"auto_now_add"`
}

func NewComment() *Comment {
	obj := new(Comment)
	return obj
}

var DBARGS = struct {
	Driver string
	Source string
	Debug  string
}{
	os.Getenv("ORM_DRIVER"),
	os.Getenv("ORM_SOURCE"),
	os.Getenv("ORM_DEBUG"),
}

var (
	IsMysql    = DBARGS.Driver == "mysql"
	IsSqlite   = DBARGS.Driver == "sqlite3"
	IsPostgres = DBARGS.Driver == "postgres"
)

var dORM Ormer

var initSQLs = map[string]string{
	"mysql": "DROP TABLE IF EXISTS `user_profile`;\n" +
		"DROP TABLE IF EXISTS `user`;\n" +
		"DROP TABLE IF EXISTS `post`;\n" +
		"DROP TABLE IF EXISTS `tag`;\n" +
		"DROP TABLE IF EXISTS `post_tags`;\n" +
		"DROP TABLE IF EXISTS `comment`;\n" +
		"DROP TABLE IF EXISTS `data`;\n" +
		"DROP TABLE IF EXISTS `data_null`;\n" +
		"CREATE TABLE `user_profile` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `age` smallint NOT NULL,\n" +
		"    `money` double precision NOT NULL\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `user` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `user_name` varchar(30) NOT NULL UNIQUE,\n" +
		"    `email` varchar(100) NOT NULL,\n" +
		"    `password` varchar(100) NOT NULL,\n" +
		"    `status` smallint NOT NULL,\n" +
		"    `is_staff` bool NOT NULL,\n" +
		"    `is_active` bool NOT NULL,\n" +
		"    `created` date NOT NULL,\n" +
		"    `updated` datetime NOT NULL,\n" +
		"    `profile_id` integer\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `post` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `user_id` integer NOT NULL,\n" +
		"    `title` varchar(60) NOT NULL,\n" +
		"    `content` longtext NOT NULL,\n" +
		"    `created` datetime NOT NULL,\n" +
		"    `updated` datetime NOT NULL\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `tag` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `name` varchar(30) NOT NULL\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `post_tags` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `post_id` integer NOT NULL,\n" +
		"    `tag_id` integer NOT NULL,\n" +
		"    UNIQUE (`post_id`, `tag_id`)\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `comment` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `post_id` integer NOT NULL,\n" +
		"    `content` longtext NOT NULL,\n" +
		"    `parent_id` integer,\n" +
		"    `created` datetime NOT NULL\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `data` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `boolean` bool NOT NULL,\n" +
		"    `char` varchar(50) NOT NULL,\n" +
		"    `text` longtext NOT NULL,\n" +
		"    `date` date NOT NULL,\n" +
		"    `date_time` datetime NOT NULL,\n" +
		"    `byte` tinyint unsigned NOT NULL,\n" +
		"    `rune` integer NOT NULL,\n" +
		"    `int` integer NOT NULL,\n" +
		"    `int8` tinyint NOT NULL,\n" +
		"    `int16` smallint NOT NULL,\n" +
		"    `int32` integer NOT NULL,\n" +
		"    `int64` bigint NOT NULL,\n" +
		"    `uint` integer unsigned NOT NULL,\n" +
		"    `uint8` tinyint unsigned NULL,\n" +
		"    `uint16` smallint unsigned NOT NULL,\n" +
		"    `uint32` integer unsigned NOT NULL,\n" +
		"    `uint64` bigint unsigned NOT NULL,\n" +
		"    `float32` double precision NOT NULL,\n" +
		"    `float64` double precision NOT NULL,\n" +
		"    `decimal` numeric(8,4) NOT NULL\n" +
		") ENGINE=INNODB;\n" +
		"CREATE TABLE `data_null` (\n" +
		"    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n" +
		"    `boolean` bool,\n" +
		"    `char` varchar(50),\n" +
		"    `text` longtext,\n" +
		"    `date` date,\n" +
		"    `date_time` datetime,\n" +
		"    `byte` tinyint unsigned,\n" +
		"    `rune` integer,\n" +
		"    `int` integer,\n" +
		"    `int8` tinyint,\n" +
		"    `int16` smallint,\n" +
		"    `int32` integer,\n" +
		"    `int64` bigint,\n" +
		"    `uint` integer unsigned,\n" +
		"    `uint8` tinyint unsigned,\n" +
		"    `uint16` smallint unsigned,\n" +
		"    `uint32` integer unsigned,\n" +
		"    `uint64` bigint unsigned,\n" +
		"    `float32` double precision,\n" +
		"    `float64` double precision,\n" +
		"    `decimal` numeric(8,4)\n" +
		") ENGINE=INNODB;\n" +
		"CREATE INDEX `user_141c6eec` ON `user` (`profile_id`);\n" +
		"CREATE INDEX `post_fbfc09f1` ON `post` (`user_id`);\n" +
		"CREATE INDEX `comment_699ae8ca` ON `comment` (`post_id`);\n" +
		"CREATE INDEX `comment_63f17a16` ON `comment` (`parent_id`);",

	"sqlite3": `
DROP TABLE IF EXISTS "user_profile";
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "post";
DROP TABLE IF EXISTS "tag";
DROP TABLE IF EXISTS "post_tags";
DROP TABLE IF EXISTS "comment";
DROP TABLE IF EXISTS "data";
DROP TABLE IF EXISTS "data_null";
CREATE TABLE "user_profile" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "age" smallint NOT NULL,
    "money" real NOT NULL
);
CREATE TABLE "user" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_name" varchar(30) NOT NULL UNIQUE,
    "email" varchar(100) NOT NULL,
    "password" varchar(100) NOT NULL,
    "status" smallint NOT NULL,
    "is_staff" bool NOT NULL,
    "is_active" bool NOT NULL,
    "created" date NOT NULL,
    "updated" datetime NOT NULL,
    "profile_id" integer
);
CREATE TABLE "post" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" integer NOT NULL,
    "title" varchar(60) NOT NULL,
    "content" text NOT NULL,
    "created" datetime NOT NULL,
    "updated" datetime NOT NULL
);
CREATE TABLE "tag" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" varchar(30) NOT NULL
);
CREATE TABLE "post_tags" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "post_id" integer NOT NULL,
    "tag_id" integer NOT NULL,
    UNIQUE ("post_id", "tag_id")
);
CREATE TABLE "comment" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "post_id" integer NOT NULL,
    "content" text NOT NULL,
    "parent_id" integer,
    "created" datetime NOT NULL
);
CREATE TABLE "data" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "boolean" bool NOT NULL,
    "char" varchar(50) NOT NULL,
    "text" text NOT NULL,
    "date" date NOT NULL,
    "date_time" datetime NOT NULL,
    "byte" tinyint unsigned NOT NULL,
    "rune" integer NOT NULL,
    "int" integer NOT NULL,
    "int8" tinyint NOT NULL,
    "int16" smallint NOT NULL,
    "int32" integer NOT NULL,
    "int64" bigint NOT NULL,
    "uint" integer unsigned NOT NULL,
    "uint8" tinyint unsigned NOT NULL,
    "uint16" smallint unsigned NOT NULL,
    "uint32" integer unsigned NOT NULL,
    "uint64" bigint unsigned NOT NULL,
    "float32" real NOT NULL,
    "float64" real NOT NULL,
    "decimal" decimal
);
CREATE TABLE "data_null" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "boolean" bool,
    "char" varchar(50),
    "text" text,
    "date" date,
    "date_time" datetime,
    "byte" tinyint unsigned,
    "rune" integer,
    "int" integer,
    "int8" tinyint,
    "int16" smallint,
    "int32" integer,
    "int64" bigint,
    "uint" integer unsigned,
    "uint8" tinyint unsigned,
    "uint16" smallint unsigned,
    "uint32" integer unsigned,
    "uint64" bigint unsigned,
    "float32" real,
    "float64" real,
    "decimal" decimal
);
CREATE INDEX "user_141c6eec" ON "user" ("profile_id");
CREATE INDEX "post_fbfc09f1" ON "post" ("user_id");
CREATE INDEX "comment_699ae8ca" ON "comment" ("post_id");
CREATE INDEX "comment_63f17a16" ON "comment" ("parent_id");
`,

	"postgres": `
DROP TABLE IF EXISTS "user_profile";
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "post";
DROP TABLE IF EXISTS "tag";
DROP TABLE IF EXISTS "post_tags";
DROP TABLE IF EXISTS "comment";
DROP TABLE IF EXISTS "data";
DROP TABLE IF EXISTS "data_null";
CREATE TABLE "user_profile" (
    "id" serial NOT NULL PRIMARY KEY,
    "age" smallint NOT NULL,
    "money" double precision NOT NULL
);
CREATE TABLE "user" (
    "id" serial NOT NULL PRIMARY KEY,
    "user_name" varchar(30) NOT NULL UNIQUE,
    "email" varchar(100) NOT NULL,
    "password" varchar(100) NOT NULL,
    "status" smallint NOT NULL,
    "is_staff" boolean NOT NULL,
    "is_active" boolean NOT NULL,
    "created" date NOT NULL,
    "updated" timestamp with time zone NOT NULL,
    "profile_id" integer
);
CREATE TABLE "post" (
    "id" serial NOT NULL PRIMARY KEY,
    "user_id" integer NOT NULL,
    "title" varchar(60) NOT NULL,
    "content" text NOT NULL,
    "created" timestamp with time zone NOT NULL,
    "updated" timestamp with time zone NOT NULL
);
CREATE TABLE "tag" (
    "id" serial NOT NULL PRIMARY KEY,
    "name" varchar(30) NOT NULL
);
CREATE TABLE "post_tags" (
    "id" serial NOT NULL PRIMARY KEY,
    "post_id" integer NOT NULL,
    "tag_id" integer NOT NULL,
    UNIQUE ("post_id", "tag_id")
);
CREATE TABLE "comment" (
    "id" serial NOT NULL PRIMARY KEY,
    "post_id" integer NOT NULL,
    "content" text NOT NULL,
    "parent_id" integer,
    "created" timestamp with time zone NOT NULL
);
CREATE TABLE "data" (
    "id" serial NOT NULL PRIMARY KEY,
    "boolean" bool NOT NULL,
    "char" varchar(50) NOT NULL,
    "text" text NOT NULL,
    "date" date NOT NULL,
    "date_time" timestamp with time zone NOT NULL,
    "byte" smallint CHECK("byte" >= 0 AND "byte" <= 255) NOT NULL,
    "rune" integer NOT NULL,
    "int" integer NOT NULL,
    "int8" smallint CHECK("int8" >= -127 AND "int8" <= 128) NOT NULL,
    "int16" smallint NOT NULL,
    "int32" integer NOT NULL,
    "int64" bigint NOT NULL,
    "uint" bigint CHECK("uint" >= 0) NOT NULL,
    "uint8" smallint CHECK("uint8" >= 0 AND "uint8" <= 255) NOT NULL,
    "uint16" integer CHECK("uint16" >= 0) NOT NULL,
    "uint32" bigint CHECK("uint32" >= 0) NOT NULL,
    "uint64" bigint CHECK("uint64" >= 0) NOT NULL,
    "float32" double precision NOT NULL,
    "float64" double precision NOT NULL,
    "decimal" numeric(8, 4)
);
CREATE TABLE "data_null" (
    "id" serial NOT NULL PRIMARY KEY,
    "boolean" bool,
    "char" varchar(50),
    "text" text,
    "date" date,
    "date_time" timestamp with time zone,
    "byte" smallint CHECK("byte" >= 0 AND "byte" <= 255),
    "rune" integer,
    "int" integer,
    "int8" smallint CHECK("int8" >= -127 AND "int8" <= 128),
    "int16" smallint,
    "int32" integer,
    "int64" bigint,
    "uint" bigint CHECK("uint" >= 0),
    "uint8" smallint CHECK("uint8" >= 0 AND "uint8" <= 255),
    "uint16" integer CHECK("uint16" >= 0),
    "uint32" bigint CHECK("uint32" >= 0),
    "uint64" bigint CHECK("uint64" >= 0),
    "float32" double precision,
    "float64" double precision,
    "decimal" numeric(8, 4)
);
CREATE INDEX "user_profile_id" ON "user" ("profile_id");
CREATE INDEX "post_user_id" ON "post" ("user_id");
CREATE INDEX "comment_post_id" ON "comment" ("post_id");
CREATE INDEX "comment_parent_id" ON "comment" ("parent_id");
`}

func init() {
	// err := os.Setenv("TZ", "+00:00")
	// fmt.Println(err)

	RegisterModel(new(Data), new(DataNull))
	RegisterModel(new(User))
	RegisterModel(new(Profile))
	RegisterModel(new(Post))
	RegisterModel(new(Tag))
	RegisterModel(new(Comment))

	Debug, _ = StrTo(DBARGS.Debug).Bool()

	if DBARGS.Driver == "" || DBARGS.Source == "" {
		fmt.Println(`need driver and source!

Default DB Drivers.

  driver: url
   mysql: https://github.com/go-sql-driver/mysql
 sqlite3: https://github.com/mattn/go-sqlite3
postgres: https://github.com/lib/pq

eg: mysql
ORM_DRIVER=mysql ORM_SOURCE="root:root@/my_db?charset=utf8" go test github.com/astaxie/beego/orm
`)
		os.Exit(2)
	}

	RegisterDataBase("default", DBARGS.Driver, DBARGS.Source, 20)

	BootStrap()

	dORM = NewOrm()

	queries := strings.Split(initSQLs[DBARGS.Driver], ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}
		_, err := dORM.Raw(query).Exec()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}
