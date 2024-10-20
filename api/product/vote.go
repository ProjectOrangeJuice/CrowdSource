package product

import (
	"context"
	"time"

	"../user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vote struct {
	ID          string
	Name        int
	Ingredients int
	Nutrition   int
	Serving     int
}

type VoteCheck struct {
	part     string
	barcode  string
	version  int64
	username string
	conn     *mongo.Database
}

func canVote(v []UserVote, user string) bool {
	for _, a := range v {
		if a.User == user {
			return true
		}
	}
	return false
}

func VoteOnProduct(v Vote, username string, conn *mongo.Database) {
	p := GetProductInfo(v.ID, username, conn)
	sec := time.Now().Unix()
	vc := VoteCheck{"NAME", v.ID, p.ProductName.Stamp, username, conn}
	if canVote(vc) {
		if v.Name > 0 {
			p.ProductName.Up++
			point := user.Point{p.ID, p.ProductName.Stamp, "NAMEUP", 1, false, sec}
			user.AddPoint(point, username, conn)
		} else if v.Name < 0 {
			p.ProductName.Down--
			point := user.Point{p.ID, p.ProductName.Stamp, "NAMEDOWN", 1, false, sec}
			user.AddPoint(point, username, conn)
		}
	}
	vc = VoteCheck{"INGREDIENTS", v.ID, p.Ingredients.Stamp, username, conn}
	if canVote(vc) {

		if v.Ingredients > 0 && canVote(vc) {
			p.Ingredients.Up++
			point := user.Point{p.ID, p.Ingredients.Stamp, "INGREDIENTSUP", 1, false, sec}
			user.AddPoint(point, username, conn)
		} else if v.Ingredients < 0 {
			p.Ingredients.Down--
			point := user.Point{p.ID, p.Ingredients.Stamp, "INGREDIENTSDOWN", 1, false, sec}
			user.AddPoint(point, username, conn)
		}
	}
	vc = VoteCheck{"NUTRITION", v.ID, p.Nutrition.Stamp, username, conn}
	if canVote(vc) {

		if v.Nutrition > 0 && canVote(vc) {
			p.Nutrition.Up++
			point := user.Point{p.ID, p.Nutrition.Stamp, "NUTRITIONUP", 1, false, sec}
			user.AddPoint(point, username, conn)
		} else if v.Nutrition < 0 {
			p.Nutrition.Down--
			point := user.Point{p.ID, p.ProductName.Stamp, "NUTRITIONDOWN", 1, false, sec}
			user.AddPoint(point, username, conn)
		}
	}
	collection := conn.Collection("products")
	filter := bson.M{"_id": p.ID}
	collection.FindOneAndReplace(context.TODO(), filter, p)

}
