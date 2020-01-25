package devicePack

import (
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
)

func init() {
	govalidator.TagMap["uuid"] = govalidator.Validator(func(str string) bool {
		_, err := uuid.FromString(str)
		return err == nil
	})
}

func ValidateDevice(u *Device) bool {
	result, err := govalidator.ValidateStruct(u)
	if err != nil {
		println("error: " + err.Error())
	}

	return result
}
