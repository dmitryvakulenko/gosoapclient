package gosoapclient

import (
	"encoding/xml"
	"reflect"
	"bytes"
)

type Marshaler struct {
	namespaces map[string]string
	typeNamespace map[string]string
	usedNamespaces map[string]string
}

func NewMarshaller(namespaceMap, typeNamespaceMap map[string]string) Marshaler {
	return Marshaler{
		namespaces: namespaceMap,
		typeNamespace: typeNamespaceMap,
		usedNamespaces: make(map[string]string)}
}


func (p *Marshaler) Encode(input interface{}) ([]byte, error) {
	res := make([]byte, 0)
	e := xml.NewEncoder(bytes.NewBuffer(res))

	p.marshal(input, e)

	return res, nil
}


func (p *Marshaler) marshal(input interface{}, e *xml.Encoder) error {
	v := reflect.ValueOf(p)

	t := xml.StartElement{
		Name: xml.Name{Local: p.getElementName(v)}}
	e.EncodeToken(t)

	fieldsNum := v.NumField()
	for i := 0; i < fieldsNum; i++ {
		field := v.Type().Field(i)
		ns := field.Tag.Get("ns")
		nsAlias := p.namespaces[ns]
		p.usedNamespaces[nsAlias] = ns
		t := xml.StartElement{
			Name: xml.Name{Local: nsAlias + ":" + field.Name}}
		e.EncodeToken(t)
		e.EncodeElement(v.Field(i).Interface(), t)
		e.EncodeToken(t.End())
	}

	e.EncodeToken(t.End())

	return nil
}


func (p *Marshaler) getElementName(val reflect.Value) string {
	t := val.Type()
	prefix := ""
	if ns, ok := p.typeNamespace[t.Name()]; ok {
		nsAlias := p.namespaces[ns]
		p.usedNamespaces[nsAlias] = ns
		prefix = nsAlias + ":"
	}

	return prefix + t.Name()
}


//func (c *Client) collectNamespaces(in interface{}) map[string]string {
//	inType := reflect.ValueOf(in)
//	res := make(map[string]string)
//
//	inTypeName := inType.Type().Name()
//	if ns, ok := c.typesNamespaces[inTypeName]; ok {
//		nsAlias := c.namespacesAlias[ns]
//		res[nsAlias] = ns
//	}
//
//	if inType.Kind() == reflect.Struct {
//		fieldsNum := inType.NumField()
//		for i := 0; i < fieldsNum; i++ {
//			val := inType.Field(i)
//			res = mergeNamespaces(res, c.collectNamespaces(val.Interface()))
//		}
//	} else if inType.Kind() == reflect.Slice {
//		val, ok := in.([]interface{})
//		if !ok {
//			return res
//		}
//		for _, v := range val {
//			res = mergeNamespaces(res, c.collectNamespaces(v))
//		}
//	} else if inType.Kind() == reflect.Ptr && !inType.IsNil() {
//		ptr := inType.Elem().Interface()
//		res = mergeNamespaces(res, c.collectNamespaces(ptr))
//	}
//
//	return res
//}

func mergeNamespaces(first map[string]string, second map[string]string) map[string]string {
	res := first
	for k, v := range second {
		res[k] = v
	}

	return res
}


func (p *Marshaler) GetUsedNamespaces() map[string]string {
	return p.usedNamespaces
}

func (p *Marshaler) Reset() {
	p.usedNamespaces = make(map[string]string)
}