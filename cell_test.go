package gotable

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cell Test Suites", func() {
	Context("DataCell", func() {
		// render
		It("render-case1", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Wordwrap, false)
			_, err := c.render(3, 1000, Console)
			Expect(errors.Is(err, ErrInsufficientColumnWidth)).Should(BeTrue())
		})
		It("render-case2", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Exception, true)
			_, err := c.render(4, 1000, Console)
			Expect(errors.Is(err, ErrInsufficientColumnWidth)).Should(BeTrue())
		})
		It("render-case3", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Truncate, true)
			out, err := c.render(4, 1, Console)
			expects := []string{
				"  ~ ",
			}
			Expect(err).Should(BeNil())
			Expect(out).Should(Equal(expects))
		})
		It("render-case4", func() {
			c := test_NewDataCell(strShort2, AlignLeft, Truncate, true)
			out, err := c.render(5, 2, Console)
			expects := []string{
				" a ~ ",
				"     ",
			}
			Expect(err).Should(BeNil())
			Expect(out).Should(Equal(expects))
		})
		It("render-case5", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Truncate, false)
			out, err := c.render(8, 2, Console)
			expects := []string{
				" 你好 ~ ",
				"        ",
			}
			Expect(err).Should(BeNil())
			Expect(out).Should(Equal(expects))
		})
		It("render-case6", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Truncate, true)
			out, err := c.render(9, 2, Console)
			expects := []string{
				" 你好  ~ ",
				"         ",
			}
			Expect(err).Should(BeNil())
			Expect(out).Should(Equal(expects))
		})
		It("render-case7", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Wordwrap, true)
			_, err := c.render(4, 4, Console)
			Expect(errors.Is(err, ErrInsufficientColumnHeight)).Should(BeTrue())
		})
		It("render-case8", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Wordwrap, false)
			out, err := c.render(4, 5, Console)
			Expect(err).Should(BeNil())
			Expect(len(out)).Should(Equal(5))
			expects := []string{" 你 ", " 好 ", " ， ", " 世 ", " 界 "}
			Expect(out).Should(Equal(expects))
		})
		It("render-case9", func() {
			c := test_NewDataCell(strHelloChinese, AlignLeft, Wordwrap, true)
			out, err := c.render(4, 6, Console)
			Expect(err).Should(BeNil())
			Expect(len(out)).Should(Equal(6))
			expects := []string{" 你 ", " 好 ", " ， ", " 世 ", " 界 ", "    "}
			Expect(out).Should(Equal(expects))
		})
		It("render-case10", func() {
			c := test_NewDataCell(strWordWrap, AlignJustify, Wordwrap, true)
			out, err := c.render(80, 6, Console)
			Expect(err).Should(BeNil())
			Expect(len(out)).Should(Equal(6))
			expects := []string{
				" Lorem     ipsum    dolor    sit    amet,    consectetur    adipiscing    elit. ",
				" Nulla        eget       mi       nec       ipsum       aliquam       pulvinar. ",
				" Aenean     id    justo    ac    diam    iaculis    gravida    nec    et    ex. ",
				" Fusce    sed    quam   hendrerit,   mollis   nisi   vitae,   porttitor   erat. ",
				"                                                                                ",
				"                                                                                ",
			}
			Expect(out).Should(Equal(expects))
		})
		It("render-case11", func() {
			c := test_NewDataCell(strWordWrap, AlignLeft, Wordwrap, false)
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			out, err := c.render(w, h, Console)
			Expect(err).Should(BeNil())
			expects := []string{
				" Lorem ipsum dolor sit amet, consectetur adipiscing elit.     ",
				" Nulla eget mi nec ipsum aliquam pulvinar.                    ",
				" Aenean id justo ac diam iaculis gravida nec et ex.           ",
				" Fusce sed quam hendrerit, mollis nisi vitae, porttitor erat. ",
			}
			Expect(out).Should(Equal(expects))
		})
		It("render-case12", func() {
			c := test_NewDataCell(strSQL, AlignLeft, Wordwrap, true)
			c.leftPadding = ""
			c.rightPadding = ""
			out, err := c.render(20, 8, Console)
			Expect(err).Should(BeNil())
			Expect(len(out)).Should(Equal(8))
			expects := []string{
				"SELECT              ",
				"    *               ",
				"FROM                ",
				"    information_sche",
				"ma.tables           ",
				"WHERE               ",
				"    TABLE_SCHEMA = '",
				"mysql'              ",
			}
			Expect(out).Should(Equal(expects))
		})
		It("render-case13", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Wordwrap, false)
			c.leftPadding = ">>"
			c.rightPadding = "<<"
			out, err := c.render(40, 9, Console)
			Expect(err).Should(BeNil())
			Expect(len(out)).Should(Equal(9))
			expects := []string{
				">>Lorem ipsum dolor sit amet, consecte<<",
				">>        tur adipiscing elit.        <<",
				">>Nulla eget mi nec ipsum aliquam pulv<<",
				">>               inar.                <<",
				">>Aenean id justo ac diam iaculis grav<<",
				">>           ida nec et ex.           <<",
				">>Fusce sed quam hendrerit, mollis nis<<",
				">>      i vitae, porttitor erat.      <<",
				">>                                    <<",
			}
			Expect(out).Should(Equal(expects))
		})
		It("render-case14", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Exception, true)
			c.leftPadding = "|"
			c.rightPadding = "|"
			out, err := c.render(250, 1, Console)
			expects := []string{
				"|                 Lorem ipsum dolor sit amet, consectetur adipiscing elit.\\nNulla eget mi nec ipsum aliquam pulvinar.\\nAenean id justo ac diam iaculis gravida nec et ex.\\nFusce sed quam hendrerit, mollis nisi vitae, porttitor erat.                  |",
			}
			Expect(err).Should(BeNil())
			Expect(out).Should(Equal(expects))
		})
		It("render-case15", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Truncate, true)
			c.leftPadding = "|"
			c.rightPadding = "|"
			out, err := c.render(80, 1, Console)
			expects := []string{
				"|Lorem ipsum dolor sit amet, consectetur adipiscing elit.\\nNulla eget mi nec  ~|",
			}
			Expect(err).Should(BeNil())
			Expect(out).Should(Equal(expects))
		})
		It("render-case16", func() {
			c := test_NewDataCell(strShort2, AlignCenter, Wordwrap, false)
			c.Style(Bold, Red, BgBlue)
			out, err := c.render(13, 1, Console)
			Expect(err).Should(BeNil())
			expects := []string{
				"\x1b[44m\x1b[31m\x1b[1m ab cd ef gh \x1b[22m\x1b[0m\x1b[0m",
			}
			Expect(out).Should(Equal(expects))
		})

		It("stats-case1", func() {
			c := test_NewDataCell(strShort2, AlignCenter, Truncate, true)
			w, h, err := c.stats(100, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(100))
			Expect(h).Should(Equal(1))
		})
		It("stats-case2", func() {
			c := test_NewDataCell(strHelloChinese, AlignCenter, Truncate, false)
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(12))
			Expect(h).Should(Equal(1))
		})
		It("stats-case3", func() {
			c := test_NewDataCell(strHelloChinese, AlignCenter, Wordwrap, true)
			w, h, err := c.stats(6, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(6))
			Expect(h).Should(Equal(3))
		})
		It("stats-case4", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Truncate, false)
			w, h, err := c.stats(20, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(20))
			Expect(h).Should(Equal(4))
		})
		It("stats-case5", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Wordwrap, true)
			w, h, err := c.stats(20, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(20))
			Expect(h).Should(Equal(14))
		})
		It("stats-case6", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Wordwrap, false)
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(62))
			Expect(h).Should(Equal(4))
		})

		It("mix-case1", func() {
			c := test_NewDataCell(strWordWrap, AlignCenter, Truncate, false)
			w, h, err := c.stats(20, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(20))
			Expect(h).Should(Equal(4))
			out, err := c.render(w, h, Console)
			Expect(err).Should(BeNil())
			expects := []string{
				" Lorem ipsum dolo ~ ",
				" Nulla eget mi ne ~ ",
				" Aenean id justo  ~ ",
				" Fusce sed quam h ~ ",
			}
			Expect(out).Should(Equal(expects))
		})
	})

	Context("TreePathCell", func() {
		It("stats-case1", func() {
			c := test_NewTreePathCell(nil, ' ', LightTreePathStyle())
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(7))
			Expect(h).Should(Equal(1))
		})
		It("stats-case2", func() {
			c := test_NewTreePathCell("路径", '*', LightTreePathStyle())
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(6))
			Expect(h).Should(Equal(1))
		})
		It("stats-case3", func() {
			c := test_NewTreePathCell("□─┬──────", '*', LightTreePathStyle())
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(11))
			Expect(h).Should(Equal(1))
		})
		It("stats-case4", func() {
			c := test_NewTreePathCell("□─", '*', LightTreePathStyle())
			w, h, err := c.stats(0, Console)
			Expect(err).Should(BeNil())
			Expect(w).Should(Equal(4))
			Expect(h).Should(Equal(1))
		})

		It("mix-case1", func() {
			sty := LightTreePathStyle()
			cells := []*Cell{
				test_NewTreePathCell("□─┬──────", ' ', sty),
				test_NewTreePathCell("  ├─┬────", ' ', sty),
				test_NewTreePathCell("  │ └─┬──", ' ', sty),
				test_NewTreePathCell("  │   └──", ' ', sty),
				test_NewTreePathCell("  └──────", ' ', sty),
				test_NewTreePathCell("□────────", ' ', sty),
			}
			expects := [][]string{
				{
					" □─┬────── ",
					"   │       ",
				},
				{
					"   ├─┬──── ",
					"   │ │     ",
				},
				{
					"   │ └─┬── ",
					"   │   │   ",
				},
				{
					"   │   └── ",
					"   │       ",
				},
				{
					"   └────── ",
					"           ",
				},
				{
					" □──────── ",
					"           ",
				},
			}
			for i, cell := range cells {
				w, _, err := cell.stats(0, Console)
				Expect(err).Should(BeNil())
				out, err := cell.render(w, 2, Console)
				Expect(err).Should(BeNil())
				Expect(out).Should(Equal(expects[i]))
			}
		})
	})
})

func test_NewDataCell(data any, al Align, ofa ColumnOverFlowAction, lf bool) *Cell {
	if ofa == Wordwrap {
		lf = false
	}
	c := &Cell{
		leftPadding:  " ",
		rightPadding: " ",
		cellRenderer: &DataCell{
			padding:        ' ',
			align:          al,
			overFlowAction: ofa,
			escapeLineFeed: lf,
		},
	}
	c.Value(data)
	return c
}

func test_NewTreePathCell(data any, padding rune, sty *TreePathStyle) *Cell {
	c := &Cell{
		leftPadding:  string(padding),
		rightPadding: string(padding),
		cellRenderer: &TreePathCell{
			style: sty,
		},
	}
	c.Value(data)
	return c
}
