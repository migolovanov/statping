package database

import (
	"fmt"
	"github.com/hunterlong/statping/types"
	"time"
)

type CheckinObj struct {
	*types.Checkin
	o *Object

	Checkiner
}

type Checkiner interface {
	Hits() *CheckinHitObj
	Failures() *FailureObj
	Model() *types.Checkin
	Object() *CheckinObj
}

func Checkin(id int64) (*CheckinObj, error) {
	var checkin types.Checkin
	query := database.Checkins().Where("id = ?", id)
	finder := query.Find(&checkin)
	return &CheckinObj{Checkin: &checkin, o: wrapObject(id, &checkin, query)}, finder.Error()
}

func CheckinByKey(api string) (*CheckinObj, error) {
	var checkin types.Checkin
	query := database.Checkins().Where("api = ?", api)
	finder := query.Find(&checkin)
	return &CheckinObj{Checkin: &checkin, o: wrapObject(checkin.Id, &checkin, query)}, finder.Error()
}

func wrapCheckins(all []*types.Checkin, db Database) []*CheckinObj {
	var arr []*CheckinObj
	for _, v := range all {
		arr = append(arr, &CheckinObj{Checkin: v, o: wrapObject(v.Id, v, db)})
	}
	return arr
}

func AllCheckins() []*CheckinObj {
	var checkins []*types.Checkin
	query := database.Checkins()
	query.Find(&checkins)
	return wrapCheckins(checkins, query)
}

func (s *CheckinObj) Service() *ServiceObj {
	var srv *types.Service
	q := database.Checkins().Where("service = ?", s.ServiceId)
	q.Find(&srv)
	return &ServiceObj{
		Service: srv,
		o:       wrapObject(srv.Id, srv, q),
	}
}

func (s *CheckinObj) Failures() *FailureObj {
	q := database.Failures().
		Where("method = 'checkin' AND id = ?", s.Id).
		Where("method = 'checkin'")
	return &FailureObj{wrapObject(s.Id, nil, q)}
}

func (s *CheckinObj) object() *Object {
	return s.o
}

func (c *CheckinObj) Model() *types.Checkin {
	return c.Checkin
}

func (c *CheckinObj) Object() *CheckinObj {
	return c
}

// Period will return the duration of the Checkin interval
func (c *CheckinObj) Period() time.Duration {
	duration, _ := time.ParseDuration(fmt.Sprintf("%vs", c.Interval))
	return duration
}

// Grace will return the duration of the Checkin Grace Period (after service hasn't responded, wait a bit for a response)
func (c *CheckinObj) Grace() time.Duration {
	duration, _ := time.ParseDuration(fmt.Sprintf("%vs", c.GracePeriod))
	return duration
}

// Expected returns the duration of when the serviec should receive a Checkin
func (c *CheckinObj) Expected() time.Duration {
	last := c.Hits().Last()
	now := time.Now().UTC()
	lastDir := now.Sub(last.CreatedAt)
	sub := time.Duration(c.Period() - lastDir)
	return sub
}

// Last returns the last checkinHit for a Checkin
func (c *CheckinObj) Hits() *CheckinHitObj {
	var checkinHits []*types.CheckinHit
	query := database.CheckinHits().Where("checkin = ?", c.Id)
	query.Find(&checkinHits)
	return &CheckinHitObj{checkinHits, wrapObject(c.Id, checkinHits, query)}
}

// Last returns the last checkinHit for a Checkin
func (c *CheckinHitObj) Last() *types.CheckinHit {
	var last types.CheckinHit
	c.o.db.Last(&last)
	return &last
}

func (c *CheckinObj) Link() string {
	return fmt.Sprintf("%v/checkin/%v", "DOMAINHERE", c.ApiKey)
}