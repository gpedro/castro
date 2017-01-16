package controllers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/database"
	"os"
	"github.com/raggaer/castro/app/util"
	"time"
)

func Signature(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get name
	name := ps.ByName("name")

	// Model to get player info
	player := models.Player{}

	// Get player information
	if err := database.DB.Find(&player, "name = ?", name).Error; err != nil {
		return
	}

	// Check if signature image needs to be updated
	info, err := os.Stat("public/img/signature/" + name + ".png")

	if err != nil {

		// Create signature image
		if err := util.CreatePlayerSignature(player); err != nil {
			return
		}

		// Serve signature file
		http.ServeFile(w, r, "public/img/signature/" + name + ".png")

		return
	}

	// Get time plus 5 minutes
	diff := info.ModTime().Add(time.Minute * 5)

	// Check if signature image is old
	if diff.Before(time.Now()) || util.Config.IsDev() {

		// Create signature image
		if err := util.CreatePlayerSignature(player); err != nil {
			return
		}
	}

	// Serve signature file
	http.ServeFile(w, r, "public/img/signature/" + name + ".png")
}