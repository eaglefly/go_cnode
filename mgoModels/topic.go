package models

import (
	//"log"
	"log"

	db "github.com/dangyanglim/go_cnode/database"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"time"
)

//"log"

type Topic struct {
	Id              bson.ObjectId `bson:"_id"`
	Title           string        `json:"title"`
	Content         string        `json:"content" `
	Author_id       bson.ObjectId `bson:"author_id" `
	Top             bool          `json:"top" `
	Good            bool          `json:"good" `
	Lock            bool          `json:"lock"`
	Reply_count     uint          `json:"reply_count"`
	Visit_count     uint          `json:"visit_count"`
	Collect_count   uint          `json:"collect_count"`
	Create_at       time.Time        `bson:"create_at"`
	Update_at       string        `json:"update_at"`
	Last_reply      bson.ObjectId          `bson:"last_reply,omitempty"`
	Last_reply_at   time.Time       `json:"last_reply_at,omitempty"`
	Content_is_html bool          `json:"content_is_html"`
	Tab             string        `json:"tab"`
	Deleted         bool          `json:"deleted"`
}
type TopicModel struct{}

var userModel = new(UserModel)
var replyModel = new(ReplyModel)
func (p *TopicModel) GetTopicByQuery(tab string, good bool) (topics []Topic, err error) {
	mgodb := db.MogSession.DB("egg_cnode")
	if tab == "" || tab == "all" {
		err = mgodb.C("topics").Find(bson.M{"good": good}).All(&topics)
	} else {
		err = mgodb.C("topics").Find(bson.M{"tab": tab, "good": good}).All(&topics)
	}

	return topics, err
}
func (p *TopicModel) GetTopicBy(tab string, good bool) (topics []Topic,topicss []byte, err error) {
	type TopciAndAuthor struct{
		Author User `json:"author"`
		Topic Topic `json:"topic"`
		Reply Reply `json:"reply"`
	}
	var temps []TopciAndAuthor 
	mgodb := db.MogSession.DB("egg_cnode")
	if tab == "" || tab == "all" {
		err = mgodb.C("topics").Find(bson.M{"good": good}).All(&topics)
	} else {
		err = mgodb.C("topics").Find(bson.M{"tab": tab, "good": good}).All(&topics)
	}
	for _,v:=range topics{
		var temp TopciAndAuthor
		temp.Topic=v
		author, _ := userModel.GetUserById(v.Author_id.Hex())
		temp.Author=author
		 if v.Last_reply.Hex()!=""{
			reply, _ := replyModel.GetReplyById(v.Last_reply.Hex())
			log.Println("dddd",reply)
			temp.Reply=reply
		 }

		temps=append(temps,temp)
	}
	//log.Println(temps)
	topicss,_=json.Marshal(temps)
	//log.Println(string(topicss))
	return topics,topicss, err
}
func (p *TopicModel) GetTopicByQueryCount(tab string, good bool) (count int, err error) {
	mgodb := db.MogSession.DB("egg_cnode")
	if tab == "" || tab == "all" {
		count,err = mgodb.C("topics").Find(bson.M{"good": good}).Count()
	} else {
		count,err = mgodb.C("topics").Find(bson.M{"tab": tab, "good": good}).Count()
	}

	return count, err
}
func (p *TopicModel) GetTopicById(id string) (topic Topic, author User, err error) {
	mgodb := db.MogSession.DB("egg_cnode")
	objectId := bson.ObjectIdHex(id)

	err = mgodb.C("topics").Find(bson.M{"_id": objectId}).One(&topic)

	author, _ = userModel.GetUserById(topic.Author_id.Hex())
	//log.Println(topic)
	//log.Println(author)

	return topic, author, err
}
func (p *TopicModel) NewAndSave(title string, tab string, id string, content string) ( topic Topic,err error) {

	
	objectId := bson.ObjectIdHex(id)
	topic = Topic{
		Id:          bson.NewObjectId(),
		Title:        title,
		Content:   content,
		Tab:        tab,
		Author_id:objectId,
		Create_at:time.Now(),
	}
	mgodb := db.MogSession.DB("egg_cnode")
	err = mgodb.C("topics").Insert(&topic)
	log.Println(topic)
	log.Println(err)
	return topic,err
}
func (p *TopicModel) GetTopicNoReply() (topics []Topic, err error) {
	mgodb := db.MogSession.DB("egg_cnode")

	err = mgodb.C("topics").Find(bson.M{"reply_count": 0}).Limit(5).All(&topics)
	

	return topics, err
}
