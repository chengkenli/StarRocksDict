/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package tools
 *@file    Try
 *@date    2024/11/13 18:04
 */

package tools

import (
	"StarRocksDict/util"
	"context"
	"fmt"
	"time"
)

// RetryOptionV2 配置选项函数
type RetryOptionV2 func(retry *RetryV2)

// RetryFunc 不带返回值的重试函数
type RetryFunc func(ctx context.Context) error

// RetryV2 重试类
type RetryV2 struct {
	interval time.Duration // 重试的间隔时长
	attempts int           // 重试次数
	timeout  time.Duration // 超时时间
}

// ObjectRetry 构造函数（实现当函数超时时，那么断开重试）
func ObjectRetry(opts ...RetryOptionV2) *RetryV2 {
	retry := &RetryV2{}
	for _, opt := range opts {
		opt(retry)
	}
	return retry
}

// WithInterval 设置重试间隔
func WithInterval(interval time.Duration) RetryOptionV2 {
	return func(retry *RetryV2) {
		retry.interval = interval
	}
}

// WithAttempts 设置重试次数
func WithAttempts(attempts int) RetryOptionV2 {
	return func(retry *RetryV2) {
		retry.attempts = attempts
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) RetryOptionV2 {
	return func(retry *RetryV2) {
		retry.timeout = timeout
	}
}

// Do 执行重试逻辑
func (r *RetryV2) Do(execFunc RetryFunc) error {
	for i := 0; i < r.attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()

		err := execFunc(ctx)
		if err == nil {
			return nil
		}
		// 检查错误是否是因为超时
		util.Loggrs.Warn(ctx.Err().Error())
		if ctx.Err() == context.DeadlineExceeded {
			util.Loggrs.Warn(fmt.Sprintf("Attempt %d: timeout, retrying...", i+1))
			time.Sleep(r.interval)
			continue
		}
		// 其他错误不重试
		return err
	}
	return fmt.Errorf("after %d attempts", r.attempts)
}

// 以下是使用示例
func test() {
	// 重构超时重试函数
	retry := ObjectRetry(
		WithInterval(1*time.Second),
		WithAttempts(3),
		WithTimeout(5*time.Second),
	)
	err := retry.Do(func(ctx context.Context) error {
		done := make(chan struct{})
		go func() {
			// 这里放置你的长时间运行函数
			//fun(){}
			fmt.Println("all tasks done.")
			done <- struct{}{}
		}()
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
}
