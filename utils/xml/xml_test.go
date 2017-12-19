package xml

import (
	"testing"
	"log"
)

func TestXmlParse(t *testing.T) {
	// error t.Errorf
	// t.Fatal
	xn := XmlNodeNew()
	if xn == nil {
		t.Fatal("create xmlnode failed!")
	}
	context := `<xml type="xxx"><person><AGE>12</AGE><address>Beijing</address><ttt>123</ttt></person></xml>`
	xn.Parse([]byte(context))
	if xn.Gets("person") == nil || len(xn.Gets("person")) != 1 {
		t.Errorf("error! %s", "get function failed!")
	}
	if len(xn.Gets("person")[0].Next) != 3{
		t.Errorf("error2! %s","get function failed!")
	}
	if len(xn.Gets("person")[0].Gets("address")) != 1 {
		t.Error("second layer error!")
	}
	if xn.Gets("person")[0].Gets("address")[0].Content != "Beijing" ||
		xn.Gets("person")[0].Gets("ttt")[0].Content != "123" ||
			xn.Gets("person")[0].Gets("age")[0].Content != "12"{
		t.Error("second layer content error !")
	}

}

func TestXmlParse2(t *testing.T) {
	// error t.Errorf
	// t.Fatal
	xn := XmlNodeNew()
	if xn == nil {
		t.Fatal("create xmlnode failed!")
	}
	context := `<xml type="xxx">
		<person><age>12</age><address>Beijing</address></person>
		<person><age>12</age><address>Beijing</address></person>
		<person><age>12</age><address>Beijing</address></person>
		</xml>`
	xn.Parse([]byte(context))
	if xn.Gets("person") == nil || len(xn.Gets("person")) != 3 {
		t.Errorf("error! %s", "get function failed!")
	}
	if len(xn.Gets("person")[1].Next) != 2 {
		t.Errorf("error2! %s", "get function failed!")
	}
	if len(xn.Gets("person")[1].Gets("address")) != 1 {
		t.Error("second layer error!")
	}
	if xn.Gets("person")[1].Gets("address")[0].Content != "Beijing" ||
		xn.Gets("person")[1].Gets("age")[0].Content != "12" {
		t.Error("second layer content error !")
	}

}

func TestXmlParse3(t *testing.T) {
	content := `<?xml version='1.0' encoding='utf-8'?><res><tranCode>HKMB000000</tranCode><returnVal></returnVal><status><value>0</value><msg></msg></status><updateStatus>NO</updateStatus><RNS>CtSxxi9Ziwh+PQZMoXyUT6euMb9cVx1/5yhacSIDAng=</RNS><RNC>AAEAJPyRdgvZu1PotkswxZMoWqJxzEvgriqeMwKWBJk=</RNC><MS>c6p94LZxZZ/idhEbU1DttGIUvN/p1Qlvl90DxSlLvjmXXdn2m9q/mGMP2okIY510HGXZQRgCotjSYOkCosS1FwRbFFw=</MS><toDownloadIOS></toDownloadIOS><toDownloadAndroid></toDownloadAndroid></res>`
	xn := XmlNodeNew()
	if xn == nil {
		t.Fatal("create xmlnode failed!")
	}
	xn.Parse([]byte(content))
	if xn.Gets("res") == nil || len(xn.Gets("res")[0].Next) != 9 {
		t.Error("len:", len(xn.Gets("res")))
	}
}

func TestXmlParse4(t *testing.T) {
	content := `<res><tranCode>HKMB000000</tranCode><tranCode>HKMB000001</tranCode><tranCode>HKMB000002</tranCode></res>`
	xn := XmlNodeNew()
	if xn == nil {
		t.Fatal("create xmlnode failed!")
	}
	xn.Parse([]byte(content))
	if xn.Gets("res") == nil || len(xn.Gets("res")[0].Next) != 3 {
		t.Error("len:", len(xn.Gets("res")))
	}
}

func TestXmlJson(t *testing.T) {
	content := `<res><tranCode>HKMB000000</tranCode><tranCode>HKMB000001</tranCode><tranCode>HKMB000002</tranCode><value>0</value></res>`
	xn := XmlNodeNew()
	if xn == nil {
		t.Fatal("create xmlnode failed!")
	}
	xn.Parse([]byte(content))

	json := xn.ToJson()
	log.Print("json:", json)
	if json == "" {
		t.Error("xml to json failed!")
	}
}

func TestXmlJson2(t *testing.T) {
	content := `<?xml version='1.0' encoding='utf-8'?><res><tranCode>HKMB000000</tranCode><returnVal></returnVal><status><value>0</value><msg></msg></status><updateStatus>NO</updateStatus><RNS>CtSxxi9Ziwh+PQZMoXyUT6euMb9cVx1/5yhacSIDAng=</RNS><RNC>AAEAJPyRdgvZu1PotkswxZMoWqJxzEvgriqeMwKWBJk=</RNC><MS>c6p94LZxZZ/idhEbU1DttGIUvN/p1Qlvl90DxSlLvjmXXdn2m9q/mGMP2okIY510HGXZQRgCotjSYOkCosS1FwRbFFw=</MS><toDownloadIOS></toDownloadIOS><toDownloadAndroid></toDownloadAndroid></res>`
	xn := XmlNodeNew()
	if xn == nil {
		t.Fatal("create xmlnode failed!")
	}
	xn.Parse([]byte(content))

	json := xn.ToJson()
	log.Print("json:", json)
	if json == "" {
		t.Error("xml to json failed!")
	}
}