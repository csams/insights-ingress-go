package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/redhatinsights/platform-go-middlewares/request_id"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TrackerResponse struct {
	Data     []Status    `json:"data"`
	Duration interface{} `json:"duration"`
}

type Status struct {
	StatusMsg   string `json:"status_msg,omitempty"`
	Date        string `json:"date,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
	Account     string `json:"account,omitempty"`
	OrgID       string `json:"org_id,omitempty"`
	InventoryID string `json:"inventory_id,omitempty"`
	Service     string `json:"service,omitempty"`
	Status      string `json:"status,omitempty"`
}

type MinimalStatus struct {
	StatusMsg   string `json:"status_msg,omitempty"`
	Date        string `json:"date,omitempty"`
	InventoryID string `json:"inventory_id,omitempty"`
	Service     string `json:"service,omitempty"`
	Status      string `json:"status,omitempty"`
}

// NewHandlers returns an http handler for tracking
func NewHandler(c CompletedConfig) http.HandlerFunc {

	logerr := func(msg string, err error) {
		c.Log.WithFields(logrus.Fields{"error": err}).Error(msg)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var id identity.XRHID
		reqID := chi.URLParam(r, "requestID")

		requestLogger := c.Log.WithFields(logrus.Fields{"source_host": c.Common.Hostname, "name": "ingress"})

		if c.Common.Auth {
			id = identity.Get(r.Context())
		}

		if !isValidUUID(reqID) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request id is not an uuid"))
			return
		}

		verbosity, _ := strconv.Atoi(r.URL.Query().Get("verbosity"))

		response, err := c.Client.Get(c.Url + reqID)
		if err != nil {
			logerr("Failed to get payload status", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logerr("Unable to read response body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var pt TrackerResponse
		var responseBody []byte
		if err = json.Unmarshal(body, &pt); err != nil {
			logerr("Failed to unmarshal tracker json", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(pt.Data) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var subjectDN string
		if id.Identity.X509.SubjectDN != "" {
			subjectSplit := strings.Split(id.Identity.X509.SubjectDN, "=")
			subjectDN = subjectSplit[len(subjectSplit)-1]
		}
		fmt.Print(id.Identity.Type)
		fmt.Print(subjectDN)
		if id.Identity.Type != "Associate" && subjectDN != "insightspipelineqe" {
			if !isIdAuthorized(id.Identity, pt.Data[0].Account, pt.Data[0].OrgID) {

				incomingRequestID := request_id.GetReqID(r.Context())

				requestLogger.WithFields(logrus.Fields{"requestID": reqID,
					"incomingRequestID":         incomingRequestID,
					"pt.Data[0].Account":        pt.Data[0].Account,
					"pt.Data[0].OrgID":          pt.Data[0].OrgID,
					"id.Identity.AccountNumber": id.Identity.AccountNumber,
					"id.Identity.OrgID":         id.Identity.OrgID,
					"id.Identity.Type":          id.Identity.Type,
					"id.Identity.AuthType":      id.Identity.AuthType,
				}).Info("Returning a 403 while retrieving a payload from payload-tracker")

				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		// Response with minimal status data by default
		latestStatus := pt.Data[0]
		ms := MinimalStatus{
			Status:      latestStatus.Status,
			Date:        latestStatus.Date,
			StatusMsg:   latestStatus.StatusMsg,
			Service:     latestStatus.Service,
			InventoryID: latestStatus.InventoryID,
		}

		responseBody, err = json.Marshal(&ms)
		if err != nil {
			logerr("Failed to marshal JSON response", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if verbosity >= 2 {
			w.Write(body)
		} else {
			w.Write(responseBody)
		}
	}
}

func isIdAuthorized(identity identity.Identity, accountNumber string, orgID string) bool {
	return identity.AccountNumber == accountNumber || identity.OrgID == orgID
}

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}