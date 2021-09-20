package dubbo_tools

import (
	"encoding/xml"
	"path/filepath"

	"code_tools/spring_tools"
	"code_tools/util"
)

type (
	// dubbo 解析
	//
	//  <dubbo:reference id="" interface="" version="" group=""/>
	//  http://dubbo.apache.org/zh-cn/docs/user/references/xml/dubbo-reference.html
	//  https://github.com/apache/dubbo-website/blob/master/docs/zh-cn/user/references/xml/dubbo-reference.md
	//
	//  <dubbo:service ref="" interface="" version="" group="" />
	//  http://dubbo.apache.org/zh-cn/docs/user/references/xml/dubbo-service.html
	//  https://github.com/apache/dubbo-website/blob/master/docs/zh-cn/user/references/xml/dubbo-service.md
	//
	//  <dubbo:service ref="" interface="" version="" group="" />
	//  http://dubbo.apache.org/zh-cn/docs/user/references/xml/dubbo-service.html
	//  https://github.com/apache/dubbo-website/blob/master/docs/zh-cn/user/references/xml/dubbo-service.md
	Dubbo struct {
		// 固定名字 XMLName
		XMLName xml.Name
		spring_tools.SpringBeans
		References     []Reference `xml:"reference"`
		Services       []Service   `xml:"service"`
		EveryReference func(result [][][]string, a Reference)
		EveryService   func(result [][][]string, a Service, bean spring_tools.Bean)
	}

	Reference struct {
		Id        string `xml:"id,attr"`
		Interface string `xml:"interface,attr"`
		Version   string `xml:"version,attr"`
		Group     string `xml:"group,attr"`
	}

	Service struct {
		Ref       string `xml:"ref,attr"`
		Interface string `xml:"interface,attr"`
		Version   string `xml:"version,attr"`
		Group     string `xml:"group,attr"`
	}
)

func (b Dubbo) FeaturesName() string {
	return "Dubbo 清单"
}

// 需要修改值，所以传指针
func (d *Dubbo) Before(i ...interface{}) (results []interface{}) {
	d.EveryReference = everyReference
	d.EveryService = everyService
	return
}

func (d Dubbo) FileNames() []string {
	return []string{"dubbo_reference_list", "dubbo_service_list"}
}

func (d Dubbo) Titles() [][]string {
	return [][]string{{
		"id",
		"giv",
	}, {
		"giv",
		"ref",
		"beanId",
		"class",
	}}
}

func (d *Dubbo) Lines(rootPath string, path string) [][][]string {
	if filepath.Ext(path) != ".xml" {
		return nil
	}

	model := Dubbo{}
	err := util.ParseXml(path, &model)
	if err != nil {
		return nil
	}
	if model.XMLName.Local != "beans" {
		return nil
	}

	result := [][][]string{{}, {}}
	for _, a := range model.References {
		d.EveryReference(result, a)
	}
	// bean 需要配在 service 的同个文件
	beanMap := make(map[string]spring_tools.Bean)
	for _, bean := range model.BeanArray {
		beanMap[bean.Id] = bean
	}
	for _, a := range model.Services {
		bean := beanMap[a.Ref]
		d.EveryService(result, a, bean)
	}
	return result
}

func everyReference(result [][][]string, a Reference) {
	result[0] = append(result[0], []string{
		a.Id,
		a.Group + ":" + a.Interface + ":" + a.Version,
	})
}

func everyService(result [][][]string, a Service, bean spring_tools.Bean) {
	result[1] = append(result[1], []string{
		a.Group + ":" + a.Interface + ":" + a.Version,
		a.Ref,
		bean.Id,
		bean.Class,
	})
}
