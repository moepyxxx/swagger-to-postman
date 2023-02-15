package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	swag "github.com/astaxie/beego/swagger"
	"github.com/rbretecher/go-postman-collection"
	"gopkg.in/yaml.v3"
)

var basicHost = "{{TENANT_URL}}%v"

func createURL(path string, params []swag.Parameter) *postman.URL {
	rePath := path
	paths := []string{}
	for _, param := range params {
		if param.In == "path" {
			paths = append(paths, param.Name)
		}
	}

	variables := []*postman.Variable{}
	for _, p := range paths {
		rePath = strings.Replace(
			rePath,
			fmt.Sprintf("{%v}", p),
			fmt.Sprintf(":%v", p),
			1,
		)
		variables = append(variables, &postman.Variable{
			Key:   p,
			Value: "1",
		})
	}

	raw := fmt.Sprintf(basicHost, rePath)
	return &postman.URL{
		Raw:       raw,
		Path:      strings.Split(string([]rune(rePath)[1:]), "/"),
		Variables: variables,
		Host:      []string{"{{TENANT_URL}}"},
	}
}

func getTag(tags []string) (string, error) {
	if len(tags) == 0 {
		return "", errors.New("invalid tags ")
	}

	return tags[0], nil
}

func main() {
	// 読み込み
	f, err := os.Open("swagger.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	swag := swag.Swagger{}
	yaml.Unmarshal(b, &swag)

	// postmanファイル作成
	c := postman.CreateCollection("casone-generated-swagger", "generate sareruyo !!")

	// カテゴリごとリクエストの塊を作る
	pis := map[string][]*postman.Items{}
	for pathKey, item := range swag.Paths {
		if item.Get != nil {
			url := createURL(pathKey, item.Get.Parameters)
			if err != nil {
				panic(err)
			}

			tag, err := getTag(item.Get.Tags)
			if err != nil {
				panic(err)
			}

			pis[tag] = append(pis[tag], &postman.Items{
				Name:        item.Get.Summary,
				Description: item.Get.Description,
				Request: &postman.Request{
					URL: url,
				},
			})
		}
	}

	// 書き込み
	for groupKey, items := range pis {
		c.AddItem(&postman.Items{
			Name:  groupKey,
			Items: items,
		})
	}

	file, err := os.Create("postman.json")
	defer file.Close()

	if err != nil {
		panic(err)
	}

	err = c.Write(file, postman.V210)
	if err != nil {
		panic(err)
	}
}
