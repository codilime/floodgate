package auth

import "testing"

func TestConfig_IsValid(t *testing.T) {
	type args struct {
		config Config
	}

	tests := []struct {
		name     string
		args     args
		wantBool bool
	}{
		{
			name: "basic is enabled",
			args: args{
				config: Config{
					Basic: Basic{
						Enabled: true,
					},
					OAuth2: OAuth2{
						Enabled: false,
					},
					X509: X509{
						Enabled: false,
					},
				},
			},
			wantBool: true,
		},
		{
			name: "oauth2 is enabled",
			args: args{
				config: Config{
					Basic: Basic{
						Enabled: false,
					},
					OAuth2: OAuth2{
						Enabled: true,
					},
					X509: X509{
						Enabled: false,
					},
				},
			},
			wantBool: true,
		},
		{
			name: "x509 auth is enabled",
			args: args{
				config: Config{
					Basic: Basic{
						Enabled: false,
					},
					OAuth2: OAuth2{
						Enabled: false,
					},
					X509: X509{
						Enabled: true,
					},
				},
			},
			wantBool: true,
		},
		{
			name: "basic and oauth2 is enabled",
			args: args{
				config: Config{
					Basic: Basic{
						Enabled: true,
					},
					OAuth2: OAuth2{
						Enabled: true,
					},
					X509: X509{
						Enabled: false,
					},
				},
			},
			wantBool: false,
		},
		{
			name: "oauth2 and x509 is enabled",
			args: args{
				config: Config{
					Basic: Basic{
						Enabled: false,
					},
					OAuth2: OAuth2{
						Enabled: true,
					},
					X509: X509{
						Enabled: true,
					},
				},
			},
			wantBool: false,
		},
		{
			name: "basic and x509 is enabled",
			args: args{
				config: Config{
					Basic: Basic{
						Enabled: true,
					},
					OAuth2: OAuth2{
						Enabled: false,
					},
					X509: X509{
						Enabled: true,
					},
				},
			},
			wantBool: false,
		},
		{
			name: "auth is not configured",
			args: args{
				config: Config{},
			},
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.args.config.IsValid()
			if valid != tt.wantBool {
				t.Errorf("Config.IsValid() want = %v, got %v", tt.wantBool, valid)
			}
		})
	}
}
