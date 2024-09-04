package main

import (
	"fmt"

	gotable "github.com/darkelf21cn/go-table"
)

var (
	titleID           = "ID"
	titleMovie        = "Movie"
	titleScore        = "Score"
	titleIntroduction = "Introduction"
	titleActors       = "Actors"

	movieRow1 = []any{1, "The Godfather", 9.2}
	movieRow2 = []any{20, "霸王别姬", 8.1}
	movieRow3 = []any{30, "東京物語", 8.2}
	movieRow4 = []any{40, "La Haine", 7.8}
	movieRow5 = []any{500, "Life of Pi", 7.9}

	movieRow11 = []any{
		1,
		"The Godfather",
		"Francis Ford Coppola's masterpiece that chronicles the Corleone mafia family, is widely regarded as one of the greatest films in world cinema.",
		"Marlon Brando\nAl Pacino\nJames Caan",
	}
	movieRow12 = []any{
		20,
		"霸王别姬",
		"A poignant drama that tells the story of two performers in the Beijing Opera during the tumultuous events of the 20th century in China.",
		"张国荣\r\n张丰毅\n巩俐\r\n葛优"}
	movieRow13 = []any{
		30,
		"La Haine",
		"A 1995 film by Mathieu Kassovitz that depicts the tensions in the suburbs of Paris, following the lives of three friends over a day.",
		"Vincent Cassel\nHubert Koundé\nSaïd Taghmaoui",
	}
)

func main() {
	simple()
	complex()
}

// simple examples use builtin methods to modify table layout
func simple() {
	fmt.Println("Everyting in default settings")
	tb := gotable.NewTable(nil)
	col1 := gotable.NewStandardColumn(titleID)
	col2 := gotable.NewStandardColumn(titleMovie)
	col3 := gotable.NewStandardColumn(titleScore)
	tb.AppendColumn(col1, col2, col3)
	tb.AppendRow(movieRow1...)
	tb.AppendRow(movieRow2...)
	tb.AppendRow(movieRow3...)
	tb.AppendRow(movieRow4...)
	tb.AppendRow(movieRow5...)
	str, err := tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Use LightTableLayout")
	tb.Layout = *gotable.LightTableLayout()
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Hide header")
	tb.Layout = *gotable.LightTableLayout().HideHeader()
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Hide outter border")
	tb.Layout = *gotable.LightTableLayout().HideOutterBorder()
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Hide header and outter border")
	tb.Layout = *gotable.LightTableLayout().HideOutterBorder().HideHeader()
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Spearate header and body")
	tb.Layout = *gotable.LightTableLayout().SplitHeaderAndBody()
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Enforce table width to 80")
	tb.Layout = *gotable.LightTableLayout()
	tb.Layout.Width = 80
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}

// complex examples applies many customize settings to the table
func complex() {
	fmt.Println("Enforce table width to 80 but I don't want column width to be adjusted on a specific column")
	layout := gotable.LightTableLayout()
	tb := gotable.NewTable(layout)
	col1 := gotable.NewStandardColumn(titleID)
	col2 := gotable.NewStandardColumn(titleMovie).Width(0, false) // disable AutoWidthControl here
	col3 := gotable.NewStandardColumn(titleIntroduction)
	col4 := gotable.NewStandardColumn(titleActors)
	tb.AppendColumn(col1, col2, col3, col4)
	tb.AppendRow(movieRow11...)
	tb.AppendRow(movieRow12...)
	tb.AppendRow(movieRow13...)
	tb.Layout.Width = 80
	str, err := tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("I want text to be truncated when it exceeds the column width")
	layout = gotable.LightTableLayout()
	tb = gotable.NewTable(layout)
	tb.Layout.Width = 80
	col1 = gotable.NewStandardColumn(titleID)
	col2 = gotable.NewStandardColumn(titleMovie)
	col3 = gotable.NewStandardColumn(titleIntroduction).BodyStyle(gotable.DefauleBodyStyle().OverFlowAction(gotable.Truncate))
	col4 = gotable.NewStandardColumn(titleActors)
	tb.AppendColumn(col1, col2, col3, col4)
	tb.AppendRow(movieRow11...)
	tb.AppendRow(movieRow12...)
	tb.AppendRow(movieRow13...)
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	fmt.Println("Don't handle line feed in text")
	layout = gotable.LightTableLayout()
	tb = gotable.NewTable(layout)
	tb.Layout.Width = 80
	col1 = gotable.NewStandardColumn(titleID)
	col2 = gotable.NewStandardColumn(titleMovie)
	col3 = gotable.NewStandardColumn(titleIntroduction).BodyStyle(gotable.DefauleBodyStyle().OverFlowAction(gotable.Truncate))
	col4 = gotable.NewStandardColumn(titleActors).BodyStyle(gotable.DefauleBodyStyle().EscapeLineFeed(true))
	tb.AppendColumn(col1, col2, col3, col4)
	tb.AppendRow(movieRow11...)
	tb.AppendRow(movieRow12...)
	tb.AppendRow(movieRow13...)
	str, err = tb.Render(gotable.Console)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
