package permissions

import (
	"errors"
	"gitee.com/quant1x/engine/models"
)

var (
	ErrAlreadyExists = errors.New("the validator already exists") // 权限验证已经存在
)

type Validator func(model models.Strategy) error

var (
	validatePermission Validator = nil
)

// RegisterValidatePermission 注册权限验证模块
func RegisterValidatePermission(f Validator) error {
	if validatePermission != nil {
		return ErrAlreadyExists
	}
	validatePermission = f
	return nil
}

// CheckPermission 权限验证
func CheckPermission(model models.Strategy) error {
	if validatePermission == nil {
		// 没有权限验证, 直接返回成功
		return nil
	}
	return validatePermission(model)
}
