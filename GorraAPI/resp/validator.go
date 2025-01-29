package resp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"reflect"
	"strings"
)

func zHTrans(locale string) (ut.Translator, error) {
	var err error
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		trans, ok := uni.GetTranslator(locale)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) err", locale)
		}

		switch locale {
		case "zh":
			err = zh_translations.RegisterDefaultTranslations(v, trans)
		case "en":
			err = en_translations.RegisterDefaultTranslations(v, trans)
		default:
			err = en_translations.RegisterDefaultTranslations(v, trans)
		}

		return trans, nil
	}

	return nil, err
}

func HandleValidatorError(err error, c *gin.Context) {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "绑定参数出错",
		})
		return
	}

	zH, _ := zHTrans("zh")

	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(zH)),
	})
}

func removeTopStruct(fields map[string]string) map[string]string {
	resp := map[string]string{}
	for field, err := range fields {
		resp[field[strings.Index(field, ".")+1:]] = err
	}
	return resp
}
