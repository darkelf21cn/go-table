package gotable

const (
	strSingle   = "a"
	strShort1   = "abcd"
	strShort2   = "ab cd ef gh"
	strLong     = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras sit amet suscipit elit. Praesent in dictum est. Aenean nec congue lectus, non gravida diam. Sed accumsan vel lectus sed fringilla. In hac habitasse platea dictumst. Duis ac justo est. Aenean bibendum tellus sed risus convallis blandit. Donec quis nunc pretium, suscipit arcu et, suscipit lectus. Etiam orci tellus, luctus in molestie molestie, efficitur et urna. Maecenas pretium tincidunt risus volutpat sodales. Nunc quis suscipit odio. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;`
	strWordWrap = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nulla eget mi nec ipsum aliquam pulvinar.
Aenean id justo ac diam iaculis gravida nec et ex.
Fusce sed quam hendrerit, mollis nisi vitae, porttitor erat.`
	strHelloChinese = "你好，世界"
	strSQL          = `SELECT
	*
FROM
	information_schema.tables
WHERE
	TABLE_SCHEMA = 'mysql'`
)

type mockTreeNode struct {
	ID       int
	Data     string
	children []TreeNodeReader
}

func (a *mockTreeNode) Children() []TreeNodeReader {
	return a.children
}

func (a *mockTreeNode) Fields() map[string]any {
	return map[string]any{
		"ID":   a.ID,
		"Data": a.Data,
	}
}
