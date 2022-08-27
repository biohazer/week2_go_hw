package homework

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

func InsertStmt(entity interface{}) (string, []interface{}, error) {

	// 解决测试1
	if entity == nil {
		return "", []interface{}{}, errInvalidEntity
	}

	val := reflect.ValueOf(entity)
	typ := val.Type()

	// 解决测试2
	if typ.Size()==0 {
		return "", []interface{}{}, errInvalidEntity
	}
	flag := 0
	// 解决测试3,4
	if typ.Kind() == reflect.Ptr {
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			val = val.Elem()
			flag += 1
		}
	}

	// 解决问题5
	if flag>1 {
		return "", []interface{}{}, errInvalidEntity
	}

	tabelName := typ.Name()
	num := typ.NumField()
	ls := []interface{}{} 
	colName := []string{}
	vals := []string{}
	for i := 0; i < num; i++ {
		fd := typ.Field(i)
		fdVal := val.Field(i)
		if fd.IsExported() {
			tmp := "`"+fd.Name+"`"
			colName = append(colName, tmp)
			ls = append(ls, fdVal.Interface())
			vals = append(vals, "?")
		} else {
			ls = append(ls, reflect.Zero(fd.Type).Interface())
		}
	}

	colName2 := strings.Join(colName, ",")
	vals2 := strings.Join(vals, ",")

	// INSERT INTO `BaseEntity`(`CreateTime`,`UpdateTime`) VALUES(?,?);"
	// INSERT INTO `BaseEntity`(`CreateTime`,`UpdateTime`) VALUES(<int64 Value>,<*int64 Value>);"

	// colName:= ([]string{"CreateTime","UpdateTime"})
	// col2 := ""
	// for _, v := range colName {
	// 	col2+=v
	// }
	// // vals := []string{"123","*456"}
	// vals2 := ""
	// for _, v := range vals {
	// 	vals2 += v
	// }
	str1 := "INSERT INTO `"+ tabelName+ "`("+colName2+")"+ " VALUES("+vals2+");"

	// 检测 entity 是否符合我们的要求
	// 我们只支持有限的几种输入

	// 使用 strings.Builder 来拼接 字符串
	// var str string = "INSERT INTO `BaseEntity`(`CreateTime`,`UpdateTime`) VALUES(?,?);"
	b := []byte(str1)
	var bd strings.Builder
	for _, s := range b {
		fmt.Fprint(&bd, s)
	}
	// 构造 INSERT INTO XXX，XXX 是你的表名，这里我们直接用结构体名字

	// 遍历所有的字段，构造出来的是 INSERT INTO XXX(col1, col2, col3)
	// 在这个遍历的过程中，你就可以把参数构造出来
	// 如果你打算支持组合，那么这里你要深入解析每一个组合的结构体
	// 并且层层深入进去

	// 拼接 VALUES，达成 INSERT INTO XXX(col1, col2, col3) VALUES

	// 再一次遍历所有的字段，要拼接成 INSERT INTO XXX(col1, col2, col3) VALUES(?,?,?)
	// 注意，在第一次遍历的时候我们就已经拿到了参数的值，所以这里就是简单拼接 ?,?,?
	
	// return bd.String(), ls, nil
	return str1, ls, nil
	// panic("implement me")

}

