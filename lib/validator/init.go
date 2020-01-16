package validator

import (
    "fmt"
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var GlobalValidate *validator.Validate

func init() {
    GlobalValidate = New()
}

func New() *validator.Validate {
    //验证器
    validate := validator.New()
    //中文翻译器
    zh_ch := zh.New()
    uni := ut.New(zh_ch)
    trans, _ := uni.GetTranslator("zh")
    if err := zh_translations.RegisterDefaultTranslations(validate, trans);err != nil {
        fmt.Println("register validate error:"+err.Error())
    }
    return validate
}
