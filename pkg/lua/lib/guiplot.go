package lib

import (
	g "github.com/AllenDang/giu"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_GUIPLOT = "guiplot"

/// @lib GUI Plots
/// @import guiplot
/// @desc
/// Extension of the GUI library for plot widgets.

func RegisterGUIPlot(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_GUIPLOT, r, r.State, lg)

	/// @func wg_plot(title) -> struct<guiplot.WidgetPlot>
	/// @arg title {string}
	/// @returns {struct<guiplot.WidgetPlot>}
	lib.CreateFunction(tab, "wg_plot",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotTable(state, args["title"].(string))

			state.Push(t)
			return 1
		})

	/// @func plot_ticker(position, label) -> struct<guiplot.PlotTicker>
	/// @arg position {float}
	/// @arg label {string}
	/// @returns {struct<guiplot.PlotTicker>}
	lib.CreateFunction(tab, "plot_ticker",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "position"},
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pos := args["position"].(float64)
			label := args["label"].(string)

			t := state.NewTable()
			t.RawSetString("position", golua.LNumber(pos))
			t.RawSetString("label", golua.LString(label))

			state.Push(t)
			return 1
		})

	/// @func pt_bar_h(title, data) -> struct<guiplot.PlotBarH>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<guiplot.PlotBarH>}
	lib.CreateFunction(tab, "pt_bar_h",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotBarHTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_bar(title, data) -> struct<guiplot.PlotBar>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<guiplot.PlotBar>}
	lib.CreateFunction(tab, "pt_bar",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotBarTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_line(title, data) -> struct<guiplot.PlotLine>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<guiplot.PlotLine>}
	lib.CreateFunction(tab, "pt_line",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotLineTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_line_xy(title, xdata, ydata) -> struct<guiplot.PlotLineXY>
	/// @arg title {string}
	/// @arg xdata {[]float}
	/// @arg ydata {[]float}
	/// @returns {struct<guiplot.PlotLineXY>}
	lib.CreateFunction(tab, "pt_line_xy",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "xdata"},
			{Type: lua.RAW_TABLE, Name: "ydata"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotLineXYTable(state, args["title"].(string), args["xdata"].(golua.LValue), args["ydata"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_pie_chart(labels, data, x, y, radius) -> struct<guiplot.PlotPieChart>
	/// @arg labels {[]string}
	/// @arg data {[]float}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg radius {float}
	/// @returns {struct<guiplot.PlotPieChart>}
	lib.CreateFunction(tab, "pt_pie_chart",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "labels"},
			{Type: lua.RAW_TABLE, Name: "data"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "radius"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			labels := args["labels"].(golua.LValue)
			data := args["data"].(golua.LValue)
			x := args["x"].(float64)
			y := args["y"].(float64)
			radius := args["radius"].(float64)
			t := plotPieTable(state, labels, data, x, y, radius)

			state.Push(t)
			return 1
		})

	/// @func pt_scatter(title, data) -> struct<guiplot.PlotScatter>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<guiplot.PlotScatter>}
	lib.CreateFunction(tab, "pt_scatter",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotScatterTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_scatter_xy(title, xdata, ydata) -> struct<guiplot.PlotScatterXY>
	/// @arg title {string}
	/// @arg xdata {[]float}
	/// @arg ydata {[]float}
	/// @returns {struct<guiplot.PlotScatterXY>}
	lib.CreateFunction(tab, "pt_scatter_xy",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "xdata"},
			{Type: lua.RAW_TABLE, Name: "ydata"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotScatterXYTable(state, args["title"].(string), args["xdata"].(golua.LValue), args["ydata"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_custom(builder) -> struct<guiplot.PlotCustom>
	/// @arg builder {function()}
	/// @returns {struct<guiplot.PlotCustom>}
	lib.CreateFunction(tab, "pt_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "builder"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotCustomTable(state, args["builder"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @constants PlotFlags {int}
	/// @const FLAGPLOT_NONE
	/// @const FLAGPLOT_NOTITLE
	/// @const FLAGPLOT_NOLEGEND
	/// @const FLAGPLOT_NOMOUSETEXT
	/// @const FLAGPLOT_NOINPUTS
	/// @const FLAGPLOT_NOMENUS
	/// @const FLAGPLOT_NOBOXSELECT
	/// @const FLAGPLOT_NOFRAME
	/// @const FLAGPLOT_EQUAL
	/// @const FLAGPLOT_CROSSHAIRS
	/// @const FLAGPLOT_CANVASONLY
	tab.RawSetString("FLAGPLOT_NONE", golua.LNumber(FLAGPLOT_NONE))
	tab.RawSetString("FLAGPLOT_NOTITLE", golua.LNumber(FLAGPLOT_NOTITLE))
	tab.RawSetString("FLAGPLOT_NOLEGEND", golua.LNumber(FLAGPLOT_NOLEGEND))
	tab.RawSetString("FLAGPLOT_NOMOUSETEXT", golua.LNumber(FLAGPLOT_NOMOUSETEXT))
	tab.RawSetString("FLAGPLOT_NOINPUTS", golua.LNumber(FLAGPLOT_NOINPUTS))
	tab.RawSetString("FLAGPLOT_NOMENUS", golua.LNumber(FLAGPLOT_NOMENUS))
	tab.RawSetString("FLAGPLOT_NOBOXSELECT", golua.LNumber(FLAGPLOT_NOBOXSELECT))
	tab.RawSetString("FLAGPLOT_NOFRAME", golua.LNumber(FLAGPLOT_NOFRAME))
	tab.RawSetString("FLAGPLOT_EQUAL", golua.LNumber(FLAGPLOT_EQUAL))
	tab.RawSetString("FLAGPLOT_CROSSHAIRS", golua.LNumber(FLAGPLOT_CROSSHAIRS))
	tab.RawSetString("FLAGPLOT_CANVASONLY", golua.LNumber(FLAGPLOT_CANVASONLY))

	/// @constants PlotAxis {int}
	/// @const PLOTAXIS_X1
	/// @const PLOTAXIS_X2
	/// @const PLOTAXIS_X3
	/// @const PLOTAXIS_Y1
	/// @const PLOTAXIS_Y2
	/// @const PLOTAXIS_Y3
	/// @const PLOTAXIS_COUNT
	tab.RawSetString("PLOTAXIS_X1", golua.LNumber(PLOTAXIS_X1))
	tab.RawSetString("PLOTAXIS_X2", golua.LNumber(PLOTAXIS_X2))
	tab.RawSetString("PLOTAXIS_X3", golua.LNumber(PLOTAXIS_X3))
	tab.RawSetString("PLOTAXIS_Y1", golua.LNumber(PLOTAXIS_Y1))
	tab.RawSetString("PLOTAXIS_Y2", golua.LNumber(PLOTAXIS_Y2))
	tab.RawSetString("PLOTAXIS_Y3", golua.LNumber(PLOTAXIS_Y3))
	tab.RawSetString("PLOTAXIS_COUNT", golua.LNumber(PLOTAXIS_COUNT))

	/// @constants PlotAxisFlags {int}
	/// @const FLAGPLOTAXIS_NONE
	/// @const FLAGPLOTAXIS_NOLABEL
	/// @const FLAGPLOTAXIS_NOGRIDLINES
	/// @const FLAGPLOTAXIS_NOTICKMARKS
	/// @const FLAGPLOTAXIS_NOTICKLABELS
	/// @const FLAGPLOTAXIS_NOINITIALFIT
	/// @const FLAGPLOTAXIS_NOMENUS
	/// @const FLAGPLOTAXIS_NOSIDESWITCH
	/// @const FLAGPLOTAXIS_NOHIGHLIGHT
	/// @const FLAGPLOTAXIS_OPPOSITE
	/// @const FLAGPLOTAXIS_FOREGROUND
	/// @const FLAGPLOTAXIS_INVERT
	/// @const FLAGPLOTAXIS_AUTOFIT
	/// @const FLAGPLOTAXIS_RANGEFIT
	/// @const FLAGPLOTAXIS_PANSTRETCH
	/// @const FLAGPLOTAXIS_LOCKMIN
	/// @const FLAGPLOTAXIS_LOCKMAX
	/// @const FLAGPLOTAXIS_LOCK
	/// @const FLAGPLOTAXIS_NODECORATIONS
	/// @const FLAGPLOTAXIS_AUXDEFAULT
	tab.RawSetString("FLAGPLOTAXIS_NONE", golua.LNumber(FLAGPLOTAXIS_NONE))
	tab.RawSetString("FLAGPLOTAXIS_NOLABEL", golua.LNumber(FLAGPLOTAXIS_NOLABEL))
	tab.RawSetString("FLAGPLOTAXIS_NOGRIDLINES", golua.LNumber(FLAGPLOTAXIS_NOGRIDLINES))
	tab.RawSetString("FLAGPLOTAXIS_NOTICKMARKS", golua.LNumber(FLAGPLOTAXIS_NOTICKMARKS))
	tab.RawSetString("FLAGPLOTAXIS_NOTICKLABELS", golua.LNumber(FLAGPLOTAXIS_NOTICKLABELS))
	tab.RawSetString("FLAGPLOTAXIS_NOINITIALFIT", golua.LNumber(FLAGPLOTAXIS_NOINITIALFIT))
	tab.RawSetString("FLAGPLOTAXIS_NOMENUS", golua.LNumber(FLAGPLOTAXIS_NOMENUS))
	tab.RawSetString("FLAGPLOTAXIS_NOSIDESWITCH", golua.LNumber(FLAGPLOTAXIS_NOSIDESWITCH))
	tab.RawSetString("FLAGPLOTAXIS_NOHIGHLIGHT", golua.LNumber(FLAGPLOTAXIS_NOHIGHLIGHT))
	tab.RawSetString("FLAGPLOTAXIS_OPPOSITE", golua.LNumber(FLAGPLOTAXIS_OPPOSITE))
	tab.RawSetString("FLAGPLOTAXIS_FOREGROUND", golua.LNumber(FLAGPLOTAXIS_FOREGROUND))
	tab.RawSetString("FLAGPLOTAXIS_INVERT", golua.LNumber(FLAGPLOTAXIS_INVERT))
	tab.RawSetString("FLAGPLOTAXIS_AUTOFIT", golua.LNumber(FLAGPLOTAXIS_AUTOFIT))
	tab.RawSetString("FLAGPLOTAXIS_RANGEFIT", golua.LNumber(FLAGPLOTAXIS_RANGEFIT))
	tab.RawSetString("FLAGPLOTAXIS_PANSTRETCH", golua.LNumber(FLAGPLOTAXIS_PANSTRETCH))
	tab.RawSetString("FLAGPLOTAXIS_LOCKMIN", golua.LNumber(FLAGPLOTAXIS_LOCKMIN))
	tab.RawSetString("FLAGPLOTAXIS_LOCKMAX", golua.LNumber(FLAGPLOTAXIS_LOCKMAX))
	tab.RawSetString("FLAGPLOTAXIS_LOCK", golua.LNumber(FLAGPLOTAXIS_LOCK))
	tab.RawSetString("FLAGPLOTAXIS_NODECORATIONS", golua.LNumber(FLAGPLOTAXIS_NODECORATIONS))
	tab.RawSetString("FLAGPLOTAXIS_AUXDEFAULT", golua.LNumber(FLAGPLOTAXIS_AUXDEFAULT))

	/// @constants PlotYAxis {int}
	/// @const PLOTYAXIS_LEFT
	/// @const PLOTYAXIS_FIRSTONRIGHT
	/// @const PLOTYAXIS_SECONDONRIGHT
	tab.RawSetString("PLOTYAXIS_LEFT", golua.LNumber(PLOTYAXIS_LEFT))
	tab.RawSetString("PLOTYAXIS_FIRSTONRIGHT", golua.LNumber(PLOTYAXIS_FIRSTONRIGHT))
	tab.RawSetString("PLOTYAXIS_SECONDONRIGHT", golua.LNumber(PLOTYAXIS_SECONDONRIGHT))

	/// @constants PlotType {string}
	/// @const PLOT_BAR_H
	/// @const PLOT_BAR
	/// @const PLOT_LINE
	/// @const PLOT_LINE_XY
	/// @const PLOT_PIE_CHART
	/// @const PLOT_SCATTER
	/// @const PLOT_SCATTER_XY
	/// @const PLOT_CUSTOM
	/// @const PLOT_STYLE
	tab.RawSetString("PLOT_BAR_H", golua.LString(PLOT_BAR_H))
	tab.RawSetString("PLOT_BAR", golua.LString(PLOT_BAR))
	tab.RawSetString("PLOT_LINE", golua.LString(PLOT_LINE))
	tab.RawSetString("PLOT_LINE_XY", golua.LString(PLOT_LINE_XY))
	tab.RawSetString("PLOT_PIE_CHART", golua.LString(PLOT_PIE_CHART))
	tab.RawSetString("PLOT_SCATTER", golua.LString(PLOT_SCATTER))
	tab.RawSetString("PLOT_SCATTER_XY", golua.LString(PLOT_SCATTER_XY))
	tab.RawSetString("PLOT_CUSTOM", golua.LString(PLOT_CUSTOM))
	tab.RawSetString("PLOT_STYLE", golua.LString(PLOT_STYLE))

	/// @constants StylePlotColorID {int}
	/// @const COLIDPLOT_LINE
	/// @const COLIDPLOT_FILL
	/// @const COLIDPLOT_MARKEROUTLINE
	/// @const COLIDPLOT_MARKERFILL
	/// @const COLIDPLOT_ERRORBAR
	/// @const COLIDPLOT_FRAMEBG
	/// @const COLIDPLOT_PLOTBG
	/// @const COLIDPLOT_PLOTBORDER
	/// @const COLIDPLOT_LEGENDBG
	/// @const COLIDPLOT_LEGENDBORDER
	/// @const COLIDPLOT_LEGENDTEXT
	/// @const COLIDPLOT_TITLETEXT
	/// @const COLIDPLOT_INLAYTEXT
	/// @const COLIDPLOT_AXISTEXT
	/// @const COLIDPLOT_AXISGRID
	/// @const COLIDPLOT_AXISTICK
	/// @const COLIDPLOT_AXISBG
	/// @const COLIDPLOT_AXISBGHOVERED
	/// @const COLIDPLOT_AXISBGACTIVE
	/// @const COLIDPLOT_SELECTION
	/// @const COLIDPLOT_CROSSHAIRS
	tab.RawSetString("COLIDPLOT_LINE", golua.LNumber(g.StylePlotColorLine))
	tab.RawSetString("COLIDPLOT_FILL", golua.LNumber(g.StylePlotColorFill))
	tab.RawSetString("COLIDPLOT_MARKEROUTLINE", golua.LNumber(g.StylePlotColorMarkerOutline))
	tab.RawSetString("COLIDPLOT_MARKERFILL", golua.LNumber(g.StylePlotColorMarkerFill))
	tab.RawSetString("COLIDPLOT_ERRORBAR", golua.LNumber(g.StylePlotColorErrorBar))
	tab.RawSetString("COLIDPLOT_FRAMEBG", golua.LNumber(g.StylePlotColorFrameBg))
	tab.RawSetString("COLIDPLOT_PLOTBG", golua.LNumber(g.StylePlotColorPlotBg))
	tab.RawSetString("COLIDPLOT_PLOTBORDER", golua.LNumber(g.StylePlotColorPlotBorder))
	tab.RawSetString("COLIDPLOT_LEGENDBG", golua.LNumber(g.StylePlotColorLegendBg))
	tab.RawSetString("COLIDPLOT_LEGENDBORDER", golua.LNumber(g.StylePlotColorLegendBorder))
	tab.RawSetString("COLIDPLOT_LEGENDTEXT", golua.LNumber(g.StylePlotColorLegendText))
	tab.RawSetString("COLIDPLOT_TITLETEXT", golua.LNumber(g.StylePlotColorTitleText))
	tab.RawSetString("COLIDPLOT_INLAYTEXT", golua.LNumber(g.StylePlotColorInlayText))
	tab.RawSetString("COLIDPLOT_AXISTEXT", golua.LNumber(g.StylePlotColorAxisText))
	tab.RawSetString("COLIDPLOT_AXISGRID", golua.LNumber(g.StylePlotColorAxisGrid))
	tab.RawSetString("COLIDPLOT_AXISTICK", golua.LNumber(g.StylePlotColorAxisTick))
	tab.RawSetString("COLIDPLOT_AXISBG", golua.LNumber(g.StylePlotColorAxisBg))
	tab.RawSetString("COLIDPLOT_AXISBGHOVERED", golua.LNumber(g.StylePlotColorAxisBgHovered))
	tab.RawSetString("COLIDPLOT_AXISBGACTIVE", golua.LNumber(g.StylePlotColorAxisBgActive))
	tab.RawSetString("COLIDPLOT_SELECTION", golua.LNumber(g.StylePlotColorSelection))
	tab.RawSetString("COLIDPLOT_CROSSHAIRS", golua.LNumber(g.StylePlotColorCrosshairs))

	/// @constants StylePlotVar {int}
	/// @const STYLEPLOTVAR_LINEWEIGHT
	/// @const STYLEPLOTVAR_MARKER
	/// @const STYLEPLOTVAR_MARKERSIZE
	/// @const STYLEPLOTVAR_FILLALPHA
	/// @const STYLEPLOTVAR_ERRORBARSIZE
	/// @const STYLEPLOTVAR_ERRORBARWEIGHT
	/// @const STYLEPLOTVAR_DIGITALBITHEIGHT
	/// @const STYLEPLOTVAR_DIGITALBITGAP
	/// @const STYLEPLOTVAR_PLOTBORDERSIZE
	/// @const STYLEPLOTVAR_MINORALPHA
	/// @const STYLEPLOTVAR_MAJORTICKLEN
	/// @const STYLEPLOTVAR_MINORTICKLEN
	/// @const STYLEPLOTVAR_MAJORTICKSIZE
	/// @const STYLEPLOTVAR_MINORTICKSIZE
	/// @const STYLEPLOTVAR_MAJORGRIDSIZE
	/// @const STYLEPLOTVAR_MINORGRIDSIZE
	/// @const STYLEPLOTVAR_PLOTPADDING
	/// @const STYLEPLOTVAR_LABELPADDING
	/// @const STYLEPLOTVAR_LEGENDPADDING
	/// @const STYLEPLOTVAR_LEGENDINNERPADDING
	/// @const STYLEPLOTVAR_LEGENDSPACING
	/// @const STYLEPLOTVAR_MOUSEPOSPADDING
	/// @const STYLEPLOTVAR_ANNOTAIONPADDING
	/// @const STYLEPLOTVAR_FITPADDING
	/// @const STYLEPLOTVAR_PLOTDEFAULTSIZE
	/// @const STYLEPLOTVAR_PLOTMINSIZE
	tab.RawSetString("STYLEPLOTVAR_LINEWEIGHT", golua.LNumber(g.StylePlotVarLineWeight))
	tab.RawSetString("STYLEPLOTVAR_MARKER", golua.LNumber(g.StylePlotVarMarker))
	tab.RawSetString("STYLEPLOTVAR_MARKERSIZE", golua.LNumber(g.StylePlotVarMarkerSize))
	tab.RawSetString("STYLEPLOTVAR_FILLALPHA", golua.LNumber(g.StylePlotVarFillAlpha))
	tab.RawSetString("STYLEPLOTVAR_ERRORBARSIZE", golua.LNumber(g.StylePlotVarErrorBarSize))
	tab.RawSetString("STYLEPLOTVAR_ERRORBARWEIGHT", golua.LNumber(g.StylePlotVarErrorBarWeight))
	tab.RawSetString("STYLEPLOTVAR_DIGITALBITHEIGHT", golua.LNumber(g.StylePlotVarDigitalBitHeight))
	tab.RawSetString("STYLEPLOTVAR_DIGITALBITGAP", golua.LNumber(g.StylePlotVarDigitalBitGap))
	tab.RawSetString("STYLEPLOTVAR_PLOTBORDERSIZE", golua.LNumber(g.StylePlotVarPlotBorderSize))
	tab.RawSetString("STYLEPLOTVAR_MINORALPHA", golua.LNumber(g.StylePlotVarMinorAlpha))
	tab.RawSetString("STYLEPLOTVAR_MAJORTICKLEN", golua.LNumber(g.StylePlotVarMajorTickLen))
	tab.RawSetString("STYLEPLOTVAR_MINORTICKLEN", golua.LNumber(g.StylePlotVarMinorTickLen))
	tab.RawSetString("STYLEPLOTVAR_MAJORTICKSIZE", golua.LNumber(g.StylePlotVarMajorTickSize))
	tab.RawSetString("STYLEPLOTVAR_MINORTICKSIZE", golua.LNumber(g.StylePlotVarMinorTickSize))
	tab.RawSetString("STYLEPLOTVAR_MAJORGRIDSIZE", golua.LNumber(g.StylePlotVarMajorGridSize))
	tab.RawSetString("STYLEPLOTVAR_MINORGRIDSIZE", golua.LNumber(g.StylePlotVarMinorGridSize))
	tab.RawSetString("STYLEPLOTVAR_PLOTPADDING", golua.LNumber(g.StylePlotVarPlotPadding))
	tab.RawSetString("STYLEPLOTVAR_LABELPADDING", golua.LNumber(g.StylePlotVarLabelPadding))
	tab.RawSetString("STYLEPLOTVAR_LEGENDPADDING", golua.LNumber(g.StylePlotVarLegendPadding))
	tab.RawSetString("STYLEPLOTVAR_LEGENDINNERPADDING", golua.LNumber(g.StylePlotVarLegendInnerPadding))
	tab.RawSetString("STYLEPLOTVAR_LEGENDSPACING", golua.LNumber(g.StylePlotVarLegendSpacing))
	tab.RawSetString("STYLEPLOTVAR_MOUSEPOSPADDING", golua.LNumber(g.StylePlotVarMousePosPadding))
	tab.RawSetString("STYLEPLOTVAR_ANNOTAIONPADDING", golua.LNumber(g.StylePlotVarAnnotationPadding))
	tab.RawSetString("STYLEPLOTVAR_FITPADDING", golua.LNumber(g.StylePlotVarFitPadding))
	tab.RawSetString("STYLEPLOTVAR_PLOTDEFAULTSIZE", golua.LNumber(g.StylePlotVarPlotDefaultSize))
	tab.RawSetString("STYLEPLOTVAR_PLOTMINSIZE", golua.LNumber(g.StylePlotVarPlotMinSize))
}

const (
	FLAGPLOT_NONE        int = 0b0000_0000_0000
	FLAGPLOT_NOTITLE     int = 0b0000_0000_0001
	FLAGPLOT_NOLEGEND    int = 0b0000_0000_0010
	FLAGPLOT_NOMOUSETEXT int = 0b0000_0000_0100
	FLAGPLOT_NOINPUTS    int = 0b0000_0000_1000
	FLAGPLOT_NOMENUS     int = 0b0000_0001_0000
	FLAGPLOT_NOBOXSELECT int = 0b0000_0010_0000
	FLAGPLOT_NOFRAME     int = 0b0000_0100_0000
	FLAGPLOT_EQUAL       int = 0b0000_1000_0000
	FLAGPLOT_CROSSHAIRS  int = 0b0001_0000_0000
	FLAGPLOT_CANVASONLY  int = 0b0000_0011_0111
)

const (
	PLOTAXIS_X1 int = iota
	PLOTAXIS_X2
	PLOTAXIS_X3
	PLOTAXIS_Y1
	PLOTAXIS_Y2
	PLOTAXIS_Y3
	PLOTAXIS_COUNT
)

const (
	FLAGPLOTAXIS_NONE          int = 0b0000_0000_0000_0000
	FLAGPLOTAXIS_NOLABEL       int = 0b0000_0000_0000_0001
	FLAGPLOTAXIS_NOGRIDLINES   int = 0b0000_0000_0000_0010
	FLAGPLOTAXIS_NOTICKMARKS   int = 0b0000_0000_0000_0100
	FLAGPLOTAXIS_NOTICKLABELS  int = 0b0000_0000_0000_1000
	FLAGPLOTAXIS_NOINITIALFIT  int = 0b0000_0000_0001_0000
	FLAGPLOTAXIS_NOMENUS       int = 0b0000_0000_0010_0000
	FLAGPLOTAXIS_NOSIDESWITCH  int = 0b0000_0000_0100_0000
	FLAGPLOTAXIS_NOHIGHLIGHT   int = 0b0000_0000_1000_0000
	FLAGPLOTAXIS_OPPOSITE      int = 0b0000_0001_0000_0000
	FLAGPLOTAXIS_FOREGROUND    int = 0b0000_0010_0000_0000
	FLAGPLOTAXIS_INVERT        int = 0b0000_0100_0000_0000
	FLAGPLOTAXIS_AUTOFIT       int = 0b0000_1000_0000_0000
	FLAGPLOTAXIS_RANGEFIT      int = 0b0001_0000_0000_0000
	FLAGPLOTAXIS_PANSTRETCH    int = 0b0010_0000_0000_0000
	FLAGPLOTAXIS_LOCKMIN       int = 0b0100_0000_0000_0000
	FLAGPLOTAXIS_LOCKMAX       int = 0b1000_0000_0000_0000
	FLAGPLOTAXIS_LOCK          int = 0b1100_0000_0000_0000
	FLAGPLOTAXIS_NODECORATIONS int = 0b0000_0000_0000_1111
	FLAGPLOTAXIS_AUXDEFAULT    int = 0b0000_0001_0000_0010
)

const (
	PLOTYAXIS_LEFT          int = 0
	PLOTYAXIS_FIRSTONRIGHT  int = 1
	PLOTYAXIS_SECONDONRIGHT int = 2
)

const (
	PLOT_BAR_H      = "plot_bar_h"
	PLOT_BAR        = "plot_bar"
	PLOT_LINE       = "plot_line"
	PLOT_LINE_XY    = "plot_line_xy"
	PLOT_PIE_CHART  = "plot_pie_chart"
	PLOT_SCATTER    = "plot_scatter"
	PLOT_SCATTER_XY = "plot_scatter_xy"
	PLOT_CUSTOM     = "plot_custom"
	PLOT_STYLE      = WIDGET_STYLE
)

var plotList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget{}

func init() {
	plotList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget{
		PLOT_BAR_H:      plotBarHBuild,
		PLOT_BAR:        plotBarBuild,
		PLOT_LINE:       plotLineBuild,
		PLOT_LINE_XY:    plotLineXYBuild,
		PLOT_PIE_CHART:  plotPieBuild,
		PLOT_SCATTER:    plotScatterBuild,
		PLOT_SCATTER_XY: plotScatterXYBuild,
		PLOT_CUSTOM:     plotCustomBuild,
		PLOT_STYLE:      plotStyleBuild,
	}
}

func plotTable(state *golua.LState, title string) *golua.LTable {
	/// @struct WidgetPlot
	/// @prop type {string<gui.WidgetType>}
	/// @prop title {string}
	/// @method axis_limits(self, xmin float, xmax float, ymin float, ymax float, cond int<gui.Condition>) -> self
	/// @method flags(self, flags int<guiplot.PlotFlags>) -> self
	/// @method set_xaxis_label(self, axis int<guiplot.PlotAxis>, label string) -> self
	/// @method set_yaxis_label(self, axis int<guiplot.PlotAxis>, label string) -> self
	/// @method size(self, width float, height float) -> self
	/// @method x_axeflags(self, flags int<guiplot.PlotAxisFlags>) -> self
	/// @method xticks(self, ticks []struct<guiplot.PlotTicker>, default bool) -> self
	/// @method y_axeflags(self, flags1 int<guiplot.PlotAxisFlags>, flags2 int<guiplot.PlotAxisFlags>, flags3 int<guiplot.PlotAxisFlags>) -> self
	/// @method yticks(self, ticks []struct<guiplot.PlotTicker>) -> self
	/// @method plots(self, []struct<guiplot.Plot>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_PLOT))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("__xmin", golua.LNil)
	t.RawSetString("__xmax", golua.LNil)
	t.RawSetString("__ymin", golua.LNil)
	t.RawSetString("__ymax", golua.LNil)
	t.RawSetString("__cond", golua.LNil)
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__xlabels", state.NewTable())
	t.RawSetString("__ylabels", state.NewTable())
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__xaxeflags", golua.LNil)
	t.RawSetString("__xticks", golua.LNil)
	t.RawSetString("__xaticksdefault", golua.LNil)
	t.RawSetString("__yaxeflags1", golua.LNil)
	t.RawSetString("__yaxeflags2", golua.LNil)
	t.RawSetString("__yaxeflags3", golua.LNil)
	t.RawSetString("__yticks", state.NewTable())
	t.RawSetString("__plots", golua.LNil)

	tableBuilderFunc(state, t, "axis_limits", func(state *golua.LState, t *golua.LTable) {
		xmin := state.CheckNumber(-5)
		xmax := state.CheckNumber(-4)
		ymin := state.CheckNumber(-3)
		ymax := state.CheckNumber(-2)
		cond := state.CheckNumber(-1)
		t.RawSetString("__xmin", xmin)
		t.RawSetString("__xmax", xmax)
		t.RawSetString("__ymin", ymin)
		t.RawSetString("__ymax", ymax)
		t.RawSetString("__cond", cond)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "set_xaxis_label", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-2)
		label := state.CheckString(-1)
		lt := state.NewTable()
		lt.RawSetString("axis", axis)
		lt.RawSetString("label", golua.LString(label))

		ft := t.RawGetString("__xlabels").(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "set_yaxis_label", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-2)
		label := state.CheckString(-1)
		lt := state.NewTable()
		lt.RawSetString("axis", axis)
		lt.RawSetString("label", golua.LString(label))

		ft := t.RawGetString("__ylabels").(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "x_axeflags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__xaxeflags", flags)
	})

	tableBuilderFunc(state, t, "xticks", func(state *golua.LState, t *golua.LTable) {
		ticks := state.CheckTable(-2)
		dflt := state.CheckBool(-1)
		t.RawSetString("__xticks", ticks)
		t.RawSetString("__xaticksdefault", golua.LBool(dflt))
	})

	tableBuilderFunc(state, t, "y_axeflags", func(state *golua.LState, t *golua.LTable) {
		flags1 := state.CheckNumber(-3)
		flags2 := state.CheckNumber(-2)
		flags3 := state.CheckNumber(-1)
		t.RawSetString("__yaxeflags1", flags1)
		t.RawSetString("__yaxeflags2", flags2)
		t.RawSetString("__yaxeflags3", flags3)
	})

	tableBuilderFunc(state, t, "yticks", func(state *golua.LState, t *golua.LTable) {
		ticks := state.CheckTable(-3)
		dflt := state.CheckBool(-2)
		axis := state.CheckNumber(-1)
		lt := state.NewTable()
		lt.RawSetString("ticks", ticks)
		lt.RawSetString("dflt", golua.LBool(dflt))
		lt.RawSetString("axis", axis)

		ft := t.RawGetString("__ylabels").(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "plots", func(state *golua.LState, t *golua.LTable) {
		plots := state.CheckTable(-1)
		t.RawSetString("__plots", plots)
	})

	return t
}

func plotBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	/// @interface Plot
	/// @prop type {string<guiplot.PlotType>}

	title := t.RawGetString("title").(golua.LString)
	p := g.Plot(string(title))

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		p.Size(int(width.(golua.LNumber)), int(height.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		p.Flags(g.PlotFlags(flags.(golua.LNumber)))
	}

	xaxeflags := t.RawGetString("__xaxeflags")
	if xaxeflags.Type() == golua.LTNumber {
		p.XAxeFlags(g.PlotAxisFlags(xaxeflags.(golua.LNumber)))
	}

	yaxeflags1 := t.RawGetString("__yaxeflags1")
	yaxeflags2 := t.RawGetString("__yaxeflags2")
	yaxeflags3 := t.RawGetString("__yaxeflags3")
	if yaxeflags1.Type() == golua.LTNumber && yaxeflags2.Type() == golua.LTNumber && yaxeflags3.Type() == golua.LTNumber {
		p.YAxeFlags(g.PlotAxisFlags(yaxeflags1.(golua.LNumber)), g.PlotAxisFlags(yaxeflags2.(golua.LNumber)), g.PlotAxisFlags(yaxeflags3.(golua.LNumber)))
	}

	xmin := t.RawGetString("__xmin")
	xmax := t.RawGetString("__xmax")
	ymin := t.RawGetString("__ymin")
	ymax := t.RawGetString("__ymax")
	cond := t.RawGetString("__cond")
	if xmin.Type() == golua.LTNumber && xmax.Type() == golua.LTNumber && ymin.Type() == golua.LTNumber && ymax.Type() == golua.LTNumber && cond.Type() == golua.LTNumber {
		p.AxisLimits(
			float64(xmin.(golua.LNumber)), float64(xmax.(golua.LNumber)),
			float64(ymin.(golua.LNumber)), float64(ymax.(golua.LNumber)),
			g.ExecCondition(cond.(golua.LNumber)),
		)
	}

	xlabels := t.RawGetString("__xlabels").(*golua.LTable)
	for i := range xlabels.Len() {
		lt := xlabels.RawGetInt(i + 1).(*golua.LTable)
		axis := lt.RawGetString("axis").(golua.LNumber)
		label := lt.RawGetString("label").(golua.LString)

		p.SetXAxisLabel(g.PlotXAxis(axis), string(label))
	}

	ylabels := t.RawGetString("__ylabels").(*golua.LTable)
	for i := range ylabels.Len() {
		lt := ylabels.RawGetInt(i + 1).(*golua.LTable)
		axis := lt.RawGetString("axis").(golua.LNumber)
		label := lt.RawGetString("label").(golua.LString)

		p.SetYAxisLabel(g.PlotYAxis(axis), string(label))
	}

	xticks := t.RawGetString("__xticks")
	xticksdefault := t.RawGetString("__xticksdefault")
	if xticks.Type() == golua.LTTable && xticksdefault.Type() == golua.LTBool {
		ticks := []g.PlotTicker{}

		xtickst := xticks.(*golua.LTable)
		for i := range xtickst.Len() {
			tick := xtickst.RawGetInt(i + 1).(*golua.LTable)
			ticks = append(ticks, plotTickerBuild(tick))
		}

		p.XTicks(ticks, bool(xticksdefault.(golua.LBool)))
	}

	yticks := t.RawGetString("__yticks").(*golua.LTable)
	for z := range yticks.Len() {
		yticksaxis := yticks.RawGetInt(z + 1).(*golua.LTable)

		ytickaxis := yticksaxis.RawGetString("ticks").(*golua.LTable)
		dflt := yticksaxis.RawGetString("dflt").(golua.LBool)
		axis := yticksaxis.RawGetString("axis").(golua.LNumber)

		ticks := []g.PlotTicker{}

		for i := range ytickaxis.Len() {
			tick := ytickaxis.RawGetInt(i + 1).(*golua.LTable)
			ticks = append(ticks, plotTickerBuild(tick))
		}

		p.YTicks(ticks, bool(dflt), g.ImPlotYAxis(axis))
	}

	plots := t.RawGetString("__plots")
	if plots.Type() == golua.LTTable {
		plist := plotsBuild(plots.(*golua.LTable), r, lg, state)

		p.Plots(plist...)
	}

	return p
}

func plotsBuild(plots *golua.LTable, r *lua.Runner, lg *log.Logger, state *golua.LState) []g.PlotWidget {
	plist := []g.PlotWidget{}

	for i := range plots.Len() {
		pt := plots.RawGetInt(i + 1).(*golua.LTable)
		plottype := pt.RawGetString("type").(golua.LString)

		build := plotList[string(plottype)]
		plist = append(plist, build(r, lg, state, pt))
	}

	return plist
}

func plotTickerBuild(t *golua.LTable) g.PlotTicker {
	/// @struct PlotTicker
	/// @prop position {float}
	/// @prop label {string}

	position := t.RawGetString("position").(golua.LNumber)
	label := t.RawGetString("label").(golua.LString)

	return g.PlotTicker{
		Position: float64(position),
		Label:    string(label),
	}
}

func plotBarHTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotBarH
	/// @prop type {string<guiplot.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method height(self, height float) -> self
	/// @method offset(self, offset float) -> self
	/// @method shift(self, shift float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_BAR_H))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__shift", golua.LNil)

	tableBuilderFunc(state, t, "height", func(state *golua.LState, t *golua.LTable) {
		height := state.CheckNumber(-1)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "shift", func(state *golua.LState, t *golua.LTable) {
		shift := state.CheckNumber(-1)
		t.RawSetString("__shift", shift)
	})

	return t
}

func plotBarHBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.BarH(string(title), dataPoints)

	height := t.RawGetString("__height")
	if height.Type() == golua.LTNumber {
		p.Height(float64(height.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	shift := t.RawGetString("__shift")
	if shift.Type() == golua.LTNumber {
		p.Shift(float64(shift.(golua.LNumber)))
	}

	return p
}

func plotBarTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotBar
	/// @prop type {string<guiplot.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method width(self, width float) -> self
	/// @method offset(self, offset float) -> self
	/// @method shift(self, shift float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_BAR))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__shift", golua.LNil)

	tableBuilderFunc(state, t, "width", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "shift", func(state *golua.LState, t *golua.LTable) {
		shift := state.CheckNumber(-1)
		t.RawSetString("__shift", shift)
	})

	return t
}

func plotBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Bar(string(title), dataPoints)

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		p.Width(float64(width.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	shift := t.RawGetString("__shift")
	if shift.Type() == golua.LTNumber {
		p.Shift(float64(shift.(golua.LNumber)))
	}

	return p
}

func plotLineTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotLine
	/// @prop type {string<guiplot.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method set_plot_y_axis(self, axis int<guiplot.PlotYAxis>) -> self
	/// @method offset(self, offset float) -> self
	/// @method x0(self, x0 float) -> self
	/// @method xscale(self, xscale float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_LINE))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__yaxis", golua.LNil)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__x0", golua.LNil)
	t.RawSetString("__xscale", golua.LNil)

	tableBuilderFunc(state, t, "set_plot_y_axis", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-1)
		t.RawSetString("__yaxis", axis)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "x0", func(state *golua.LState, t *golua.LTable) {
		x0 := state.CheckNumber(-1)
		t.RawSetString("__x0", x0)
	})

	tableBuilderFunc(state, t, "xscale", func(state *golua.LState, t *golua.LTable) {
		xscale := state.CheckNumber(-1)
		t.RawSetString("__xscale", xscale)
	})

	return t
}

func plotLineBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Line(string(title), dataPoints)

	yaxis := t.RawGetString("__yaxis")
	if yaxis.Type() == golua.LTNumber {
		p.SetPlotYAxis(g.ImPlotYAxis(yaxis.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	x0 := t.RawGetString("__x0")
	if x0.Type() == golua.LTNumber {
		p.X0(float64(x0.(golua.LNumber)))
	}

	xscale := t.RawGetString("__xscale")
	if xscale.Type() == golua.LTNumber {
		p.XScale(float64(xscale.(golua.LNumber)))
	}

	return p
}

func plotLineXYTable(state *golua.LState, title string, xdata, ydata golua.LValue) *golua.LTable {
	/// @struct PlotLineXY
	/// @prop type {string<guiplot.PlotType>}
	/// @prop title {string}
	/// @prop xdata {[]float}
	/// @prop ydata {[]float}
	/// @method set_plot_y_axis(self, axis int<guiplot.PlotYAxis>) -> self
	/// @method offset(self, offset float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_LINE_XY))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("xdata", xdata)
	t.RawSetString("ydata", ydata)
	t.RawSetString("__yaxis", golua.LNil)
	t.RawSetString("__offset", golua.LNil)

	tableBuilderFunc(state, t, "set_plot_y_axis", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-1)
		t.RawSetString("__yaxis", axis)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	return t
}

func plotLineXYBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	xdata := t.RawGetString("xdata").(*golua.LTable)
	ydata := t.RawGetString("ydata").(*golua.LTable)

	xdataPoints := []float64{}
	for i := range xdata.Len() {
		point := xdata.RawGetInt(i + 1).(golua.LNumber)
		xdataPoints = append(xdataPoints, float64(point))
	}

	ydataPoints := []float64{}
	for i := range ydata.Len() {
		point := ydata.RawGetInt(i + 1).(golua.LNumber)
		ydataPoints = append(ydataPoints, float64(point))
	}

	p := g.LineXY(string(title), xdataPoints, ydataPoints)

	yaxis := t.RawGetString("__yaxis")
	if yaxis.Type() == golua.LTNumber {
		p.SetPlotYAxis(g.ImPlotYAxis(yaxis.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	return p
}

func plotPieTable(state *golua.LState, labels golua.LValue, data golua.LValue, x, y, radius float64) *golua.LTable {
	/// @struct PlotPieChart
	/// @prop type {string<guiplot.PlotType>}
	/// @prop labels {[]string}
	/// @prop data {[]float}
	/// @prop x {float}
	/// @prop y {float}
	/// @prop radius {float}
	/// @method angle0(self, angle0 float) -> self
	/// @method label_format(self, format string) -> self
	/// @method normalize(self, bool) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_PIE_CHART))
	t.RawSetString("labels", labels)
	t.RawSetString("data", data)
	t.RawSetString("x", golua.LNumber(x))
	t.RawSetString("y", golua.LNumber(y))
	t.RawSetString("radius", golua.LNumber(radius))
	t.RawSetString("__angle0", golua.LNil)
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__normalize", golua.LNil)

	tableBuilderFunc(state, t, "angle0", func(state *golua.LState, t *golua.LTable) {
		angle0 := state.CheckNumber(-1)
		t.RawSetString("__angle0", angle0)
	})

	tableBuilderFunc(state, t, "label_format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "normalize", func(state *golua.LState, t *golua.LTable) {
		normalize := state.CheckBool(-1)
		t.RawSetString("__normalize", golua.LBool(normalize))
	})

	return t
}

func plotPieBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	labels := t.RawGetString("labels").(*golua.LTable)
	data := t.RawGetString("data").(*golua.LTable)
	x := t.RawGetString("x").(golua.LNumber)
	y := t.RawGetString("y").(golua.LNumber)
	radius := t.RawGetString("radius").(golua.LNumber)

	labelPoints := []string{}
	for i := range labels.Len() {
		point := labels.RawGetInt(i + 1).(golua.LString)
		labelPoints = append(labelPoints, string(point))
	}

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.PieChart(labelPoints, dataPoints, float64(x), float64(y), float64(radius))

	angle0 := t.RawGetString("__angle0")
	if angle0.Type() == golua.LTNumber {
		p.Angle0(float64(angle0.(golua.LNumber)))
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		p.LabelFormat(string(format.(golua.LString)))
	}

	normalize := t.RawGetString("__normalize")
	if normalize.Type() == golua.LTBool {
		p.Normalize(bool(normalize.(golua.LBool)))
	}

	return p
}

func plotScatterTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotScatter
	/// @prop type {string<guiplot.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method offset(self, offset float) -> self
	/// @method x0(self, x0 float) -> self
	/// @method xscale(self, xscale float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_SCATTER))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__x0", golua.LNil)
	t.RawSetString("__xscale", golua.LNil)

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "x0", func(state *golua.LState, t *golua.LTable) {
		x0 := state.CheckNumber(-1)
		t.RawSetString("__x0", x0)
	})

	tableBuilderFunc(state, t, "xscale", func(state *golua.LState, t *golua.LTable) {
		xscale := state.CheckNumber(-1)
		t.RawSetString("__xscale", xscale)
	})

	return t
}

func plotScatterBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Scatter(string(title), dataPoints)

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	x0 := t.RawGetString("__x0")
	if x0.Type() == golua.LTNumber {
		p.X0(float64(x0.(golua.LNumber)))
	}

	xscale := t.RawGetString("__xscale")
	if xscale.Type() == golua.LTNumber {
		p.XScale(float64(xscale.(golua.LNumber)))
	}

	return p
}

func plotScatterXYTable(state *golua.LState, title string, xdata, ydata golua.LValue) *golua.LTable {
	/// @struct PlotScatterXY
	/// @prop type {string<guiplot.PlotType>}
	/// @prop title {string}
	/// @prop xdata {[]float}
	/// @prop ydata {[]float}
	/// @method offset(self, offset float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_SCATTER_XY))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("xdata", xdata)
	t.RawSetString("ydata", ydata)
	t.RawSetString("__offset", golua.LNil)

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	return t
}

func plotScatterXYBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	xdata := t.RawGetString("xdata").(*golua.LTable)
	ydata := t.RawGetString("ydata").(*golua.LTable)

	xdataPoints := []float64{}
	for i := range xdata.Len() {
		point := xdata.RawGetInt(i + 1).(golua.LNumber)
		xdataPoints = append(xdataPoints, float64(point))
	}

	ydataPoints := []float64{}
	for i := range ydata.Len() {
		point := ydata.RawGetInt(i + 1).(golua.LNumber)
		ydataPoints = append(ydataPoints, float64(point))
	}

	p := g.ScatterXY(string(title), xdataPoints, ydataPoints)

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	return p
}

func plotCustomTable(state *golua.LState, builder *golua.LFunction) *golua.LTable {
	/// @struct PlotCustom
	/// @prop type {string<guiplot.PlotType>}
	/// @prop builder {function()}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_CUSTOM))
	t.RawSetString("builder", builder)

	return t
}

func plotCustomBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	builder := t.RawGetString("builder").(*golua.LFunction)

	c := g.Custom(func() {
		state.Push(builder)
		state.Call(0, 0)
	})

	return c
}

func plotStyleBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	s := styleBuild(r, lg, state, t)
	return s.(*g.StyleSetter)
}
