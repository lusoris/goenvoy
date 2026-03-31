package flaresolverr

// Response is the standard FlareSolverr API response.
type Response struct {
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	Solution       *Solution `json:"solution"`
	StartTimestamp int64     `json:"startTimestamp"`
	EndTimestamp   int64     `json:"endTimestamp"`
	Version        string    `json:"version"`
}

// Solution contains the resolved page data.
type Solution struct {
	URL            string            `json:"url"`
	Status         int               `json:"status"`
	Headers        map[string]string `json:"headers"`
	Response       string            `json:"response"`
	Cookies        []Cookie          `json:"cookies"`
	UserAgent      string            `json:"userAgent"`
	TurnstileToken string            `json:"turnstile_token"`
}

// Cookie represents an HTTP cookie from the solved challenge.
type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	Size     int     `json:"size"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	Session  bool    `json:"session"`
	SameSite string  `json:"sameSite"`
}

// RequestOptions configures optional parameters for a request.
type RequestOptions struct {
	Session           string   `json:"session,omitempty"`
	MaxTimeout        int      `json:"maxTimeout,omitempty"`
	Cookies           []Cookie `json:"cookies,omitempty"`
	ReturnOnlyCookies bool     `json:"returnOnlyCookies,omitempty"`
	Proxy             *Proxy   `json:"proxy,omitempty"`
	WaitInSeconds     int      `json:"waitInSeconds,omitempty"`
	DisableMedia      bool     `json:"disableMedia,omitempty"`
}

// Proxy configures proxy settings for a request.
type Proxy struct {
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
