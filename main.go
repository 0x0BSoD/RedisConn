package main

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type RedisConn struct {
	Port int
	Adders string
	pool *redis.Pool
}

func (r *RedisConn) Init() {
	r.pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", r.Adders, r.Port))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func (r *RedisConn) Close() error {
	if r.pool != nil {
		err := r.pool.Close()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("pool is nil")
}

func (r *RedisConn) Check(keys ...interface{}) (interface{}, error) {
	data, err := r.doAction("EXISTS", keys)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RedisConn) Set(key, value string) error {
	data, err := r.doAction("SET", key, value)
	if err != nil {
		return err
	}
	// cast to string
	if data.(string) != "OK" {
		return errors.New(fmt.Sprintf("set failed, %s", data.(string)))
	}
	return nil
}

func (r *RedisConn) initCon() (redis.Conn, error) {
	if r.pool != nil {
		return r.pool.Get(), nil
	}
	return nil, errors.New("pool is nil")
}

func (r *RedisConn) doAction(action string, params ...interface{}) (interface{}, error) {
	if r.pool != nil {
		c, err := r.initCon()
		if err != nil {
			return nil, err
		}

		defer func() {
			err := c.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		data, err := c.Do(action, params)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("pool is nil")
}

func (r *RedisConn) clientClose() error {
	if r.pool != nil {
		err := r.pool.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("pool is nil")
}

func main() {
	var rc RedisConn
	rc.Adders = ""
	rc.Port = 6379
	rc.Init()

	defer func() {
		err := rc.Close()
		if err != nil {
			panic(err)
		}
	}()
}