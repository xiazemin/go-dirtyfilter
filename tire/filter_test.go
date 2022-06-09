package tire

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
)

func BenchmarkFilterResult(t *testing.B) {
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
				texts: []string{"中国", "人", "外国", "中国人", "中", "ab", "ef", "cd", "好"}, //[]string
			},
			want: [][]string{{"中国"}, {"中国", "人"}, {"人", "外国"}, {"ab", "cd", "ef"}, {"中国", "人", "好"}}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"中国"}, {"中国", "人"}, {"人", "中国"}, {"外国", "人"}, {"不是", "中国", "人"}, {"中国", "好", "人"}, {"cd", "ef", "ab"}}, true)
				fields.root = fileter.GetRoot()
				//fileter.Print()
				gomock.InOrder()
			},
		},
		{
			name: "case4",
			fields: fields{
				ignoreOrder: true, // bool
				root:        nil,  // *node
			},
			args: args{
				texts: []string{"劫掠", "斯大林", "长春", "劫掠", "斯大林", "沈阳", "劫掠", "哈尔滨", "斯大林", "东北", "搜刮", "斯大林", "东三省", "搜刮", "斯大林", "搜刮", "斯大林", "长春", "搜刮", "斯大林", "沈阳", "哈尔滨", "搜刮", "斯大林", "东北", "搬运", "斯大林", "东三省", "搬运", "斯大林", "搬运"}, //[]string
			},
			want: [][]string{{"劫掠", "斯大林", "沈阳"}, {"劫掠", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}, {"搜刮", "斯大林", "长春"}}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"劫掠", "斯大林", "沈阳"}, {"劫掠", "斯大林", "长春"}, {"搜刮", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}}, true)
				fields.root = fileter.GetRoot()
				//fileter.Print()
				gomock.InOrder()
			},
		},
	}
	for i := 0; i < t.N; i++ {
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
				lock:        &sync.RWMutex{},
			}
			if got := tf.FilterResult(tt.args.texts); !reflect.DeepEqual(got, tt.want) {
				j1, _ := json.Marshal(got)
				j2, _ := json.Marshal(tt.want)
				t.Errorf("%q. tireFilter.FilterResult() = %#v, want %#v\n", tt.name, string(j1), string(j2))
			}
			httpmock.DeactivateAndReset()
		}
	}
}
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
				texts: []string{"中国", "人", "外国", "中国人", "中", "ab", "ef", "cd", "好"}, //[]string
			},
			want: [][]string{{"中国"}, {"中国", "人"}, {"人", "外国"}, {"ab", "cd", "ef"}, {"中国", "人", "好"}}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"中国"}, {"中国", "人"}, {"人", "中国"}, {"外国", "人"}, {"不是", "中国", "人"}, {"中国", "好", "人"}, {"cd", "ef", "ab"}}, true)
				fields.root = fileter.GetRoot()
				fileter.Print()
				gomock.InOrder()
			},
		},
		{
			name: "case4",
			fields: fields{
				ignoreOrder: true, // bool
				root:        nil,  // *node
			},
			args: args{
				texts: []string{"劫掠", "斯大林", "长春", "劫掠", "斯大林", "沈阳", "劫掠", "哈尔滨", "斯大林", "东北", "搜刮", "斯大林", "东三省", "搜刮", "斯大林", "搜刮", "斯大林", "长春", "搜刮", "斯大林", "沈阳", "哈尔滨", "搜刮", "斯大林", "东北", "搬运", "斯大林", "东三省", "搬运", "斯大林", "搬运"}, //[]string
			},
			want: [][]string{{"劫掠", "斯大林", "沈阳"}, {"劫掠", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}, {"搜刮", "斯大林", "长春"}}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"劫掠", "斯大林", "沈阳"}, {"劫掠", "斯大林", "长春"}, {"搜刮", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}}, true)
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
			lock:        &sync.RWMutex{},
		}
		if got := tf.FilterResult(tt.args.texts); !reflect.DeepEqual(got, tt.want) {
			j1, _ := json.Marshal(got)
			j2, _ := json.Marshal(tt.want)
			t.Errorf("%q. tireFilter.FilterResult() = %#v, want %#v\n", tt.name, string(j1), string(j2))
		}
		httpmock.DeactivateAndReset()
	}
}

func BenchmarkFilterResultCount(t *testing.B) {
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
		want    map[string]int
	}{
		{
			name: "case1",
			fields: fields{
				ignoreOrder: true, // bool
				root:        nil,  // *node
			},
			args: args{
				texts: []string{"劫掠", "劫掠", "劫掠", "哈尔滨", "哈尔滨", "哈尔滨", "哈尔滨", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "沈阳", "沈阳", "沈阳", "沈阳", "苏联", "苏联", "苏联", "苏联", "苏联", "苏联", "长春", "长春", "长春", "长春"}, //[]string
			},
			want: map[string]int{"劫掠,斯大林,沈阳": 3, "劫掠,斯大林,长春": 3, "搜刮,斯大林,沈阳": 4, "搜刮,斯大林,长春": 4}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"劫掠", "斯大林", "长春"}, {"劫掠", "斯大林", "沈阳"}, {"搜刮", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}}, true)
				fields.root = fileter.GetRoot()
				fileter.Print()
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
				texts: []string{"劫掠", "斯大林", "长春", "劫掠", "劫掠", "哈尔滨", "哈尔滨", "哈尔滨", "哈尔滨", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "沈阳", "沈阳", "沈阳", "沈阳", "苏联", "苏联", "苏联", "苏联", "苏联", "苏联", "长春", "长春", "长春"}, //[]string
			},
			want: map[string]int{"劫掠,斯大林,长春": 1}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"劫掠", "斯大林", "沈阳"}, {"劫掠", "斯大林", "长春"}, {"搜刮", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}}, false)
				fields.root = fileter.GetRoot()
				fileter.Print()
				gomock.InOrder()
			},
		},
	}
	for i := 0; i < t.N; i++ {
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
				lock:        &sync.RWMutex{},
			}
			if got := tf.FilterResultCount(tt.args.texts); !reflect.DeepEqual(got, tt.want) {
				j1, _ := json.Marshal(got)
				j2, _ := json.Marshal(tt.want)
				t.Errorf("%q. tireFilter.FilterResultCount() = %#v, want %#v", tt.name, string(j1), string(j2))
			}
			httpmock.DeactivateAndReset()
		}
	}
}

func Test_tireFilter_FilterResultCount(t *testing.T) {
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
		want    map[string]int
	}{
		{
			name: "case1",
			fields: fields{
				ignoreOrder: true, // bool
				root:        nil,  // *node
			},
			args: args{
				texts: []string{"劫掠", "劫掠", "劫掠", "哈尔滨", "哈尔滨", "哈尔滨", "哈尔滨", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "沈阳", "沈阳", "沈阳", "沈阳", "苏联", "苏联", "苏联", "苏联", "苏联", "苏联", "长春", "长春", "长春", "长春"}, //[]string
			},
			want: map[string]int{"劫掠,斯大林,沈阳": 3, "劫掠,斯大林,长春": 3, "搜刮,斯大林,沈阳": 4, "搜刮,斯大林,长春": 4}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"劫掠", "斯大林", "长春"}, {"劫掠", "斯大林", "沈阳"}, {"搜刮", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}}, true)
				fields.root = fileter.GetRoot()
				//fileter.Print()
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
				texts: []string{"劫掠", "斯大林", "长春", "劫掠", "劫掠", "哈尔滨", "哈尔滨", "哈尔滨", "哈尔滨", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搜刮", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "搬运", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "斯大林", "沈阳", "沈阳", "沈阳", "沈阳", "苏联", "苏联", "苏联", "苏联", "苏联", "苏联", "长春", "长春", "长春"}, //[]string
			},
			want: map[string]int{"劫掠,斯大林,长春": 1}, //[][]string,
			prepare: func(fields *fields, args *args) {
				//httpmock.RegisterResponder("GET", "https://mytest.com/httpmock", httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Test"}]`))
				fileter := NewTireFilter([][]string{{"劫掠", "斯大林", "沈阳"}, {"劫掠", "斯大林", "长春"}, {"搜刮", "斯大林", "长春"}, {"搜刮", "斯大林", "沈阳"}}, false)
				fields.root = fileter.GetRoot()
				//fileter.Print()
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
			lock:        &sync.RWMutex{},
		}
		if got := tf.FilterResultCount(tt.args.texts); !reflect.DeepEqual(got, tt.want) {
			j1, _ := json.Marshal(got)
			j2, _ := json.Marshal(tt.want)
			t.Errorf("%q. tireFilter.FilterResultCount() = %#v, want %#v", tt.name, string(j1), string(j2))
		}
		httpmock.DeactivateAndReset()
	}
}
