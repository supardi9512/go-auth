package libraries

import (
	"database/sql"
	"go-auth/config"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validation struct {
	conn *sql.DB
}

func NewValidation() *Validation {
	conn, err := config.DBConn()

	if err != nil {
		panic(err)
	}

	return &Validation{
		conn: conn,
	}
}

func (v *Validation) Init() (*validator.Validate, ut.Translator) {
	// get translator package
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()

	// register default translation (english)
	en_translations.RegisterDefaultTranslations(validate, trans)

	// change default label
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")

		return labelName
	})

	// change default text for required
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} cannot be empty", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// custom validation for unique value
	validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		split_params := strings.Split(params, "-")
		tableName := split_params[0]
		fieldName := split_params[1]
		fieldValue := fl.Field().String()

		return v.checkIsUnique(tableName, fieldName, fieldValue)
	})

	// change default text for isunique
	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} is already in use", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isunique", fe.Field())
		return t
	})

	return validate, trans
}

func (v *Validation) Struct(s interface{}) interface{} {
	validate, trans := v.Init()

	vErrors := make(map[string]interface{})

	err := validate.Struct(s)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			vErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}

// check is unique function
func (v *Validation) checkIsUnique(tableName, fieldName, fieldValue string) bool {

	row, _ := v.conn.Query("SELECT "+fieldName+" FROM "+tableName+" WHERE "+fieldName+" = ?", fieldValue)

	defer row.Close()

	var result string

	for row.Next() {
		row.Scan(&result)
	}

	return result != fieldValue
}
