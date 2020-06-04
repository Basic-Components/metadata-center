package errs

import (
	"errors"
)

// ErrConfigParams 配置参数错误
var ErrConfigParams = errors.New("config params error")

// ErrAlgoType 算法不支持
var ErrAlgoType = errors.New("unknown algo type key")

// ErrExpOutOfRange 过期时间超出范围
var ErrExpOutOfRange = errors.New("exp is out of range")

// ErrNotFindMatchKey 未找到key
var ErrNotFindMatchKey = errors.New("not find match key")
