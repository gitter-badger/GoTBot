package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/3stadt/GoTBot/src/structs"
	"encoding/json"
	"github.com/imdario/mergo"
	"github.com/3stadt/GoTBot/src/globals"
	"fmt"
)

var db *bolt.DB

func CreateOrUpdateUser(updateUser structs.User) error {
	baseUser := GetUser(updateUser.Name)
	if baseUser != nil {
		if err := mergo.MergeWithOverwrite(baseUser, &updateUser); err != nil {
			panic(err)
		}
	} else {
		baseUser = &updateUser
	}
	open()
	dberr := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(globals.UserbucketName))
		err := b.Put([]byte(baseUser.Name), marshalUser(*baseUser))
		return err
	})
	db.Close()
	return dberr
}

func GetUser(username string) *structs.User {
	open()
	var user *structs.User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(globals.UserbucketName))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(username))
		var err error
		err, user = unmarshalUser(v)
		return err
	})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	db.Close()
	return user
}

func marshalUser(user structs.User) []byte {
	jUser, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	return jUser
}

func unmarshalUser(bytes []byte) (error, *structs.User) {
	var objmap map[string]*json.RawMessage
	var user *structs.User
	if len(bytes) == 0 {
		return nil, nil
	}
	if err:= json.Unmarshal(bytes, &objmap); err != nil {
		return err, nil
	}
	json.Unmarshal(*objmap["Name"], &user)
	json.Unmarshal(*objmap["LastJoin"], &user)
	json.Unmarshal(*objmap["LastPart"], &user)
	json.Unmarshal(*objmap["LastActive"], &user)
	json.Unmarshal(*objmap["FirstSeen"], &user)
	return nil, user
}

func open() {
	var err error
	db, err = bolt.Open("gotbot.db", 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(globals.UserbucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
