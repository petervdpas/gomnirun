package main

import (
	"flag"
	"fmt"
	"gomnirun/core/config"
	"gomnirun/core/executor"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func stripAnsiCodes(input string) string {
	ansiRegex := regexp.MustCompile(`\x1B\[[0-9;]*[mK]`)
	return ansiRegex.ReplaceAllString(input, "")
}

func loadConfig(configFile string) config.Config {
	conf, err := config.Load(configFile)
	if err != nil {
		fmt.Println("Failed to load config, using defaults.")
		conf = config.Config{}
	}
	return conf
}

func createCommandEntry(conf config.Config) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(conf.CommandTemplate)
	entry.MultiLine = true
	return entry
}

func createVariableInputs(w fyne.Window, conf config.Config) (map[string]fyne.CanvasObject, *fyne.Container) {
	varEntries := make(map[string]fyne.CanvasObject)
	varBox := container.NewVBox()
	for key, variable := range conf.Variables {
		label := widget.NewLabel(key)
		var inputWidget fyne.CanvasObject

		switch variable.Type {
		case "string":
			entry := widget.NewEntry()
			entry.SetText(variable.Value)
			varEntries[key] = entry
			inputWidget = container.NewGridWithColumns(1, entry)
		case "file":
			entry := widget.NewEntry()
			entry.SetText(variable.Value)
			filePicker := widget.NewButton("Select File", func() {
				dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
					if reader != nil {
						entry.SetText(reader.URI().Path())
					}
				}, w)
			})
			varEntries[key] = entry
			inputWidget = container.NewBorder(nil, nil, nil, filePicker, entry)
		case "directory":
			entry := widget.NewEntry()
			entry.SetText(variable.Value)
			dirPicker := widget.NewButton("Select Directory", func() {
				dialog.ShowFolderOpen(func(reader fyne.ListableURI, err error) {
					if reader != nil {
						entry.SetText(reader.Path())
					}
				}, w)
			})
			varEntries[key] = entry
			inputWidget = container.NewBorder(nil, nil, nil, dirPicker, entry)
		case "bool":
			check := widget.NewCheck("", nil)
			check.SetChecked(variable.Value == "true")
			varEntries[key] = check
			inputWidget = check
		}
		varBox.Add(container.NewVBox(label, inputWidget))
	}
	return varEntries, varBox
}

func createOutputArea() (*widget.Label, *widget.RichText, *fyne.Container) {
	outputLabel := widget.NewLabel("Output:")
	outputText := widget.NewRichTextFromMarkdown("**Script output will appear here...**")
	outputText.Wrapping = fyne.TextWrapWord

	bgColor := theme.BackgroundColor()
	outputBg := canvas.NewRectangle(bgColor)
	outputBg.FillColor = bgColor

	outputContainer := container.NewStack(outputBg, container.NewVScroll(outputText))
	return outputLabel, outputText, outputContainer
}

func createRunButton(w fyne.Window, conf config.Config, commandEntry *widget.Entry, varEntries map[string]fyne.CanvasObject, outputText *widget.RichText) *widget.Button {
	return widget.NewButton("Run Script", func() {
		conf.CommandTemplate = commandEntry.Text
		for key, obj := range varEntries {
			switch widget := obj.(type) {
			case *widget.Entry:
				conf.Variables[key] = config.Variable{Type: conf.Variables[key].Type, Value: widget.Text}
			case *widget.Check:
				conf.Variables[key] = config.Variable{Type: "bool", Value: fmt.Sprintf("%v", widget.Checked)}
			}
		}

		// Get the final command (as []string) and join it into a single string
		finalCommandParts := executor.ReplacePlaceholders(conf.CommandTemplate, conf.Variables)
		finalCommand := strings.Join(finalCommandParts, " ") // Convert []string to a single string

		fmt.Println("Executing command:", finalCommand)

		dialog.ShowConfirm("Execute Command", "Executing: "+finalCommand, func(confirmed bool) {
			if confirmed {
				output, err := executor.RunScript(conf.CommandTemplate, conf.Variables)
				if err != nil {
					dialog.ShowError(err, w)
					outputText.ParseMarkdown(fmt.Sprintf("**Error:** %s", err.Error()))
				} else {
					cleanOutput := stripAnsiCodes(output)
					outputText.ParseMarkdown(fmt.Sprintf("**Executed:** `%s`\n\n```\n%s\n```", finalCommand, cleanOutput))
				}
			}
		}, w)

		if conf.Overwrite {
			config.Save("config.json", conf)
			fmt.Println("Configuration updated and saved.")
		} else {
			fmt.Println("Overwrite is disabled. Changes will not be saved.")
		}
	})
}

func main() {
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	fmt.Printf("Starting GomniRun GUI Mode with config: %s...\n", *configFile)

	a := app.New()
	w := a.NewWindow("GomniRun - Script Runner")
	w.Resize(fyne.NewSize(800, 500))

	conf := loadConfig(*configFile)
	commandEntry := createCommandEntry(conf)
	varEntries, varBox := createVariableInputs(w, conf)
	outputLabel, outputText, outputContainer := createOutputArea()
	runButton := createRunButton(w, conf, commandEntry, varEntries, outputText)

	leftPanel := container.NewVBox(
		widget.NewLabel("Command Template:"),
		commandEntry,
		widget.NewLabel("Variables:"),
		varBox,
		layout.NewSpacer(),
		runButton,
	)
	leftPanelWithPadding := container.NewPadded(leftPanel)

	rightPanel := container.NewBorder(outputLabel, nil, nil, nil, outputContainer)
	rightPanelWithPadding := container.NewPadded(rightPanel)

	splitView := container.NewHSplit(leftPanelWithPadding, rightPanelWithPadding)
	splitView.SetOffset(0.4)

	w.SetContent(splitView)
	w.ShowAndRun()
}
