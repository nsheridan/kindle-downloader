package api_response

type GetContentOwnershipDataResponse struct {
	Success                 bool                    `json:"success"`
	GetContentOwnershipData GetContentOwnershipData `json:"GetContentOwnershipData"`
}

type BookProducerDetails struct {
	Role string `json:"role"`
	Name string `json:"name"`
	Asin string `json:"asin"`
}

type ContentItem struct {
	ReadStatus                  string                `json:"readStatus"`
	TargetDevices               any                   `json:"targetDevices"`
	BookProducerDetails         []BookProducerDetails `json:"bookProducerDetails"`
	OrderID                     string                `json:"orderId"`
	IsNellOptimized             bool                  `json:"isNellOptimized"`
	Title                       string                `json:"title"`
	IsGiftOption                bool                  `json:"isGiftOption"`
	SortableAuthors             string                `json:"sortableAuthors"`
	IsPurchased                 bool                  `json:"isPurchased"`
	ExcludedDeviceMap           any                   `json:"excludedDeviceMap"`
	GetOrderDetails             bool                  `json:"getOrderDetails"`
	ExpiredPublicLibraryLending bool                  `json:"expiredPublicLibraryLending"`
	ProductImage                string                `json:"productImage"`
	AcquiredDate                string                `json:"acquiredDate"`
	IsDeleteRestrictionEnabled  bool                  `json:"isDeleteRestrictionEnabled"`
	IsGift                      bool                  `json:"isGift"`
	CollectionList              []any                 `json:"collectionList"`
	ContentCategoryType         string                `json:"contentCategoryType"`
	OrderDetailURL              string                `json:"orderDetailURL"`
	ShowProductDetails          bool                  `json:"showProductDetails"`
	IsContentValid              bool                  `json:"isContentValid"`
	CanLoan                     bool                  `json:"canLoan"`
	StatusFromPlatformSearch    string                `json:"statusFromPlatformSearch"`
	UdlCategory                 string                `json:"udlCategory"`
	RenderDownloadElements      bool                  `json:"renderDownloadElements"`
	AcquiredTime                float64               `json:"acquiredTime"`
	SortableTitle               string                `json:"sortableTitle"`
	RefundEligibility           bool                  `json:"refundEligibility"`
	OriginType                  string                `json:"originType"`
	CapabilityList              []string              `json:"capabilityList"`
	DpURL                       string                `json:"dpURL"`
	IsInstitutionalRental       bool                  `json:"isInstitutionalRental"`
	CollectionCount             int                   `json:"collectionCount"`
	IsAudibleOwned              bool                  `json:"isAudibleOwned"`
	Asin                        string                `json:"asin"`
	IsKCRSupported              bool                  `json:"isKCRSupported"`
	ContentIdentifier           string                `json:"contentIdentifier"`
	Category                    string                `json:"category"`
	IsPrimeShared               bool                  `json:"isPrimeShared"`
	Actions                     []string              `json:"actions"`
	Authors                     string                `json:"authors"`
	ReadAlongSupport            string                `json:"readAlongSupport"`
}

type GetContentOwnershipData struct {
	HasMoreItems             bool          `json:"hasMoreItems"`
	NumberOfItems            int           `json:"numberOfItems"`
	Success                  bool          `json:"success"`
	ContentItems             []ContentItem `json:"items"`
	ContentCategoryReference string        `json:"contentCategoryReference"`
}
