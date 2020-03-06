// Statping
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statping
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package notifications

import (
	"fmt"
	"github.com/hunterlong/statping/database"
	"github.com/hunterlong/statping/types/failures"
	"github.com/hunterlong/statping/types/null"
	"github.com/hunterlong/statping/types/services"
	"github.com/hunterlong/statping/types/users"
	"github.com/hunterlong/statping/utils"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	dir    string
	METHOD = "example"
)

var service = &services.Service{
	Name:               "Interpol - All The Rage Back Home",
	Domain:             "https://www.youtube.com/watch?v=-u6DvRyyKGU",
	ExpectedStatus:     200,
	Interval:           30,
	Type:               "http",
	Method:             "GET",
	Timeout:            20,
	AllowNotifications: null.NewNullBool(true),
}

var failure = &failures.Failure{
	Issue: "testing",
}

var user = &users.User{
	Username: "admin",
	Email:    "info@email.com",
}

func injectDatabase() {
	sqlPath := dir + "/notifier.db"
	utils.DeleteFile(sqlPath)
	db, _ = database.Openw("sqlite3", sqlPath)
	db.CreateTable(&Notification{})
}

func TestIsBasicType(t *testing.T) {
	assert.True(t, utils.IsType(example, new(Notifier)))
	assert.True(t, utils.IsType(example, new(BasicEvents)))
	assert.True(t, utils.IsType(example, new(ServiceEvents)))
	assert.True(t, utils.IsType(example, new(UserEvents)))
	assert.True(t, utils.IsType(example, new(NotifierEvents)))
	assert.True(t, utils.IsType(example, new(Tester)))
}

func TestSelectNotification(t *testing.T) {
	notifier, notif, err := SelectNotifier(example.Method)
	assert.Nil(t, err)
	assert.NotNil(t, notifier)
	assert.NotNil(t, notif)
	assert.Equal(t, "example", notifier.Method)
	assert.False(t, notifier.Enabled.Bool)
	assert.False(t, notifier.IsRunning())
}

func TestAddQueue(t *testing.T) {
	msg := "this is a test in the queue!"
	example.AddQueue(fmt.Sprintf("service_%v", 0), msg)
	assert.Equal(t, 1, len(example.Queue))
}

func TestNotification_Update(t *testing.T) {
	notifier, err := SelectNotification(example)
	assert.Nil(t, err)
	notifier.Host = "http://demo.statping.com/api"
	notifier.Port = 9090
	notifier.Username = "admin"
	notifier.Password = "password123"
	notifier.Var1 = "var1_is_here"
	notifier.Var2 = "var2_is_here"
	notifier.ApiKey = "USBdu82HDiiuw9327yGYDGw"
	notifier.ApiSecret = "PQopncow929hUIDHGwiud"
	notifier.Limits = 10
	err = notifier.Update()
	assert.Nil(t, err)

	selected, err := SelectNotification(example)
	assert.Nil(t, err)
	assert.Equal(t, "http://demo.statping.com/api", selected.GetValue("host"))
	assert.Equal(t, "http://demo.statping.com/api", example.Notification.Host)
	assert.Equal(t, "http://demo.statping.com/api", example.Host)
	assert.Equal(t, "USBdu82HDiiuw9327yGYDGw", selected.GetValue("api_key"))
	assert.Equal(t, "USBdu82HDiiuw9327yGYDGw", example.ApiKey)
	assert.False(t, selected.Enabled.Bool)
	assert.False(t, selected.IsRunning())
}

func TestEnableNotification(t *testing.T) {
	notifier, err := SelectNotification(example)
	assert.Nil(t, err)
	notifier.Enabled = null.NewNullBool(true)
	err = notifier.Update()
	assert.Nil(t, err)
	assert.True(t, notifier.Enabled.Bool)
	assert.True(t, notifier.IsRunning())
}

func TestIsEnabled(t *testing.T) {
	assert.True(t, isEnabled(example))
}

func TestIsRunning(t *testing.T) {
	assert.True(t, example.IsRunning())
}

func TestLastSent(t *testing.T) {
	notifier, err := SelectNotification(example)
	assert.Nil(t, err)
	assert.Equal(t, "0s", notifier.LastSent().String())
}

func TestWithinLimits(t *testing.T) {
	notifier, err := SelectNotification(example)
	assert.Nil(t, err)
	assert.Equal(t, 10, notifier.Limits)
	assert.True(t, inLimits(example))
}

func TestNotification_GetValue(t *testing.T) {
	notifier, err := SelectNotification(example)
	assert.Nil(t, err)
	val := notifier.GetValue("Host")
	assert.Equal(t, "http://demo.statping.com/api", val)
}

func TestOnSave(t *testing.T) {
	err := example.OnSave()
	assert.Equal(t, "onsave triggered", err.Error())
}

func TestOnSuccess(t *testing.T) {
	OnSuccess(service)
	assert.Equal(t, 2, len(example.Queue))
}

func TestOnFailure(t *testing.T) {
	OnFailure(service, failure)
	assert.Equal(t, 3, len(example.Queue))
}

func TestOnNewService(t *testing.T) {
	OnNewService(service)
	assert.Equal(t, 4, len(example.Queue))
}

func TestOnUpdatedService(t *testing.T) {
	OnUpdatedService(service)
	assert.Equal(t, 5, len(example.Queue))
}

func TestOnDeletedService(t *testing.T) {
	OnDeletedService(service)
	assert.Equal(t, 6, len(example.Queue))
}

func TestOnNewUser(t *testing.T) {
	OnNewUser(user)
	assert.Equal(t, 7, len(example.Queue))
}

func TestOnUpdatedUser(t *testing.T) {
	OnUpdatedUser(user)
	assert.Equal(t, 8, len(example.Queue))
}

func TestOnDeletedUser(t *testing.T) {
	OnDeletedUser(user)
	assert.Equal(t, 9, len(example.Queue))
}

func TestOnUpdatedNotifier(t *testing.T) {
	OnUpdatedNotifier(example.Select())
	assert.Equal(t, 11, len(example.Queue))
}

func TestRunAllQueueAndStop(t *testing.T) {
	assert.True(t, example.IsRunning())
	assert.Equal(t, 11, len(example.Queue))
	go Queue(example)
	time.Sleep(13 * time.Second)
	assert.NotZero(t, len(example.Queue))
	example.Close()
	assert.False(t, example.IsRunning())
	assert.NotZero(t, len(example.Queue))
}