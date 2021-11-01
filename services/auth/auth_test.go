package auth

import (
	"github.com/google/uuid"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"
	"testing"
)

func TestAuth_Login(t *testing.T) {
	appconf.LoadTestConf()

	type fields struct {
		Token uuid.UUID
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test login success - no password",
			fields: fields{Token: uuid.UUID{}},
			args: args{
				username: "root",
				password: "",
			},
			wantErr: false,
		},
		{
			name:   "test login success - with password",
			fields: fields{Token: uuid.UUID{}},
			args: args{
				username: "john",
				password: "sentinel",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Auth{
				Token: tt.fields.Token,
			}
			if err := s.Login(tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
