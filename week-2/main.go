package main

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

var dao userDAO

func init() {
	var (
		err error
	)
	dao, err = NewUserDao()
	if err != nil {
		log.Fatalf("init dao err: %v \n", err)
	}
}

// 模拟dao
type userDAO struct {
}

func (d userDAO) findUserByName(name string) (string, error) {
	return "", sql.ErrNoRows
}

//我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
// dao不应该包装 sql.ErrNoRows 错误
//	如果这个查询应该有结果而数据库中没有数据，这应该被看做是一个错误，调用方处理后向上抛出
//	如果这只是一个普通的查询结果，没有数据是正常现象，我们直接可以处理掉这个 error

func main() {
	name := "wk"
	_, err := mustFindUserByName(name)
	if err != nil {
		log.Printf("[mustFindUserByName] err: %v\n", err)
		//return
	}

	_, err = findUserByName(name)
	if err != nil {
		log.Printf("[findUserByName] err: %v\n", err)
		//return
	}

}

// 如果应该有结果，将错误sql.ErrNoRows包装后抛出
func mustFindUserByName(name string) (string, error) {

	findUserByName, err := dao.findUserByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.Wrapf(err, "通过名字未找到对应用户. params: %s", name)
		}
		return "", errors.Wrapf(err, "其他错误. params: %s", name)
	}
	return findUserByName, nil
}

// 如果不是必须有结果，不需要特殊处理sql.ErrNoRows将错误包装后抛出
func findUserByName(name string) (string, error) {
	findUserByName, err := dao.findUserByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[findUserByName] 未找到用户: %v, 后续进行降级处理", name)
			// 降级处理
			return "降级处理后的结果", nil
		} else {
			return "", errors.Wrapf(err, "异常. params: %s", name)
		}
	}
	return findUserByName, nil
}

func NewUserDao() (userDAO, error) {
	return userDAO{}, nil
}
