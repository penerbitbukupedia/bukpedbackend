package hub

import (
	"context"
	"time"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HubHandler(Profile itmodel.Profile, msg itmodel.IteungMessage, db *mongo.Database) string {
	//check apakah ada sesssion hub
	var shub SessionHub
	shub, res, err := CheckHubSessionUser(msg.Phone_number, db)
	if err != nil {
		return err.Error()
	}
	if !res {
		shub, res, err := CheckHubSessionAdmin(msg.Phone_number, db)
		if err != nil {
			return err.Error()
		}
		if res { //klo session admin masih ada
			msgstr := msg.Message + "\n> _" + shub.AdminName + "_"
			dt := &itmodel.TextMessage{
				To:       shub.UserPhone,
				IsGroup:  false,
				Messages: msgstr,
			}
			go atapi.PostStructWithToken[itmodel.Response]("Token", Profile.Token, dt, Profile.URLAPIText)
			//return "> _ⓘ terkirim ke:" + shub.UserPhone + "_"
			return ""
		}
		//kalo ga ada session
		return ""
	}
	//kalo session user masih ada
	msgstr := msg.Message + "\n> _" + shub.UserName + "_"
	dt := &itmodel.TextMessage{
		To:       shub.AdminPhone,
		IsGroup:  false,
		Messages: msgstr,
	}
	go atapi.PostStructWithToken[itmodel.Response]("Token", Profile.Token, dt, Profile.URLAPIText)
	//return "> _ⓘ terkirim ke:" + shub.AdminPhone + "_"
	return ""

}

// check session hub buat baru atau refresh session
func CheckHubSessionUser(userphone string, db *mongo.Database) (session SessionHub, result bool, err error) {
	session, err = atdb.GetOneDoc[SessionHub](db, "hub", bson.M{"userphone": userphone})
	session.CreatedAt = time.Now()
	if err != nil { //insert session klo belum ada
		return
	} else { //jika sesssion udah ada
		//refresh waktu session dengan waktu sekarang
		_, err = atdb.DeleteManyDocs(db, "hub", bson.M{"userphone": userphone})
		if err != nil {
			return
		}
		_, err = db.Collection("hub").InsertOne(context.TODO(), session)
		if err != nil {
			return
		}
		result = true
	}
	return
}

// check session hub buat baru atau refresh session
func CheckHubSessionAdmin(adminphone string, db *mongo.Database) (session SessionHub, result bool, err error) {
	session, err = atdb.GetOneDoc[SessionHub](db, "hub", bson.M{"adminphone": adminphone})
	session.CreatedAt = time.Now()
	if err != nil { //return session klo belum ada
		return
	} else { //jika sesssion udah ada
		//refresh waktu session dengan waktu sekarang
		_, err = atdb.DeleteManyDocs(db, "hub", bson.M{"adminphone": adminphone})
		if err != nil {
			return
		}
		_, err = db.Collection("hub").InsertOne(context.TODO(), session)
		if err != nil {
			return
		}
		result = true
	}
	return
}

// check session hub buat baru atau refresh session
func CheckHubSession(userphone, username, adminphone, adminname string, db *mongo.Database) (session SessionHub, result bool, err error) {
	session, err = atdb.GetOneDoc[SessionHub](db, "hub", bson.M{"userphone": userphone})
	session.CreatedAt = time.Now()
	if err != nil { //insert session klo belum ada
		session.UserPhone = userphone
		session.UserName = username
		session.AdminPhone = adminphone
		session.AdminName = adminname
		_, err = db.Collection("hub").InsertOne(context.TODO(), session)
		if err != nil {
			return
		}
	} else { //jika sesssion udah ada
		//refresh waktu session dengan waktu sekarang
		_, err = atdb.DeleteManyDocs(db, "hub", bson.M{"userphone": userphone})
		if err != nil {
			return
		}
		_, err = db.Collection("hub").InsertOne(context.TODO(), session)
		if err != nil {
			return
		}
		result = true
	}
	return
}