// Package module are modules in ghip
package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// IP is a model saving the url and it's remote ip. both from local or ipaddress.com
type IP struct {
	URL       string `gorm:"primaryKey; not null"`
	LocalIP   string
	ADDRCOMIP datatypes.JSON
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Write will write a ip record. If the record has been updated, then return true, otherwise false.
// If ip has not been in database previous, record is recogonized as updated. If ip has some changes, it
// will also update. If recorded ip is same as inserting one, it will not update.
func (ip *IP) Write(db *gorm.DB) (bool, error) {
	if ip.URL == "" {
		return false, fmt.Errorf("invalid empty url")
	}
	prev := &IP{
		URL: ip.URL,
	}
	if err := db.First(prev).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			return false, err
		}
		return true, db.Create(ip).Error
	}
	if prev.Equal(ip) {
		// no update, pass
		return false, nil
	}
	return true, db.Save(ip).Error
}

// LoadFrom will load the ip record with given url from given database. If record not exist,
// it will create an enpty one with initialized url
func (ip *IP) LoadFrom(db *gorm.DB, url string) error {
	ip.URL = url
	return db.FirstOrCreate(ip).Error
}

// IPLocal will change IP fetched from local end system
func (ip *IP) IPLocal(local string) {
	ip.LocalIP = local
}

// IPADDRCOM will change IP fetched from ipaddress.com, saved to Raw JSON after sort
func (ip *IP) IPADDRCOM(ips []string) {
	// sort
	sort.Strings(ips)
	d, _ := json.Marshal(ips)
	ip.ADDRCOMIP = datatypes.JSON(d)
}

// Equal returns whether ip is equal with another one
func (ip *IP) Equal(another *IP) bool {
	return ip.URL == another.URL &&
		ip.LocalIP == another.LocalIP &&
		reflect.DeepEqual(ip.ADDRCOMIP, another.ADDRCOMIP)
}
