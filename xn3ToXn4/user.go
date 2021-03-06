package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/golib"
	"log"
)

type user struct {
	db3str,
	db4str dbstr
	fields userFields
}

type userFields struct {
	uid, gid, email, username, password, salt, threads, posts, credits, create_ip, create_date, avatar string
}

func (this *user) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"user") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "user 失败: " + err.Error())
	}

	fmt.Printf("转换 %suser 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *user) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "uid,gid,email,username,password,salt,threads,posts,credits,create_ip,create_date,avatar"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
	xn3 := fmt.Sprintf("SELECT %s FROM %suser", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %suser (%s) VALUES (%s)", xn4pre, fields, qmark)

	data, err := xiuno3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4Clear := "TRUNCATE `" + xn4pre + "user`"
	_, err = xiuno4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %suser 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %suser 表成功\r\n", xn4pre)

	tx, err := xiuno4db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(xn4)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %suser 表\r\n", xn4pre)

	var field userFields
	for data.Next() {
		err = data.Scan(
			&field.uid,
			&field.gid,
			&field.email,
			&field.username,
			&field.password,
			&field.salt,
			&field.threads,
			&field.posts,
			&field.credits,
			&field.create_ip,
			&field.create_date,
			&field.avatar)

		_, err = stmt.Exec(
			&field.uid,
			&field.gid,
			&field.email,
			&field.username,
			&field.password,
			&field.salt,
			&field.threads,
			&field.posts,
			&field.credits,
			&field.create_ip,
			&field.create_date,
			&field.avatar)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条 user", count), 0)
		}
	}
	defer data.Close()

	if err = data.Err(); err != nil {
		log.Fatalln("user error: ", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return count, err
}
