package contactmanaging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"infrastructure/myerror"
	"infrastructure/myhttp"
	"time"
)

func ListenHTTP(s Service) {
	r := gin.Default()

	r.POST("/users/:userID/contacts", makeHTTPEndpointCreateContact(s))
	r.PUT("/users/:userID/contacts/:contactID", makeHTTPEndpointUpdateContact(s))

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
	UpdateAtVersion time.Time `json:"updateAt"`
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
