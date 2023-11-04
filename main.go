package main

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/gen2brain/beeep"
	"github.com/sqweek/dialog"
)

func writeConfig(conf Config) {
	file, _ := json.MarshalIndent(conf, "", " ")
	_ = os.WriteFile("config.json", file, 0644)
}

func readConfig() Config {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		file, _ := json.MarshalIndent(Config{LastBrowseFolder: ""}, "", " ")
		_ = os.WriteFile("config.json", file, 0644)
		return Config{LastBrowseFolder: ""}
	} else {
		data := Config{}
		json.Unmarshal(file, &data)
		return data
	}
}

func parseImageToGrid(g *fyne.Container, path string) {
	images := GetImages(path)

	g.Objects = []fyne.CanvasObject{}

	for _, img := range images {
		imgFix := img

		im := canvas.NewImageFromFile(imgFix.Path)
		im.FillMode = canvas.ImageFillOriginal

		button := widget.NewButton("", func() {
			err := beeep.Notify("Quick Image Browser", "Image copied to clipboard!", "./Icon.png")
			if err != nil {
				panic(err)
			}
			CopyImageToClipboard(imgFix.Path)
		})

		container := container.New(
			layout.NewStackLayout(),
			button,
			im,
		)

		g.Objects = append(g.Objects, container)
	}

	g.Refresh()

	curConfig := readConfig()
	curConfig.LastBrowseFolder = path
	writeConfig(curConfig)
}

func main() {
	// currentDirectory := ""
	currentConfig := readConfig()

	a := app.New()
	w := a.NewWindow("Image Browser")
	w.SetIcon(theme.FyneLogo())
	w.Resize(fyne.NewSize(500, 500))

	g := container.NewGridWrap(fyne.NewSize(100, 100))
	s := container.NewVScroll(g)
	s.SetMinSize(fyne.NewSize(100, 500))

	input := widget.NewEntry()
	input.SetText(currentConfig.LastBrowseFolder)
	// input.Resize(fyne.NewSize(400, 50))

	input.SetPlaceHolder("Enter path...")

	goButton := widget.NewButton("Go", func() {
		parseImageToGrid(g, input.Text)
	})

	browseButton := widget.NewButton("Choose Folder", func() {
		directory, err := dialog.Directory().Title("Load images").Browse()
		if err != nil {
			panic(err)
		}
		input.SetText(directory)
		parseImageToGrid(g, directory)
	})

	inputContainer := container.NewBorder(input, nil, nil, container.NewHBox(goButton, browseButton))

	// list := widget.NewTable(
	// 	func() (int, int) {
	// 		return len(data), len(data[0])
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return widget.NewLabel("wide content")
	// 	},
	// 	func(i widget.TableCellID, o fyne.CanvasObject) {
	// 		o.(*widget.Label).SetText(data[i.Row][i.Col])
	// 	})

	w.SetContent(container.NewVBox(

		inputContainer,
		s,
	))

	go func() {
		if currentConfig.LastBrowseFolder != "" {
			parseImageToGrid(g, currentConfig.LastBrowseFolder)
		}
	}()

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(theme.FileIcon())

		m := fyne.NewMenu("Quick Image Browser",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))

		desk.SetSystemTrayMenu(m)
	}

	w.SetCloseIntercept(func() {
		w.Hide()
	})

	w.ShowAndRun()
}
