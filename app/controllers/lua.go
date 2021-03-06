package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"net/http"
	"path/filepath"
	"strings"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Create application paypal REST client
	lua.CreatePaypalClient(util.Config.Configuration.PayPal.SandBox)

	// Check if request is POST
	if r.Method == http.MethodPost {

		// Parse POST form
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			return
		}
	}

	// If development mode reload pages, widgets and config file
	if util.Config.Configuration.IsDev() {

		// Reload config file
		if err := util.LoadConfig("config.toml"); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload config file: %v", err)

			return
		}

		// Reload pages
		if err := lua.PageList.Load("pages"); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload subtopic %v: %v", ps.ByName("page"), err)

			return
		}

		// Reload widgets
		if err := lua.WidgetList.Load("widgets"); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload widgets when executing %v subtopic: %v", ps.ByName("page"), err)

			return
		}
	}

	// Get session
	session, ok := r.Context().Value("session").(map[string]interface{})

	if !ok {
		// Set error header
		w.WriteHeader(500)
		util.Logger.Logger.Error("Cannot get session as map")

		return
	}

	// Set LUA file name
	pageName := ps.ByName("filepath")

	// If there is no subtopic request index
	if pageName == "" {
		pageName = "index"
	}

	// Get state from the pool
	s, err := lua.PageList.Get(filepath.Join("pages", pageName, r.Method+".lua"))

	if err != nil {

		// Set not found header
		w.WriteHeader(404)
		util.Logger.Logger.Errorf("Cannot get %v subtopic source: %v", pageName, err)

		return
	}

	// Return state to the pool
	defer lua.PageList.Put(s, filepath.Join("pages", pageName, r.Method+".lua"))

	// Create HTTP metatable
	lua.SetHTTPMetaTable(s)

	// Set the state user data
	lua.SetHTTPUserData(s, w, r)

	// Set session user data
	lua.SetSessionMetaTableUserData(s, session)

	// Call file function
	if err := s.CallByParam(
		glua.P{
			Fn:      s.GetGlobal(strings.ToLower(r.Method)),
			NRet:    0,
			Protect: !util.Config.Configuration.IsDev(),
		},
	); err != nil {

		// Set error header
		w.WriteHeader(500)
		util.Logger.Logger.Errorf("Cannot execute %v subtopic: %v", pageName, err)
	}
}
