package babynamer

import (
	"github.com/AndyNortrup/baby-namer/settings"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"net/http"
)

const partnerEmailFormID = "PartnerEmail"

type SettingsPage struct {
	Partner      *settings.Partner
	PartnerEmail string
	ctx          context.Context
}

func NewSettingsPage(ctx context.Context) *SettingsPage {
	page := &SettingsPage{ctx: ctx}
	page.Partner, _ = settings.GetPartner(ctx, user.Current(ctx))
	return page
}

func (sp *SettingsPage) Render(w http.ResponseWriter) {
	temp, err := getSettingTemplate()
	if err != nil {
		log.Errorf(sp.ctx, "action=SettingsPage.Render error=%v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	temp.Execute(w, *sp)
}

func handleSettingsPage(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := r.ParseForm()

	if err != nil {
		log.Errorf(ctx, "action=handleSettingsPage error=%v", err)
	}

	partner := r.PostFormValue(partnerEmailFormID)
	log.Infof(ctx, "Partner: %v", r.FormValue(partnerEmailFormID))

	if partner != "" {
		err := settings.SetPartner(ctx, user.Current(ctx), partner)
		if err != nil {
			log.Errorf(ctx, "action=handleSettingsPage error=%v", err)
		}
	}

	sp := NewSettingsPage(ctx)
	sp.PartnerEmail = r.PostFormValue(partnerEmailFormID)
	sp.Render(w)
}
