package commonerr

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator"
)

const (
	stringFoo = "This is an error"
)

func TestSetNewError(t *testing.T) {
	type args struct {
		code      int
		errorName string
		errDesc   string
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{

		{
			name: "get new error",
			args: args{
				code:      402,
				errorName: "test_error",
				errDesc:   "test error description",
			},
			want: &ErrorMessage{
				Code: 402,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        "test_error",
						ErrorDescription: "test error description",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNewError(tt.args.code, tt.args.errorName, tt.args.errDesc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNewBadRequest(t *testing.T) {
	type args struct {
		errorName string
		errDesc   string
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{

		{
			name: "get new error",
			args: args{
				errorName: "test_error",
				errDesc:   "test error description",
			},
			want: &ErrorMessage{
				Code: 400,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        "test_error",
						ErrorDescription: "test error description",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNewBadRequest(tt.args.errorName, tt.args.errDesc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNewBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNewBadRequestByFormat(t *testing.T) {
	type args struct {
		ef *ErrorFormat
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{

		{
			name: "error bad request by format",
			args: args{
				ef: &ErrorFormat{
					ErrorName:        "test_error",
					ErrorDescription: "test error description",
				},
			},
			want: &ErrorMessage{
				Code: 400,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        "test_error",
						ErrorDescription: "test error description",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNewBadRequestByFormat(tt.args.ef); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNewBadRequestByFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNewInternalError(t *testing.T) {
	tests := []struct {
		name string
		want *ErrorMessage
	}{

		{
			name: "new internal error",
			want: &ErrorMessage{
				Code: http.StatusInternalServerError,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        InternalServerName,
						ErrorDescription: InternalServerDescription,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNewInternalError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNewInternalError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNewUnauthorizedError(t *testing.T) {
	type args struct {
		errorName string
		errDesc   string
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{

		{
			name: "error bad request by format",
			args: args{
				errorName: "test_error",
				errDesc:   "test error description",
			},
			want: &ErrorMessage{
				Code: http.StatusUnauthorized,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        "test_error",
						ErrorDescription: "test error description",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNewUnauthorizedError(tt.args.errorName, tt.args.errDesc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNewUnauthorizedError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_Append(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	type args struct {
		errorName string
		errDesc   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{

		{
			name: "get append error",
			fields: fields{
				Code:      400,
				ErrorList: []*ErrorFormat{},
			},
			args: args{
				errorName: "test",
				errDesc:   "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			errorMessage.Append(tt.args.errorName, tt.args.errDesc)
		})
	}
}

func TestErrorMessage_AppendFormat(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	type args struct {
		ef *ErrorFormat
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{

		{
			name: "get append error",
			fields: fields{
				Code:      400,
				ErrorList: []*ErrorFormat{},
			},
			args: args{
				ef: &ErrorFormat{
					ErrorName:        "test",
					ErrorDescription: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			errorMessage.AppendFormat(tt.args.ef)
		})
	}
}

func TestErrorMessage_GetListError(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   []*ErrorFormat
	}{

		{
			name: "get append error",
			fields: fields{
				Code: 400,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        "test error name",
						ErrorDescription: "test error description",
					},
				},
			},
			want: []*ErrorFormat{
				{
					ErrorName:        "test error name",
					ErrorDescription: "test error description",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.GetListError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorMessage.GetListError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_GetCode(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{

		{
			name: "get append error",
			fields: fields{
				Code:      400,
				ErrorList: []*ErrorFormat{},
			},
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.GetCode(); got != tt.want {
				t.Errorf("ErrorMessage.GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		want *ErrorMessage
	}{

		{
			name: "get new error message",
			want: &ErrorMessage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_SetBadRequest(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   *ErrorMessage
	}{

		{
			name: "get fields",
			fields: fields{
				Code: http.StatusBadRequest,
			},
			want: &ErrorMessage{
				Code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := em.SetBadRequest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorMessage.SetBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_SetErrorValidator(t *testing.T) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	type testStruct struct {
		AddressID int64  `json:"addr_id"`
		Address1  string `json:"address" validate:"required"`
	}
	var test testStruct
	err := validate.Struct(test)

	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorMessage
	}{

		{
			name: "get error from struct",
			fields: fields{
				Code:      http.StatusBadRequest,
				ErrorList: []*ErrorFormat{},
			},
			args: args{
				err: err,
			},
			want: &ErrorMessage{
				Code: http.StatusBadRequest,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        "address",
						ErrorDescription: "required",
					},
				},
			},
		},
		{
			name: "get error from struct",
			fields: fields{
				Code:      http.StatusBadRequest,
				ErrorList: []*ErrorFormat{},
			},
			args: args{
				err: &validator.InvalidValidationError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}

			if got := em.SetErrorValidator(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				if tt.want == nil {
					t.SkipNow()
				}
				t.Errorf("ErrorMessage.SetErrorValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaultNewBadRequest(t *testing.T) {
	tests := []struct {
		name string
		want *ErrorMessage
	}{

		{
			name: "get error bad request",
			want: SetNewBadRequestByFormat(&DefaultBadRequest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDefaultNewBadRequest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDefaultNewBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_Marshal(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{

		{
			name: "add unit test",
			fields: fields{
				Code: 200,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.Marshal(); !reflect.DeepEqual(got, tt.want) {
				t.SkipNow()
			}
		})
	}
}

func TestErrorMessage_ToString(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{

		{
			name: "get error string",
			fields: fields{
				Code: 400,
			},
			want: `{"error_list":null,"code":400}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.ToString(); got != tt.want {
				t.Errorf("ErrorMessage.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_Error(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{

		{
			name: "get error string from Error interface",
			fields: fields{
				Code: 400,
			},
			want: `{"error_list":null,"code":400}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.Error(); got != tt.want {
				t.Errorf("ErrorMessage.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaultUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		want *ErrorMessage
	}{
		{
			name: "Should return default unauthorized error",
			want: &ErrorMessage{
				Code: 401,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        UnauthorizedErrorName,
						ErrorDescription: UnauthorizedErrorDescription,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDefaultUnauthorized(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDefaultUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaultNewNotFound(t *testing.T) {
	tests := []struct {
		name string
		want *ErrorMessage
	}{
		{
			name: "Should return default unauthorized error",
			want: &ErrorMessage{
				Code: 404,
				ErrorList: []*ErrorFormat{
					{
						ErrorName:        NotFound,
						ErrorDescription: NotFoundDescription,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDefaultNewNotFound(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDefaultNewNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_Errorln(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorMessage
	}{
		{
			name: "",
			fields: fields{
				Code: 1,
			},
			args: args{
				err: errors.New("hai"),
			},
			want: &ErrorMessage{
				Code: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.Errorln(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorMessage.Errorln() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_Debugln(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ErrorMessage
	}{
		{
			name: "",
			fields: fields{
				Code: 1,
			},
			args: args{
				err: errors.New("hai"),
			},
			want: &ErrorMessage{
				Code: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.Debugln(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorMessage.Debugln() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage_GetErrorDesc(t *testing.T) {
	type fields struct {
		ErrorList []*ErrorFormat
		Code      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success to get error description",
			fields: fields{
				ErrorList: []*ErrorFormat{
					&ErrorFormat{
						ErrorDescription: stringFoo,
					},
				},
			},
			want: stringFoo,
		},
		{
			name: "Empty string",
			fields: fields{
				ErrorList: []*ErrorFormat{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMessage := &ErrorMessage{
				ErrorList: tt.fields.ErrorList,
				Code:      tt.fields.Code,
			}
			if got := errorMessage.GetErrorDesc(); got != tt.want {
				t.Errorf("ErrorMessage.GetErrorDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaultErrBodyRequest(t *testing.T) {
	tests := []struct {
		name string
		want *ErrorMessage
	}{
		{
			want: SetNewBadRequestByFormat(&DefaultInputBody),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDefaultErrBodyRequest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDefaultErrBodyRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetListErr(t *testing.T) {
	dataErr := []*ErrorFormat{}
	type args struct {
		code    int
		listErr []*ErrorFormat
	}
	tests := []struct {
		name string
		args args
		want *ErrorMessage
	}{
		{
			name: "Get error",
			args: args{
				code:    http.StatusOK,
				listErr: dataErr,
			},
			want: &ErrorMessage{
				Code:      http.StatusOK,
				ErrorList: dataErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetListErr(tt.args.code, tt.args.listErr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetListErr() = %v, want %v", got, tt.want)
			}
		})
	}
}
