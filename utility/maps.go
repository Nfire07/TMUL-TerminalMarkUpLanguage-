package utility

type Map map[string]string

func (attributes Map) HasAttr(attrName string) bool {
	_, exists := attributes[attrName]
	return exists
}
