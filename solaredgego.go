package solaredgego

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Options struct {
	Username   string
	Password   string
	HttpClient *http.Client
}

const (
	LOGIN_URL            = "https://api.solaredge.com/solaredge-apigw/api/login"
	SITES_URL            = "https://api.solaredge.com/services/m/so/sites/?status=ACTIVE%2CPENDING&sortName=name&sortOrder=ASC&page=0&pageSize=1&isDemoSite=false"
	POWERFLOW_URL        = "https://api.solaredge.com/services/m/so/dashboard/v2/site/%d/powerflow/latest/?components=consumption,grid,storage"
	BATTERIES_URL        = "https://ha.monitoring.solaredge.com/api/homeautomation/v1.0/storage/%d/getBatteries?triggerHF=false"
	GET_BATTERY_MODE_URL = "https://ha.monitoring.solaredge.com/api/homeautomation/v1.0/storage/%d/batteryMode"
)

type SolarEdge interface {
	Login() (*SitesResponse, error)
	GetData(siteId int) (*PowerflowLatestResponse, *GetBatteriesResponse, error)
	GetBatteryMode(siteId int) (*GetBatteryModeResponse, error)
	PutBatteryMode(siteId int, req *PutBatteryModeRequest) error
}

func NewSolarEdge(options *Options) SolarEdge {
	cj, err := cookiejar.New(nil)
	if err != nil {
		return nil
	}

	client := resty.New().
		SetCookieJar(cj).
		SetHeader("client-version", "3.12").
		SetHeader("user-agent", "okhttp/4.10.0").
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			if r.StatusCode() >= 400 {
				return fmt.Errorf("bad status code: %d, body: %s", r.StatusCode(), r.String())
			}
			return nil
		})
	if options == nil {
		options = &Options{}
	}
	if options.HttpClient != nil {
		client = resty.NewWithClient(options.HttpClient)
	}

	return &solarEdge{
		cj:      cj,
		client:  client,
		options: options,
	}
}

type solarEdge struct {
	options *Options
	client  *resty.Client
	cj      *cookiejar.Jar
}

func (n *solarEdge) PutBatteryMode(siteId int, req *PutBatteryModeRequest) error {
	retries := 5
	for retries > 0 {
		url := fmt.Sprintf("https://ha.monitoring.solaredge.com/api/homeautomation/v1.0/storage/%d/batteryMode", siteId)
		data, _ := json.Marshal(req)
		var resp PutBatteryModeResponse
		_, err := n.client.R().SetBody(req).SetResult(&resp).Put(url)
		if err == nil && resp.HttpStatus == 200 {
			return nil
		}

		if err != nil {
			fmt.Printf("PutBattery error: %s\n", err.Error())
		} else if resp.HttpStatus != 200 {
			fmt.Printf("PutBattery http error: %s %d (%+v)\n", resp.Status, resp.HttpStatus, resp)
		}

		time.Sleep(time.Second * 20)
		retries--
	}

	return fmt.Errorf("PutBattery failed")
}

func (n *solarEdge) GetBatteryMode(siteId int) (*GetBatteryModeResponse, error) {
	{
		x, err := n.client.R().Get(fmt.Sprintf("https://ha.monitoring.solaredge.com/api/homeautomation/v1.0/sites/%d/excessPvPrioritiesV2", siteId))
		if err != nil {
			panic(err)
		}
		fmt.Println(string(x.Body()))
	}

	var resp GetBatteryModeResponse
	x, err := n.client.R().SetResult(&resp).Get(fmt.Sprintf(GET_BATTERY_MODE_URL, siteId))
	if err != nil {
		return nil, fmt.Errorf("Failed to get powerflow data: %s", err.Error())
	}
	return &resp, nil
}

func (n *solarEdge) GetData(siteId int) (*PowerflowLatestResponse, *GetBatteriesResponse, error) {
	var powerflowLatestResponse PowerflowLatestResponse
	var getBatteriesResponse GetBatteriesResponse

	_, err := n.client.R().SetResult(&powerflowLatestResponse).Get(fmt.Sprintf(POWERFLOW_URL, siteId))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get powerflow data: %s", err.Error())
	}

	_, err = n.client.R().SetResult(&getBatteriesResponse).Get(fmt.Sprintf(BATTERIES_URL, siteId))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get battery data: %s", err.Error())
	}

	return &powerflowLatestResponse, &getBatteriesResponse, nil
}

func (n *solarEdge) Login() (*SitesResponse, error) {
	{
		_, err := n.client.R().
			SetQueryParam("j_username", n.options.Username).
			SetQueryParam("j_password", n.options.Password).
			Post(LOGIN_URL)
		if err != nil {
			return nil, fmt.Errorf("Failed to login: %s", err.Error())
		}

		// Login response gives us cookies with restricted path but then needs to use them against a different path and a different domain
		// So hack the cookie URL
		u1, err := url.Parse(LOGIN_URL)
		if err != nil {
			return nil, fmt.Errorf("URL parse error: %s", err.Error())
		}
		u2, err := url.Parse(fmt.Sprintf(BATTERIES_URL, 1))
		if err != nil {
			return nil, fmt.Errorf("URL parse error: %s", err.Error())
		}
		for _, cookie := range n.cj.Cookies(u1) {
			{
				newCookie := *cookie
				newCookie.Domain = u1.Host
				newCookie.Path = "/"
				n.cj.SetCookies(u1, []*http.Cookie{&newCookie})
			}

			{
				newCookie := *cookie
				newCookie.Domain = u2.Host
				newCookie.Path = "/"
				n.cj.SetCookies(u2, []*http.Cookie{&newCookie})
			}
		}

		// Note this 302s to /user/details and returns some user XML but we ignore it
	}

	{
		var sitesResponse SitesResponse
		_, err := n.client.R().SetResult(&sitesResponse).Get(SITES_URL)
		if err != nil {
			return nil, fmt.Errorf("Failed to get sites: %s", err.Error())
		}
		return &sitesResponse, nil
	}
}

func BatteryModeMsc() *PutBatteryModeRequest {
	return createBatteryRequest(BatteryModeMSC, false, false, false)
}

func BatteryModeCharge() *PutBatteryModeRequest {
	return createBatteryRequest(BatteryModeManualToU, true, true, false)
}

func BatteryModeDischarge() *PutBatteryModeRequest {
	return createBatteryRequest(BatteryModeManualToU, true, false, false)
}

func BatteryModeDisable() *PutBatteryModeRequest {
	return createBatteryRequest(BatteryModeManualToU, true, false, true) // Disable battery by saying discharge 2 hours ago
}

func createBatteryRequest(batteryMode BatteryMode, touEnabled bool, touCharging bool, touPause bool) *PutBatteryModeRequest {
	req := &PutBatteryModeRequest{
		BatteryMode: batteryMode,
		ManualTouConfiguration: PutBatteryModeManualTouConfiguration{
			TouConfiguration: PutBatteryModeTouConfiguration{
				TouPlan:            nil,
				OwnerConfiguration: []PutBatteryModeOwnerConfiguration{},
			},
		},
		TouConfiguration: PutBatteryModeTouConfiguration{},
	}

	addOwnerConfig := func(t time.Time, touMode string) {
		req.ManualTouConfiguration.TouConfiguration.OwnerConfiguration = append(req.ManualTouConfiguration.TouConfiguration.OwnerConfiguration, PutBatteryModeOwnerConfiguration{
			Months: []string{strings.ToUpper(t.Month().String())},
			DaysSegments: []PutBatteryModeDaySegment{{
				Days: []string{strings.ToUpper(t.Weekday().String())},
				HoursSegments: []PutBatteryModeHoursSegment{{
					From: t.Hour() * 60,
					To:   (t.Hour() + 1) * 60,
					ManualTouData: PutBatteryModeManualTouData{
						BatteryMode: touMode,
					},
				}},
			}},
		})
	}

	if touEnabled {
		const Mode_Charging = "CHARGING"
		const Mode_Discharging = "DISCHARGING"
		const Mode_Paused = "PAUSE"
		now := time.Now()

		if touCharging {
			addOwnerConfig(now, Mode_Charging)
			addOwnerConfig(now.Add(time.Hour*2), Mode_Discharging)
		} else if touPause {
			addOwnerConfig(now, Mode_Paused)
		} else {
			addOwnerConfig(now, Mode_Discharging)
			addOwnerConfig(now.Add(time.Hour*2), Mode_Charging)
		}
	}

	return req
}
