package contactmanaging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"infrastructure/myerror"
	"infrastructure/myhttp"
	"strconv"
	"strings"
	"time"
)

const (
	createContactURL                  = "/users/:userID/contacts"
	updateContactURL                  = "/users/:userID/contacts/:contactID"
	getContactURL                     = "/users/:userID/contacts/:contactID"
	deleteContactURL                  = "/users/:userID/contacts/:contactID"
	searchContactsURL                 = "/users/:userID/contacts"
	searchContactsPaginationFormatURL = "%s?phone=%s&firstName=%s&lastName=%s&address=%s&limit=%d&offset=%d"
)

func ListenHTTP(s Service) {
	r := gin.Default()

	r.POST(createContactURL, makeHTTPEndpointCreateContact(s))
	r.PUT(updateContactURL, makeHTTPEndpointUpdateContact(s))
	r.GET(getContactURL, makeHTTPEndpointGetContact(s))
	r.GET(searchContactsURL, makeHTTPEndpointSearchContacts(s))
	r.DELETE(deleteContactURL, makeHTTPEndpointDeleteContact(s))

	fmt.Println("Server listening on port 8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

// Create
type createContactHTTPRequest struct {
	UserID    string
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
}

func (r createContactHTTPRequest) ToCreateContactRequest() createContactRequest {
	return createContactRequest{
		UserID:    r.UserID,
		Phone:     r.Phone,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Address:   r.Address,
	}
}

type createContactHTTPResponse struct {
	ID string `json:"id"`
}

func makeHTTPEndpointCreateContact(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := decodeCreateContactHTTPRequest(c)
		if err != nil {
			encodeCreateContactResponse(c, createContactResponse{}, err)
			return
		}

		resp, err := endpointCreateContact(c, s, req)
		encodeCreateContactResponse(c, resp, err)
	}
}

func decodeCreateContactHTTPRequest(c *gin.Context) (createContactRequest, error) {
	var req createContactHTTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return createContactRequest{}, myerror.Wrap(err, "decodeCreateContactHTTPRequest")
	}

	req.UserID = c.Param("userID")

	return req.ToCreateContactRequest(), nil
}

func encodeCreateContactResponse(c *gin.Context, resp createContactResponse, err error) {
	jsonResponse := createContactHTTPResponse{
		ID: resp.ID,
	}

	if err != nil {
		myhttp.EncodeJSONError(c, err)
	} else {
		myhttp.EncodeJSONSuccess(c, jsonResponse)
	}
}

// Update
type updateContactHTTPRequest struct {
	UserID          string
	ContactID       string
	Phone           string    `json:"phone"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Address         string    `json:"address"`
	UpdateAtVersion time.Time `json:"updatedAt"`
}

func (r updateContactHTTPRequest) ToUpdateContactRequest() updateContactRequest {
	return updateContactRequest{
		UserID:          r.UserID,
		ContactID:       r.ContactID,
		Phone:           r.Phone,
		FirstName:       r.FirstName,
		LastName:        r.LastName,
		Address:         r.Address,
		UpdateAtVersion: r.UpdateAtVersion,
	}
}

func makeHTTPEndpointUpdateContact(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := decodeUpdateContactHTTPRequest(c)
		if err != nil {
			myhttp.EncodeJSONError(c, err)
			return
		}

		err = endpointUpdateContact(c, s, req)
		if err != nil {
			myhttp.EncodeJSONError(c, err)
			return
		}

		myhttp.EncodeJSONSuccess(c, struct{}{})
	}
}

func decodeUpdateContactHTTPRequest(c *gin.Context) (updateContactRequest, error) {
	var req updateContactHTTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return updateContactRequest{}, myerror.Wrap(err, "decodeUpdateContactHTTPRequest")
	}

	req.UserID = c.Param("userID")
	req.ContactID = c.Param("contactID")

	return req.ToUpdateContactRequest(), nil
}

// Get
type getContactHTTPRequest struct {
	UserID    string
	ContactID string
}

func (r getContactHTTPRequest) ToGetContactRequest() getContactRequest {
	return getContactRequest{
		UserID:    r.UserID,
		ContactID: r.ContactID,
	}
}

type getContactHTTPResponse struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func makeHTTPEndpointGetContact(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := decodeGetContactHTTPRequest(c)
		if err != nil {
			myhttp.EncodeJSONError(c, err)
			return
		}

		resp, err := endpointGetContact(c, s, req)
		encodeGetContactResponse(c, resp, err)
	}
}

func decodeGetContactHTTPRequest(c *gin.Context) (getContactRequest, error) {
	var req getContactHTTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return getContactRequest{}, myerror.Wrap(err, "decodeGetContactHTTPRequest")
	}

	req.UserID = c.Param("userID")
	req.ContactID = c.Param("contactID")

	return req.ToGetContactRequest(), nil
}

func encodeGetContactResponse(c *gin.Context, resp getContactResponse, err error) {
	jsonResponse := getContactResponseToJSON(resp)

	if err != nil {
		myhttp.EncodeJSONError(c, err)
	} else {
		myhttp.EncodeJSONSuccess(c, jsonResponse)
	}
}

func getContactResponseToJSON(resp getContactResponse) getContactHTTPResponse {
	return getContactHTTPResponse{
		ID:        resp.ID,
		Phone:     resp.Phone,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Address:   resp.Address,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

// Search
type searchContactsHTTPRequest struct {
	UserID    string
	Phone     string
	FirstName string
	LastName  string
	Address   string
	Limit     int
	Offset    int
}

func (r searchContactsHTTPRequest) ToSearchContactsRequest() searchContactsRequest {
	return searchContactsRequest{
		UserID:    r.UserID,
		Phone:     r.Phone,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Address:   r.Address,
		Limit:     r.Limit,
		Offset:    r.Offset,
	}
}

type searchContactsHTTPResponse struct {
	Contacts   []getContactHTTPResponse `json:"contacts"`
	Pagination myhttp.Pagination        `json:"pagination"`
}

func makeHTTPEndpointSearchContacts(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := decodeSearchContactsHTTPRequest(c)
		if err != nil {
			encodeSearchContactsResponse(c, searchContactsResponse{}, err)
			return
		}

		resp, err := endpointSearchContacts(c, s, req.ToSearchContactsRequest())
		encodeSearchContactsResponse(c, resp, err)
	}
}

func decodeSearchContactsHTTPRequest(c *gin.Context) (searchContactsHTTPRequest, error) {
	req := searchContactsHTTPRequest{
		UserID:    c.Param("userID"),
		Phone:     c.Query("phone"),
		FirstName: c.Query("firstName"),
		LastName:  c.Query("lastName"),
		Address:   c.Query("address"),
	}

	var err error
	if limitStr := c.Query("limit"); limitStr != "" {
		req.Limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return searchContactsHTTPRequest{}, myerror.Wrap(err, "decodeSearchContactsHTTPRequest")
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		req.Offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return searchContactsHTTPRequest{}, myerror.Wrap(err, "decodeSearchContactsHTTPRequest")
		}
	}

	return req, nil
}

func formatSearchContactsURL(userID, phone, firstName, lastName, address string, limit, offset int) string {
	return fmt.Sprintf(searchContactsPaginationFormatURL,
		strings.Replace(searchContactsURL, ":userID", userID, 1),
		phone,
		firstName,
		lastName,
		address,
		limit,
		offset,
	)
}

func encodeSearchContactsResponse(c *gin.Context, resp searchContactsResponse, err error) {
	req, err := decodeSearchContactsHTTPRequest(c)
	if err != nil {
		myhttp.EncodeJSONError(c, err)
		return
	}

	var nextURL, prevURL string
	if req.Offset > 0 {
		prevURL = formatSearchContactsURL(req.UserID, req.Phone, req.FirstName, req.LastName, req.Address, req.Limit, req.Offset-req.Limit)
	}
	if len(resp.Contacts) > 0 {
		nextURL = formatSearchContactsURL(req.UserID, req.Phone, req.FirstName, req.LastName, req.Address, req.Limit, req.Offset+req.Limit)
	}

	contacts := make([]getContactHTTPResponse, 0, len(resp.Contacts))
	for _, c := range resp.Contacts {
		contacts = append(contacts, getContactResponseToJSON(c))
	}

	jsonResponse := searchContactsHTTPResponse{
		Contacts: contacts,
		Pagination: myhttp.Pagination{
			Previous: prevURL,
			Next:     nextURL,
		},
	}

	if err != nil {
		myhttp.EncodeJSONError(c, err)
	} else {
		myhttp.EncodeJSONSuccess(c, jsonResponse)
	}
}

// Delete
type deleteContactHTTPRequest struct {
	UserID    string
	ContactID string
}

func (r deleteContactHTTPRequest) ToDeleteContactRequest() deleteContactRequest {
	return deleteContactRequest{
		UserID:    r.UserID,
		ContactID: r.ContactID,
	}
}

func makeHTTPEndpointDeleteContact(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := decodeDeleteContactHTTPRequest(c)
		if err := endpointDeleteContact(c, s, req.ToDeleteContactRequest()); err != nil {
			myhttp.EncodeJSONError(c, err)
			return
		}

		myhttp.EncodeJSONSuccess(c, struct{}{})
	}
}

func decodeDeleteContactHTTPRequest(c *gin.Context) deleteContactHTTPRequest {
	return deleteContactHTTPRequest{
		UserID:    c.Param("userID"),
		ContactID: c.Param("contactID"),
	}
}
