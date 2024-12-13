package solaredgego

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Options struct {
	Username   string
	Password   string
	HttpClient *http.Client
}

const (
	LOGIN_URL     = "https://api.solaredge.com/solaredge-apigw/api/login"
	SITES_URL     = "https://api.solaredge.com/services/m/so/sites/?status=ACTIVE%2CPENDING&sortName=name&sortOrder=ASC&page=0&pageSize=1&isDemoSite=false"
	POWERFLOW_URL = "https://api.solaredge.com/services/m/so/dashboard/v2/site/%d/powerflow/latest/?components=consumption,grid,storage"
	BATTERIES_URL = "https://ha.monitoring.solaredge.com/api/homeautomation/v1.0/storage/%d/getBatteries?triggerHF=false"
)

type SolarEdge interface {
	Login() (*SitesResponse, error)
	GetData(siteId int) (*PowerflowLatestResponse, *GetBatteriesResponse, error)
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

func (n *solarEdge) GetData(siteId int) (*PowerflowLatestResponse, *GetBatteriesResponse, error) {
	var powerflowLatestResponse PowerflowLatestResponse
	var getBatteriesResponse GetBatteriesResponse

	_, err := n.client.R().SetResult(&powerflowLatestResponse).Get(fmt.Sprintf(POWERFLOW_URL, siteId))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get powerflow data: %s", err.Error())
	}

	_, err = n.client.R().SetResult(&getBatteriesResponse).Get(fmt.Sprintf(BATTERIES_URL, siteId))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get powerflow data: %s", err.Error())
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
