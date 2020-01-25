package userPack

import (
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
)

func init() {
	govalidator.TagMap["duck"] = govalidator.Validator(func(str string) bool {
		return str == "duck"
	})

	govalidator.TagMap["uuid"] = govalidator.Validator(func(str string) bool {
		_, err := uuid.FromString(str)
		return err == nil
	})
}

func ValidateUser(u *User) bool {
	result, err := govalidator.ValidateStruct(u)
	if err != nil {
		println("error: " + err.Error())
	}

	return result
}

func SignValidate(u *User) (bool, error) {
	return govalidator.ValidateStruct(u)
}
