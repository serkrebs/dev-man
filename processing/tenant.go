package processing

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"bytes"

	"github.com/gorilla/mux"
	"se.com/data-processing/devdevman/cors"
)

const tenantsPath = "tenant"
const devicesPath = "device"

// Tenant holds the GUID for the device owner for the environment and tenant
type Tenant struct {
	TenantName    string `json:"tenant,omitempty"`
	Environment   string `json:"environment"`
	DeviceOwnerID string `json:"deviceOwner"`
}

func getAllTenants() ([]Tenant, error) {
	log.Println("loading tenants...")
	var tenants []Tenant
	resp, err := http.Get("https://processing.dev.struxurewarecloud.com/api/processing/environment/8fa5c77f-5c2c-4a81-929b-92efe8f876f0/tenant")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tenants)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(tenants, func(i, j int) bool {
		if tenants[i].Environment == tenants[j].Environment {
			return tenants[i].TenantName < tenants[j].TenantName
		}
		return tenants[i].Environment < tenants[j].Environment
	})
	return tenants, nil
}

func handleTenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tenants, err := getAllTenants()
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(tenants)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Printf("Response error: %s", err.Error())
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getOneTenant(tenantId string) (*Tenant, error) {
	log.Println("loading tenant ", tenantId, "...")
	var tenant Tenant
	resp, err := http.Get(fmt.Sprintf("https://processing.dev.struxurewarecloud.com/api/processing/environment/8fa5c77f-5c2c-4a81-929b-92efe8f876f0/tenant/%s", tenantId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func handleOneTenant(w http.ResponseWriter, r *http.Request) {
	vmap := mux.Vars(r)
	tenant, ok := vmap["tenant"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("tenant not recognized"))
	}
	log.Println("loading 1 tenant ", tenant, "...")
	switch r.Method {
	case http.MethodGet:
		tenantReg, err := getOneTenant(tenant)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(tenantReg)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Printf("Response error: %s", err.Error())
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


// DeviceMetaSummary holds common device metadata fields for display in a list
type DeviceMetaSummary struct {
	Environment  string `json:"environment"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"modelNumber"`
	Serial       string `json:"serialNumber"`
}

// DeviceSummary holds common device registration fields for display in a list
type DeviceSummary struct {
	DeviceID    string            `json:"deviceId"`
	DeviceType  string            `json:"deviceType"`
	Tenant      string            `json:"tenant"`
	DeviceOwner string            `json:"deviceOwner"`
	Meta        DeviceMetaSummary `json:"metadata"`
}

func getAllTenantDevices(tenant string) ([]DeviceSummary, error) {
	log.Println("loading tenants...")
	var devices []DeviceSummary
	resp, err := http.Get(fmt.Sprintf("https://processing.dev.struxurewarecloud.com/api/processing/metadata/%s/device", tenant))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].DeviceID < devices[j].DeviceID
	})
	return devices, nil
}

// DeviceSummary holds common device registration fields for display in a list
type DeviceRegistration struct {
	DeviceID    string            `json:"deviceId"`
	DeviceType  string            `json:"deviceType"`
	Tenant      string            `json:"tenant"`
	DeviceOwner string            `json:"deviceOwner"`
	Meta        map[string]string `json:"metadata"`
}

func handleDevicesForTenant(w http.ResponseWriter, r *http.Request) {
	vmap := mux.Vars(r)
	tenant, ok := vmap["tenant"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("tenant not recognized: " + tenant))
	}
	log.Println("loading devices ", tenant, "...")
	switch r.Method {
	case http.MethodGet:
		devices, err := getAllTenantDevices(tenant)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(devices)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Printf("Response error: %s", err.Error())
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getOneTenantDevice(tenant, deviceId string) (*DeviceRegistration, error) {
	log.Printf("loading device '%s'...", deviceId)
	var deviceReg DeviceRegistration
	resp, err := http.Get(fmt.Sprintf("https://processing.dev.struxurewarecloud.com/api/processing/metadata/%s/device/%s", tenant, deviceId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &deviceReg)
	if err != nil {
		return nil, err
	}

	return &deviceReg, nil
}

func updateOneTenantDevice(tenant, deviceId string, device *DeviceRegistration) (*DeviceRegistration, error) {
	log.Printf("update device '%s'...", deviceId)

	jBytes, err := json.Marshal(device)
    if err != nil {
        return nil, err
    }

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("https://processing.dev.struxurewarecloud.com/api/processing/environment/8fa5c77f-5c2c-4a81-929b-92efe8f876f0/tenant/%s/device/%s", tenant, deviceId),
		bytes.NewBuffer(jBytes),
	)
	if err != nil {
		return nil, err
	}

	
    req.Header.Set("Content-Type", "application/json; charset=utf-8")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        panic(err)
    }

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var deviceReg DeviceRegistration
	err = json.Unmarshal(body, &deviceReg)
	if err != nil {
		return nil, err
	}

	return &deviceReg, nil
}

func handleOneDeviceForTenant(w http.ResponseWriter, r *http.Request) {
	vmap := mux.Vars(r)
	tenant, ok := vmap["tenant"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("tenant not recognized"))
	}
	deviceId, ok := vmap["device"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("device not recognized"))
	}
	switch r.Method {
	case http.MethodGet:
		log.Println("loading device ", tenant, " - ", deviceId, "...")
		deviceReg, err := getOneTenantDevice(tenant, deviceId)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(deviceReg)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Printf("Response error: %s", err.Error())
		}

	case http.MethodPut:
		log.Println("updating device ", tenant, " - ", deviceId, "...")
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		var devUpdate DeviceRegistration
		err = json.Unmarshal(body, &devUpdate)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		deviceReg, err := updateOneTenantDevice(tenant, deviceId, &devUpdate)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		j, err := json.Marshal(deviceReg)
		if err != nil {
			log.Printf("Server error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Printf("Response error: %s", err.Error())
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SetupRoutes :
func SetupRoutes(apiBasePath string, r *mux.Router) {
	spa := spaHandler{staticPath: "./public", indexPath: "index.html"}
	tenantsHandler := http.HandlerFunc(handleTenants)
	oneTenantHandler := http.HandlerFunc(handleOneTenant)
	devicesHandler := http.HandlerFunc(handleDevicesForTenant)
	oneDeviceHandler := http.HandlerFunc(handleOneDeviceForTenant)
	r.Handle(fmt.Sprintf("%s/%s", apiBasePath, tenantsPath), cors.Middleware(tenantsHandler))
	r.Handle(fmt.Sprintf("%s/%s/{tenant}", apiBasePath, tenantsPath), cors.Middleware(oneTenantHandler))
	r.Handle(fmt.Sprintf("%s/%s/{tenant}/%s", apiBasePath, tenantsPath, devicesPath), cors.Middleware(devicesHandler))
	r.Handle(fmt.Sprintf("%s/%s/{tenant}/%s/{device}", apiBasePath, tenantsPath, devicesPath), cors.Middleware(oneDeviceHandler))
	r.PathPrefix("/").Handler(spa)
}
