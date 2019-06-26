package main

// PeeringDB, Comcast example: https://beta.peeringdb.com/api/net?asn=7922

/*
{"meta": {}, "data": [{"id": 822, "org_id": 1061, "name": "Comcast", "aka": "Comcast Backbone - 7922",
	"website": "", "asn": 7922, "looking_glass": "", "route_server": "telnet://route-server.newyork.ny.ibone.comcast.net",
	"irr_as_set": "AS-COMCAST-IBONE", "info_type": "Cable/DSL/ISP", "info_prefixes4": 14000, "info_prefixes6": 600,
	"info_traffic": "1 Tbps+", "info_ratio": "Balanced", "info_scope": "North America", "info_unicast": true,
	"info_multicast": false, "info_ipv6": true, "notes": "peering@comcast.com is the best email address to send requests
	for IPv6 & settled peering after reviewing http://www.comcast.com/peering.  It is not a 24x7 operations contact and
	should not be assumed appropriate for that purpose.\n\nWe do not offer peering, paid or otherwise, on the shared fabric
	public switches at any IX.\n\nREMINDER:  Max prefix settings:\nIPv4 set to 30K\nIPv6 set to 2K\n\nPlease ensure that
	you are set up to accept these prefixes from AS7922 for server l.root-servers.net -\n199.7.83.0/24  &  2001:500:3::/48",
	"policy_url": "http://www.comcast.com/peering", "policy_general": "Selective", "policy_locations": "Required - US",
	"policy_ratio": true, "policy_contracts": "Required", "created": "2010-05-24T13:13:54Z", "updated": "2017-03-28T16:08:20Z",
	"status": "ok"}]}
*/

// NetworkSerializer is the https://peeringdb.com/apidocs/#!/net/Network_list data model
/* type NetworkSerializer struct {
	ID              int               `json:"id"`
	OrgID           int               `json:"org_id"`
	Org             string            `json:"org"`
	Name            string            `json:"name"`
	Aka             string            `json:"aka"`
	Website         string            `json:"website"`
	ASN             int               `json:"asn"`
	LookingGlass    string            `json:"looking_glass"`
	RouteServer     string            `json:"route_server"`
	IrrAsSet        string            `json:"irr_as_set"`
	InfoType        string            `json:"info_type"`
	InfoPrefixes4   int               `json:"info_prefixes4"`
	InfoPrefixes6   int               `json:"info_prefixes6"`
	InfoTraffic     string            `json:"info_traffic"`
	InfoRatio       string            `json:"info_ratio"`
	InfoScope       string            `json:"info_scope"`
	InfoUnicast     bool              `json:"info_unicast"`
	InfoMulticast   bool              `json:"info_multicast"`
	InfoIPv6        bool              `json:"info_ipv6"`
	Notes           string            `json:"notes"`
	PolicyURL       string            `json:"policy_url"`
	PolicyGeneral   string            `json:"policy_general"`
	PolicyLocations string            `json:"policy_locations"`
	PolicyRatio     bool              `json:"policy_ratio"`
	PolicyContracts string            `json:"policy_contracts"`
	NetfacSet       map[string]string `json:"netfac_set"`
	NetixlanSet     map[string]string `json:"netixlan_set"`
	PocSet          map[string]string `json:"poc_set"`
	Created         string            `json:"created"`
	Updated         string            `json:"updated"`
	Status          string            `json:"status"`
} */

// NetworkSerializer is the https://peeringdb.com/apidocs/#!/net/Network_list data model
//
// InfoType:
// ['' or 'Not Disclosed' or 'NSP' or 'Content' or 'Cable/DSL/ISP' or 'Enterprise'
// or 'Educational/Research' or 'Non-Profit' or 'Route Server']
// InfoTraffic:
// ['' or '0-20 Mbps' or '20-100Mbps' or '100-1000Mbps' or '1-5Gbps' or '5-10Gbps'
// or '10-20Gbps' or '20-50 Gbps' or '50-100 Gbps' or '100+ Gbps' or '100-200 Gbps'
// or '200-300 Gbps' or '300-500 Gbps' or '500-1000 Gbps' or '1 Tbps+' or '10 Tbps+']
// InfoRatio:
// ['' or 'Not Disclosed' or 'Heavy Outbound' or 'Mostly Outbound' or 'Balanced'
// or 'Mostly Inbound' or 'Heavy Inbound']
// InfoScope:
// ['' or 'Not Disclosed' or 'Regional' or 'North America' or 'Asia Pacific' or 'Europe'
// or 'South America' or 'Africa' or 'Australia' or 'Middle East' or 'Global']
// PolicyGeneral:
// ['Open' or 'Selective' or 'Restrictive' or 'No']
// PolicyLocations:
// ['Not Required' or 'Preferred' or 'Required - US' or 'Required - EU' or 'Required - International']
// PolicyContracts:
// ['Not Required' or 'Private Only' or 'Required']
// NetfacSet:
// (array[NetworkFacilitySerializer])
// NetixlanSet:
// (array[NetworkIXLanSerializer])
// PocSet:
// (array[NetworkContactSerializer])
//
