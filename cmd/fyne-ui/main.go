package main

import (
	"fmt"
	"gomnirun/core/config"
	"gomnirun/core/executor"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const configFile = "config.json"

// Function to strip ANSI escape codes
func stripAnsiCodes(input string) string {
	ansiRegex := regexp.MustCompile(`\x1B\[[0-9;]*[mK]`)
	return ansiRegex.ReplaceAllString(input, "")
}

func main() {
	a := app.New()
	w := a.NewWindow("GomniRun - Script Runner")
	w.Resize(fyne.NewSize(800, 500)) // Wider window for split layout

	// Load configuration
	conf, _ := config.Load(configFile)

	// UI Elements
	commandEntry := widget.NewEntry()
	commandEntry.SetText(conf.CommandTemplate)
	commandEntry.MultiLine = true

	// Variable Input Fields (Matching Type)
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

	// Output Area with High Contrast
	outputLabel := widget.NewLabel("Output:")
	outputText := widget.NewRichTextFromMarkdown("**Script output will appear here...**")
	outputText.Wrapping = fyne.TextWrapWord

	// Dynamically set background and text colors based on theme
	bgColor := theme.BackgroundColor()

	// Background Rectangle
	outputBg := canvas.NewRectangle(bgColor)
	outputBg.FillColor = bgColor

	// Use a Stack to layer the background and text
	outputContainer := container.NewStack(outputBg, container.NewVScroll(outputText))

	// Run Script Button
	runButton := widget.NewButton("Run Script", func() {
		conf.CommandTemplate = commandEntry.Text
		for key, obj := range varEntries {
			switch widget := obj.(type) {
			case *widget.Entry:
				conf.Variables[key] = config.Variable{Type: conf.Variables[key].Type, Value: widget.Text}
			case *widget.Check:
				conf.Variables[key] = config.Variable{Type: "bool", Value: fmt.Sprintf("%v", widget.Checked)}
			}
		}

		// Generate final command with variables replaced
		finalCommand := executor.ReplacePlaceholders(conf.CommandTemplate, conf.Variables)
		fmt.Println("Executing command:", finalCommand)

		// Show final command in dialog before execution
		dialog.ShowConfirm("Execute Command", "Executing: "+finalCommand, func(confirmed bool) {
			if confirmed {
				output, err := executor.RunScript(conf.CommandTemplate, conf.Variables)
				if err != nil {
					dialog.ShowError(err, w)
					outputText.ParseMarkdown(fmt.Sprintf("**Error:** %s", err.Error()))
				} else {
					cleanOutput := stripAnsiCodes(output) // Strip ANSI color codes
					outputText.ParseMarkdown(fmt.Sprintf("**Executed:** `%s`\n\n```\n%s\n```", finalCommand, cleanOutput))
				}
			}
		}, w)

		// Save config only if overwrite is enabled
		if conf.Overwrite {
			config.Save(configFile, conf)
			fmt.Println("Configuration updated and saved.")
		} else {
			fmt.Println("Overwrite is disabled. Changes will not be saved.")
		}
	})

	// Left Panel (Command + Parameters)
	leftPanel := container.NewVBox(
		widget.NewLabel("Command Template:"),
		commandEntry,
		widget.NewLabel("Variables:"),
		varBox,
		layout.NewSpacer(), // Push Run button to the bottom
		runButton,
	)
	leftPanelWithPadding := container.NewPadded(leftPanel) // Add padding

	// Right Panel (Output, spans remaining height)
	rightPanel := container.NewBorder(outputLabel, nil, nil, nil, outputContainer)
	rightPanelWithPadding := container.NewPadded(rightPanel) // Add padding

	// Split View (Left: Inputs, Right: Output)
	splitView := container.NewHSplit(leftPanelWithPadding, rightPanelWithPadding)
	splitView.SetOffset(0.4) // Adjusts split size (40% inputs, 60% output)

	// Set Content
	w.SetContent(splitView)
	w.ShowAndRun()
}
