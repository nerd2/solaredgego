package solaredgego

import "time"

// https://api.solaredge.com/services/m/so/sites/?status=ACTIVE%2CPENDING&sortName=name&sortOrder=ASC&page=0&pageSize=1&isDemoSite=false
type SitesResponse struct {
	Count int `json:"count"`
	Sites []struct {
		Address          string      `json:"address"`
		Address2         interface{} `json:"address2"`
		City             string      `json:"city"`
		Country          interface{} `json:"country"`
		FieldCity        interface{} `json:"fieldCity"`
		Id               int         `json:"id"`
		Image            interface{} `json:"image"`
		InstallationDate time.Time   `json:"installationDate"`
		Name             string      `json:"name"`
		PeakPower        float64     `json:"peakPower"`
		State            string      `json:"state"`
		Status           string      `json:"status"`
		Zip              interface{} `json:"zip"`
	} `json:"sites"`
}

// https://api.solaredge.com/services/m/so/dashboard/v2/site/4417661/powerflow/latest/?components=consumption,grid,storage
type PowerflowLatestResponse struct {
	BatteryConsumers []string `json:"batteryConsumers"`
	Consumption      struct {
		CurrentPower float64 `json:"currentPower"`
		IsActive     bool    `json:"isActive"`
		IsConsuming  bool    `json:"isConsuming"`
	} `json:"consumption"`
	EnergyConsumers []string `json:"energyConsumers"`
	Grid            struct {
		CurrentPower   float64 `json:"currentPower"`
		HasPowerOutage bool    `json:"hasPowerOutage"`
		IsActive       bool    `json:"isActive"`
		Status         string  `json:"status"`
	} `json:"grid"`
	IsCommunicating bool      `json:"isCommunicating"`
	IsRealTime      bool      `json:"isRealTime"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
	SolarProduction struct {
		CurrentPower float64 `json:"currentPower"`
		IsActive     bool    `json:"isActive"`
		IsProducing  bool    `json:"isProducing"`
	} `json:"solarProduction"`
	Storage struct {
		ChargeLevel  int     `json:"chargeLevel"`
		CurrentPower float64 `json:"currentPower"`
		IsActive     bool    `json:"isActive"`
		Status       string  `json:"status"`
		StoragePlan  string  `json:"storagePlan"`
	} `json:"storage"`
	Unit              string `json:"unit"`
	UpdateRefreshRate int    `json:"updateRefreshRate"`
}

// https://ha.monitoring.solaredge.com/api/homeautomation/v1.0/storage/4417661/getBatteries?triggerHF=false
type GetBatteriesResponse struct {
	C                string        `json:"@c"`
	BillingProviders []interface{} `json:"billingProviders"`
	Devices          []interface{} `json:"devices"`
	DevicesByType    struct {
		BATTERY []struct {
			BatteryState           string  `json:"batteryState"`
			ChargeEnergy           float64 `json:"chargeEnergy"`
			CommunicationStatus    string  `json:"communicationStatus"`
			DischargeEnergy        float64 `json:"dischargeEnergy"`
			EstimatedRemainingTime int     `json:"estimatedRemainingTime"`
			FullPackEnergy         float64 `json:"fullPackEnergy"`
			Info                   struct {
				DeviceId     int64  `json:"deviceId"`
				Manufacturer string `json:"manufacturer"`
				Model        string `json:"model"`
				Name         string `json:"name"`
				PortiaSN     string `json:"portiaSN"`
				SerialNumber string `json:"serialNumber"`
				SwVersion    string `json:"swVersion"`
			} `json:"info"`
			LastCommunicationTime int64   `json:"lastCommunicationTime"`
			LastUpdated           int64   `json:"lastUpdated"`
			NamePlateEnergy       float64 `json:"namePlateEnergy"`
			PowerSavingMode       bool    `json:"powerSavingMode"`
			RemainingEnergy       float64 `json:"remainingEnergy"`
		} `json:"BATTERY"`
		POLESTAR []struct {
			CommunicationStatus string  `json:"communicationStatus"`
			ConnectedDevices    []int64 `json:"connectedDevices"`
			Info                struct {
				DeviceId     int64  `json:"deviceId"`
				SerialNumber string `json:"serialNumber"`
				SwVersion    string `json:"swVersion"`
			} `json:"info"`
			LastCommunicationTime int64       `json:"lastCommunicationTime"`
			LastUpdated           interface{} `json:"lastUpdated"`
		} `json:"POLESTAR"`
	} `json:"devicesByType"`
	ErrorMessages             []interface{} `json:"errorMessages"`
	EssentialDevicesSupported bool          `json:"essentialDevicesSupported"`
	ExcessPvSupported         bool          `json:"excessPvSupported"`
	FieldLastUpdateTS         int64         `json:"fieldLastUpdateTS"`
	GlobalCommunicationStatus string        `json:"globalCommunicationStatus"`
	GridProgram               struct {
		OnBoard bool        `json:"onBoard"`
		Program interface{} `json:"program"`
	} `json:"gridProgram"`
	HttpStatus             int           `json:"httpStatus"`
	InfoMessages           []interface{} `json:"infoMessages"`
	IsUpdated              bool          `json:"isUpdated"`
	RfIdCardsListSupported string        `json:"rfIdCardsListSupported"`
	Status                 string        `json:"status"`
	StorageInfo            struct {
		BackupMinSOE                               float64     `json:"backupMinSOE"`
		BackupReserveControlUpToLimit              int         `json:"backupReserveControlUpToLimit"`
		BatteryProfileManualTouSupported           bool        `json:"batteryProfileManualTouSupported"`
		BatteryProfileTouAllowedByNumberOfInverter bool        `json:"batteryProfileTouAllowedByNumberOfInverter"`
		ChargingInitiator                          string      `json:"chargingInitiator"`
		ExportImportMeterExist                     bool        `json:"exportImportMeterExist"`
		NotManagedByGridServices                   bool        `json:"notManagedByGridServices"`
		PortiaBackupReserveEditCapable             bool        `json:"portiaBackupReserveEditCapable"`
		PortiaBackupReserveViewCapable             bool        `json:"portiaBackupReserveViewCapable"`
		PreviousBackupReserve                      interface{} `json:"previousBackupReserve"`
		StorageFullPackEnergy                      float64     `json:"storageFullPackEnergy"`
		StoragePower                               float64     `json:"storagePower"`
		StorageRemainingEnergy                     float64     `json:"storageRemainingEnergy"`
		TouInfo                                    struct {
			IsSupported bool          `json:"isSupported"`
			Reasons     []interface{} `json:"reasons"`
		} `json:"touInfo"`
		UserBackupReserveCommandExpiration interface{} `json:"userBackupReserveCommandExpiration"`
		UserBackupReserveValue             float64     `json:"userBackupReserveValue"`
		UserBatteryProfileSupported        bool        `json:"userBatteryProfileSupported"`
		WeatherGuardServiceConfiguration   string      `json:"weatherGuardServiceConfiguration"`
		WeatherGuardServiceStatus          string      `json:"weatherGuardServiceStatus"`
		WeatherGuardStartTime              interface{} `json:"weatherGuardStartTime"`
	} `json:"storageInfo"`
	UpdateRefreshRate int           `json:"updateRefreshRate"`
	WarningMessages   []interface{} `json:"warningMessages"`
}

type GetBatteryModeResponse struct {
	Status           string        `json:"status"`
	HttpStatus       int           `json:"httpStatus"`
	WarningMessages  []interface{} `json:"warningMessages"`
	InfoMessages     []interface{} `json:"infoMessages"`
	ErrorMessages    []interface{} `json:"errorMessages"`
	BatteryMode      string        `json:"batteryMode"`
	TouConfiguration struct {
		OwnerConfiguration []interface{} `json:"ownerConfiguration"`
		TouPlan            interface{}   `json:"touPlan"`
	} `json:"touConfiguration"`
	ManualTouConfiguration struct {
		TouConfiguration struct {
			OwnerConfiguration []struct {
				Months       []string `json:"months"`
				DaysSegments []struct {
					Days          []string `json:"days"`
					HoursSegments []struct {
						From          int `json:"from"`
						To            int `json:"to"`
						ManualTouData struct {
							BatteryMode string `json:"batteryMode"`
						} `json:"manualTouData"`
					} `json:"hoursSegments"`
				} `json:"daysSegments"`
			} `json:"ownerConfiguration"`
		} `json:"touConfiguration"`
	} `json:"manualTouConfiguration"`
	LastUpdateDate time.Time `json:"lastUpdateDate"`
}

type BatteryMode string

const BatteryModeManualToU BatteryMode = "MANUAL_TOU"
const BatteryModeMSC BatteryMode = "MSC" // Maximize Self Consumption

type TouBatteryMode string

type PutBatteryModeManualTouData struct {
	BatteryMode string `json:"batteryMode"`
}

type PutBatteryModeHoursSegment struct {
	From          int                         `json:"from"`
	To            int                         `json:"to"`
	ManualTouData PutBatteryModeManualTouData `json:"manualTouData"`
}

type PutBatteryModeDaySegment struct {
	Days          []string                     `json:"days"`
	HoursSegments []PutBatteryModeHoursSegment `json:"hoursSegments"`
}

type PutBatteryModeOwnerConfiguration struct {
	Months       []string                   `json:"months"`
	DaysSegments []PutBatteryModeDaySegment `json:"daysSegments"`
}

type PutBatteryModeTouConfiguration struct {
	TouPlan            interface{}                        `json:"touPlan,omitempty"`
	OwnerConfiguration []PutBatteryModeOwnerConfiguration `json:"ownerConfiguration"`
}

type PutBatteryModeManualTouConfiguration struct {
	TouConfiguration PutBatteryModeTouConfiguration `json:"touConfiguration"`
}

type PutBatteryModeRequest struct {
	BatteryMode            BatteryMode                          `json:"batteryMode"`
	ManualTouConfiguration PutBatteryModeManualTouConfiguration `json:"manualTouConfiguration"`
	TouConfiguration       PutBatteryModeTouConfiguration       `json:"touConfiguration"`
}

type PutBatteryModeResponse struct {
	Status          string        `json:"status"`
	HttpStatus      int           `json:"httpStatus"`
	WarningMessages []interface{} `json:"warningMessages"`
	InfoMessages    []interface{} `json:"infoMessages"`
	ErrorMessages   []struct {
		Message        string        `json:"message"`
		MessageCode    interface{}   `json:"messageCode"`
		Params         []interface{} `json:"params"`
		TranslationKey string        `json:"translationKey"`
	} `json:"errorMessages"`
}
