package settings_test

import (
	"github.com/AndyNortrup/baby-namer/settings"
	"google.golang.org/appengine/user"
	"testing"
)

func TestGetSetPartner(t *testing.T) {
	ctx := newTestContext()
	self := &user.User{Email: "self@gmail.com"}
	partnerEmail := "partner@gmail.com"
	err := settings.SetPartner(ctx, self, partnerEmail)
	if err != nil {
		t.Fatalf("action=TestGetSetPartner error=%v", err)
	}

	resp, err := settings.GetPartner(ctx, self)
	if err != nil {
		t.Fatalf("action=TestGetSetPartner error=%v", err)
	}
	if resp == nil {
		t.Log("action=TestGetSetPartner recieved null response.")
		t.FailNow()
	}
	if resp.PartnerEmail != partnerEmail {
		t.Logf("action=TestGetSetPartner wrong email returned "+"expected=%v recieved=%v", partnerEmail, resp.PartnerEmail)
		t.Fail()
	}
}
