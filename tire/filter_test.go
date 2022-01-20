package tire

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
)

func Test_tireFilter_FilterResult(t *testing.T) {
	type fields struct {
		ignoreOrder bool
		root        *node
	}
	type args struct {
		texts []string
	}
	tests := []struct {
		name    string
		fields  fields
		prepare func(fields *fields, args *args)
		args    args
		want    [][]string
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				ignoreOrder: false, // bool
				root:        nil,   // *node
			},
			args: args{
				texts: []string{}, //[]string
			},
			want: nil, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				gomock.InOrder()
			},
		},
		{
			name: "case2",
			fields: fields{
				ignoreOrder: false, // bool
				root:        nil,   // *node
			},
			args: args{
				texts: []string{"人", "中国", "外国", "中国人", "中"}, //[]string
			},
			want: [][]string{{"人", "中国"}, {"中国"}}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"中国"}, {"中国", "人"}, {"人", "中国"}, {"外国", "人"}, {"不是", "中国", "人"}, {"cd", "ef", "ab"}}, false)
				fields.root = fileter.GetRoot()
				//fileter.Print()
				gomock.InOrder()
			},
		},
		{
			name: "case3",
			fields: fields{
				ignoreOrder: true, // bool
				root:        nil,  // *node
			},
			args: args{
				texts: []string{"中国", "人", "外国", "中国人", "中", "ab", "ef", "cd"}, //[]string
			},
			want: [][]string{{"ab", "cd", "ef"}, {"中国"}, {"中国", "人"}, {"人", "外国"}}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"中国"}, {"中国", "人"}, {"人", "中国"}, {"外国", "人"}, {"不是", "中国", "人"}, {"cd", "ef", "ab"}}, true)
				fields.root = fileter.GetRoot()
				fileter.Print()
				gomock.InOrder()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpmock.Activate()
		if tt.prepare != nil {
			tt.prepare(&tt.fields, &tt.args)
		}

		tf := &tireFilter{
			ignoreOrder: tt.fields.ignoreOrder,
			root:        tt.fields.root,
		}
		if got := tf.FilterResult(tt.args.texts); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. tireFilter.FilterResult() = %#v, want %#v", tt.name, got, tt.want)
		}
		httpmock.DeactivateAndReset()
	}
}
