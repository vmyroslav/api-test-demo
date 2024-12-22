/*
FakeRESTApi.Web V1

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the CoverPhoto type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CoverPhoto{}

// CoverPhoto struct for CoverPhoto
type CoverPhoto struct {
	Id *int32 `json:"id,omitempty"`
	IdBook *int32 `json:"idBook,omitempty"`
	Url NullableString `json:"url,omitempty"`
}

// NewCoverPhoto instantiates a new CoverPhoto object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCoverPhoto() *CoverPhoto {
	this := CoverPhoto{}
	return &this
}

// NewCoverPhotoWithDefaults instantiates a new CoverPhoto object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCoverPhotoWithDefaults() *CoverPhoto {
	this := CoverPhoto{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CoverPhoto) GetId() int32 {
	if o == nil || IsNil(o.Id) {
		var ret int32
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CoverPhoto) GetIdOk() (*int32, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CoverPhoto) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int32 and assigns it to the Id field.
func (o *CoverPhoto) SetId(v int32) {
	o.Id = &v
}

// GetIdBook returns the IdBook field value if set, zero value otherwise.
func (o *CoverPhoto) GetIdBook() int32 {
	if o == nil || IsNil(o.IdBook) {
		var ret int32
		return ret
	}
	return *o.IdBook
}

// GetIdBookOk returns a tuple with the IdBook field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CoverPhoto) GetIdBookOk() (*int32, bool) {
	if o == nil || IsNil(o.IdBook) {
		return nil, false
	}
	return o.IdBook, true
}

// HasIdBook returns a boolean if a field has been set.
func (o *CoverPhoto) HasIdBook() bool {
	if o != nil && !IsNil(o.IdBook) {
		return true
	}

	return false
}

// SetIdBook gets a reference to the given int32 and assigns it to the IdBook field.
func (o *CoverPhoto) SetIdBook(v int32) {
	o.IdBook = &v
}

// GetUrl returns the Url field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CoverPhoto) GetUrl() string {
	if o == nil || IsNil(o.Url.Get()) {
		var ret string
		return ret
	}
	return *o.Url.Get()
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CoverPhoto) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Url.Get(), o.Url.IsSet()
}

// HasUrl returns a boolean if a field has been set.
func (o *CoverPhoto) HasUrl() bool {
	if o != nil && o.Url.IsSet() {
		return true
	}

	return false
}

// SetUrl gets a reference to the given NullableString and assigns it to the Url field.
func (o *CoverPhoto) SetUrl(v string) {
	o.Url.Set(&v)
}
// SetUrlNil sets the value for Url to be an explicit nil
func (o *CoverPhoto) SetUrlNil() {
	o.Url.Set(nil)
}

// UnsetUrl ensures that no value is present for Url, not even an explicit nil
func (o *CoverPhoto) UnsetUrl() {
	o.Url.Unset()
}

func (o CoverPhoto) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CoverPhoto) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.IdBook) {
		toSerialize["idBook"] = o.IdBook
	}
	if o.Url.IsSet() {
		toSerialize["url"] = o.Url.Get()
	}
	return toSerialize, nil
}

type NullableCoverPhoto struct {
	value *CoverPhoto
	isSet bool
}

func (v NullableCoverPhoto) Get() *CoverPhoto {
	return v.value
}

func (v *NullableCoverPhoto) Set(val *CoverPhoto) {
	v.value = val
	v.isSet = true
}

func (v NullableCoverPhoto) IsSet() bool {
	return v.isSet
}

func (v *NullableCoverPhoto) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCoverPhoto(val *CoverPhoto) *NullableCoverPhoto {
	return &NullableCoverPhoto{value: val, isSet: true}
}

func (v NullableCoverPhoto) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCoverPhoto) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


