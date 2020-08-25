package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/unix-streamdeck/api"
)

type button struct {
	widget.BaseWidget
	size int

	keyID int
	key   api.Key
}

func newButton(key api.Key, id int, size int) *button {
	b := &button{key: key, keyID: id, size: size}
	b.ExtendBaseWidget(b)
	return b
}

func (b *button) CreateRenderer() fyne.WidgetRenderer {
	icon := canvas.NewImageFromFile(b.key.Icon)
	text := canvas.NewText(b.key.Text, color.White)
	text.Alignment = fyne.TextAlignCenter

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 2
	border.SetMinSize(fyne.NewSize(b.size, b.size))

	bg := canvas.NewRectangle(color.Black)
	render := &buttonRenderer{border: border, text: text, icon: icon, bg: bg,
		objects: []fyne.CanvasObject{bg, icon, text, border}, b: b}
	render.Refresh()
	return render
}

func (b *button) Tapped(ev *fyne.PointEvent) {
	editButton(b)
}

func (b *button) updateKey() {
	if len(config.Pages) == 0 {
		config.Pages = append(config.Pages, api.Page{api.Key{}})
	}
	config.Pages[0][b.keyID] = b.key
	err := conn.SetConfig(config)
	if err != nil {
		dialog.ShowError(err, win)
	}
}

const (
	buttonInset = 5
)

type buttonRenderer struct {
	border, bg *canvas.Rectangle
	text       *canvas.Text
	icon       *canvas.Image

	objects []fyne.CanvasObject

	b *button
}

func (r *buttonRenderer) Layout(s fyne.Size) {
	size := s.Subtract(fyne.NewSize(buttonInset*2, buttonInset*2))
	offset := fyne.NewPos(buttonInset, buttonInset)

	for _, obj := range r.objects {
		obj.Move(offset)
		obj.Resize(size)
	}
}

func (r *buttonRenderer) MinSize() fyne.Size {
	iconSize := fyne.NewSize(r.b.size, r.b.size)
	return iconSize.Add(fyne.NewSize(buttonInset*2, buttonInset*2))
}

func (r *buttonRenderer) Refresh() {
	if currentButton == r.b {
		r.border.StrokeColor = theme.FocusColor()
	} else {
		r.border.StrokeColor = &color.Gray{128}
	}

	if r.b.key.Text != r.text.Text {
		r.text.Text = r.b.key.Text
		r.text.Refresh()
	}
	if r.b.key.Icon != r.icon.File {
		r.icon.File = r.b.key.Icon
		r.icon.Refresh()
	}

	r.border.Refresh()
}

func (r *buttonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *buttonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *buttonRenderer) Destroy() {
	// nothing
}