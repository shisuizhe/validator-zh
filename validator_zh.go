package validatorzh

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(i interface{}) []string {
	validate := validator.New()
	// register mobile
	if err := validate.RegisterValidation("mobile", mobile); err != nil {
		fmt.Println(err)
		return nil
	}
	// register idcard
	if err := validate.RegisterValidation("idcard", idcard); err != nil {
		fmt.Println(err)
		return nil
	}
	// register label for better prompt
	validate.RegisterTagNameFunc(func(filed reflect.StructField) string {
		name := filed.Tag.Get("label")
		return name
	})

	// i18n
	e := en.New()
	universalTranslator := ut.New(e, e, zh.New())
	translator, found := universalTranslator.GetTranslator("zh")
	if found {
		err := zhTranslation.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			fmt.Println("register zh translotor failed")
			return nil
		}
	} else {
		fmt.Println("not found zh translator")
		return nil
	}

	// register mobile field translation
	validate.RegisterTranslation("mobile", translator, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0}格式错误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field(), fe.Field())
		return t
	})
	// register idcard field translation
	validate.RegisterTranslation("idcard", translator, func(ut ut.Translator) error {
		return ut.Add("idcard", "{0}格式错误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("idcard", fe.Field(), fe.Field())
		return t
	})

	err := validate.Struct(i)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			fmt.Println("invalid field")
			return nil
		}
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errMsgs := make([]string, 0)
			for _, fieldErr := range validationErrors {
				errMsg := fieldErr.Translate(translator)
				errMsgs = append(errMsgs, errMsg)
			}
			return errMsgs
		}
		fmt.Println("type assertion failed")
		return nil
	}
	return nil
}

// mobile 验证手机号码
func mobile(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^(13|14|15|17|18|19)[0-9]{9}$`, fl.Field().String())
	return ok
}

// idcard 验证身份证号码
func idcard(fl validator.FieldLevel) bool {
	id := fl.Field().String()

	var a1Map = map[int]int{
		0:  1,
		1:  0,
		2:  10,
		3:  9,
		4:  8,
		5:  7,
		6:  6,
		7:  5,
		8:  4,
		9:  3,
		10: 2,
	}

	var idStr = strings.ToUpper(string(id))
	var reg, err = regexp.Compile(`^[0-9]{17}[0-9X]$`)
	if err != nil {
		return false
	}
	if !reg.Match([]byte(idStr)) {
		return false
	}
	var sum int
	var signChar = ""
	for index, c := range idStr {
		var i = 18 - index
		if i != 1 {
			if v, err := strconv.Atoi(string(c)); err == nil {
				var weight = int(math.Pow(2, float64(i-1))) % 11
				sum += v * weight
			} else {
				return false
			}
		} else {
			signChar = string(c)
		}
	}
	var a1 = a1Map[sum%11]
	var a1Str = fmt.Sprintf("%d", a1)
	if a1 == 10 {
		a1Str = "X"
	}
	return a1Str == signChar
}
