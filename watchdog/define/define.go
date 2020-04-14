// type
package define

const (
	Port                = ":50501"
	AlertUrl            = "http://122.10.232.11:8888/alert"
	ElectMasterInterval = 1
	DiskAlertThreshold  = 80
	PublicKeyString     = `-----BEGIN PRIVATE KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAor5fLr86uql4XCc1lTV1
q+wMiDrCkdtrDjHedYr95iPOarrQMlgGHX6D1YMzAkyfbCzHu3VldtSz02iXvlXU
GAMY5o3ME1Cph+oNjPp9WewZv0lbvcuoD2P8kueEMb607VAcuAfQL8gWeaCakkJX
FNaydifAfKt2fzTBnH7Z8hBaPlKG2xskbNaOKVQGByI8eJNsXusD1+L/P4WWZaMV
2GXH8OhtuWZn/QyFVwYfXt0RGfJkg7X6w0WIdOJ1mJNWdwWSiFIunHTVp8ebQ2rs
LjGv1RMn27dOqJRCHfKSxOl9j7BDH4F82pJriqkwPJVSLvrql43RBaWek2dkYcPy
DwIDAQAB
-----END PRIVATE KEY-----`
)

var ClusterHosts []string
var Master string
var IsMaster bool
var Hostname string
var Vote []string

type AlertMessageStruct struct {
	Shipid  string            `json:"shipid"`
	Title   string            `json:"title"`
	Message map[string]string `json:"message"`
	Time    string            `json:"time"`
}

const HelpInfo = `  get:
    master
    status
  test:
    alert`
