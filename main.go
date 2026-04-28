package main

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var version = "0.1.0"
var selfUpdateRepo = ""

func normalizeVersion(v string) string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "v") {
		return trimmed
	}
	return "v" + trimmed
}

func main() {
	a := app.New()
	w := a.NewWindow("Simple Fyne App")
	w.Resize(fyne.NewSize(520, 300))

	label := widget.NewLabel("Hello from Go + Fyne!")
	button := widget.NewButton("Click me", func() {
		label.SetText("Button clicked")
	})

	versionLabel := widget.NewLabel("Current version: " + normalizeVersion(version))
	updateStatus := widget.NewLabel("Self-update: ready")

	updateButton := widget.NewButton("Check updates", func() {
		repoSlug := strings.TrimSpace(selfUpdateRepo)
		if repoSlug == "" {
			updateStatus.SetText("SELFUPDATE_REPO is not embedded (check build ldflags)")
			return
		}

		currentVersion, err := semver.ParseTolerant(version)
		if err != nil {
			updateStatus.SetText("Invalid app version: " + err.Error())
			return
		}

		latest, found, err := selfupdate.DetectLatest(repoSlug)
		if err != nil {
			updateStatus.SetText("Update check failed: " + err.Error())
			return
		}
		if !found {
			updateStatus.SetText("No suitable release found for this platform")
			return
		}
		if latest.Version.LTE(currentVersion) {
			updateStatus.SetText("Already up to date: " + normalizeVersion(currentVersion.String()))
			return
		}

		updateStatus.SetText(fmt.Sprintf("Updating to %s...", normalizeVersion(latest.Version.String())))
		updatedRelease, err := selfupdate.UpdateSelf(currentVersion, repoSlug)
		if err != nil {
			updateStatus.SetText("Update failed: " + err.Error())
			return
		}

		newVersion := normalizeVersion(latest.Version.String())
		if updatedRelease != nil {
			newVersion = normalizeVersion(updatedRelease.Version.String())
		}
		updateStatus.SetText("Updated successfully to " + newVersion)
		dialog.ShowInformation(
			"Update completed",
			"Application was updated to "+newVersion+". Please restart the app.",
			w,
		)
	})

	content := container.NewVBox(
		widget.NewLabel("Minimal desktop app"),
		versionLabel,
		label,
		button,
		widget.NewSeparator(),
		widget.NewLabel("GitHub self-update"),
		widget.NewLabel("SELFUPDATE_REPO can be embedded at build stage"),
		updateStatus,
		updateButton,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
