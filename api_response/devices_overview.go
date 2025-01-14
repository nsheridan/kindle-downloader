package api_response

type GetDevicesOverviewResponse struct {
	Success            bool               `json:"success"`
	GetDevicesOverview GetDevicesOverview `json:"GetDevicesOverview"`
}

type GetDevicesOverview struct {
	DeviceList                   []Device `json:"deviceList"`
	DevicePromotionOfferDataList []any    `json:"devicePromotionOfferDataList"`
	Successful                   bool     `json:"successful"`
}

type Device struct {
	ChildDevices                any                `json:"childDevices"`
	DeviceTypeID                string             `json:"deviceTypeID"`
	DeviceImageURL              any                `json:"deviceImageURL"`
	Metadata                    Metadata           `json:"metadata,omitempty"`
	IsChildDevice               bool               `json:"isChildDevice"`
	DisplayCategoryFriendlyName string             `json:"displayCategoryFriendlyName"`
	LastUsedDate                any                `json:"lastUsedDate"`
	ChildDevice                 bool               `json:"childDevice"`
	FormattedLastRegisteredDate string             `json:"formattedLastRegisteredDate"`
	DeviceFamily                string             `json:"deviceFamily"`
	DeviceIdentificationNumber  any                `json:"deviceIdentificationNumber"`
	DeviceName                  string             `json:"deviceName"`
	DeviceSerialNumber          string             `json:"deviceSerialNumber"`
	LastRegisteredDate          LastRegisteredDate `json:"lastRegisteredDate"`
	IsDefaultDevice             bool               `json:"isDefaultDevice"`
	DisplayCategoryImage        string             `json:"displayCategoryImage"`
	DeviceTypeString            string             `json:"deviceTypeString"`
	DefaultDevice               bool               `json:"defaultDevice"`
	CustomerID                  string             `json:"customerID"`
	ParentDevice                any                `json:"parentDevice"`
	DeviceClassification        string             `json:"deviceClassification"`
	DeviceAccountID             string             `json:"deviceAccountID"`
	DeviceGroup                 any                `json:"deviceGroup"`
	AlexaOnDevice               bool               `json:"alexaOnDevice"`
}

type Zone struct {
	Fixed bool   `json:"fixed"`
	ID    string `json:"id"`
}
type Chronology struct {
	Zone Zone `json:"zone"`
}

type LastRegisteredDate struct {
	Year           int        `json:"year"`
	DayOfYear      int        `json:"dayOfYear"`
	EqualNow       bool       `json:"equalNow"`
	Weekyear       int        `json:"weekyear"`
	Chronology     Chronology `json:"chronology"`
	WeekOfWeekyear int        `json:"weekOfWeekyear"`
	SecondOfMinute int        `json:"secondOfMinute"`
	MillisOfDay    int        `json:"millisOfDay"`
	MonthOfYear    int        `json:"monthOfYear"`
	DayOfWeek      int        `json:"dayOfWeek"`
	BeforeNow      bool       `json:"beforeNow"`
	MinuteOfDay    int        `json:"minuteOfDay"`
	DayOfMonth     int        `json:"dayOfMonth"`
	Era            int        `json:"era"`
	Zone           Zone       `json:"zone"`
	YearOfCentury  int        `json:"yearOfCentury"`
	HourOfDay      int        `json:"hourOfDay"`
	CenturyOfEra   int        `json:"centuryOfEra"`
	SecondOfDay    int        `json:"secondOfDay"`
	YearOfEra      int        `json:"yearOfEra"`
	Millis         int64      `json:"millis"`
	MinuteOfHour   int        `json:"minuteOfHour"`
	MillisOfSecond int        `json:"millisOfSecond"`
	AfterNow       bool       `json:"afterNow"`
}

type Metadata struct {
	DeviceImageURL                    string   `json:"deviceImageURL"`
	DeviceCustomerSupportStatus       string   `json:"deviceCustomerSupportStatus"`
	SecurityUpdateSupportedThrough    string   `json:"securityUpdateSupportedThrough"`
	DeviceTypeString                  string   `json:"deviceTypeString"`
	Namespace                         string   `json:"namespace"`
	DeviceFamily                      string   `json:"deviceFamily"`
	Actions                           []string `json:"actions"`
	CustomDeleteVoiceRecordingMessage string   `json:"customDeleteVoiceRecordingMessage"`
}
