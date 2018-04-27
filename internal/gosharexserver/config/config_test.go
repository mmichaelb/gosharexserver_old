package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestMainConfig(t *testing.T) {
	err := loadViper("../../../test/test-config.toml")
	if err != nil {
		t.Fatalf("Could not load test config file, %T: %v", err, err)
	}
	if webServerAddress := viper.GetString("webserver.address"); webServerAddress != ":80" {
		t.Fatalf(`Invalid value for "webserver.webserver_address": %s`, strconv.Quote(webServerAddress))
	}
	if reverseProxyHeader := viper.GetString("webserver.reverse_proxy_header"); reverseProxyHeader != "This-Header-Contains-The-Real-IP" {
		t.Fatalf(`Invalid value for "webserver.reverse_proxy_header": %s`, strconv.Quote(reverseProxyHeader))
	}
	if whitelistedContentTypes := viper.GetStringSlice("webserver.whitelisted_content_types"); !reflect.DeepEqual(whitelistedContentTypes, []string{"first-ct", "a-mime-type", "sp€ci4l"}) {
		t.Fatalf(`Invalid value for "webserver.whitelisted_content_types": %s`, strconv.Quote(fmt.Sprintf("%+v", whitelistedContentTypes)))
	}
	if authorizationToken := viper.GetString("webserver.authorization_token"); authorizationToken != "123456" {
		t.Fatalf(`Invalid value for "webserver.authorization_token": %s`, strconv.Quote(authorizationToken))
	}
	testMongoConfig(t)
}

func testMongoConfig(t *testing.T) {
	if address := viper.GetString("mongodb.address"); address != "0.0.0.0:1337" {
		t.Fatalf(`Invalid value for "mongodb.address": %s`, strconv.Quote(address))
	}
	if connectTimeout := viper.GetDuration("mongodb.connect_timeout"); connectTimeout != time.Minute+time.Second*30 {
		t.Fatalf(`Invalid value for "mongodb.connect_timeout": %s`, strconv.Quote(connectTimeout.String()))
	}
	if authDB := viper.GetString("mongodb.auth_db"); authDB != "sharex-admin-db" {
		t.Fatalf(`Invalid value for "mongodb.auth_db": %s`, strconv.Quote(authDB))
	}
	if authUser := viper.GetString("mongodb.auth_user"); authUser != "l_torvalds" {
		t.Fatalf(`Invalid value for "mongodb.auth_user": %s`, strconv.Quote(authUser))
	}
	if authPasswd := viper.GetString("mongodb.auth_passwd"); authPasswd != "MySuperSecurePassword+!#" {
		t.Fatalf(`Invalid value for "mongodb.auth_passwd": %s`, strconv.Quote(authPasswd))
	}
	if db := viper.GetString("mongodb.db"); db != "sharex-upload-metadata" {
		t.Fatalf(`Invalid value for "mongodb.db": %s`, strconv.Quote(db))
	}
	if gridFSPrefix := viper.GetString("mongodb.gridfs_prefix"); gridFSPrefix != "gridfsisnice" {
		t.Fatalf(`Invalid value for "mongodb.gridfs_prefix": %s`, strconv.Quote(gridFSPrefix))
	}
	if gridFSChunkSize := viper.GetInt("mongodb.grids_chunk_size"); gridFSChunkSize != 1337 {
		t.Fatalf(`Invalid value for "mongodb.grids_chunk_size": %d`, gridFSChunkSize)
	}
}
