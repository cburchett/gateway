package ontap

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type IscsiService struct {
	Target  IscsiTarget `json:"target,omitempty"`
	Svm     SvmRef      `json:"svm,omitempty"`
	Enabled *bool       `json:"enabled,omitempty"`
}

type IscsiTarget struct {
	Alias string `json:"alias,omitempty"`
}

const returnIscsiRecords string = "?return_records=true"

func (c *Client) GetIscsiServiceBySvmUuid(uuid string) (iscsiService IscsiService, err error) {
	uri := "/api/protocols/san/iscsi/services/" + uuid

	data, err := c.clientGet(uri)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find iSCSI service") {
			return iscsiService, errors.NewNotFound(schema.GroupResource{Group: "gateway.netapp.com", Resource: "StorageVirtualMachine"}, "no iscsi")
		}
		return iscsiService, &apiError{1, err.Error()}
	}

	var resp IscsiService
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, &apiError{2, err.Error()}
	}

	return resp, nil
}

func (c *Client) CreateIscsiService(jsonPayload []byte) (err error) {
	uri := "/api/protocols/san/iscsi/services" + returnIscsiRecords
	_, err = c.clientPost(uri, jsonPayload)
	if err != nil {
		//fmt.Println("Error: " + err.Error())
		return &apiError{1, err.Error()}
	}

	return nil
}

func (c *Client) PatchIscsiService(uuid string, jsonPayload []byte) (err error) {
	uri := "/api/protocols/san/iscsi/services/" + uuid

	_, err = c.clientPatch(uri, jsonPayload)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return &apiError{404, fmt.Sprintf("SVM with UUID \"%s\" not found", uuid)}
		}
		if strings.Contains(err.Error(), "not running") {
			return &apiError{3276916, err.Error()}
		}
		//miscellaneous errror
		return &apiError{1, err.Error()}
	}

	return nil
}

func (c *Client) DeleteIscsiService(uuid string) (err error) {
	uri := "/api/protocols/san/iscsi/services/" + uuid

	_, err = c.clientDelete(uri)
	if err != nil {
		return &apiError{1, err.Error()}
	}

	return nil
}

func (c *Client) GetIscsiInterfacesBySvmUuid(uuid string, servicePolicy string) (lifs IpInterfacesResponse, err error) {
	uri := "/api/network/ip/interfaces" + returnNFSRecords + "&service_policy.name=" + servicePolicy + "&svm.uuid=" + uuid

	data, err := c.clientGet(uri)
	if err != nil {
		return lifs, &apiError{1, err.Error()}
	}

	var resp IpInterfacesResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, &apiError{2, err.Error()}
	}

	return resp, nil
}

func (c *Client) GetIscsiServicePolicyByName(servicePolicy string) (err error) {
	uri := "/api/network/ip/service-policies?name=" + servicePolicy

	_, err = c.clientGet(uri)
	if err != nil {
		return &apiError{1, err.Error()}
	}

	return nil
}

func (c *Client) CreateIscsiServicePolicy(jsonPayload []byte) (err error) {
	uri := "/api/network/ip/service-policies"
	_, err = c.clientPost(uri, jsonPayload)
	if err != nil {
		//fmt.Println("Error: " + err.Error())
		return &apiError{1, err.Error()}
	}
	return nil
}
