package validdate

import (
    "errors"
    "fmt"
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
    GlobalValidate *validator.Validate
    Trans ut.Translator
)

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
    }else {
        Trans = trans
    }
    return validate
}

func Struct(i interface{}) error {
    if GlobalValidate == nil {
        return errors.New("GlobalValidate need init")
    }
    if err := GlobalValidate.Struct(i);err != nil {
        for _,paramErr := range err.(validator.ValidationErrors) {
            return errors.New(paramErr.Translate(Trans))
        }
    }
    return nil
}
