package main

import (
	"fmt"
	"log"
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
var githubToken = ""

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
	log.Printf("app starting, version=%s", normalizeVersion(version))

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
		token := strings.TrimSpace(githubToken)
		log.Printf("check updates clicked, repo=%q", repoSlug)

		if repoSlug == "" {
			log.Print("self update repo is empty")
			updateStatus.SetText("SELFUPDATE_REPO is not embedded (check build ldflags)")
			return
		}

		currentVersion, err := semver.ParseTolerant(version)
		if err != nil {
			log.Printf("invalid current version: %v", err)
			updateStatus.SetText("Invalid app version: " + err.Error())
			return
		}
		log.Printf("current version parsed: %s", normalizeVersion(currentVersion.String()))

		updater, err := selfupdate.NewUpdater(selfupdate.Config{
			APIToken: token,
		})
		if err != nil {
			log.Printf("updater init failed: %v", err)
			updateStatus.SetText("Updater initialization failed: " + err.Error())
			return
		}

		latest, found, err := updater.DetectLatest(repoSlug)
		if err != nil {
			log.Printf("update check failed: %v", err)
			updateStatus.SetText("Update check failed: " + err.Error())
			return
		}
		if !found {
			log.Print("no suitable release found for current platform")
			updateStatus.SetText("No suitable release found for this platform")
			return
		}
		log.Printf("latest version found: %s", normalizeVersion(latest.Version.String()))
		if latest.Version.LTE(currentVersion) {
			log.Print("already up to date")
			updateStatus.SetText("Already up to date: " + normalizeVersion(currentVersion.String()))
			return
		}

		updateStatus.SetText(fmt.Sprintf("Updating to %s...", normalizeVersion(latest.Version.String())))
		log.Printf("starting self update to: %s", normalizeVersion(latest.Version.String()))
		updatedRelease, err := updater.UpdateSelf(currentVersion, repoSlug)
		if err != nil {
			log.Printf("self update failed: %v", err)
			updateStatus.SetText("Update failed: " + err.Error())
			return
		}

		newVersion := normalizeVersion(latest.Version.String())
		if updatedRelease != nil {
			newVersion = normalizeVersion(updatedRelease.Version.String())
		}
		log.Printf("self update completed, new version=%s", newVersion)
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
