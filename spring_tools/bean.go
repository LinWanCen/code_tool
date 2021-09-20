package spring_tools

type (
	// Spring bean 解析
	//  <bean id="" class="">
	//    <property name="" ref="" />
	//  </bean>
	SpringBeans struct {
		// 有beans标签在所以换名字
		BeanArray []Bean `xml:"bean"`
	}

	Bean struct {
		Id         string     `xml:"id,attr"`
		Name       string     `xml:"name,attr"`
		Class      string     `xml:"class,attr"`
		Properties []Property `xml:"property"`
	}

	Property struct {
		Name string `xml:"name,attr"`
		Ref  string `xml:"ref,attr"`
	}
)
