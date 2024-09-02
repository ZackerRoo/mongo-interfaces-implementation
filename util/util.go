package util

// IsListField 检查字段是否是一个列表类型的字段
func IsListField(fieldName string) bool {
	listFields := []string{"tags", "knowledgeType", "tacticsId", "techniquesId", "subTechniquesId"}
	for _, field := range listFields {
		if field == fieldName {
			return true
		}
	}
	return false
}
