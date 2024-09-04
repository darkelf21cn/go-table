package gotable

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common Test Suites", func() {
	// splitStringByWidth
	It("splitStringByWidth-case1", func() {
		out := splitStringByWidth(strWordWrap, 10000)
		Expect(len(out)).Should(Equal(4))
	})
	It("splitStringByWidth-case2", func() {
		out := splitStringByWidth(strShort2, 5)
		Expect(len(out)).Should(Equal(3))
		Expect(out[0]).Should(Equal("ab cd"))
		Expect(out[1]).Should(Equal(" ef g"))
		Expect(out[2]).Should(Equal("h"))
	})
	It("splitStringByWidth-case3", func() {
		out := splitStringByWidth(strHelloChinese, 10)
		Expect(len(out)).Should(Equal(1))
		Expect(out[0]).Should(Equal("你好，世界"))
	})
	It("splitStringByWidth-case4", func() {
		out := splitStringByWidth(strHelloChinese, 9)
		Expect(len(out)).Should(Equal(2))
		Expect(out[0]).Should(Equal("你好，世"))
		Expect(out[1]).Should(Equal("界"))
	})

	// getTreeStatistics
	It("getTreeStatistics-t1", func() {
		stat := getTreeStatistics(nil)
		Expect(stat.MaxDepth).Should(Equal(0))
	})
	It("getTreeStatistics-t2", func() {
		nodes := []TreeNodeReader{
			&mockTreeNode{},
			&mockTreeNode{},
		}
		stat := getTreeStatistics(nodes)
		Expect(stat.MaxDepth).Should(Equal(1))
	})
	It("getTreeStatistics-t3", func() {
		nodes := []TreeNodeReader{
			&mockTreeNode{
				children: []TreeNodeReader{
					&mockTreeNode{},
					&mockTreeNode{},
				}},
			&mockTreeNode{},
		}
		stat := getTreeStatistics(nodes)
		Expect(stat.MaxDepth).Should(Equal(2))
	})
	It("getTreeStatistics-t4", func() {
		nodes := []TreeNodeReader{
			&mockTreeNode{},
			&mockTreeNode{
				children: []TreeNodeReader{
					&mockTreeNode{},
					&mockTreeNode{},
				}},
			&mockTreeNode{},
		}
		stat := getTreeStatistics(nodes)
		Expect(stat.MaxDepth).Should(Equal(2))
	})

	// formatAlignment
	It("formatAlignment-case1", func() {
		out := formatAlignment(strSingle, 4, ' ', AlignLeft)
		Expect(out).Should(Equal("a   "))
	})
	It("formatAlignment-case2", func() {
		out := formatAlignment(strSingle, 4, ' ', AlignRight)
		Expect(out).Should(Equal("   a"))
	})
	It("formatAlignment-case3", func() {
		out := formatAlignment(strSingle, 4, ' ', AlignCenter)
		Expect(out).Should(Equal(" a  "))
	})
	It("formatAlignment-case4", func() {
		out := formatAlignment(strSingle, 3, ' ', AlignCenter)
		Expect(out).Should(Equal(" a "))
	})
	It("formatAlignment-case5", func() {
		out := formatAlignment(strShort2, 16, ' ', AlignJustify)
		Expect(out).Should(Equal("ab   cd   ef  gh"))
	})
	It("formatAlignment-case6", func() {
		out := formatAlignment(strShort2, 17, ' ', AlignJustify)
		Expect(out).Should(Equal("ab   cd   ef   gh"))
	})
	It("formatAlignment-case7", func() {
		out := formatAlignment(strShort1, 4, ' ', AlignJustify)
		Expect(out).Should(Equal(strShort1))
	})
})
