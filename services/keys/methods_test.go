package keys

import (
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"
	"os"
	"testing"
)

func TestStore_RecoverAddressFromMnemonic(t *testing.T) {
	appconf.LoadTestConf()

	// load sentinel key storage
	if err := Load(appconf.Paths.SentinelPath()); err != nil {
		panic(err)
	}

	defer os.RemoveAll(appconf.Paths.SentinelPath())

	type fields struct {
		homeDir string
	}
	type args struct {
		mnemonic string
		name     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "recover address from mnemonic successfully",
			fields: struct{ homeDir string }{homeDir: appconf.Paths.SentinelDir},
			args: args{
				mnemonic: "identify tree express horse insect sure vendor remove spare bracket average chuckle tube actor habit system clock gas virtual motion use hero afford come",
				name:     "test1",
			},
			want:    "sent1sduag3y5dv4ffry3fy0sk8kahfsnc456n7j25v",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Wallet.RecoverAddressFromMnemonic(tt.args.mnemonic, tt.args.name)

			if (err != nil) != tt.wantErr {
				t.Errorf("RecoverAddressFromMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecoverAddressFromMnemonic() got = %v, want %v", got, tt.want)
			}
		})
	}
}
