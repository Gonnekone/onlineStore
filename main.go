package main

import (
	"net/http"

	"github.com/Gonnekone/onlineStore/handlers"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	http.Handle("/", initRoutes())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Failed to start the server")
	}
}

func initRoutes() http.Handler {
    mux := http.NewServeMux()

    api := "/api"
    user := "/user"
    typeRoute := "/type"
    brand := "/brand"
    device := "/device"

    mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, api, http.StatusFound)
    }))

    mux.Handle(api+user+"/registration", http.HandlerFunc(handlers.Registration))
    mux.Handle(api+user+"/login", http.HandlerFunc(handlers.Login))
    mux.Handle(api+user+"/auth", middleware.Auth(http.HandlerFunc(handlers.Validate)))

    mux.Handle(api+typeRoute+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handlers.ReadAllTypes(w, r)
        case http.MethodPost:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.CreateType))).ServeHTTP(w, r)
        case http.MethodPut:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.UpdateType))).ServeHTTP(w, r)
        case http.MethodDelete:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.DeleteType))).ServeHTTP(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))

    mux.Handle(api+brand+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handlers.ReadAllBrands(w, r)
        case http.MethodPost:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.CreateBrand))).ServeHTTP(w, r)
        case http.MethodPut:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.UpdateBrand))).ServeHTTP(w, r)
        case http.MethodDelete:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.DeleteBrand))).ServeHTTP(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))

    mux.Handle(api+device+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handlers.ReadAllDevices(w, r)
        case http.MethodPost:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.CreateDevice))).ServeHTTP(w, r)
        case http.MethodPut:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.UpdateDevice))).ServeHTTP(w, r)
        case http.MethodDelete:
            middleware.Auth(middleware.CheckRoleAdmin(http.HandlerFunc(handlers.DeleteDevice))).ServeHTTP(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }))

    mux.Handle(api+device+"/:id", http.HandlerFunc(handlers.ReadOneDevice))

    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    return mux
}

