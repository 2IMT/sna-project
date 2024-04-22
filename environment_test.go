package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func unsetEnvs() {
	os.Unsetenv("BOT_TOKEN")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
}

func setEnvs() {
	os.Setenv("BOT_TOKEN", "token")
	os.Setenv("DB_HOST", "db_host")
	os.Setenv("DB_PORT", "db_port")
	os.Setenv("DB_NAME", "db_name")
	os.Setenv("DB_USER", "db_user")
	os.Setenv("DB_PASS", "db_pass")
}

func Test(t *testing.T) {
	unsetEnvs()
	setEnvs()

	env, err := LoadEnvironment()
    assert.Nil(t, err);
    
    assert.Equal(t, env.BotToken, "token")
    assert.Equal(t, env.DbHost, "db_host")
    assert.Equal(t, env.DbPort, "db_port")
    assert.Equal(t, env.DbName, "db_name")
    assert.Equal(t, env.DbUser, "db_user")
    assert.Equal(t, env.DbPass, "db_pass")
}
