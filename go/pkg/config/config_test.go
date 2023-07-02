package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	type args struct {
		environmentVariables map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "Test LoadConfig",
			args: args{
				environmentVariables: map[string]string{
					"BOOKKEEPER_DATABASE_HOST":     "localhost",
					"BOOKKEEPER_DATABASE_PORT":     "5432",
					"BOOKKEEPER_DATABASE_USER":     "postgres",
					"BOOKKEEPER_DATABASE_PASSWORD": "postgres",
					"BOOKKEEPER_DATABASE_NAME":     "bookkeeper",
					"BOOKKEEPER_DATABASE_SSLMODE":  "disabled",
				},
			},
			wantConfig: &Config{
				DatabaseConfig: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "postgres",
					Password: "postgres",
					Name:     "bookkeeper",
					SSLMode:  "disabled",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.args.environmentVariables {
				os.Setenv(key, value)
			}
			gotConfig, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfig.DatabaseConfig.Host, tt.wantConfig.DatabaseConfig.Host) {
				t.Errorf("LoadConfig() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
