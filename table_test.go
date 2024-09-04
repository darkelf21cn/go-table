package gotable

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtilities(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils test suite")
}

var _ = Describe("Table Test Suites", func() {
	Context("render-row", func() {
		It("t1", func() {
			tb := NewTable(nil)
			tb.AppendColumn(test_NewStdColumn("ID"))
			tb.AppendColumn(test_NewStdColumn("Data"))
			tb.AppendRow(1, strShort1)
			tb.AppendRow(2, strShort2)
			tb.AppendRow(3, strHelloChinese)
			out, err := tb.Render(Console)
			Expect(err).Should(BeNil())
			expects := []string{
				"+----+-------------+",
				"| ID |    Data     |",
				"+----+-------------+",
				"| 1  | abcd        |",
				"| 2  | ab cd ef gh |",
				"| 3  | 你好，世界  |",
				"+----+-------------+",
				"",
			}
			for i, line := range strings.Split(out, "\n") {
				Expect(line).Should(Equal(expects[i]))
			}
		})
		It("t2", func() {
			ly := LightTableLayout()
			ly.Width = 40
			tb := NewTable(ly)
			tb.AppendColumn(test_NewStdColumn("ID"))
			tb.AppendColumn(test_NewStdColumn("Data"))
			tb.AppendRow(1, strShort1)
			tb.AppendRow(2, strShort2)
			tb.AppendRow(3, strHelloChinese)
			out, err := tb.Render(Console)
			Expect(err).Should(BeNil())
			expects := []string{
				"┌────┬─────────────────────────────────┐",
				"│ ID │              Data               │",
				"├────┼─────────────────────────────────┤",
				"│ 1  │ abcd                            │",
				"│ 2  │ ab cd ef gh                     │",
				"│ 3  │ 你好，世界                      │",
				"└────┴─────────────────────────────────┘",
				"",
			}
			for i, line := range strings.Split(out, "\n") {
				Expect(line).Should(Equal(expects[i]))
			}
		})
		It("t3", func() {
			ly := LightTableLayout()
			ly.Width = 60
			tb := NewTable(ly)
			tb.AppendColumn(test_NewStdColumn("ID"))
			tb.AppendColumn(test_NewStdColumn("Data1"))
			tb.AppendColumn(test_NewStdColumn("Data2"))
			tb.AppendRow(1, strShort1, strShort2)
			tb.AppendRow(2, strSQL, strHelloChinese)
			tb.AppendRow(3, strHelloChinese, strWordWrap)
			out, err := tb.Render(Console)
			Expect(err).Should(BeNil())
			expects := []string{
				"┌────┬──────────────────────────┬──────────────────────────┐",
				"│ ID │          Data1           │          Data2           │",
				"├────┼──────────────────────────┼──────────────────────────┤",
				"│ 1  │ abcd                     │ ab cd ef gh              │",
				"│ 2  │ SELECT                   │ 你好，世界               │",
				"│    │     *                    │                          │",
				"│    │ FROM                     │                          │",
				"│    │     information_schema.t │                          │",
				"│    │ ables                    │                          │",
				"│    │ WHERE                    │                          │",
				"│    │     TABLE_SCHEMA = 'mysq │                          │",
				"│    │ l'                       │                          │",
				"│ 3  │ 你好，世界               │ Lorem ipsum dolor sit am │",
				"│    │                          │ et, consectetur adipisci │",
				"│    │                          │ ng elit.                 │",
				"│    │                          │ Nulla eget mi nec ipsum  │",
				"│    │                          │ aliquam pulvinar.        │",
				"│    │                          │ Aenean id justo ac diam  │",
				"│    │                          │ iaculis gravida nec et e │",
				"│    │                          │ x.                       │",
				"│    │                          │ Fusce sed quam hendrerit │",
				"│    │                          │ , mollis nisi vitae, por │",
				"│    │                          │ ttitor erat.             │",
				"└────┴──────────────────────────┴──────────────────────────┘",
				"",
			}
			for i, line := range strings.Split(out, "\n") {
				Expect(line).Should(Equal(expects[i]))
			}
		})
		It("t4", func() {
			ly := LightTableLayout()
			ly.Width = 40
			tb := NewTable(ly)
			tb.AppendColumn(test_NewStdColumn("ID"))
			tb.AppendColumn(test_NewStdColumn("Data1").Width(10, false))
			tb.AppendColumn(test_NewStdColumn("Data2").Hidden(true))
			tb.AppendColumn(test_NewStdColumn("Data3"))
			rows := []map[string]any{
				{"ID": 1, "Data1": strShort1, "Data2": strShort2, "Data3": strHelloChinese},
				{"ID": 2, "Data1": strShort2, "Data2": strHelloChinese, "Data3": strShort1},
				{"ID": 3, "Data1": strHelloChinese, "Data2": strShort1, "Data3": strShort2},
			}
			for _, row := range rows {
				err := tb.AppendRowM(row)
				Expect(err).Should(BeNil())
			}
			out, err := tb.Render(Console)
			Expect(err).Should(BeNil())
			expects := []string{
				"┌────┬──────────┬──────────────────────┐",
				"│ ID │  Data1   │        Data3         │",
				"├────┼──────────┼──────────────────────┤",
				"│ 1  │ abcd     │ 你好，世界           │",
				"│ 2  │ ab cd ef │ abcd                 │",
				"│    │  gh      │                      │",
				"│ 3  │ 你好，世 │ ab cd ef gh          │",
				"│    │ 界       │                      │",
				"└────┴──────────┴──────────────────────┘",
				"",
			}
			for i, line := range strings.Split(out, "\n") {
				Expect(line).Should(Equal(expects[i]))
			}
		})
		It("t5", func() {
			ly := LightTableLayout()
			ly.Width = 10
			tb := NewTable(ly)
			tb.AppendColumn(test_NewStdColumn("ID"))
			tb.AppendColumn(test_NewStdColumn("Data"))
			tb.AppendRow(1, strShort1)
			tb.AppendRow(2, strShort2)
			tb.AppendRow(3, strHelloChinese)
			_, err := tb.Render(Console)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("render-tree", func() {
		It("t1", func() {
			tb := NewTable(nil)
			tb.AppendColumn(test_NewStdColumn("ID"))
			tb.AppendColumn(test_NewStdColumn("Data"))
			nodes := []TreeNodeReader{
				&mockTreeNode{
					ID:   1,
					Data: strShort1,
					children: []TreeNodeReader{
						&mockTreeNode{
							ID:   2,
							Data: strShort2,
							children: []TreeNodeReader{
								&mockTreeNode{
									ID:   4,
									Data: strHelloChinese,
									children: []TreeNodeReader{
										&mockTreeNode{
											ID:   5,
											Data: strShort1,
										},
									},
								},
							},
						},
						&mockTreeNode{
							ID:   3,
							Data: strShort2,
						},
					},
				},
				&mockTreeNode{
					ID:   6,
					Data: strHelloChinese,
				},
			}
			sty := DefaultTreePathStyle().Header()
			err := tb.AppendTrees(*sty, nodes...)
			Expect(err).To(BeNil())
			out, err := tb.Render(Console)
			Expect(err).Should(BeNil())
			expects := []string{
				`+-----------+----+-------------+`,
				`|   Path    | ID |    Data     |`,
				`+-----------+----+-------------+`,
				`| >-+------ | 1  | abcd        |`,
				`|   +-+---- | 2  | ab cd ef gh |`,
				`|   | \-+-- | 4  | 你好，世界  |`,
				`|   |   \-- | 5  | abcd        |`,
				`|   \------ | 3  | ab cd ef gh |`,
				`| >-------- | 6  | 你好，世界  |`,
				`+-----------+----+-------------+`,
				``,
			}
			for i, line := range strings.Split(out, "\n") {
				Expect(line).Should(Equal(expects[i]))
			}
		})
	})
})

func test_NewStdColumn(name string) *Column {
	// clean text style to simplify unit test
	return NewStandardColumn(name).HeaderStyle(DefauleHeaderStyle().Text()).BodyStyle(DefauleBodyStyle().Text())
}
