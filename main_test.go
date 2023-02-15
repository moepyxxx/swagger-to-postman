package main

import (
	"testing"

	"github.com/K0kubun/pp"
	swag "github.com/astaxie/beego/swagger"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rbretecher/go-postman-collection"
)

func Test_CreateURL(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		params []swag.Parameter
		want   *postman.URL
	}{
		{
			name: "pathパラメータが指定されていた場合は1つの動的パスに対応できる",
			path: "/path/to/{tenant_id}",
			params: []swag.Parameter{
				{
					In:   "path",
					Name: "tenant_id",
				},
			},
			want: &postman.URL{
				Raw: "{{TENANT_URL}}/path/to/:tenant_id",
				Path: []string{
					"path",
					"to",
					":tenant_id",
				},
				Variables: []*postman.Variable{
					{
						Key:   "tenant_id",
						Value: "1",
					},
				},
				Host: []string{"{{TENANT_URL}}"},
			},
		},
		{
			name: "pathパラメータが指定されていた場合は複数の動的パスに対応できる",
			path: "/path/to/{tenant_id}/to/staff/{staff_id}",
			params: []swag.Parameter{
				{
					In:   "path",
					Name: "tenant_id",
				},
				{
					In:   "path",
					Name: "staff_id",
				},
			},
			want: &postman.URL{
				Raw: "{{TENANT_URL}}/path/to/:tenant_id/to/staff/:staff_id",
				Path: []string{
					"path",
					"to",
					":tenant_id",
					"to",
					"staff",
					":staff_id",
				},
				Variables: []*postman.Variable{
					{
						Key:   "tenant_id",
						Value: "1",
					},
					{
						Key:   "staff_id",
						Value: "1",
					},
				},
				Host: []string{"{{TENANT_URL}}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createURL(tt.path, tt.params)
			pp.Print(tt.want)
			pp.Print(got)
			diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(postman.URL{}))
			if diff != "" {
				t.Errorf("createURL() = %v, want %v", got, tt.want)
			}
		})
	}

}
